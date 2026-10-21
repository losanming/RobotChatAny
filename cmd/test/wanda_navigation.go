package main

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"new-robot-self/global"
	"new-robot-self/internal/tools"
	amap "new-robot-self/pkg/location/provider/amap"
)

// 使用我们 tools 系统中真实的工具来完成需求
// - around_search: 周边搜索（已对接高德真实 API）
// - navigation: 路径规划（刚刚修复，现在返回真实数据）
func main() {
	if global.AmapKey == "" {
		fmt.Println("高德地图API密钥未配置")
		return
	}
	ctx := context.Background()

	// 用 amap 客户端做地理编码，把地址转成 lng,lat
	cli, err := amap.NewClient(global.AmapKey, 5000)
	if err != nil {
		fmt.Println("创建客户端失败:", err)
		return
	}

	originLoc, originLat, originLng, err := geocode(ctx, cli, "徐州市铜山区保利学府", "徐州")
	if err != nil {
		fmt.Println("起点解析失败:", err)
		return
	}
	destLoc, destLat, destLng, err := geocode(ctx, cli, "徐州市铜山万达", "徐州")
	if err != nil {
		fmt.Println("终点解析失败:", err)
		return
	}
	_ = originLat
	_ = originLng

	fmt.Printf("起点: 徐州市铜山区保利学府 -> %s\n", originLoc)
	fmt.Printf("终点: 徐州市铜山万达 -> %s\n\n", destLoc)

	// === 1. 用 around_search 工具搜索周边 ===
	aroundTool := tools.NewAroundSearchTool()

	fmt.Println("=== [around_search 工具] 铜山万达 2km 内娱乐设施 ===")
	callAround(ctx, aroundTool, destLat, destLng, "娱乐|电影|KTV", 2000, 10)

	fmt.Println("\n=== [around_search 工具] 铜山万达 1km 内餐饮 ===")
	callAround(ctx, aroundTool, destLat, destLng, "餐饮", 1000, 8)

	fmt.Println("\n=== [around_search 工具] 铜山万达 1km 内购物 ===")
	callAround(ctx, aroundTool, destLat, destLng, "购物|商场|超市", 1000, 8)

	// === 2. 用 navigation 工具规划路线（现在返回真实数据）===
	navTool := tools.NewNavigationTool()

	fmt.Println("\n=== [navigation 工具] 驾车路线 ===")
	callNav(ctx, navTool, originLoc, destLoc, "driving")

	fmt.Println("\n=== [navigation 工具] 步行路线 ===")
	callNav(ctx, navTool, originLoc, destLoc, "walking")

	fmt.Println("\n=== [navigation 工具] 骑行路线 ===")
	callNav(ctx, navTool, originLoc, destLoc, "bicycling")

	fmt.Println("\n=== [navigation 工具] 公交路线 ===")
	callNav(ctx, navTool, originLoc, destLoc, "transit")
}

func geocode(ctx context.Context, cli *amap.Client, addr, city string) (string, float64, float64, error) {
	r, err := cli.Geocode(ctx, amap.GeoParams{Address: addr, City: city})
	if err != nil {
		return "", 0, 0, err
	}
	if len(r.Items) == 0 {
		return "", 0, 0, fmt.Errorf("没有结果")
	}
	loc := r.Items[0].Location
	parts := strings.Split(loc, ",")
	lng, _ := strconv.ParseFloat(strings.TrimSpace(parts[0]), 64)
	lat, _ := strconv.ParseFloat(strings.TrimSpace(parts[1]), 64)
	return loc, lat, lng, nil
}

func callAround(ctx context.Context, tool *tools.AroundSearchTool, lat, lng float64, keywords string, radius, pageSize int) {
	args := fmt.Sprintf(`{"lat":%f,"lng":%f,"keywords":%q,"radius":%d,"page_size":%d}`, lat, lng, keywords, radius, pageSize)
	out, err := tool.InvokableRun(ctx, args)
	if err != nil {
		fmt.Println("调用失败:", err)
		return
	}
	// 解析返回的 JSON
	var r struct {
		Count int `json:"count"`
		POIs  []struct {
			Name    string `json:"name,omitempty"`
			Address string `json:"address,omitempty"`
			Type    string `json:"type,omitempty"`
		} `json:"pois"`
	}
	if err := jsonUnmarshal(out, &r); err != nil {
		fmt.Println("解析失败:", err)
		fmt.Println("原始返回:", out)
		return
	}
	fmt.Printf("找到 %d 个 (展示前%d个)\n", r.Count, len(r.POIs))
	for i, p := range r.POIs {
		fmt.Printf("  %d. %s - %s\n", i+1, p.Name, p.Address)
	}
}

func callNav(ctx context.Context, tool *tools.NavigationTool, origin, dest, mode string) {
	args := fmt.Sprintf(`{"origin":%q,"destination":%q,"mode":%q}`, origin, dest, mode)
	out, err := tool.InvokableRun(ctx, args)
	if err != nil {
		fmt.Println("调用失败:", err)
		return
	}
	var r struct {
		Distance int `json:"distance"`
		Duration int `json:"duration"`
		Steps    []struct {
			Instruction string `json:"instruction"`
			Distance    int    `json:"distance"`
		} `json:"steps"`
	}
	if err := jsonUnmarshal(out, &r); err != nil {
		fmt.Println("解析失败:", err)
		fmt.Println("原始返回:", out)
		return
	}
	fmt.Printf("距离: %.2f 公里, 预计耗时: %d 分钟\n", float64(r.Distance)/1000, r.Duration/60)
	if len(r.Steps) > 0 {
		fmt.Println("路段指引:")
		for i, s := range r.Steps {
			if i >= 5 {
				fmt.Printf("  ... 还有 %d 个步骤\n", len(r.Steps)-5)
				break
			}
			fmt.Printf("  %d. %s (%d米)\n", i+1, s.Instruction, s.Distance)
		}
	}
}

func jsonUnmarshal(s string, v interface{}) error {
	return json.Unmarshal([]byte(s), v)
}
