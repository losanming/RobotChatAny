package visualization

import (
	"context"
	"fmt"
	"time"

	"github.com/bytedance/sonic"
	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/schema"
)

// DataVisualizerTool 数据可视化工具
// 支持健康数据图表展示、趋势分析
type DataVisualizerTool struct {
	EnableLearning bool
	Category       string
	ToolName       string
}

// GetCategory 获取工具分类
func (t *DataVisualizerTool) GetCategory() string {
	return t.Category
}

// GetToolName 获取工具名称
func (t *DataVisualizerTool) GetToolName() string {
	return t.ToolName
}

// AfterRun 工具调用后的后处理
func (t *DataVisualizerTool) AfterRun(ctx context.Context, result string, err error) {
	// 简化实现，不进行学习集成以避免循环依赖
}

type dataVisualizerParams struct {
	Action     string  `json:"action"`       // generate_chart, generate_report
	ChartType  string  `json:"chart_type,omitempty"`  // line, bar, pie, scatter
	DataType   string  `json:"data_type,omitempty"`   // blood_sugar, weight, exercise, diet
	StartDate  int64   `json:"start_date,omitempty"`  // 开始时间(时间戳)
	EndDate    int64   `json:"end_date,omitempty"`    // 结束时间(时间戳)
	Resolution string  `json:"resolution,omitempty"`  // 数据分辨率: daily, weekly, monthly
}

type chartData struct {
	Type        string         `json:"type"`
	Title       string         `json:"title"`
	XAxis       axisConfig     `json:"x_axis"`
	YAxis       axisConfig     `json:"y_axis"`
	Datasets    []dataset      `json:"datasets"`
	Options     chartOptions  `json:"options"`
	GeneratedAt int64         `json:"generated_at"`
}

type axisConfig struct {
	Label    string   `json:"label"`
	Format   string   `json:"format"`
	Category []string `json:"category,omitempty"`
	Min      float64  `json:"min,omitempty"`
	Max      float64  `json:"max,omitempty"`
}

type dataset struct {
	Label        string    `json:"label"`
	Data         []float64 `json:"data"`
	BorderColor  string    `json:"border_color"`
	BackgroundColor string  `json:"background_color"`
	Fill         bool      `json:"fill"`
}

type chartOptions struct {
	Responsive  bool     `json:"responsive"`
	MaintainAspectRatio bool `json:"maintain_aspect_ratio"`
	Plugins     []string `json:"plugins"`
}

type visualReport struct {
	Title       string      `json:"title"`
	Period      string      `json:"period"`
	Summary     summaryData `json:"summary"`
	Charts      []chartData `json:"charts"`
	Insights    []string    `json:"insights"`
	Recommendations []string `json:"recommendations"`
	GeneratedAt int64       `json:"generated_at"`
}

type summaryData struct {
	TotalDataPoints int     `json:"total_data_points"`
	AverageValue   float64 `json:"average_value"`
	MinValue       float64 `json:"min_value"`
	MaxValue       float64 `json:"max_value"`
	Trend          string  `json:"trend"` // improving, stable, declining
	TrendPercent  float64 `json:"trend_percent"`
}

func NewDataVisualizerTool() *DataVisualizerTool {
	return &DataVisualizerTool{
		Category:       "visualization",
		EnableLearning: true,
		ToolName:       "data_visualizer",
	}
}

func (t *DataVisualizerTool) Info(_ context.Context) (*schema.ToolInfo, error) {
	return &schema.ToolInfo{
		Name: "data_visualizer",
		Desc: "数据可视化工具，支持健康数据图表展示、趋势分析和报告生成",
		ParamsOneOf: schema.NewParamsOneOfByParams(map[string]*schema.ParameterInfo{
			"action":     {Desc: "操作类型: generate_chart(生成图表), generate_report(生成报告)", Type: schema.String, Required: true},
			"chart_type":  {Desc: "图表类型: line(折线图), bar(柱状图), pie(饼图), scatter(散点图)", Type: schema.String, Required: false},
			"data_type":  {Desc: "数据类型: blood_sugar(血糖), weight(体重), exercise(运动), diet(饮食)", Type: schema.String, Required: false},
			"start_date": {Desc: "开始时间(时间戳)", Type: schema.Integer, Required: false},
			"end_date":   {Desc: "结束时间(时间戳)", Type: schema.Integer, Required: false},
			"resolution": {Desc: "数据分辨率: daily, weekly, monthly", Type: schema.String, Required: false},
		}),
	}, nil
}

func (t *DataVisualizerTool) InvokableRun(ctx context.Context, argumentsInJSON string, _ ...tool.Option) (string, error) {
	p := dataVisualizerParams{}
	if err := sonic.UnmarshalString(argumentsInJSON, &p); err != nil {
		return "", err
	}

	// 验证参数
	if p.Action == "" {
		return "", fmt.Errorf("action is required")
	}

	// 设置默认时间范围（最近30天）
	if p.StartDate == 0 {
		p.StartDate = time.Now().AddDate(0, 0, -30).Unix()
	}
	if p.EndDate == 0 {
		p.EndDate = time.Now().Unix()
	}

	// 设置默认分辨率
	if p.Resolution == "" {
		p.Resolution = "daily"
	}

	// 根据操作类型执行相应功能
	switch p.Action {
	case "generate_chart":
		return t.generateChart(p)
	case "generate_report":
		return t.generateReport(p)
	default:
		return "", fmt.Errorf("unknown action: %s", p.Action)
	}
}

// generateChart 生成图表
func (t *DataVisualizerTool) generateChart(p dataVisualizerParams) (string, error) {
	if p.DataType == "" {
		return "", fmt.Errorf("data_type is required for generate_chart action")
	}
	if p.ChartType == "" {
		p.ChartType = "line" // 默认折线图
	}

	// 根据数据类型生成图表数据
	chart, err := t.generateChartData(p)
	if err != nil {
		return "", err
	}

	b, _ := sonic.Marshal(chart)
	return string(b), nil
}

// generateReport 生成报告
func (t *DataVisualizerTool) generateReport(p dataVisualizerParams) (string, error) {
	// 设置默认报告时间范围（最近30天）
	if p.StartDate == 0 {
		p.StartDate = time.Now().AddDate(0, 0, -30).Unix()
	}
	if p.EndDate == 0 {
		p.EndDate = time.Now().Unix()
	}

	// 生成多个图表
	charts := []chartData{}
	
	// 血糖趋势图
	bloodSugarChart, _ := t.generateChartData(dataVisualizerParams{
		Action:     "generate_chart",
		ChartType:  "line",
		DataType:   "blood_sugar",
		StartDate:  p.StartDate,
		EndDate:    p.EndDate,
		Resolution: p.Resolution,
	})
	charts = append(charts, bloodSugarChart)
	
	// 体重变化图
	weightChart, _ := t.generateChartData(dataVisualizerParams{
		Action:     "generate_chart",
		ChartType:  "line",
		DataType:   "weight",
		StartDate:  p.StartDate,
		EndDate:    p.EndDate,
		Resolution: p.Resolution,
	})
	charts = append(charts, weightChart)
	
	// 运动统计图
	exerciseChart, _ := t.generateChartData(dataVisualizerParams{
		Action:     "generate_chart",
		ChartType:  "bar",
		DataType:   "exercise",
		StartDate:  p.StartDate,
		EndDate:    p.EndDate,
		Resolution: p.Resolution,
	})
	charts = append(charts, exerciseChart)

	report := visualReport{
		Title:  "健康数据综合报告",
		Period: fmt.Sprintf("%s - %s", 
			time.Unix(p.StartDate, 0).Format("2006-01-02"),
			time.Unix(p.EndDate, 0).Format("2006-01-02")),
		Summary: summaryData{
			TotalDataPoints: 90,
			AverageValue:    6.2,
			MinValue:        5.8,
			MaxValue:        7.2,
			Trend:           "improving",
			TrendPercent:    -8.5,
		},
		Charts: charts,
		Insights: []string{
			"血糖水平整体呈下降趋势，控制效果良好",
			"体重稳步下降，符合减肥目标",
			"运动量充足，建议保持当前运动频率",
			"饮食控制得当，营养摄入均衡",
		},
		Recommendations: []string{
			"继续保持当前的健康管理计划",
			"建议增加力量训练以提高代谢",
			"定期监测血糖变化",
			"保持规律的作息时间",
		},
		GeneratedAt: time.Now().Unix(),
	}

	b, _ := sonic.Marshal(report)
	return string(b), nil
}

// generateChartData 生成图表数据
func (t *DataVisualizerTool) generateChartData(p dataVisualizerParams) (chartData, error) {
	// TODO: 从数据库获取实际数据
	// 这里模拟数据生成
	
	var chart chartData
	
	switch p.DataType {
	case "blood_sugar":
		chart = t.generateBloodSugarChart(p)
	case "weight":
		chart = t.generateWeightChart(p)
	case "exercise":
		chart = t.generateExerciseChart(p)
	case "diet":
		chart = t.generateDietChart(p)
	default:
		return chartData{}, fmt.Errorf("unknown data_type: %s", p.DataType)
	}
	
	return chart, nil
}

// generateBloodSugarChart 生成血糖图表
func (t *DataVisualizerTool) generateBloodSugarChart(p dataVisualizerParams) chartData {
	// 模拟30天的血糖数据
	days := 30
	labels := make([]string, days)
	data := make([]float64, days)
	
	for i := 0; i < days; i++ {
		date := time.Unix(p.StartDate, 0).AddDate(0, 0, i)
		labels[i] = date.Format("01-02")
		// 模拟血糖数据：6.0-7.2之间波动，整体呈下降趋势
		baseValue := 7.0 - (float64(i) / float64(days) * 0.5)
		data[i] = baseValue + (float64(i%3) * 0.2)
	}

	return chartData{
		Type:  p.ChartType,
		Title: "血糖趋势图",
		XAxis: axisConfig{
			Label:    "日期",
			Category: labels,
		},
		YAxis: axisConfig{
			Label: "血糖 (mmol/L)",
			Min:   4.0,
			Max:   10.0,
		},
		Datasets: []dataset{
			{
				Label:        "空腹血糖",
				Data:         data,
				BorderColor:  "#FF6B6B",
				BackgroundColor: "rgba(255, 107, 107, 0.1)",
				Fill:         true,
			},
		},
		Options: chartOptions{
			Responsive:  true,
			MaintainAspectRatio: false,
			Plugins:     []string{"legend", "tooltip"},
		},
		GeneratedAt: time.Now().Unix(),
	}
}

// generateWeightChart 生成体重图表
func (t *DataVisualizerTool) generateWeightChart(p dataVisualizerParams) chartData {
	// 模拟30天的体重数据
	days := 30
	labels := make([]string, days)
	data := make([]float64, days)
	
	for i := 0; i < days; i++ {
		date := time.Unix(p.StartDate, 0).AddDate(0, 0, i)
		labels[i] = date.Format("01-02")
		// 模拟体重数据：从77.8kg降到76.6kg
		startWeight := 77.8
		endWeight := 76.6
		data[i] = startWeight - ((startWeight - endWeight) * float64(i) / float64(days))
	}

	return chartData{
		Type:  p.ChartType,
		Title: "体重变化图",
		XAxis: axisConfig{
			Label:    "日期",
			Category: labels,
		},
		YAxis: axisConfig{
			Label: "体重 (kg)",
			Min:   75.0,
			Max:   80.0,
		},
		Datasets: []dataset{
			{
				Label:        "体重",
				Data:         data,
				BorderColor:  "#4ECDC4",
				BackgroundColor: "rgba(78, 205, 196, 0.1)",
				Fill:         true,
			},
		},
		Options: chartOptions{
			Responsive:  true,
			MaintainAspectRatio: false,
			Plugins:     []string{"legend", "tooltip"},
		},
		GeneratedAt: time.Now().Unix(),
	}
}

// generateExerciseChart 生成运动图表
func (t *DataVisualizerTool) generateExerciseChart(p dataVisualizerParams) chartData {
	// 模拟30天的运动数据（按周汇总）
	labels := []string{"第1周", "第2周", "第3周", "第4周"}
	data := []float64{180, 200, 220, 240} // 每周运动分钟数

	return chartData{
		Type:  p.ChartType,
		Title: "运动统计图",
		XAxis: axisConfig{
			Label:    "周期",
			Category: labels,
		},
		YAxis: axisConfig{
			Label: "运动时长 (分钟)",
			Min:   0,
			Max:   300,
		},
		Datasets: []dataset{
			{
				Label:        "运动时长",
				Data:         data,
				BorderColor:  "#95E1D3",
				BackgroundColor: "rgba(149, 225, 211, 0.5)",
				Fill:         false,
			},
		},
		Options: chartOptions{
			Responsive:  true,
			MaintainAspectRatio: false,
			Plugins:     []string{"legend", "tooltip"},
		},
		GeneratedAt: time.Now().Unix(),
	}
}

// generateDietChart 生成饮食图表
func (t *DataVisualizerTool) generateDietChart(p dataVisualizerParams) chartData {
	// 模拟营养摄入分布
	data := []float64{20, 60, 20} // 百分比

	return chartData{
		Type:  "pie",
		Title: "营养摄入分布",
		XAxis: axisConfig{
			Label: "营养成分",
		},
		YAxis: axisConfig{
			Label: "占比 (%)",
			Min:   0,
			Max:   100,
		},
		Datasets: []dataset{
			{
				Label:           "营养占比",
				Data:            data,
				BorderColor:     "#FF6B6B",
				BackgroundColor: "rgba(255, 107, 107, 0.8)",
				Fill:            false,
			},
		},
		Options: chartOptions{
			Responsive:           true,
			MaintainAspectRatio:  false,
			Plugins:              []string{"legend", "tooltip"},
		},
		GeneratedAt: time.Now().Unix(),
	}
}