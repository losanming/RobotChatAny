package health

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/schema"
	"new-robot-self/model"
	"new-robot-self/public"
	"xorm.io/xorm"
)

// DietManagerTool 饮食管理工具
// 支持卡路里计算、糖尿病饮食建议、膳食记录和营养分析
type DietManagerTool struct {
	EnableLearning bool
	Category       string
	ToolName       string
	db             *xorm.Engine // 数据库实例
}

// SetDB 设置数据库实例
func (t *DietManagerTool) SetDB(db *xorm.Engine) {
	t.db = db
}

// getDB 获取数据库实例
func (t *DietManagerTool) getDB() *xorm.Engine {
	if t.db == nil {
		return public.DB
	}
	return t.db
}

// GetCategory 获取工具分类
func (t *DietManagerTool) GetCategory() string {
	return t.Category
}

// GetToolName 获取工具名称
func (t *DietManagerTool) GetToolName() string {
	return t.ToolName
}

// AfterRun 工具调用后的后处理
func (t *DietManagerTool) AfterRun(ctx context.Context, result string, err error) {
	// 简化实现，不进行学习集成以避免循环依赖
}

type dietManagerParams struct {
	Action      string                 `json:"action"`        // calculate, record, query, suggest, analyze
	FoodName    string                 `json:"food_name,omitempty"`
	Amount      float64                `json:"amount,omitempty"`       // 食物数量
	AmountUnit  string                 `json:"amount_unit,omitempty"`  // 数量单位: g, ml, portion
	MealType    string                 `json:"meal_type,omitempty"`    // 餐次: breakfast, lunch, dinner, snack
	Calories    float64                `json:"calories,omitempty"`     // 卡路里
	Serving     string                 `json:"serving,omitempty"`      // 份量描述
	Notes       string                 `json:"notes,omitempty"`        // 备注
	RecordedAt  int64                  `json:"recorded_at,omitempty"`  // 记录时间(时间戳)
	StartDate   int64                  `json:"start_date,omitempty"`
	EndDate     int64                  `json:"end_date,omitempty"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
}

type foodInfo struct {
	Name             string  `json:"name"`
	CaloriesPer100g  float64 `json:"calories_per_100g"`
	ProteinPer100g   float64 `json:"protein_per_100g"`
	CarbsPer100g     float64 `json:"carbs_per_100g"`
	FatPer100g       float64 `json:"fat_per_100g"`
	GI               float64 `json:"gi"`
	Category         string  `json:"category"` // low_gi, medium_gi, high_gi
}

type dietAnalysis struct {
	Period           string         `json:"period"`
	TotalCalories    float64        `json:"total_calories"`
	AvgDailyCalories float64        `json:"avg_daily_calories"`
	Macronutrients   macronutrients `json:"macronutrients"`
	MealDistribution map[string]float64 `json:"meal_distribution"`
	GIAnalysis       giAnalysis     `json:"gi_analysis"`
	Recommendations  []string       `json:"recommendations"`
}

type macronutrients struct {
	TotalProtein float64 `json:"total_protein"`
	TotalCarbs   float64 `json:"total_carbs"`
	TotalFat     float64 `json:"total_fat"`
	ProteinRatio float64 `json:"protein_ratio"` // 蛋白质占比
	CarbsRatio   float64 `json:"carbs_ratio"`   // 碳水占比
	FatRatio     float64 `json:"fat_ratio"`     // 脂肪占比
}

type giAnalysis struct {
	LowGIFoods     []string `json:"low_gi_foods"`
	MediumGIFoods  []string `json:"medium_gi_foods"`
	HighGIFoods    []string `json:"high_gi_foods"`
	AvgGI          float64  `json:"avg_gi"`
	Recommendation string   `json:"recommendation"`
}

type dietSuggestion struct {
	TargetCalories      float64            `json:"target_calories"`
	MealPlan            map[string]mealPlan `json:"meal_plan"`
	FoodRecommendations []string           `json:"food_recommendations"`
	AvoidFoods          []string           `json:"avoid_foods"`
	Tips                []string           `json:"tips"`
}

type mealPlan struct {
	RecommendedCalories float64  `json:"recommended_calories"`
	FoodSuggestions     []string `json:"food_suggestions"`
	SampleMenu          []string `json:"sample_menu"`
}

func NewDietManagerTool() *DietManagerTool {
	return &DietManagerTool{
		Category:       "health",
		EnableLearning: true,
		ToolName:       "diet_manager",
	}
}

func (t *DietManagerTool) Info(_ context.Context) (*schema.ToolInfo, error) {
	return &schema.ToolInfo{
		Name: "diet_manager",
		Desc: "饮食管理工具，支持卡路里计算、糖尿病饮食建议、膳食记录和营养分析",
		ParamsOneOf: schema.NewParamsOneOfByParams(map[string]*schema.ParameterInfo{
			"action":      {Desc: "操作类型: calculate(计算卡路里), record(记录膳食), query(查询记录), suggest(饮食建议), analyze(营养分析)", Type: schema.String, Required: true},
			"food_name":   {Desc: "食物名称", Type: schema.String, Required: false},
			"amount":      {Desc: "食物数量(g)", Type: schema.Number, Required: false},
			"amount_unit": {Desc: "数量单位: g, ml, portion", Type: schema.String, Required: false},
			"meal_type":   {Desc: "餐次: breakfast, lunch, dinner, snack", Type: schema.String, Required: false},
			"calories":    {Desc: "卡路里", Type: schema.Number, Required: false},
			"serving":     {Desc: "份量描述", Type: schema.String, Required: false},
			"notes":       {Desc: "备注", Type: schema.String, Required: false},
			"recorded_at": {Desc: "记录时间(时间戳)", Type: schema.Integer, Required: false},
			"start_date":  {Desc: "查询开始时间(时间戳)", Type: schema.Integer, Required: false},
			"end_date":    {Desc: "查询结束时间(时间戳)", Type: schema.Integer, Required: false},
			"metadata":    {Desc: "额外元数据", Type: schema.Object, Required: false},
		}),
	}, nil
}

func (t *DietManagerTool) InvokableRun(ctx context.Context, argumentsInJSON string, _ ...tool.Option) (string, error) {
	var p dietManagerParams
	if err := json.Unmarshal([]byte(argumentsInJSON), &p); err != nil {
		return "", fmt.Errorf("参数解析失败: %w", err)
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
	case "calculate":
		return t.calculateCalories(p)
	case "record":
		return t.recordDiet(p)
	case "query":
		return t.queryDiet(p)
	case "suggest":
		return t.getDietSuggestion(p)
	case "analyze":
		return t.analyzeDiet(p)
	default:
		return "", fmt.Errorf("unknown action: %s", p.Action)
	}
}

// calculateCalories 计算卡路里
func (t *DietManagerTool) calculateCalories(p dietManagerParams) (string, error) {
	if p.FoodName == "" {
		return "", fmt.Errorf("food_name is required for calculate action")
	}
	if p.Amount == 0 {
		return "", fmt.Errorf("amount is required for calculate action")
	}

	// 获取食物信息（从数据库或API）
	info := t.getFoodInfo(p.FoodName)
	if info == nil {
		return "", fmt.Errorf("food not found: %s", p.FoodName)
	}

	// 计算营养值
	calories := (info.CaloriesPer100g / 100) * p.Amount
	protein := (info.ProteinPer100g / 100) * p.Amount
	carbs := (info.CarbsPer100g / 100) * p.Amount
	fat := (info.FatPer100g / 100) * p.Amount

	result := map[string]interface{}{
		"success": true,
		"food":    p.FoodName,
		"amount":  p.Amount,
		"unit":    p.AmountUnit,
		"nutrition": map[string]interface{}{
			"calories": calories,
			"protein":  protein,
			"carbs":    carbs,
			"fat":      fat,
			"gi":       info.GI,
		},
		"diabetes_friendly": info.GI < 55,
	}

	return toJSON(result), nil
}

// recordDiet 记录膳食
func (t *DietManagerTool) recordDiet(p dietManagerParams) (string, error) {
	if p.FoodName == "" {
		return "", fmt.Errorf("food_name is required for record action")
	}
	if p.MealType == "" {
		return "", fmt.Errorf("meal_type is required for record action")
	}

	// 获取食物信息
	info := t.getFoodInfo(p.FoodName)
	if info == nil {
		return "", fmt.Errorf("food not found: %s", p.FoodName)
	}

	// 计算营养值
	calories := p.Calories
	if calories == 0 && p.Amount > 0 {
		calories = (info.CaloriesPer100g / 100) * p.Amount
	}

	db := t.getDB()
	if db == nil {
		return "", fmt.Errorf("database not initialized")
	}

	record := &model.DietRecord{
		UserID:     "weimeng_xiansheng",
		MealType:   p.MealType,
		FoodName:   p.FoodName,
		Calories:   calories,
		Protein:    (info.ProteinPer100g / 100) * p.Amount,
		Carbs:      (info.CarbsPer100g / 100) * p.Amount,
		Fat:        (info.FatPer100g / 100) * p.Amount,
		Serving:    p.Serving,
		Notes:      p.Notes,
		RecordedAt: p.RecordedAt,
	}

	_, err := db.Insert(record)
	if err != nil {
		return "", fmt.Errorf("failed to record diet: %w", err)
	}

	result := map[string]interface{}{
		"success": true,
		"message": "膳食记录成功",
		"record": map[string]interface{}{
			"id":          record.Id,
			"food_name":   record.FoodName,
			"meal_type":   record.MealType,
			"calories":    record.Calories,
			"protein":     record.Protein,
			"carbs":       record.Carbs,
			"fat":         record.Fat,
			"serving":     record.Serving,
			"recorded_at": record.RecordedAt,
		},
	}

	return toJSON(result), nil
}

// queryDiet 查询膳食记录
func (t *DietManagerTool) queryDiet(p dietManagerParams) (string, error) {
	db := t.getDB()
	if db == nil {
		return "", fmt.Errorf("database not initialized")
	}

	// 设置默认查询时间范围（最近7天）
	if p.StartDate == 0 {
		p.StartDate = time.Now().AddDate(0, 0, -7).Unix()
	}
	if p.EndDate == 0 {
		p.EndDate = time.Now().Unix()
	}

	query := db.Where("user_id = ? AND recorded_at >= ? AND recorded_at <= ?",
		"weimeng_xiansheng", p.StartDate, p.EndDate)

	// 可选的餐次过滤
	if p.MealType != "" && p.MealType != "all" {
		query = query.And("meal_type = ?", p.MealType)
	}

	var records []model.DietRecord
	err := query.OrderBy("recorded_at DESC").Find(&records)
	if err != nil {
		return "", fmt.Errorf("failed to query diet records: %w", err)
	}

	resultRecords := make([]map[string]interface{}, 0, len(records))
	for _, record := range records {
		resultRecords = append(resultRecords, map[string]interface{}{
			"id":          record.Id,
			"food_name":   record.FoodName,
			"meal_type":   record.MealType,
			"calories":    record.Calories,
			"protein":     record.Protein,
			"carbs":       record.Carbs,
			"fat":         record.Fat,
			"serving":     record.Serving,
			"notes":       record.Notes,
			"recorded_at": record.RecordedAt,
		})
	}

	result := map[string]interface{}{
		"success": true,
		"records": resultRecords,
		"count":   len(records),
	}

	return toJSON(result), nil
}

// getDietSuggestion 获取饮食建议
func (t *DietManagerTool) getDietSuggestion(p dietManagerParams) (string, error) {
	// 基于威猛先生的糖尿病情况，生成个性化建议
	suggestion := dietSuggestion{
		TargetCalories: 1800, // 糖尿病患者推荐每日摄入
		MealPlan: map[string]mealPlan{
			"breakfast": {
				RecommendedCalories: 400,
				FoodSuggestions: []string{
					"低脂牛奶或豆浆",
					"全麦面包",
					"鸡蛋",
					"燕麦片",
				},
				SampleMenu: []string{
					"200ml低脂牛奶 + 1片全麦面包 + 1个鸡蛋",
				},
			},
			"lunch": {
				RecommendedCalories: 600,
				FoodSuggestions: []string{
					"糙米饭或杂粮饭（替代白米饭）",
					"瘦肉（鸡胸肉、鱼肉）",
					"大量蔬菜",
				},
				SampleMenu: []string{
					"100g糙米饭 + 100g鸡胸肉 + 200g蔬菜",
				},
			},
			"dinner": {
				RecommendedCalories: 500,
				FoodSuggestions: []string{
					"蔬菜沙拉",
					"清蒸鱼或豆腐",
					"少量主食",
				},
				SampleMenu: []string{
					"蔬菜沙拉 + 清蒸鱼 + 少量糙米",
				},
			},
			"snack": {
				RecommendedCalories: 300,
				FoodSuggestions: []string{
					"坚果（少量）",
					"低糖水果",
					"无糖酸奶",
				},
				SampleMenu: []string{
					"1根黄瓜 + 1杯黑咖啡",
				},
			},
		},
		FoodRecommendations: []string{
			"选择低GI食物（GI<55）",
			"增加膳食纤维摄入",
			"控制碳水化合物总量",
			"选择优质蛋白质",
			"健康脂肪来源（坚果、橄榄油）",
		},
		AvoidFoods: []string{
			"精制糖和甜食",
			"白米饭、白面包",
			"油炸食品",
			"加工肉类",
			"含糖饮料",
		},
		Tips: []string{
			"定时定量，少食多餐",
			"细嚼慢咽，控制进食速度",
			"餐后适量运动有助于控制血糖",
			"保持饮食记录，便于分析调整",
		},
	}

	return toJSON(suggestion), nil
}

// analyzeDiet 分析饮食
func (t *DietManagerTool) analyzeDiet(p dietManagerParams) (string, error) {
	db := t.getDB()
	if db == nil {
		return "", fmt.Errorf("database not initialized")
	}

	// 设置默认分析时间范围（最近7天）
	if p.StartDate == 0 {
		p.StartDate = time.Now().AddDate(0, 0, -7).Unix()
	}
	if p.EndDate == 0 {
		p.EndDate = time.Now().Unix()
	}

	// 查询记录
	var records []model.DietRecord
	err := db.Where("user_id = ? AND recorded_at >= ? AND recorded_at <= ?",
		"weimeng_xiansheng", p.StartDate, p.EndDate).
		Find(&records)
	if err != nil {
		return "", fmt.Errorf("failed to query records: %w", err)
	}

	// 计算统计数据
	totalCalories := 0.0
	totalProtein := 0.0
	totalCarbs := 0.0
	totalFat := 0.0
	mealDist := map[string]float64{
		"breakfast": 0,
		"lunch":     0,
		"dinner":    0,
		"snack":     0,
	}

	lowGIFoods := make([]string, 0)
	mediumGIFoods := make([]string, 0)
	highGIFoods := make([]string, 0)
	totalGI := 0.0
	giCount := 0

	for _, record := range records {
		totalCalories += record.Calories
		totalProtein += record.Protein
		totalCarbs += record.Carbs
		totalFat += record.Fat

		if cal, ok := mealDist[record.MealType]; ok {
			mealDist[record.MealType] = cal + record.Calories
		}

		// GI分类
		info := t.getFoodInfo(record.FoodName)
		if info != nil {
			totalGI += info.GI
			giCount++
			switch info.Category {
			case "low_gi":
				lowGIFoods = append(lowGIFoods, record.FoodName)
			case "medium_gi":
				mediumGIFoods = append(mediumGIFoods, record.FoodName)
			case "high_gi":
				highGIFoods = append(highGIFoods, record.FoodName)
			}
		}
	}

	days := float64((p.EndDate - p.StartDate) / 86400)
	if days < 1 {
		days = 1
	}
	avgDailyCalories := totalCalories / days

	avgGI := 0.0
	if giCount > 0 {
		avgGI = totalGI / float64(giCount)
	}

	// 计算宏量营养素比例
	totalMacro := totalProtein*4 + totalCarbs*4 + totalFat*9
	proteinRatio := 0.0
	carbsRatio := 0.0
	fatRatio := 0.0
	if totalMacro > 0 {
		proteinRatio = (totalProtein * 4) / totalMacro
		carbsRatio = (totalCarbs * 4) / totalMacro
		fatRatio = (totalFat * 9) / totalMacro
	}

	analysis := dietAnalysis{
		Period: fmt.Sprintf("%s - %s",
			time.Unix(p.StartDate, 0).Format("2006-01-02"),
			time.Unix(p.EndDate, 0).Format("2006-01-02")),
		TotalCalories:    totalCalories,
		AvgDailyCalories: avgDailyCalories,
		Macronutrients: macronutrients{
			TotalProtein: totalProtein,
			TotalCarbs:   totalCarbs,
			TotalFat:     totalFat,
			ProteinRatio: proteinRatio,
			CarbsRatio:   carbsRatio,
			FatRatio:     fatRatio,
		},
		MealDistribution: mealDist,
		GIAnalysis: giAnalysis{
			LowGIFoods:     uniqueStrings(lowGIFoods),
			MediumGIFoods:  uniqueStrings(mediumGIFoods),
			HighGIFoods:    uniqueStrings(highGIFoods),
			AvgGI:          avgGI,
			Recommendation: generateGIRecommendation(avgGI),
		},
		Recommendations: generateDietRecommendations(proteinRatio, carbsRatio, fatRatio, avgDailyCalories),
	}

	return toJSON(analysis), nil
}

// getFoodInfo 获取食物营养信息
func (t *DietManagerTool) getFoodInfo(foodName string) *foodInfo {
	// 简化的食物数据库
	foodDatabase := map[string]foodInfo{
		"牛奶": {
			Name:            "牛奶",
			CaloriesPer100g: 60,
			ProteinPer100g:  3.0,
			CarbsPer100g:    5.0,
			FatPer100g:      2.0,
			GI:              27,
			Category:        "low_gi",
		},
		"米饭": {
			Name:            "米饭",
			CaloriesPer100g: 130,
			ProteinPer100g:  2.7,
			CarbsPer100g:    28.0,
			FatPer100g:      0.3,
			GI:              83,
			Category:        "high_gi",
		},
		"黄瓜": {
			Name:            "黄瓜",
			CaloriesPer100g: 16,
			ProteinPer100g:  0.8,
			CarbsPer100g:    3.6,
			FatPer100g:      0.1,
			GI:              15,
			Category:        "low_gi",
		},
		"鸡蛋": {
			Name:            "鸡蛋",
			CaloriesPer100g: 155,
			ProteinPer100g:  13.0,
			CarbsPer100g:    1.1,
			FatPer100g:      11.0,
			GI:              42,
			Category:        "low_gi",
		},
		"鸡胸肉": {
			Name:            "鸡胸肉",
			CaloriesPer100g: 165,
			ProteinPer100g:  31.0,
			CarbsPer100g:    0.0,
			FatPer100g:      3.6,
			GI:              0,
			Category:        "low_gi",
		},
		"糙米": {
			Name:            "糙米",
			CaloriesPer100g: 111,
			ProteinPer100g:  2.6,
			CarbsPer100g:    23.0,
			FatPer100g:      0.9,
			GI:              55,
			Category:        "medium_gi",
		},
		"燕麦": {
			Name:            "燕麦",
			CaloriesPer100g: 389,
			ProteinPer100g:  16.9,
			CarbsPer100g:    66.3,
			FatPer100g:      6.9,
			GI:              55,
			Category:        "medium_gi",
		},
		"全麦面包": {
			Name:            "全麦面包",
			CaloriesPer100g: 247,
			ProteinPer100g:  13.0,
			CarbsPer100g:    41.0,
			FatPer100g:      3.4,
			GI:              51,
			Category:        "low_gi",
		},
	}

	if info, ok := foodDatabase[foodName]; ok {
		return &info
	}
	return nil
}

// uniqueStrings 去重字符串数组
func uniqueStrings(strs []string) []string {
	seen := make(map[string]bool)
	result := make([]string, 0)
	for _, s := range strs {
		if !seen[s] {
			seen[s] = true
			result = append(result, s)
		}
	}
	return result
}

// generateGIRecommendation 生成GI建议
func generateGIRecommendation(avgGI float64) string {
	if avgGI < 55 {
		return "整体GI控制良好，继续保持低GI饮食习惯"
	} else if avgGI < 70 {
		return "GI水平中等，建议增加低GI食物比例"
	}
	return "GI水平偏高，建议减少高GI食物摄入，增加蔬菜和全谷物"
}

// generateDietRecommendations 生成饮食建议
func generateDietRecommendations(proteinRatio, carbsRatio, fatRatio, avgDailyCalories float64) []string {
	recommendations := make([]string, 0)

	// 蛋白质建议
	if proteinRatio < 0.15 {
		recommendations = append(recommendations, "蛋白质摄入不足，建议增加优质蛋白质来源")
	} else if proteinRatio > 0.25 {
		recommendations = append(recommendations, "蛋白质摄入充足，继续保持")
	}

	// 碳水建议
	if carbsRatio > 0.60 {
		recommendations = append(recommendations, "碳水化合物比例偏高，建议适当减少")
	}

	// 脂肪建议
	if fatRatio > 0.30 {
		recommendations = append(recommendations, "脂肪摄入偏高，建议选择健康脂肪来源")
	}

	// 热量建议
	if avgDailyCalories > 2000 {
		recommendations = append(recommendations, "日均热量摄入偏高，建议适当控制")
	} else if avgDailyCalories < 1500 {
		recommendations = append(recommendations, "日均热量摄入偏低，确保营养均衡")
	}

	if len(recommendations) == 0 {
		recommendations = append(recommendations, "饮食结构合理，继续保持")
	}

	return recommendations
}

// toJSON 将对象转换为JSON字符串
func toJSON(obj interface{}) string {
	data, _ := json.Marshal(obj)
	return string(data)
}
