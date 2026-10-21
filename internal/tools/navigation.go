package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/bytedance/sonic"
	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/schema"
)

// NavigationTool 导航路径工具
// 提供路径规划、实时导航、交通信息、多种出行方式
type NavigationTool struct {
	BaseToolImpl
	client *http.Client
}

type navigationParams struct {
	Origin      string  `json:"origin"`                // 起点坐标或地址
	Destination string  `json:"destination"`           // 终点坐标或地址
	Strategy    string  `json:"strategy,omitempty"`    // 路径策略: fast(最快), short(最短), economy(经济)
	Waypoints   string  `json:"waypoints,omitempty"`   // 途经点，多个用|分隔
	Mode        string  `json:"mode,omitempty"`        // 出行方式: driving, walking, bicycling, transit
}

type navigationResult struct {
	Origin      string           `json:"origin"`
	Destination string           `json:"destination"`
	Distance    int              `json:"distance"`     // 距离(米)
	Duration    int              `json:"duration"`     // 时间(秒)
	Steps       []navigationStep `json:"steps"`        // 路径步骤
	Strategy    string           `json:"strategy"`
	Mode        string           `json:"mode"`
	PlannedAt   time.Time        `json:"planned_at"`
}

type navigationStep struct {
	Instruction string  `json:"instruction"` // 指导说明
	Distance    int     `json:"distance"`    // 距离(米)
	Duration    int     `json:"duration"`    // 时间(秒)
	Polyline    string  `json:"polyline"`    // 路径编码
	Action      string  `json:"action"`      // 动作: left, right, straight, uturn
}

func NewNavigationTool() *NavigationTool {
	return &NavigationTool{
		BaseToolImpl: BaseToolImpl{
			Category:       "location",
			EnableLearning: true,
			ToolName:       "navigation",
		},
		client: &http.Client{Timeout: 10 * time.Second},
	}
}

func (t *NavigationTool) Info(_ context.Context) (*schema.ToolInfo, error) {
	return &schema.ToolInfo{
		Name: "navigation",
		Desc: "提供路径规划、实时导航、交通信息。支持多种出行方式和路径策略",
		ParamsOneOf: schema.NewParamsOneOfByParams(map[string]*schema.ParameterInfo{
			"origin":       {Desc: "起点坐标(经度,纬度)或地址", Type: schema.String, Required: true},
			"destination":  {Desc: "终点坐标(经度,纬度)或地址", Type: schema.String, Required: true},
			"strategy":     {Desc: "路径策略: fast(最快), short(最短), economy(经济)", Type: schema.String, Required: false},
			"waypoints":    {Desc: "途经点，多个用|分隔", Type: schema.String, Required: false},
			"mode":         {Desc: "出行方式: driving(驾车), walking(步行), bicycling(骑行), transit(公交)", Type: schema.String, Required: false},
		}),
	}, nil
}

func (t *NavigationTool) InvokableRun(ctx context.Context, argumentsInJSON string, _ ...tool.Option) (string, error) {
	p := navigationParams{}
	if err := sonic.UnmarshalString(argumentsInJSON, &p); err != nil {
		return "", err
	}

	// 验证参数
	if p.Origin == "" {
		return "", fmt.Errorf("origin is required")
	}
	if p.Destination == "" {
		return "", fmt.Errorf("destination is required")
	}

	// 设置默认值
	if p.Strategy == "" {
		p.Strategy = "fast"
	}
	if p.Mode == "" {
		p.Mode = "driving"
	}

	// 获取高德地图API密钥
	key, err := GetAMapKey()
	if err != nil {
		return "", err
	}

	// 调用高德地图路径规划API
	result, err := t.planRoute(ctx, key, p)
	if err != nil {
		return "", fmt.Errorf("route planning failed: %w", err)
	}

	b, _ := sonic.Marshal(result)
	return string(b), nil
}

// planRoute 调用高德地图路径规划API
func (t *NavigationTool) planRoute(ctx context.Context, key string, p navigationParams) (*navigationResult, error) {
	// 构建API URL
	var baseURL string
	switch p.Mode {
	case "walking":
		baseURL = "https://restapi.amap.com/v3/direction/walking"
	case "bicycling":
		baseURL = "https://restapi.amap.com/v4/direction/bicycling"
	case "transit":
		baseURL = "https://restapi.amap.com/v3/direction/transit/integrated"
	default: // driving
		baseURL = "https://restapi.amap.com/v3/direction/driving"
	}

	params := url.Values{}
	params.Set("key", key)
	params.Set("origin", p.Origin)
	params.Set("destination", p.Destination)

	if p.Strategy != "" {
		params.Set("strategy", t.convertStrategy(p.Strategy))
	}
	if p.Waypoints != "" {
		params.Set("waypoints", p.Waypoints)
	}

	apiURL := fmt.Sprintf("%s?%s", baseURL, params.Encode())

	req, err := http.NewRequestWithContext(ctx, "GET", apiURL, nil)
	if err != nil {
		return nil, err
	}

	resp, err := t.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned status: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// 根据出行方式解析真实API响应
	switch p.Mode {
	case "bicycling":
		return t.parseBicyclingResponse(body, p)
	case "transit":
		return t.parseTransitResponse(body, p)
	default: // driving, walking
		return t.parseDirectionResponse(body, p)
	}
}

// parseDirectionResponse 解析 v3 驾车/步行响应（数字字段为字符串）
// 注意：高德 action 字段在没有数据时返回 []（空数组），有数据时返回字符串；用 json.RawMessage 容错
func (t *NavigationTool) parseDirectionResponse(body []byte, p navigationParams) (*navigationResult, error) {
	var resp struct {
		Status string `json:"status"`
		Info   string `json:"info"`
		Route  struct {
			Paths []struct {
				Distance string `json:"distance"`
				Duration string `json:"duration"`
				Tolls    string `json:"tolls"`
				Steps    []struct {
					Instruction string          `json:"instruction"`
					Distance    string          `json:"distance"`
					Duration    string          `json:"duration"`
					Action      json.RawMessage `json:"action"`
					Polyline    string          `json:"polyline"`
				} `json:"steps"`
			} `json:"paths"`
		} `json:"route"`
	}
	if err := sonic.Unmarshal(body, &resp); err != nil {
		return nil, err
	}
	if resp.Status != "1" {
		return nil, fmt.Errorf("amap error: %s", resp.Info)
	}
	if len(resp.Route.Paths) == 0 {
		return nil, fmt.Errorf("no route found")
	}
	path := resp.Route.Paths[0]
	steps := make([]navigationStep, 0, len(path.Steps))
	for _, s := range path.Steps {
		steps = append(steps, navigationStep{
			Instruction: s.Instruction,
			Distance:    atoiSafe(s.Distance),
			Duration:    atoiSafe(s.Duration),
			Polyline:    s.Polyline,
			Action:      parseRawAction(s.Action),
		})
	}
	return &navigationResult{
		Origin:      p.Origin,
		Destination: p.Destination,
		Distance:    atoiSafe(path.Distance),
		Duration:    atoiSafe(path.Duration),
		Steps:       steps,
		Strategy:    p.Strategy,
		Mode:        p.Mode,
		PlannedAt:   time.Now(),
	}, nil
}

// parseRawAction 兼容 action 可能是 string 或 [] 的情况
func parseRawAction(raw json.RawMessage) string {
	if len(raw) == 0 {
		return ""
	}
	var s string
	if err := json.Unmarshal(raw, &s); err == nil {
		return s
	}
	return ""
}

// parseBicyclingResponse 解析 v4 骑行响应（数字字段为 int）
func (t *NavigationTool) parseBicyclingResponse(body []byte, p navigationParams) (*navigationResult, error) {
	var resp struct {
		Errcode int    `json:"errcode"`
		Errmsg  string `json:"errmsg"`
		Data    struct {
			Paths []struct {
				Distance int `json:"distance"`
				Duration int `json:"duration"`
				Steps    []struct {
					Instruction string `json:"instruction"`
					Distance    int    `json:"distance"`
					Action      string `json:"action"`
				} `json:"steps"`
			} `json:"paths"`
		} `json:"data"`
	}
	if err := sonic.Unmarshal(body, &resp); err != nil {
		return nil, err
	}
	if resp.Errcode != 0 {
		return nil, fmt.Errorf("amap bicycling error: %s", resp.Errmsg)
	}
	if len(resp.Data.Paths) == 0 {
		return nil, fmt.Errorf("no route found")
	}
	path := resp.Data.Paths[0]
	steps := make([]navigationStep, 0, len(path.Steps))
	for _, s := range path.Steps {
		steps = append(steps, navigationStep{
			Instruction: s.Instruction,
			Distance:    s.Distance,
			Action:      s.Action,
		})
	}
	return &navigationResult{
		Origin:      p.Origin,
		Destination: p.Destination,
		Distance:    path.Distance,
		Duration:    path.Duration,
		Steps:       steps,
		Strategy:    p.Strategy,
		Mode:        p.Mode,
		PlannedAt:   time.Now(),
	}, nil
}

// parseTransitResponse 解析 v3 公交响应
func (t *NavigationTool) parseTransitResponse(body []byte, p navigationParams) (*navigationResult, error) {
	var resp struct {
		Status string `json:"status"`
		Info   string `json:"info"`
		Route  struct {
			Distance   string `json:"distance"`
			Transits   []struct {
				Duration string `json:"duration"`
				Distance string `json:"distance"`
				Segments []struct {
					Walking struct {
						Distance string `json:"distance"`
						Duration string `json:"duration"`
					} `json:"walking"`
					Bus struct {
						Buslines []struct {
							Name     string `json:"name"`
							Distance string `json:"distance"`
							Duration string `json:"duration"`
						} `json:"buslines"`
					} `json:"bus"`
				} `json:"segments"`
			} `json:"transits"`
		} `json:"route"`
	}
	if err := sonic.Unmarshal(body, &resp); err != nil {
		return nil, err
	}
	if resp.Status != "1" {
		return nil, fmt.Errorf("amap transit error: %s", resp.Info)
	}
	if len(resp.Route.Transits) == 0 {
		return nil, fmt.Errorf("no transit route found")
	}
	t0 := resp.Route.Transits[0]
	steps := make([]navigationStep, 0)
	for _, seg := range t0.Segments {
		if d := atoiSafe(seg.Walking.Distance); d > 0 {
			steps = append(steps, navigationStep{
				Instruction: fmt.Sprintf("步行 %d 米", d),
				Distance:    d,
				Duration:    atoiSafe(seg.Walking.Duration),
				Action:      "walking",
			})
		}
		for _, bl := range seg.Bus.Buslines {
			steps = append(steps, navigationStep{
				Instruction: fmt.Sprintf("乘坐 %s", bl.Name),
				Distance:    atoiSafe(bl.Distance),
				Duration:    atoiSafe(bl.Duration),
				Action:      "bus",
			})
		}
	}
	return &navigationResult{
		Origin:      p.Origin,
		Destination: p.Destination,
		Distance:    atoiSafe(t0.Distance),
		Duration:    atoiSafe(t0.Duration),
		Steps:       steps,
		Strategy:    p.Strategy,
		Mode:        p.Mode,
		PlannedAt:   time.Now(),
	}, nil
}

func atoiSafe(s string) int {
	n, _ := strconv.Atoi(s)
	return n
}

// convertStrategy 转换策略参数
func (t *NavigationTool) convertStrategy(strategy string) string {
	switch strategy {
	case "fast":
		return "0" // 速度优先
	case "short":
		return "1" // 距离优先
	case "economy":
		return "2" // 费用优先
	default:
		return "0"
	}
}

// AfterRun 工具调用后的后处理
// 集成学习引擎，保存导航路径到知识库
func (t *NavigationTool) AfterRun(ctx context.Context, result string, err error) {
	// 使用默认的BaseToolImpl后处理
	t.BaseToolImpl.AfterRun(ctx, result, err)
}
