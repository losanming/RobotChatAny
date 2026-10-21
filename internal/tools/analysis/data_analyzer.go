package analysis

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"sort"

	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/schema"
)

// DataAnalyzer 数据分析工具
type DataAnalyzer struct {
	enableLearning bool
	category       string
	toolName       string
}

// NewDataAnalyzer 创建数据分析工具
func NewDataAnalyzer() *DataAnalyzer {
	return &DataAnalyzer{
		enableLearning: true,
		category:       "analysis",
		toolName:       "data_analyzer",
	}
}

// Info 返回工具信息
func (d *DataAnalyzer) Info(ctx context.Context) (*schema.ToolInfo, error) {
	return &schema.ToolInfo{
		Name: "data_analyzer",
		Desc: "数据分析工具，支持基本统计分析、数据聚合和分组、趋势分析、异常检测",
		ParamsOneOf: schema.NewParamsOneOfByParams(map[string]*schema.ParameterInfo{
			"action": {
				Desc:     "操作类型",
				Type:     schema.String,
				Required: true,
				Enum:     []string{"statistics", "aggregate", "trend", "anomaly", "correlation"},
			},
			"data": {
				Desc:     "数据数组",
				Type:     schema.Array,
				Required: false,
			},
			"analysis_type": {
				Desc:     "分析类型",
				Type:     schema.String,
				Required: false,
				Enum:     []string{"mean", "median", "mode", "stddev", "variance", "min", "max", "range", "sum", "count", "all"},
			},
			"group_by": {
				Desc:     "分组字段",
				Type:     schema.String,
				Required: false,
			},
			"aggregate_function": {
				Desc:     "聚合函数",
				Type:     schema.String,
				Required: false,
				Enum:     []string{"sum", "avg", "count", "min", "max", "stddev"},
			},
			"time_field": {
				Desc:     "时间字段",
				Type:     schema.String,
				Required: false,
			},
			"value_field": {
				Desc:     "值字段",
				Type:     schema.String,
				Required: false,
			},
			"threshold": {
				Desc:     "异常检测阈值",
				Type:     schema.Number,
				Required: false,
			},
			"method": {
				Desc:     "异常检测方法",
				Type:     schema.String,
				Required: false,
				Enum:     []string{"zscore", "iqr", "isolation"},
			},
		}),
	}, nil
}

// InvokableRun 执行工具调用
func (d *DataAnalyzer) InvokableRun(ctx context.Context, args string, opts ...tool.Option) (string, error) {
	var params map[string]interface{}
	if err := json.Unmarshal([]byte(args), &params); err != nil {
		return "", fmt.Errorf("参数解析失败: %w", err)
	}

	action := params["action"].(string)
	switch action {
	case "statistics":
		return d.analyzeStatistics(params)
	case "aggregate":
		return d.analyzeAggregate(params)
	case "trend":
		return d.analyzeTrend(params)
	case "anomaly":
		return d.detectAnomaly(params)
	case "correlation":
		return d.analyzeCorrelation(params)
	default:
		return "", fmt.Errorf("不支持的操作类型: %s", action)
	}
}

// analyzeStatistics 统计分析
func (d *DataAnalyzer) analyzeStatistics(params map[string]interface{}) (string, error) {
	analysisType := getStringParam(params, "analysis_type", "all")

	// 获取数据
	data, err := d.extractNumericData(params)
	if err != nil {
		return "", fmt.Errorf("数据提取失败: %w", err)
	}

	if len(data) == 0 {
		return "", fmt.Errorf("数据不能为空")
	}

	stats := make(map[string]interface{})

	switch analysisType {
	case "mean":
		stats["mean"] = d.calculateMean(data)
	case "median":
		stats["median"] = d.calculateMedian(data)
	case "mode":
		stats["mode"] = d.calculateMode(data)
	case "stddev":
		stats["stddev"] = d.calculateStdDev(data)
	case "variance":
		stats["variance"] = d.calculateVariance(data)
	case "min":
		stats["min"] = d.calculateMin(data)
	case "max":
		stats["max"] = d.calculateMax(data)
	case "range":
		stats["range"] = d.calculateMax(data) - d.calculateMin(data)
	case "sum":
		stats["sum"] = d.calculateSum(data)
	case "count":
		stats["count"] = len(data)
	case "all":
		stats["count"] = len(data)
		stats["sum"] = d.calculateSum(data)
		stats["mean"] = d.calculateMean(data)
		stats["median"] = d.calculateMedian(data)
		stats["mode"] = d.calculateMode(data)
		stats["stddev"] = d.calculateStdDev(data)
		stats["variance"] = d.calculateVariance(data)
		stats["min"] = d.calculateMin(data)
		stats["max"] = d.calculateMax(data)
		stats["range"] = d.calculateMax(data) - d.calculateMin(data)
	}

	result := map[string]interface{}{
		"success":       true,
		"analysis_type": analysisType,
		"statistics":    stats,
		"data_count":    len(data),
	}

	return toJSON(result), nil
}

// analyzeAggregate 聚合分析
func (d *DataAnalyzer) analyzeAggregate(params map[string]interface{}) (string, error) {
	groupBy := getStringParam(params, "group_by", "")
	aggregateFunc := getStringParam(params, "aggregate_function", "sum")

	// 简化实现：对数值数组进行聚合
	data, err := d.extractNumericData(params)
	if err != nil {
		return "", fmt.Errorf("数据提取失败: %w", err)
	}

	if len(data) == 0 {
		return "", fmt.Errorf("数据不能为空")
	}

	var result float64
	switch aggregateFunc {
	case "sum":
		result = d.calculateSum(data)
	case "avg":
		result = d.calculateMean(data)
	case "count":
		result = float64(len(data))
	case "min":
		result = d.calculateMin(data)
	case "max":
		result = d.calculateMax(data)
	case "stddev":
		result = d.calculateStdDev(data)
	default:
		return "", fmt.Errorf("不支持的聚合函数: %s", aggregateFunc)
	}

	analysisResult := map[string]interface{}{
		"aggregate_function": aggregateFunc,
		"result":             result,
		"group_by":           groupBy,
		"data_count":         len(data),
	}

	finalResult := map[string]interface{}{
		"success": true,
		"result":  analysisResult,
	}

	return toJSON(finalResult), nil
}

// analyzeTrend 趋势分析
func (d *DataAnalyzer) analyzeTrend(params map[string]interface{}) (string, error) {
	// 获取数据
	data, err := d.extractNumericData(params)
	if err != nil {
		return "", fmt.Errorf("数据提取失败: %w", err)
	}

	if len(data) < 2 {
		return "", fmt.Errorf("数据点数量不足，至少需要2个数据点")
	}

	// 计算趋势
	trend := d.calculateTrend(data)
	growthRate := d.calculateGrowthRate(data)

	// 趋势判断
	trendDirection := "stable"
	if growthRate > 0.05 {
		trendDirection = "increasing"
	} else if growthRate < -0.05 {
		trendDirection = "decreasing"
	}

	result := map[string]interface{}{
		"success":         true,
		"trend":           trend,
		"growth_rate":     growthRate,
		"trend_direction": trendDirection,
		"data_points":     len(data),
		"first_value":     data[0],
		"last_value":      data[len(data)-1],
	}

	return toJSON(result), nil
}

// detectAnomaly 异常检测
func (d *DataAnalyzer) detectAnomaly(params map[string]interface{}) (string, error) {
	method := getStringParam(params, "method", "zscore")
	threshold := getFloatParam(params, "threshold", 2.0)

	// 获取数据
	data, err := d.extractNumericData(params)
	if err != nil {
		return "", fmt.Errorf("数据提取失败: %w", err)
	}

	if len(data) == 0 {
		return "", fmt.Errorf("数据不能为空")
	}

	var anomalies []map[string]interface{}

	switch method {
	case "zscore":
		anomalies = d.detectZScoreAnomalies(data, threshold)
	case "iqr":
		anomalies = d.detectIQRAnomalies(data, threshold)
	default:
		return "", fmt.Errorf("不支持的异常检测方法: %s", method)
	}

	result := map[string]interface{}{
		"success":      true,
		"method":       method,
		"threshold":    threshold,
		"anomalies":    anomalies,
		"anomaly_count": len(anomalies),
		"total_count":  len(data),
	}

	return toJSON(result), nil
}

// analyzeCorrelation 相关性分析
func (d *DataAnalyzer) analyzeCorrelation(params map[string]interface{}) (string, error) {
	// 简化实现：计算自相关性
	data, err := d.extractNumericData(params)
	if err != nil {
		return "", fmt.Errorf("数据提取失败: %w", err)
	}

	if len(data) < 2 {
		return "", fmt.Errorf("数据点数量不足，至少需要2个数据点")
	}

	// 计算自相关性
	correlation := d.calculateAutoCorrelation(data)

	result := map[string]interface{}{
		"success":     true,
		"correlation": correlation,
		"data_count":  len(data),
	}

	return toJSON(result), nil
}

// extractNumericData 提取数值数据
func (d *DataAnalyzer) extractNumericData(params map[string]interface{}) ([]float64, error) {
	dataInterface, ok := params["data"]
	if !ok {
		return nil, fmt.Errorf("缺少data参数")
	}

	dataArray, ok := dataInterface.([]interface{})
	if !ok {
		return nil, fmt.Errorf("data参数必须是数组")
	}

	var result []float64
	for _, item := range dataArray {
		switch v := item.(type) {
		case float64:
			result = append(result, v)
		case float32:
			result = append(result, float64(v))
		case int:
			result = append(result, float64(v))
		case int64:
			result = append(result, float64(v))
		default:
			// 忽略非数值类型
		}
	}

	return result, nil
}

// 统计计算函数
func (d *DataAnalyzer) calculateMean(data []float64) float64 {
	if len(data) == 0 {
		return 0
	}
	return d.calculateSum(data) / float64(len(data))
}

func (d *DataAnalyzer) calculateSum(data []float64) float64 {
	sum := 0.0
	for _, v := range data {
		sum += v
	}
	return sum
}

func (d *DataAnalyzer) calculateMedian(data []float64) float64 {
	if len(data) == 0 {
		return 0
	}

	sorted := make([]float64, len(data))
	copy(sorted, data)
	sort.Float64s(sorted)

	n := len(sorted)
	if n%2 == 0 {
		return (sorted[n/2-1] + sorted[n/2]) / 2
	}
	return sorted[n/2]
}

func (d *DataAnalyzer) calculateMode(data []float64) float64 {
	if len(data) == 0 {
		return 0
	}

	frequency := make(map[float64]int)
	for _, v := range data {
		frequency[v]++
	}

	mode := data[0]
	maxCount := 0
	for value, count := range frequency {
		if count > maxCount {
			maxCount = count
			mode = value
		}
	}

	return mode
}

func (d *DataAnalyzer) calculateVariance(data []float64) float64 {
	if len(data) == 0 {
		return 0
	}

	mean := d.calculateMean(data)
	sum := 0.0
	for _, v := range data {
		diff := v - mean
		sum += diff * diff
	}
	return sum / float64(len(data))
}

func (d *DataAnalyzer) calculateStdDev(data []float64) float64 {
	return math.Sqrt(d.calculateVariance(data))
}

func (d *DataAnalyzer) calculateMin(data []float64) float64 {
	if len(data) == 0 {
		return 0
	}
	min := data[0]
	for _, v := range data {
		if v < min {
			min = v
		}
	}
	return min
}

func (d *DataAnalyzer) calculateMax(data []float64) float64 {
	if len(data) == 0 {
		return 0
	}
	max := data[0]
	for _, v := range data {
		if v > max {
			max = v
		}
	}
	return max
}

// 趋势分析函数
func (d *DataAnalyzer) calculateTrend(data []float64) float64 {
	if len(data) < 2 {
		return 0
	}

	// 简单线性回归
	n := float64(len(data))
	sumX := 0.0
	sumY := 0.0
	sumXY := 0.0
	sumX2 := 0.0

	for i, y := range data {
		x := float64(i)
		sumX += x
		sumY += y
		sumXY += x * y
		sumX2 += x * x
	}

	slope := (n*sumXY - sumX*sumY) / (n*sumX2 - sumX*sumX)
	return slope
}

func (d *DataAnalyzer) calculateGrowthRate(data []float64) float64 {
	if len(data) < 2 {
		return 0
	}

	first := data[0]
	last := data[len(data)-1]

	if first == 0 {
		return 0
	}

	return (last - first) / first
}

// 异常检测函数
func (d *DataAnalyzer) detectZScoreAnomalies(data []float64, threshold float64) []map[string]interface{} {
	var anomalies []map[string]interface{}

	mean := d.calculateMean(data)
	stddev := d.calculateStdDev(data)

	if stddev == 0 {
		return anomalies
	}

	for i, value := range data {
		zscore := math.Abs((value - mean) / stddev)
		if zscore > threshold {
			anomalies = append(anomalies, map[string]interface{}{
				"index":   i,
				"value":   value,
				"zscore":  zscore,
				"reason":  fmt.Sprintf("Z-score %.2f exceeds threshold %.2f", zscore, threshold),
			})
		}
	}

	return anomalies
}

func (d *DataAnalyzer) detectIQRAnomalies(data []float64, threshold float64) []map[string]interface{} {
	var anomalies []map[string]interface{}

	if len(data) < 4 {
		return anomalies
	}

	sorted := make([]float64, len(data))
	copy(sorted, data)
	sort.Float64s(sorted)

	n := len(sorted)
	q1 := sorted[n/4]
	q3 := sorted[3*n/4]
	iqr := q3 - q1

	if iqr == 0 {
		return anomalies
	}

	lowerBound := q1 - threshold*iqr
	upperBound := q3 + threshold*iqr

	for i, value := range data {
		if value < lowerBound || value > upperBound {
			anomalies = append(anomalies, map[string]interface{}{
				"index":       i,
				"value":       value,
				"lower_bound": lowerBound,
				"upper_bound": upperBound,
				"reason":      fmt.Sprintf("Value %.2f outside IQR bounds [%.2f, %.2f]", value, lowerBound, upperBound),
			})
		}
	}

	return anomalies
}

// 相关性分析函数
func (d *DataAnalyzer) calculateAutoCorrelation(data []float64) float64 {
	if len(data) < 2 {
		return 0
	}

	// 计算滞后1的自相关性
	n := len(data)
	mean := d.calculateMean(data)

	var numerator float64
	var denominator float64

	for i := 0; i < n-1; i++ {
		numerator += (data[i] - mean) * (data[i+1] - mean)
	}

	for i := 0; i < n; i++ {
		denominator += (data[i] - mean) * (data[i] - mean)
	}

	if denominator == 0 {
		return 0
	}

	return numerator / denominator
}

// AfterRun 工具调用后的后处理
func (d *DataAnalyzer) AfterRun(ctx context.Context, result string, err error) {
	// 简化实现，暂时不做学习集成
}

// 辅助函数
func getFloatParam(params map[string]interface{}, key string, defaultValue float64) float64 {
	if val, ok := params[key]; ok {
		return val.(float64)
	}
	return defaultValue
}