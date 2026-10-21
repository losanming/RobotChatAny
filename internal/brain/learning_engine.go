package brain

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"new-robot-self/internal/resource"
)

// InformationLearningEngine 信息学习引擎
type InformationLearningEngine struct {
	brain           *Brain
	resourceFetcher *resource.Fetcher
	knowledgeBase   *KnowledgeBase
	preferenceLearner *PreferenceLearner
	evolutionEngine *EvolutionEngine
	autoScheduler   *AutoScheduler
	contentEnhancer *ContentEnhancer
	
	// 配置
	config          LearningConfig
	
	// 状态
	isLearning      bool
	learningMutex   sync.RWMutex
	lastLearningTime int64
}

// LearningConfig 学习配置
type LearningConfig struct {
	// 自动学习间隔（小时）
	AutoLearnInterval int `json:"auto_learn_interval"`
	
	// 每次学习的最大条目数
	MaxItemsPerLearn int `json:"max_items_per_learn"`
	
	// 启用的领域
	EnabledTypes []resource.ResourceType `json:"enabled_types"`
	
	// 是否启用自动学习
	EnableAutoLearn bool `json:"enable_auto_learn"`
	
	// 资源配置文件路径
	ResourceConfigPath string `json:"resource_config_path"`
}

// DefaultLearningConfig 默认学习配置
func DefaultLearningConfig() LearningConfig {
	return LearningConfig{
		AutoLearnInterval:  6, // 每6小时学习一次
		MaxItemsPerLearn:   20, // 每次最多学习20条
		EnabledTypes:       resource.GetAllType(),
		EnableAutoLearn:    true,
		ResourceConfigPath: "asset/resources.json",
	}
}

// NewInformationLearningEngine 创建信息学习引擎
func NewInformationLearningEngine(brain *Brain, config LearningConfig) (*InformationLearningEngine, error) {
	// 创建资源获取器
	resourceFetcher, err := resource.NewFetcher(config.ResourceConfigPath)
	if err != nil {
		return nil, fmt.Errorf("创建资源获取器失败: %v", err)
	}
	
	// 创建知识库
	knowledgeBase := NewKnowledgeBase()
	
	// 创建偏好学习器
	preferenceLearner := NewPreferenceLearner()
	
	// 创建进化引擎
	evolutionEngine := NewEvolutionEngine()
	
	// 创建自动调度器
	autoScheduler := NewAutoScheduler(config.AutoLearnInterval)
	
	// 创建内容增强器
	contentEnhancer := NewContentEnhancer()
	
	engine := &InformationLearningEngine{
		brain:             brain,
		resourceFetcher:   resourceFetcher,
		knowledgeBase:     knowledgeBase,
		preferenceLearner: preferenceLearner,
		evolutionEngine:   evolutionEngine,
		autoScheduler:     autoScheduler,
		contentEnhancer:  contentEnhancer,
		config:           config,
		lastLearningTime:  0,
	}
	
	// 启动自动学习
	if config.EnableAutoLearn {
		go engine.startAutoLearning()
	}
	
	return engine, nil
}

// LearnOnce 执行一次学习
func (e *InformationLearningEngine) LearnOnce(ctx context.Context) (*LearningResult, error) {
	e.learningMutex.Lock()
	defer e.learningMutex.Unlock()
	
	startTime := time.Now()
	log.Println("开始执行信息学习...")
	
	result := &LearningResult{
		StartTime: startTime,
		Status:    "learning",
		ByType:    make(map[resource.ResourceType]int),
	}
	
	// 1. 获取各领域资源
	allResources := make(map[resource.ResourceType][]resource.Resource)
	for _, resType := range e.config.EnabledTypes {
		resources, err := e.resourceFetcher.GetRandomResources([]resource.ResourceType{resType}, e.config.MaxItemsPerLearn/len(e.config.EnabledTypes))
		if err != nil {
			log.Printf("获取%s类型资源失败: %v", resType, err)
			continue
		}
		if len(resources) > 0 {
			allResources[resType] = resources
			result.TotalFetched += len(resources)
		}
	}
	
	// 2. 处理和学习知识
	for resType, resources := range allResources {
		processedCount := 0
		for _, res := range resources {
			// 质量评估
			if !e.evaluateQuality(res) {
				continue
			}
			
			// 转换为知识条目
			knowledgeItem := e.convertToKnowledgeItem(res, resType)
			
			// 添加到知识库
			if err := e.knowledgeBase.AddItem(knowledgeItem); err != nil {
				log.Printf("添加知识条目失败: %v", err)
				continue
			}
			
			processedCount++
			
			// 学习用户偏好
			e.preferenceLearner.LearnFromResource(res)
		}
		
		result.ProcessedCount += processedCount
		result.ByType[resType] = processedCount
	}
	
	// 3. 更新进化引擎
	e.evolutionEngine.RecordLearning(result)
	
	// 4. 更新状态
	e.lastLearningTime = time.Now().Unix()
	result.EndTime = time.Now()
	result.Duration = result.EndTime.Sub(startTime)
	result.Status = "completed"
	
	log.Printf("信息学习完成: 获取%d条，处理%d条，耗时%v", 
		result.TotalFetched, result.ProcessedCount, result.Duration)
	
	return result, nil
}

// evaluateQuality 评估信息质量
func (e *InformationLearningEngine) evaluateQuality(res resource.Resource) bool {
	// 基本质量检查
	if res.Title == "" || res.Content == "" {
		return false
	}
	
	// 内容长度检查
	if len(res.Content) < 10 || len(res.Content) > 10000 {
		return false
	}
	
	// 热度检查（可选）
	if res.HotValue > 0 && res.HotValue < 100 {
		// 过滤低热度内容
		return false
	}
	
	return true
}

// convertToKnowledgeItem 转换为知识条目
func (e *InformationLearningEngine) convertToKnowledgeItem(res resource.Resource, resType resource.ResourceType) *KnowledgeItem {
	// 初始内容
	content := res.Content
	
	// 如果内容太短，尝试增强
	if len(content) < 100 && res.URL != "" {
		enhancedContent, err := e.contentEnhancer.EnhanceContent(res.Title, res.URL, content)
		if err == nil && enhancedContent != "" {
			content = enhancedContent
		}
	}
	
	return &KnowledgeItem{
		ID:          generateKnowledgeID(res.Title, res.URL),
		Title:       res.Title,
		Content:     content,
		URL:         res.URL,
		Category:    string(resType),
		Source:      "auto_learning",
		HotValue:    res.HotValue,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		Tags:        extractTags(res.Title, res.Content),
		Metadata: map[string]interface{}{
			"types":      res.Types,
			"timestamp":  res.Timestamp,
		},
	}
}

// startAutoLearning 启动自动学习
func (e *InformationLearningEngine) startAutoLearning() {
	ticker := time.NewTicker(time.Duration(e.config.AutoLearnInterval) * time.Hour)
	defer ticker.Stop()
	
	for {
		select {
		case <-ticker.C:
			ctx := context.Background()
			if _, err := e.LearnOnce(ctx); err != nil {
				log.Printf("自动学习失败: %v", err)
			}
		}
	}
}

// GetKnowledgeBase 获取知识库
func (e *InformationLearningEngine) GetKnowledgeBase() *KnowledgeBase {
	return e.knowledgeBase
}

// GetPreferences 获取用户偏好
func (e *InformationLearningEngine) GetPreferences() *UserPreferences {
	return e.preferenceLearner.GetPreferences()
}

// GetEvolutionReport 获取进化报告
func (e *InformationLearningEngine) GetEvolutionReport() *EvolutionReport {
	return e.evolutionEngine.GenerateReport()
}

// GetContextBuilder 获取上下文构建器
func (e *InformationLearningEngine) GetContextBuilder() *ContextBuilder {
	return NewContextBuilder(e.knowledgeBase, e.preferenceLearner.GetPreferences())
}

// LearningResult 学习结果
type LearningResult struct {
	StartTime       time.Time
	EndTime         time.Time
	Duration        time.Duration
	Status          string
	TotalFetched    int
	ProcessedCount  int
	ByType         map[resource.ResourceType]int
	Error          error
}
