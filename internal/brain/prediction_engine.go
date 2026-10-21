package brain

import (
	"context"
	"fmt"
	"math"
	"sort"
	"time"

	"new-robot-self/global"
	"new-robot-self/pkg/rag"
)

// Prediction 预测结果
type Prediction struct {
	Type         string    `json:"type"`          // 预测类型: health, work, life
	Confidence   float64   `json:"confidence"`     // 置信度 0-1
	Content      string    `json:"content"`       // 预测内容
	Action       string    `json:"action"`        // 建议动作
	Priority     int       `json:"priority"`      // 优先级 1-5
	Timing       time.Time `json:"timing"`        // 建议时机
	Reason       string    `json:"reason"`        // 预测理由
	RelatedData  map[string]interface{} `json:"related_data"`
}

// PredictionEngine 智能预测引擎
type PredictionEngine struct {
	brain        *Brain
	ragService   *rag.Service
	currentState *WorldState
}

// NewPredictionEngine 创建预测引擎
func NewPredictionEngine(brain *Brain) *PredictionEngine {
	return &PredictionEngine{
		brain:      brain,
		ragService: global.RAGService,
	}
}

// PredictNextActions 预测下一步可能的行动
func (p *PredictionEngine) PredictNextActions(ctx context.Context) ([]*Prediction, error) {
	p.brain.mu.RLock()
	p.currentState = p.brain.state
	p.brain.mu.RUnlock()

	predictions := []*Prediction{}

	// 1. 健康相关预测
	healthPredictions := p.predictHealthNeeds(ctx)
	predictions = append(predictions, healthPredictions...)

	// 2. 工作相关预测
	workPredictions := p.predictWorkNeeds(ctx)
	predictions = append(predictions, workPredictions...)

	// 3. 生活相关预测
	lifePredictions := p.predictLifeNeeds(ctx)
	predictions = append(predictions, lifePredictions...)

	// 4. 按优先级排序
	p.sortPredictionsByPriority(predictions)

	return predictions, nil
}

// predictHealthNeeds 预测健康相关需求
func (p *PredictionEngine) predictHealthNeeds(ctx context.Context) []*Prediction {
	if p.currentState.Health == nil {
		return nil
	}

	predictions := []*Prediction{}
	now := time.Now()
	hour := now.Hour()
	health := p.currentState.Health

	// 1. 运动时间预测
	if p.shouldExerciseNow(hour) {
		// 使用RAG增强运动建议
		content := p.enhanceWithRAG(ctx, "运动建议", "根据您的运动计划，现在是运动时间。建议进行室内慢跑1小时。")
		
		predictions = append(predictions, &Prediction{
			Type:       "health",
			Confidence: 0.8,
			Content:    content,
			Action:     "exercise_reminder",
			Priority:   4,
			Timing:     now,
			Reason:     "运动计划时间匹配",
			RelatedData: map[string]interface{}{
				"exercise_plan": health.ExecDesc,
				"current_time":  hour,
			},
		})
	} else {
		// 测试用：如果不是运动时间，也返回一个低优先级的运动建议
		content := p.enhanceWithRAG(ctx, "运动习惯", fmt.Sprintf("虽然现在不是最佳运动时间，但建议您保持运动习惯。您的运动计划：%s", health.ExecDesc))
		
		predictions = append(predictions, &Prediction{
			Type:       "health",
			Confidence: 0.5,
			Content:    content,
			Action:     "exercise_suggestion",
			Priority:   1,
			Timing:     now,
			Reason:     "运动习惯提醒",
			RelatedData: map[string]interface{}{
				"exercise_plan": health.ExecDesc,
				"current_time":  hour,
			},
		})
	}

	// 2. 血糖监测预测
	if p.shouldMonitorBloodSugar(health) {
		content := p.enhanceWithRAG(ctx, "血糖管理", "根据您的健康数据，建议监测血糖水平。您的血糖控制需要关注。")
		
		predictions = append(predictions, &Prediction{
			Type:       "health",
			Confidence: 0.7,
			Content:    content,
			Action:     "blood_sugar_reminder",
			Priority:   3,
			Timing:     now,
			Reason:     "血糖管理需要",
			RelatedData: map[string]interface{}{
				"blood_sugar": health.BloodSugar,
				"last_updated": health.LastUpdatedTS,
			},
		})
	}

	// 3. 饮食建议预测
	if p.shouldProvideDietAdvice(hour) {
		predictions = append(predictions, &Prediction{
			Type:       "health",
			Confidence: 0.6,
			Content:    "根据您的糖尿病管理需求，建议注意饮食控制。避免高糖食物，增加蔬菜摄入。",
			Action:     "diet_advice",
			Priority:   3,
			Timing:     now,
			Reason:     "饮食管理建议",
			RelatedData: map[string]interface{}{
				"diet_preference": "diabetes_friendly",
				"current_time": hour,
			},
		})
	}

	// 4. 体重管理预测
	if p.shouldCheckWeight(health) {
		targetWeight := 68.0 // 目标体重
		currentWeight := health.WeightKg
		weightLoss := currentWeight - targetWeight
		
		predictions = append(predictions, &Prediction{
			Type:       "health",
			Confidence: 0.75,
			Content:    fmt.Sprintf("您的当前体重%.1fkg，距离目标体重%.1fkg还差%.1fkg。建议继续保持运动和饮食控制。", 
				currentWeight, targetWeight, weightLoss),
			Action:     "weight_management",
			Priority:   2,
			Timing:     now,
			Reason:     "体重管理跟踪",
			RelatedData: map[string]interface{}{
				"current_weight": currentWeight,
				"target_weight": targetWeight,
				"weight_loss": weightLoss,
			},
		})
	}

	return predictions
}

// predictWorkNeeds 预测工作相关需求
func (p *PredictionEngine) predictWorkNeeds(ctx context.Context) []*Prediction {
	predictions := []*Prediction{}
	now := time.Now()
	hour := now.Hour()
	weekday := now.Weekday()

	// 1. 工作专注时间预测
	if p.isWorkTime(hour, weekday) {
		if p.shouldStartFocusTime(hour) {
			predictions = append(predictions, &Prediction{
				Type:       "work",
				Confidence: 0.7,
				Content:    "现在是工作专注时间。建议开始25分钟的专注工作，然后休息5分钟。",
				Action:     "focus_time_reminder",
				Priority:   3,
				Timing:     now,
				Reason:     "工作专注时间",
				RelatedData: map[string]interface{}{
					"current_hour": hour,
					"weekday": weekday,
				},
			})
		}

		// 2. 休息提醒预测
		if p.shouldTakeBreak(hour) {
			predictions = append(predictions, &Prediction{
				Type:       "work",
				Confidence: 0.8,
				Content:    "您已经工作较长时间了。建议起身活动5-10分钟，保护颈椎和腰椎。",
				Action:     "break_reminder",
				Priority:   4,
				Timing:     now,
				Reason:     "长时间工作需要休息",
				RelatedData: map[string]interface{}{
					"work_duration": "long",
					"current_hour": hour,
				},
			})
		}
	}

	// 3. 任务规划预测
	if p.shouldPlanDailyTasks(hour) {
		predictions = append(predictions, &Prediction{
			Type:       "work",
			Confidence: 0.6,
			Content:    "建议规划今天的三个重要任务。优先处理最重要的事情。",
			Action:     "task_planning",
			Priority:   2,
			Timing:     now,
			Reason:     "每日任务规划",
			RelatedData: map[string]interface{}{
				"current_hour": hour,
				"weekday": weekday,
			},
		})
	}

	return predictions
}

// predictLifeNeeds 预测生活相关需求
func (p *PredictionEngine) predictLifeNeeds(ctx context.Context) []*Prediction {
	predictions := []*Prediction{}
	hour := time.Now().Hour()

	// 1. 天气相关预测
	if p.currentState.Env != nil && p.currentState.Env.Awareness.Weather != nil {
		weather := p.currentState.Env.Awareness.Weather
		if p.shouldCheckWeather(weather) {
			predictions = append(predictions, &Prediction{
				Type:       "life",
				Confidence: 0.7,
				Content:    "根据当前天气情况，建议注意出行安全。天气变化较大时请及时增减衣物。",
				Action:     "weather_advice",
				Priority:   2,
				Timing:     time.Now(),
				Reason:     "天气变化提醒",
				RelatedData: map[string]interface{}{
					"weather": weather,
				},
			})
		}
	}

	// 2. 作息管理预测
	if p.shouldCheckSleepSchedule(hour) {
		predictions = append(predictions, &Prediction{
			Type:       "life",
			Confidence: 0.65,
			Content:    "建议保持规律作息。充足的睡眠对控制血糖和体重都很重要。",
			Action:     "sleep_schedule_advice",
			Priority:   2,
			Timing:     time.Now(),
			Reason:     "作息管理建议",
			RelatedData: map[string]interface{}{
				"current_hour": hour,
			},
		})
	}

	// 3. 周末活动建议
	if time.Now().Weekday() == time.Saturday || time.Now().Weekday() == time.Sunday {
		predictions = append(predictions, &Prediction{
			Type:       "life",
			Confidence: 0.5,
			Content:    "周末是放松和调整的好时机。建议安排一些户外活动或兴趣爱好。",
			Action:     "weekend_activity_suggestion",
			Priority:   1,
			Timing:     time.Now(),
			Reason:     "周末活动建议",
			RelatedData: map[string]interface{}{
				"is_weekend": true,
			},
		})
	}

	return predictions
}

// shouldExerciseNow 判断是否应该运动
func (p *PredictionEngine) shouldExerciseNow(hour int) bool {
	// 假设运动时间：晚上19-21点
	return hour >= 19 && hour <= 21
}

// shouldExerciseNowForTest 测试用：强制返回true
func (p *PredictionEngine) shouldExerciseNowForTest(hour int) bool {
	return true
}

// enhanceWithRAG 使用RAG增强建议内容
func (p *PredictionEngine) enhanceWithRAG(ctx context.Context, category, baseContent string) string {
	if p.ragService == nil {
		return baseContent
	}

	// 构建查询
	query := fmt.Sprintf("%s 威猛先生", category)
	
	// 搜索相关文档
	searchCtx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()
	
	results, err := p.ragService.Search(searchCtx, rag.SearchQuery{
		Query: query,
		TopK:  2,
	})
	
	if err != nil || len(results) == 0 {
		// RAG失败时返回基础内容
		return baseContent
	}
	
	// 如果找到相关文档，增强内容
	enhancedContent := baseContent
	if len(results) > 0 && results[0].Score > 0.6 {
		enhancedContent = fmt.Sprintf("%s\n\n💡 个人参考: %s", baseContent, results[0].Content)
	}
	
	return enhancedContent
}

// shouldMonitorBloodSugar 判断是否应该监测血糖
func (p *PredictionEngine) shouldMonitorBloodSugar(health *HealthState) bool {
	// 如果距离上次血糖监测超过8小时
	if health.LastUpdatedTS > 0 {
		lastUpdate := time.Unix(health.LastUpdatedTS, 0)
		hoursSinceUpdate := time.Since(lastUpdate).Hours()
		return hoursSinceUpdate > 8
	}
	return true
}

// shouldProvideDietAdvice 判断是否应该提供饮食建议
func (p *PredictionEngine) shouldProvideDietAdvice(hour int) bool {
	// 在饭点附近提供饮食建议
	return hour == 7 || hour == 12 || hour == 18
}

// shouldCheckWeight 判断是否应该检查体重
func (p *PredictionEngine) shouldCheckWeight(health *HealthState) bool {
	// 每周检查一次体重
	if health.LastUpdatedTS > 0 {
		lastUpdate := time.Unix(health.LastUpdatedTS, 0)
		daysSinceUpdate := int(time.Since(lastUpdate).Hours() / 24)
		return daysSinceUpdate >= 7
	}
	return true
}

// isWorkTime 判断是否是工作时间
func (p *PredictionEngine) isWorkTime(hour int, weekday time.Weekday) bool {
	// 周一到周五，9-18点为工作时间
	if weekday == time.Saturday || weekday == time.Sunday {
		return false
	}
	return hour >= 9 && hour <= 18
}

// shouldStartFocusTime 判断是否应该开始专注时间
func (p *PredictionEngine) shouldStartFocusTime(hour int) bool {
	// 整点时刻开始专注时间
	return hour%1 == 0
}

// shouldTakeBreak 判断是否应该休息
func (p *PredictionEngine) shouldTakeBreak(hour int) bool {
	// 工作2小时后提醒休息
	return hour == 11 || hour == 15 || hour == 17
}

// shouldPlanDailyTasks 判断是否应该规划每日任务
func (p *PredictionEngine) shouldPlanDailyTasks(hour int) bool {
	// 早上9点规划任务
	return hour == 9
}

// shouldCheckWeather 判断是否应该检查天气
func (p *PredictionEngine) shouldCheckWeather(weather interface{}) bool {
	// 天气有变化时提醒
	return weather != nil
}

// shouldCheckSleepSchedule 判断是否应该检查作息
func (p *PredictionEngine) shouldCheckSleepSchedule(hour int) bool {
	// 睡前和起床时提醒
	return hour == 22 || hour == 7
}

// sortPredictionsByPriority 按优先级排序预测
func (p *PredictionEngine) sortPredictionsByPriority(predictions []*Prediction) {
	sort.Slice(predictions, func(i, j int) bool {
		return predictions[i].Priority > predictions[j].Priority
	})
}

// EnhanceWithRAG 使用RAG增强预测
func (p *PredictionEngine) EnhanceWithRAG(ctx context.Context, prediction *Prediction) (*Prediction, error) {
	if p.ragService == nil {
		return prediction, nil
	}

	// 基于预测内容进行RAG检索
	query := fmt.Sprintf("%s %s", prediction.Type, prediction.Content)
	results, err := p.ragService.Search(ctx, rag.SearchQuery{
		Query: query,
		TopK:  3,
	})
	if err != nil {
		return prediction, err
	}

	if len(results) > 0 {
		// 增强预测内容
		enhancedContent := fmt.Sprintf("%s\n\n参考信息: %s", prediction.Content, results[0].Content)
		prediction.Content = enhancedContent
		prediction.Confidence = math.Min(prediction.Confidence+0.1, 1.0)
	}

	return prediction, nil
}

// ExecutePrediction 执行预测建议
func (p *PredictionEngine) ExecutePrediction(ctx context.Context, prediction *Prediction) error {
	// 这里可以根据预测类型执行不同的操作
	// 目前只是打印，后续可以集成到通知系统
	
	fmt.Printf("🤖 主动建议 [%s]: %s\n", prediction.Type, prediction.Content)
	fmt.Printf("   置信度: %.1f%%, 优先级: %d\n", prediction.Confidence*100, prediction.Priority)
	fmt.Printf("   理由: %s\n", prediction.Reason)
	
	return nil
}
