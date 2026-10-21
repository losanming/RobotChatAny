package tools

import (
	"fmt"
	"time"

	"new-robot-self/internal/brain"
	"new-robot-self/internal/tools/analysis"
	"new-robot-self/internal/tools/automation"
	"new-robot-self/internal/tools/cloud"
	"new-robot-self/internal/tools/communication"
	"new-robot-self/internal/tools/finance"
	"new-robot-self/internal/tools/health"
	"new-robot-self/internal/tools/intelligence"
	"new-robot-self/internal/tools/mobile"
	"new-robot-self/internal/tools/multimedia"
	"new-robot-self/internal/tools/personal"
	"new-robot-self/internal/tools/security"
	"new-robot-self/internal/tools/smart_schedule"
	"new-robot-self/internal/tools/social"
	"new-robot-self/internal/tools/utility"
	"new-robot-self/internal/tools/visualization"
	"new-robot-self/internal/tools/voice"
	"new-robot-self/internal/tools/work"
)

// RegisterAllTools 注册所有工具到全局注册中心
// knowledgeBase: 知识库实例，用于学习集成
func RegisterAllTools(knowledgeBase *brain.KnowledgeBase) error {
	// 设置全局知识收集器
	if knowledgeBase != nil {
		SetKnowledgeCollector(&knowledgeBaseAdapter{kb: knowledgeBase})
	}

	// 注册位置服务类工具
	if err := registerLocationTools(); err != nil {
		return fmt.Errorf("failed to register location tools: %w", err)
	}

	// 注册健康管理类工具
	if err := registerHealthTools(); err != nil {
		return fmt.Errorf("failed to register health tools: %w", err)
	}

	// 注册信息获取类工具
	if err := registerInfoFetcherTools(); err != nil {
		return fmt.Errorf("failed to register info fetcher tools: %w", err)
	}

	// 注册学习进化类工具
	if err := registerLearningTools(knowledgeBase); err != nil {
		return fmt.Errorf("failed to register learning tools: %w", err)
	}

	// 注册个人管理类工具
	if err := registerPersonalTools(); err != nil {
		return fmt.Errorf("failed to register personal tools: %w", err)
	}

	// 注册信息分析类工具
	if err := registerAnalysisTools(); err != nil {
		return fmt.Errorf("failed to register analysis tools: %w", err)
	}

	// 注册通用能力类工具
	if err := registerUtilityTools(); err != nil {
		return fmt.Errorf("failed to register utility tools: %w", err)
	}

	// 注册通信类工具
	if err := registerCommunicationTools(); err != nil {
		return fmt.Errorf("failed to register communication tools: %w", err)
	}

	// 注册多媒体类工具
	if err := registerMultimediaTools(); err != nil {
		return fmt.Errorf("failed to register multimedia tools: %w", err)
	}

	// 注册安全类工具
	if err := registerSecurityTools(); err != nil {
		return fmt.Errorf("failed to register security tools: %w", err)
	}

	// 注册金融管理类工具
	if err := registerFinanceTools(); err != nil {
		return fmt.Errorf("failed to register finance tools: %w", err)
	}

	// 注册自动化类工具
	if err := registerAutomationTools(); err != nil {
		return fmt.Errorf("failed to register automation tools: %w", err)
	}

	// 注册云服务类工具
	if err := registerCloudTools(); err != nil {
		return fmt.Errorf("failed to register cloud tools: %w", err)
	}

	// 注册社交媒体类工具
	if err := registerSocialTools(); err != nil {
		return fmt.Errorf("failed to register social tools: %w", err)
	}

	// 注册智能化类工具
	if err := registerIntelligenceTools(); err != nil {
		return fmt.Errorf("failed to register intelligence tools: %w", err)
	}

	// 注册移动端类工具
	if err := registerMobileTools(); err != nil {
		return fmt.Errorf("failed to register mobile tools: %w", err)
	}

	// 注册健康增强类工具
	if err := registerHealthEnhancedTools(); err != nil {
		return fmt.Errorf("failed to register health enhanced tools: %w", err)
	}

	// 注册工作助手类工具
	if err := registerWorkTools(); err != nil {
		return fmt.Errorf("failed to register work tools: %w", err)
	}

	// 注册数据可视化类工具
	if err := registerVisualizationTools(); err != nil {
		return fmt.Errorf("failed to register visualization tools: %w", err)
	}

	// 注册智能日程类工具
	if err := registerSmartScheduleTools(); err != nil {
		return fmt.Errorf("failed to register smart schedule tools: %w", err)
	}

	// 注册语音类工具
	if err := registerVoiceTools(); err != nil {
		return fmt.Errorf("failed to register voice tools: %w", err)
	}

	return nil
}

// registerLocationTools 注册位置服务类工具
func registerLocationTools() error {
	tools := map[string]BaseTool{
		"reverse_geocode": NewReverseGeocodeTool(),
		"around_search":   NewAroundSearchTool(),
		"weather_query":  NewWeatherQueryTool(),
		"city_adcode":     NewCityAdcodeTool(),
		"navigation":      NewNavigationTool(),
	}

	config := LocationConfig()

	for name, tool := range tools {
		if err := Register(name, tool, config); err != nil {
			return err
		}
	}

	return nil
}

// registerHealthTools 注册健康管理类工具
func registerHealthTools() error {
	tools := map[string]BaseTool{
		"health_manager":   NewHealthManagerTool(),
		"exercise_advisor": NewExerciseAdvisorTool(),
	}

	config := HealthConfig()

	for name, tool := range tools {
		if err := Register(name, tool, config); err != nil {
			return err
		}
	}

	return nil
}

// registerInfoFetcherTools 注册信息获取类工具
func registerInfoFetcherTools() error {
	tools := map[string]BaseTool{
		"rss_hot_news":    NewRSSHotNewsTool(),
		"smart_search":    NewSmartSearchTool(),
		"content_fetcher": NewContentFetcherTool(),
	}

	config := InfoFetcherConfig()

	for name, tool := range tools {
		if err := Register(name, tool, config); err != nil {
			return err
		}
	}

	return nil
}

// registerLearningTools 注册学习进化类工具
func registerLearningTools(knowledgeBase *brain.KnowledgeBase) error {
	tools := map[string]BaseTool{
		"knowledge_collector":  NewKnowledgeCollectorTool(knowledgeBase),
		"preference_learner":   NewPreferenceLearnerTool(),
	}

	config := LearningConfig()

	for name, tool := range tools {
		if err := Register(name, tool, config); err != nil {
			return err
		}
	}

	return nil
}

// knowledgeBaseAdapter 知识库适配器
// 实现KnowledgeCollector接口，避免循环导入
type knowledgeBaseAdapter struct {
	kb *brain.KnowledgeBase
}

func (a *knowledgeBaseAdapter) AddItem(title, content, category, source string, tags []string, metadata map[string]interface{}) error {
	item := &brain.KnowledgeItem{
		ID:        GenerateUUID(),
		Title:     title,
		Content:   content,
		Category:  category,
		Source:    source,
		HotValue:  1,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Tags:      tags,
		Metadata:  metadata,
	}
	return a.kb.AddItem(item)
}

// registerPersonalTools 注册个人管理类工具
func registerPersonalTools() error {
	tools := map[string]BaseTool{
		"calendar_manager": personal.NewCalendarManager(),
		"todo_manager":     personal.NewTodoManager(),
		"time_manager":     personal.NewTimeManager(),
	}

	config := PersonalConfig()

	for name, tool := range tools {
		if err := Register(name, tool, config); err != nil {
			return err
		}
	}

	return nil
}

// registerAnalysisTools 注册信息分析类工具
func registerAnalysisTools() error {
	tools := map[string]BaseTool{
		"file_processor":   analysis.NewFileProcessor(),
		"data_analyzer":    analysis.NewDataAnalyzer(),
		"report_generator": analysis.NewReportGenerator(),
	}

	config := AnalysisConfig()

	for name, tool := range tools {
		if err := Register(name, tool, config); err != nil {
			return err
		}
	}

	return nil
}

// registerUtilityTools 注册通用能力类工具
func registerUtilityTools() error {
	tools := map[string]BaseTool{
		"code_executor":     utility.NewCodeExecutor(),
		"translator":        utility.NewTranslator(),
		"unit_converter":    utility.NewUnitConverter(),
		"calculator":        utility.NewCalculator(),
		"timezone_converter": utility.NewTimezoneConverter(),
	}

	config := UtilityConfig()

	for name, tool := range tools {
		if err := Register(name, tool, config); err != nil {
			return err
		}
	}

	return nil
}

// registerCommunicationTools 注册通信类工具
func registerCommunicationTools() error {
	tools := map[string]BaseTool{
		"email_manager":        communication.NewEmailManager(),
		"notification_manager": communication.NewNotificationManager(),
	}

	config := CommunicationConfig()

	for name, tool := range tools {
		if err := Register(name, tool, config); err != nil {
			return err
		}
	}

	return nil
}

// registerMultimediaTools 注册多媒体类工具
func registerMultimediaTools() error {
	tools := map[string]BaseTool{
		"image_processor": multimedia.NewImageProcessor(),
		"audio_processor": multimedia.NewAudioProcessor(),
		"ocr_processor":   multimedia.NewOCRProcessor(),
	}

	config := MultimediaConfig()

	for name, tool := range tools {
		if err := Register(name, tool, config); err != nil {
			return err
		}
	}

	return nil
}

// registerSecurityTools 注册安全类工具
func registerSecurityTools() error {
	tools := map[string]BaseTool{
		"security_manager":   security.NewSecurityManager(),
		"encryption_manager": security.NewEncryptionManager(),
	}

	config := SecurityConfig()

	for name, tool := range tools {
		if err := Register(name, tool, config); err != nil {
			return err
		}
	}

	return nil
}

// registerFinanceTools 注册金融管理类工具
func registerFinanceTools() error {
	tools := map[string]BaseTool{
		"finance_manager": finance.NewFinanceManager(),
		"budget_manager":   finance.NewBudgetManager(),
	}

	config := FinanceConfig()

	for name, tool := range tools {
		if err := Register(name, tool, config); err != nil {
			return err
		}
	}

	return nil
}

// registerAutomationTools 注册自动化类工具
func registerAutomationTools() error {
	tools := map[string]BaseTool{
		"workflow_manager":   automation.NewWorkflowManager(),
		"scheduler_manager": automation.NewSchedulerManager(),
	}

	config := AutomationConfig()

	for name, tool := range tools {
		if err := Register(name, tool, config); err != nil {
			return err
		}
	}

	return nil
}

// registerCloudTools 注册云服务类工具
func registerCloudTools() error {
	tools := map[string]BaseTool{
		"cloud_storage_manager": cloud.NewCloudStorageManager(),
		"api_integrator":        cloud.NewAPIIntegrator(),
	}

	config := CloudConfig()

	for name, tool := range tools {
		if err := Register(name, tool, config); err != nil {
			return err
		}
	}

	return nil
}

// registerSocialTools 注册社交媒体类工具
func registerSocialTools() error {
	tools := map[string]BaseTool{
		"social_media_manager": social.NewSocialMediaManager(),
		"content_publisher":    social.NewContentPublisher(),
	}

	config := SocialConfig()

	for name, tool := range tools {
		if err := Register(name, tool, config); err != nil {
			return err
		}
	}

	return nil
}

// registerIntelligenceTools 注册智能化类工具
func registerIntelligenceTools() error {
	tools := map[string]BaseTool{
		"knowledge_graph_manager": intelligence.NewKnowledgeGraphManager(),
		"reasoner":                intelligence.NewReasoner(),
		"collaboration_manager":   intelligence.NewCollaborationManager(),
		"sharing_manager":         intelligence.NewSharingManager(),
		"emotion_analyzer":        intelligence.NewEmotionAnalyzer(),
		"sentiment_processor":     intelligence.NewSentimentProcessor(),
		"recommendation_engine":   intelligence.NewRecommendationEngine(),
		"behavior_predictor":      intelligence.NewBehaviorPredictor(),
	}

	config := IntelligenceConfig()

	for name, tool := range tools {
		if err := Register(name, tool, config); err != nil {
			return err
		}
	}

	return nil
}

// registerMobileTools 注册移动端类工具
func registerMobileTools() error {
	tools := map[string]BaseTool{
		"mobile_adapter":    mobile.NewMobileAdapter(),
		"push_notification": mobile.NewPushNotification(),
	}

	config := MobileConfig()

	for name, tool := range tools {
		if err := Register(name, tool, config); err != nil {
			return err
		}
	}

	return nil
}

// registerHealthEnhancedTools 注册健康增强类工具
func registerHealthEnhancedTools() error {
	tools := map[string]BaseTool{
		"health_tracker":  health.NewHealthTrackerTool(),
		"diet_manager":    health.NewDietManagerTool(),
		"smart_reminder":  health.NewSmartReminderTool(),
	}

	config := HealthConfig()

	for name, tool := range tools {
		if err := Register(name, tool, config); err != nil {
			return err
		}
	}

	return nil
}

// registerWorkTools 注册工作助手类工具
func registerWorkTools() error {
	tools := map[string]BaseTool{
		"work_assistant": work.NewWorkAssistantTool(),
	}

	config := WorkConfig()

	for name, tool := range tools {
		if err := Register(name, tool, config); err != nil {
			return err
		}
	}

	return nil
}

// registerVisualizationTools 注册数据可视化类工具
func registerVisualizationTools() error {
	tools := map[string]BaseTool{
		"data_visualizer": visualization.NewDataVisualizerTool(),
	}

	config := VisualizationConfig()

	for name, tool := range tools {
		if err := Register(name, tool, config); err != nil {
			return err
		}
	}

	return nil
}

// registerSmartScheduleTools 注册智能日程类工具
func registerSmartScheduleTools() error {
	tools := map[string]BaseTool{
		"smart_scheduler": smart_schedule.NewSmartSchedulerTool(),
	}

	config := SmartScheduleConfig()

	for name, tool := range tools {
		if err := Register(name, tool, config); err != nil {
			return err
		}
	}

	return nil
}

// registerVoiceTools 注册语音类工具
func registerVoiceTools() error {
	tools := map[string]BaseTool{
		"speech_to_text": voice.NewSpeechToTextTool(),
	}

	config := VoiceConfig()

	for name, tool := range tools {
		if err := Register(name, tool, config); err != nil {
			return err
		}
	}

	return nil
}
