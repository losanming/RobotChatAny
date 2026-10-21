package tools

import (
	"context"
	"fmt"
	"time"

	"github.com/bytedance/sonic"
	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/schema"
	"github.com/google/uuid"

	"new-robot-self/internal/brain"
)

// KnowledgeCollectorTool 知识收集工具
// 收集工具调用结果，提取关键信息，保存到知识库
type KnowledgeCollectorTool struct {
	BaseToolImpl
	knowledgeBase *brain.KnowledgeBase
}

type collectParams struct {
	Title    string                 `json:"title"`              // 标题
	Content  string                 `json:"content"`            // 内容
	URL      string                 `json:"url,omitempty"`      // URL
	Category string                 `json:"category"`           // 分类
	Source   string                 `json:"source"`             // 来源
	Tags     []string               `json:"tags,omitempty"`     // 标签
	Metadata map[string]interface{} `json:"metadata,omitempty"` // 元数据
}

type collectResult struct {
	Success    bool   `json:"success"`
	ID         string `json:"id"`
	Message    string `json:"message"`
	CollectedAt time.Time `json:"collected_at"`
}

func NewKnowledgeCollectorTool(knowledgeBase *brain.KnowledgeBase) *KnowledgeCollectorTool {
	return &KnowledgeCollectorTool{
		BaseToolImpl: BaseToolImpl{
			Category:       "learning",
			EnableLearning: false, // 知识收集工具本身不需要学习
			ToolName:       "knowledge_collector",
		},
		knowledgeBase: knowledgeBase,
	}
}

func (t *KnowledgeCollectorTool) Info(_ context.Context) (*schema.ToolInfo, error) {
	return &schema.ToolInfo{
		Name: "knowledge_collector",
		Desc: "收集信息到知识库，支持分类、标签和元数据。用于保存工具调用结果和用户交互信息",
		ParamsOneOf: schema.NewParamsOneOfByParams(map[string]*schema.ParameterInfo{
			"title":    {Desc: "信息标题", Type: schema.String, Required: true},
			"content":  {Desc: "信息内容", Type: schema.String, Required: true},
			"url":      {Desc: "相关URL（可选）", Type: schema.String, Required: false},
			"category": {Desc: "分类：tech, life, health, location, news", Type: schema.String, Required: true},
			"source":   {Desc: "来源工具名称", Type: schema.String, Required: true},
			"tags":     {Desc: "标签数组（可选）", Type: schema.Array, Required: false},
			"metadata": {Desc: "扩展元数据（可选）", Type: schema.Object, Required: false},
		}),
	}, nil
}

func (t *KnowledgeCollectorTool) InvokableRun(ctx context.Context, argumentsInJSON string, _ ...tool.Option) (string, error) {
	if t.knowledgeBase == nil {
		return "", fmt.Errorf("knowledge_collector: knowledge base not initialized")
	}

	p := collectParams{}
	if err := sonic.UnmarshalString(argumentsInJSON, &p); err != nil {
		return "", err
	}

	// 验证必填字段
	if p.Title == "" {
		return "", fmt.Errorf("title is required")
	}
	if p.Content == "" {
		return "", fmt.Errorf("content is required")
	}
	if p.Category == "" {
		return "", fmt.Errorf("category is required")
	}
	if p.Source == "" {
		return "", fmt.Errorf("source is required")
	}

	// 生成知识ID
	id := uuid.New().String()

	// 提取标签（如果没有提供）
	if len(p.Tags) == 0 {
		p.Tags = t.extractTags(p.Title, p.Content)
	}

	// 创建知识条目
	item := &brain.KnowledgeItem{
		ID:        id,
		Title:     p.Title,
		Content:   p.Content,
		URL:       p.URL,
		Category:  p.Category,
		Source:    p.Source,
		HotValue:  1, // 初始热度
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Tags:      p.Tags,
		Metadata:  p.Metadata,
	}

	// 添加到知识库
	if err := t.knowledgeBase.AddItem(item); err != nil {
		return "", fmt.Errorf("failed to add item to knowledge base: %w", err)
	}

	result := collectResult{
		Success:    true,
		ID:         id,
		Message:    "知识收集成功",
		CollectedAt: time.Now(),
	}

	b, _ := sonic.Marshal(result)
	return string(b), nil
}

// extractTags 提取标签
// 简化的标签提取逻辑，基于关键词匹配
func (t *KnowledgeCollectorTool) extractTags(title, content string) []string {
	tags := make([]string, 0)

	// 关键词字典
	keywords := map[string][]string{
		"tech":    {"科技", "AI", "人工智能", "互联网", "数码", "手机", "电脑", "芯片", "软件"},
		"health":  {"健康", "医疗", "运动", "饮食", "睡眠", "体重", "血糖", "血压"},
		"finance": {"财经", "股票", "投资", "经济", "金融", "基金", "银行", "保险"},
		"life":    {"生活", "美食", "旅游", "娱乐", "电影", "音乐", "购物", "家居"},
		"news":    {"新闻", "热点", "时事", "政策", "社会", "国际", "国内"},
	}

	// 检查标题和内容中的关键词
	text := title + " " + content
	for category, words := range keywords {
		for _, word := range words {
			if containsString(text, word) {
				// 避免重复标签
				if !containsStringSlice(tags, category) {
					tags = append(tags, category)
				}
				break
			}
		}
	}

	return tags
}

// containsString 检查字符串是否包含子串（简单实现）
func containsString(text, substr string) bool {
	return len(text) >= len(substr) && findSubstring(text, substr)
}

// containsStringSlice 检查字符串切片是否包含某个字符串
func containsStringSlice(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

// findSubstring 查找子串（简单实现）
func findSubstring(text, substr string) bool {
	for i := 0; i <= len(text)-len(substr); i++ {
		if text[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

// AfterRun 工具调用后的后处理
// 知识收集工具不需要额外的后处理
func (t *KnowledgeCollectorTool) AfterRun(ctx context.Context, result string, err error) {
	// 知识收集工具本身就是学习工具，不需要额外的学习集成
}
