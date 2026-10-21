package health

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/bytedance/sonic"
	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/schema"
	"new-robot-self/model"
	"new-robot-self/public"
	"xorm.io/xorm"
)

// HealthTrackerTool 健康数据追踪工具
// 支持记录血糖、体重、运动数据,生成趋势报告和异常预警
type HealthTrackerTool struct {
	EnableLearning bool
	Category       string
	ToolName       string
	db             *xorm.Engine // 数据库实例
}

// GetCategory 获取工具分类
func (t *HealthTrackerTool) GetCategory() string {
	return t.Category
}

// GetToolName 获取工具名称
func (t *HealthTrackerTool) GetToolName() string {
	return t.ToolName
}

// AfterRun 工具调用后的后处理
func (t *HealthTrackerTool) AfterRun(ctx context.Context, result string, err error) {
	// 简化实现，不进行学习集成以避免循环依赖
}

// SetDB 设置数据库实例
func (t *HealthTrackerTool) SetDB(db *xorm.Engine) {
	t.db = db
}

// getDB 获取数据库实例
func (t *HealthTrackerTool) getDB() *xorm.Engine {
	if t.db == nil {
		return public.DB
	}
	return t.db
}

type healthTrackerParams struct {
	Action      string                 `json:"action"`      // record, query, report, alert
	RecordType  string                 `json:"record_type,omitempty"`  // blood_sugar, weight, exercise
	Value       float64                `json:"value,omitempty"`       // 数值
	Unit        string                 `json:"unit,omitempty"`        // 单位: mg/dL, kg, min, kcal
	Notes       string                 `json:"notes,omitempty"`       // 备注
	RecordedAt  int64                  `json:"recorded_at,omitempty"` // 记录时间(时间戳)
	StartDate   int64                  `json:"start_date,omitempty"`  // 查询开始时间
	EndDate     int64                  `json:"end_date,omitempty"`    // 查询结束时间
	Metadata    map[string]interface{} `json:"metadata,omitempty"`     // 额外元数据
}

type healthRecord struct {
	ID         string                 `json:"id"`
	UserID     string                 `json:"user_id"`
	RecordType string                 `json:"record_type"` // blood_sugar, weight, exercise
	Value      float64                `json:"value"`
	Unit       string                 `json:"unit"`
	Notes      string                 `json:"notes"`
	RecordedAt int64                  `json:"recorded_at"` // 时间戳
	Metadata   map[string]interface{} `json:"metadata"`
}

type healthReport struct {
	Period         string          `json:"period"`
	Records        []healthRecord  `json:"records"`
	Summary        healthSummary   `json:"summary"`
	Trends         map[string]trendData `json:"trends"`
	Recommendations []string       `json:"recommendations"`
}

type healthSummary struct {
	TotalRecords      int     `json:"total_records"`
	AvgBloodSugar     float64 `json:"avg_blood_sugar"`
	WeightChange      float64 `json:"weight_change"`
	TotalExerciseTime int     `json:"total_exercise_time"`
	TotalCaloriesBurned int   `json:"total_calories_burned"`
}

type trendData struct {
	CurrentValue float64   `json:"current_value"`
	PreviousValue float64 `json:"previous_value"`
	Change       float64   `json:"change"`
	ChangePercent float64  `json:"change_percent"`
	Trend        string    `json:"trend"` // up, down, stable
}

type healthAlert struct {
	AlertType   string  `json:"alert_type"`
	Severity    string  `json:"severity"` // low, medium, high
	Message     string  `json:"message"`
	Value       float64 `json:"value"`
	Threshold   float64 `json:"threshold"`
	RecordedAt  int64   `json:"recorded_at"`
}

func NewHealthTrackerTool() *HealthTrackerTool {
	return &HealthTrackerTool{
		Category:       "health",
		EnableLearning: true,
		ToolName:       "health_tracker",
	}
}

func (t *HealthTrackerTool) Info(_ context.Context) (*schema.ToolInfo, error) {
	return &schema.ToolInfo{
		Name: "health_tracker",
		Desc: "健康数据追踪工具，支持记录血糖、体重、运动数据，生成趋势报告和异常预警",
		ParamsOneOf: schema.NewParamsOneOfByParams(map[string]*schema.ParameterInfo{
			"action":       {Desc: "操作类型: record(记录), query(查询), report(报告), alert(预警)", Type: schema.String, Required: true},
			"record_type":  {Desc: "记录类型: blood_sugar(血糖), weight(体重), exercise(运动)", Type: schema.String, Required: false},
			"value":        {Desc: "数值", Type: schema.Number, Required: false},
			"unit":         {Desc: "单位: mg/dL(血糖), kg(体重), min(运动时长), kcal(卡路里)", Type: schema.String, Required: false},
			"notes":        {Desc: "备注信息", Type: schema.String, Required: false},
			"recorded_at":  {Desc: "记录时间(时间戳)", Type: schema.Integer, Required: false},
			"start_date":   {Desc: "查询开始时间(时间戳)", Type: schema.Integer, Required: false},
			"end_date":     {Desc: "查询结束时间(时间戳)", Type: schema.Integer, Required: false},
			"metadata":     {Desc: "额外元数据", Type: schema.Object, Required: false},
		}),
	}, nil
}

func (t *HealthTrackerTool) InvokableRun(ctx context.Context, argumentsInJSON string, _ ...tool.Option) (string, error) {
	p := healthTrackerParams{}
	if err := sonic.UnmarshalString(argumentsInJSON, &p); err != nil {
		return "", err
	}

	// 验证参数
	if p.Action == "" {
		return "", fmt.Errorf("action is required")
	}

	// 设置默认时间
	if p.RecordedAt == 0 {
		p.RecordedAt = time.Now().Unix()
	}

	// 根据操作类型执行相应功能
	switch p.Action {
	case "record":
		return t.recordHealthData(p)
	case "query":
		return t.queryHealthData(p)
	case "report":
		return t.generateHealthReport(p)
	case "alert":
		return t.checkHealthAlerts(p)
	default:
		return "", fmt.Errorf("unknown action: %s", p.Action)
	}
}

// recordHealthData 记录健康数据
func (t *HealthTrackerTool) recordHealthData(p healthTrackerParams) (string, error) {
	if p.RecordType == "" {
		return "", fmt.Errorf("record_type is required for record action")
	}
	if p.Value == 0 {
		return "", fmt.Errorf("value is required for record action")
	}

	// 验证单位
	if err := t.validateUnit(p.RecordType, p.Unit); err != nil {
		return "", err
	}

	// 构建数据库记录
	metadataJSON := ""
	if p.Metadata != nil && len(p.Metadata) > 0 {
		metadataBytes, _ := json.Marshal(p.Metadata)
		metadataJSON = string(metadataBytes)
	}

	record := &model.HealthRecord{
		UserID:     "weimeng_xiansheng", // 威猛先生
		RecordType: p.RecordType,
		Value:      p.Value,
		Unit:       p.Unit,
		Notes:      p.Notes,
		Metadata:   metadataJSON,
		RecordedAt: p.RecordedAt,
	}

	// 存储到数据库
	db := t.getDB()
	if db == nil {
		return "", fmt.Errorf("database not initialized")
	}

	_, err := db.Insert(record)
	if err != nil {
		return "", fmt.Errorf("failed to save health record: %w", err)
	}

	result := map[string]interface{}{
		"success": true,
		"message": "健康数据记录成功",
		"record": map[string]interface{}{
			"id":          record.Id,
			"user_id":     record.UserID,
			"record_type": record.RecordType,
			"value":       record.Value,
			"unit":        record.Unit,
			"notes":       record.Notes,
			"recorded_at": record.RecordedAt,
		},
	}

	b, _ := sonic.Marshal(result)
	return string(b), nil
}

// queryHealthData 查询健康数据
func (t *HealthTrackerTool) queryHealthData(p healthTrackerParams) (string, error) {
	// 设置默认查询时间范围（最近7天）
	if p.StartDate == 0 {
		p.StartDate = time.Now().AddDate(0, 0, -7).Unix()
	}
	if p.EndDate == 0 {
		p.EndDate = time.Now().Unix()
	}

	db := t.getDB()
	if db == nil {
		return "", fmt.Errorf("database not initialized")
	}

	// 构建查询
	query := db.Where("user_id = ? AND recorded_at >= ? AND recorded_at <= ?",
		"weimeng_xiansheng", p.StartDate, p.EndDate)

	// 如果指定了记录类型，添加过滤条件
	if p.RecordType != "" {
		query = query.And("record_type = ?", p.RecordType)
	}

	// 查询数据
	var records []model.HealthRecord
	err := query.OrderBy("recorded_at DESC").Find(&records)
	if err != nil {
		return "", fmt.Errorf("failed to query health records: %w", err)
	}

	// 转换为响应格式
	resultRecords := make([]map[string]interface{}, 0, len(records))
	for _, r := range records {
		metadata := make(map[string]interface{})
		if r.Metadata != "" {
			json.Unmarshal([]byte(r.Metadata), &metadata)
		}
		resultRecords = append(resultRecords, map[string]interface{}{
			"id":          r.Id,
			"user_id":     r.UserID,
			"record_type": r.RecordType,
			"value":       r.Value,
			"unit":        r.Unit,
			"notes":       r.Notes,
			"recorded_at": r.RecordedAt,
			"metadata":    metadata,
		})
	}

	result := map[string]interface{}{
		"success": true,
		"records": resultRecords,
		"count":   len(records),
	}

	b, _ := sonic.Marshal(result)
	return string(b), nil
}

// generateHealthReport 生成健康报告
func (t *HealthTrackerTool) generateHealthReport(p healthTrackerParams) (string, error) {
	// 设置默认报告时间范围（最近30天）
	if p.StartDate == 0 {
		p.StartDate = time.Now().AddDate(0, 0, -30).Unix()
	}
	if p.EndDate == 0 {
		p.EndDate = time.Now().Unix()
	}

	db := t.getDB()
	if db == nil {
		return "", fmt.Errorf("database not initialized")
	}

	// 查询时间段内的所有记录
	var records []model.HealthRecord
	err := db.Where("user_id = ? AND recorded_at >= ? AND recorded_at <= ?",
		"weimeng_xiansheng", p.StartDate, p.EndDate).
		OrderBy("recorded_at DESC").
		Find(&records)
	if err != nil {
		return "", fmt.Errorf("failed to query health records: %w", err)
	}

	// 计算统计信息
	summary := t.calculateHealthSummary(records)

	// 计算趋势
	trends := t.calculateHealthTrends(records, p.StartDate, p.EndDate)

	// 生成建议
	recommendations := t.generateRecommendations(summary, trends)

	report := healthReport{
		Period: fmt.Sprintf("%s - %s",
			time.Unix(p.StartDate, 0).Format("2006-01-02"),
			time.Unix(p.EndDate, 0).Format("2006-01-02")),
		Records:       []healthRecord{},
		Summary:       summary,
		Trends:        trends,
		Recommendations: recommendations,
	}

	b, _ := sonic.Marshal(report)
	return string(b), nil
}

// calculateHealthSummary 计算健康摘要
func (t *HealthTrackerTool) calculateHealthSummary(records []model.HealthRecord) healthSummary {
	summary := healthSummary{
		TotalRecords: len(records),
	}

	var bloodSugarSum float64
	var bloodSugarCount int
	var firstWeight, lastWeight float64
	var weightRecords []model.HealthRecord

	for _, r := range records {
		switch r.RecordType {
		case "blood_sugar":
			bloodSugarSum += r.Value
			bloodSugarCount++
		case "weight":
			weightRecords = append(weightRecords, r)
		case "exercise":
			if r.Unit == "min" {
				summary.TotalExerciseTime += int(r.Value)
			} else if r.Unit == "kcal" {
				summary.TotalCaloriesBurned += int(r.Value)
			}
		}
	}

	// 计算平均血糖
	if bloodSugarCount > 0 {
		summary.AvgBloodSugar = bloodSugarSum / float64(bloodSugarCount)
	}

	// 计算体重变化
	if len(weightRecords) >= 2 {
		// 按时间排序
		for i := 0; i < len(weightRecords); i++ {
			for j := i + 1; j < len(weightRecords); j++ {
				if weightRecords[i].RecordedAt > weightRecords[j].RecordedAt {
					weightRecords[i], weightRecords[j] = weightRecords[j], weightRecords[i]
				}
			}
		}
		firstWeight = weightRecords[0].Value
		lastWeight = weightRecords[len(weightRecords)-1].Value
		summary.WeightChange = lastWeight - firstWeight
	}

	return summary
}

// calculateHealthTrends 计算健康趋势
func (t *HealthTrackerTool) calculateHealthTrends(records []model.HealthRecord, startDate, endDate int64) map[string]trendData {
	trends := make(map[string]trendData)

	// 按类型分组
	typeRecords := make(map[string][]model.HealthRecord)
	for _, r := range records {
		typeRecords[r.RecordType] = append(typeRecords[r.RecordType], r)
	}

	// 计算每种类型的趋势
	for recordType, recs := range typeRecords {
		if len(recs) < 2 {
			continue
		}

		// 按时间排序
		for i := 0; i < len(recs); i++ {
			for j := i + 1; j < len(recs); j++ {
				if recs[i].RecordedAt > recs[j].RecordedAt {
					recs[i], recs[j] = recs[j], recs[i]
				}
			}
		}

		// 计算前半段和后半段的平均值
		mid := len(recs) / 2
		var firstHalfSum, secondHalfSum float64
		for _, r := range recs[:mid] {
			firstHalfSum += r.Value
		}
		for _, r := range recs[mid:] {
			secondHalfSum += r.Value
		}

		firstHalfAvg := firstHalfSum / float64(mid)
		secondHalfAvg := secondHalfSum / float64(len(recs)-mid)

		change := secondHalfAvg - firstHalfAvg
		changePercent := 0.0
		if firstHalfAvg != 0 {
			changePercent = (change / firstHalfAvg) * 100
		}

		trend := "stable"
		if change > 0.5 {
			trend = "up"
		} else if change < -0.5 {
			trend = "down"
		}

		trends[recordType] = trendData{
			CurrentValue:  secondHalfAvg,
			PreviousValue: firstHalfAvg,
			Change:        change,
			ChangePercent: changePercent,
			Trend:         trend,
		}
	}

	return trends
}

// generateRecommendations 生成健康建议
func (t *HealthTrackerTool) generateRecommendations(summary healthSummary, trends map[string]trendData) []string {
	recommendations := []string{}

	// 血糖建议
	if summary.AvgBloodSugar > 0 {
		if summary.AvgBloodSugar > 7.0 {
			recommendations = append(recommendations, "血糖偏高，建议控制饮食并增加运动")
		} else if summary.AvgBloodSugar < 4.0 {
			recommendations = append(recommendations, "血糖偏低，注意及时补充能量")
		} else {
			recommendations = append(recommendations, "血糖控制良好，继续保持")
		}
	}

	// 体重建议
	if summary.WeightChange != 0 {
		if summary.WeightChange > 2 {
			recommendations = append(recommendations, "体重上升明显，建议调整饮食结构并增加运动")
		} else if summary.WeightChange < -2 {
			recommendations = append(recommendations, "体重下降趋势明显，请确保营养均衡")
		} else {
			recommendations = append(recommendations, "体重稳定，保持当前生活方式")
		}
	}

	// 运动建议
	if summary.TotalExerciseTime < 150 {
		recommendations = append(recommendations, fmt.Sprintf("本周运动时间%d分钟，建议每周至少150分钟中等强度运动", summary.TotalExerciseTime))
	} else {
		recommendations = append(recommendations, "运动量充足，建议结合力量训练和有氧运动")
	}

	return recommendations
}

// checkHealthAlerts 检查健康预警
func (t *HealthTrackerTool) checkHealthAlerts(p healthTrackerParams) (string, error) {
	db := t.getDB()
	if db == nil {
		return "", fmt.Errorf("database not initialized")
	}

	alerts := []healthAlert{}

	// 查询最新的血糖记录
	var bloodSugarRecords []model.HealthRecord
	err := db.Where("user_id = ? AND record_type = ?", "weimeng_xiansheng", "blood_sugar").
		OrderBy("recorded_at DESC").
		Limit(1).
		Find(&bloodSugarRecords)
	if err == nil && len(bloodSugarRecords) > 0 {
		latestBloodSugar := bloodSugarRecords[0].Value
		// 空腹血糖正常范围: 3.9-6.1 mmol/L
		if latestBloodSugar > 7.0 {
			alerts = append(alerts, healthAlert{
				AlertType:  "blood_sugar_high",
				Severity:   "medium",
				Message:    fmt.Sprintf("血糖%.1f mmol/L偏高（正常<7.0），建议控制饮食并增加运动", latestBloodSugar),
				Value:      latestBloodSugar,
				Threshold:  7.0,
				RecordedAt: bloodSugarRecords[0].RecordedAt,
			})
		} else if latestBloodSugar < 3.9 {
			alerts = append(alerts, healthAlert{
				AlertType:  "blood_sugar_low",
				Severity:   "high",
				Message:    fmt.Sprintf("血糖%.1f mmol/L偏低（正常>3.9），请及时补充能量", latestBloodSugar),
				Value:      latestBloodSugar,
				Threshold:  3.9,
				RecordedAt: bloodSugarRecords[0].RecordedAt,
			})
		}
	}

	// 查询最新的体重记录
	var weightRecords []model.HealthRecord
	err = db.Where("user_id = ? AND record_type = ?", "weimeng_xiansheng", "weight").
		OrderBy("recorded_at DESC").
		Limit(1).
		Find(&weightRecords)
	if err == nil && len(weightRecords) > 0 {
		latestWeight := weightRecords[0].Value
		targetWeight := 70.0 // 目标体重
		if latestWeight > targetWeight+5 {
			alerts = append(alerts, healthAlert{
				AlertType:  "weight_above_target",
				Severity:   "medium",
				Message:    fmt.Sprintf("体重%.1fkg，超过目标体重%.1fkg较多，建议加强运动和控制饮食", latestWeight, targetWeight),
				Value:      latestWeight,
				Threshold:  targetWeight,
				RecordedAt: weightRecords[0].RecordedAt,
			})
		}
	}

	// 检查今日运动量
	todayStart := time.Now().Truncate(24 * time.Hour).Unix()
	exerciseCount, _ := db.Where("user_id = ? AND record_type = ? AND recorded_at >= ?",
		"weimeng_xiansheng", "exercise", todayStart).
		Count(&model.HealthRecord{})
	if exerciseCount == 0 {
		alerts = append(alerts, healthAlert{
			AlertType:  "no_exercise_today",
			Severity:   "low",
			Message:    "今日暂无运动记录，建议进行适量运动",
			Value:      0,
			Threshold:  1,
			RecordedAt: time.Now().Unix(),
		})
	}

	result := map[string]interface{}{
		"success": true,
		"alerts":  alerts,
		"count":   len(alerts),
	}

	b, _ := sonic.Marshal(result)
	return string(b), nil
}

// validateUnit 验证单位是否正确
func (t *HealthTrackerTool) validateUnit(recordType, unit string) error {
	validUnits := map[string][]string{
		"blood_sugar": {"mmol/L", "mg/dL"},
		"weight":      {"kg", "lb"},
		"exercise":    {"min", "kcal"},
	}

	if units, ok := validUnits[recordType]; ok {
		for _, validUnit := range units {
			if unit == validUnit {
				return nil
			}
		}
		return fmt.Errorf("invalid unit %s for record_type %s, valid units are: %v", unit, recordType, units)
	}

	return fmt.Errorf("unknown record_type: %s", recordType)
}