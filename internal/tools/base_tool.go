package tools

import (
	"context"
	"fmt"
	"time"

	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/schema"
	"github.com/google/uuid"
)

// KnowledgeCollector 知识收集器接口
// 用于避免循环导入
type KnowledgeCollector interface {
	AddItem(title, content, category, source string, tags []string, metadata map[string]interface{}) error
}

// globalKnowledgeCollector 全局知识收集器
var globalKnowledgeCollector KnowledgeCollector

// SetKnowledgeCollector 设置全局知识收集器
func SetKnowledgeCollector(collector KnowledgeCollector) {
	globalKnowledgeCollector = collector
}

// BaseTool 统一工具基础接口
// 所有工具都必须实现此接口，提供统一的工具管理和学习集成
type BaseTool interface {
	// Info 返回工具信息，包括名称、描述、参数等
	Info(ctx context.Context) (*schema.ToolInfo, error)

	// InvokableRun 执行工具调用
	// args: JSON格式的参数字符串
	// opts: 工具选项
	// 返回: 执行结果字符串和错误
	InvokableRun(ctx context.Context, args string, opts ...tool.Option) (string, error)

	// AfterRun 工具调用后的后处理
	// 用于集成学习引擎，收集工具调用结果和用户反馈
	// result: 工具执行结果
	// err: 执行错误（如果有）
	AfterRun(ctx context.Context, result string, err error)
}

// BaseToolImpl 提供默认的BaseTool实现
// 工具可以嵌入此结构体来获得默认行为
type BaseToolImpl struct {
	// EnableLearning 是否启用学习集成
	EnableLearning bool
	// Category 工具分类 (info_fetcher, location, health, learning, lifecycle)
	Category string
	// ToolName 工具名称
	ToolName string
}

// GetCategory 获取工具分类
func (b *BaseToolImpl) GetCategory() string {
	return b.Category
}

// GetToolName 获取工具名称
func (b *BaseToolImpl) GetToolName() string {
	return b.ToolName
}

// AfterRun 默认的后处理实现
// 子类可以重写此方法来实现自定义的学习逻辑
func (b *BaseToolImpl) AfterRun(ctx context.Context, result string, err error) {
	// 默认实现：如果启用了学习，将结果保存到知识库
	if !b.EnableLearning {
		return
	}

	// 如果执行失败，不保存结果
	if err != nil {
		return
	}

	// 通过全局知识收集器收集到知识库
	collectToKnowledgeBase(b.ToolName, b.Category, result)
}

// collectToKnowledgeBase 收集结果到知识库
func collectToKnowledgeBase(toolName, category, result string) {
	if globalKnowledgeCollector == nil {
		return
	}

	title := fmt.Sprintf("%s 调用结果", toolName)
	tags := []string{category}
	metadata := map[string]interface{}{
		"tool_name": toolName,
		"category":  category,
		"collected_at": time.Now().Format(time.RFC3339),
	}

	if err := globalKnowledgeCollector.AddItem(title, result, category, toolName, tags, metadata); err != nil {
		// 静默处理错误，避免影响主流程
		fmt.Printf("Failed to collect to knowledge base: %v\n", err)
	}
}

// GenerateUUID 生成UUID
func GenerateUUID() string {
	return uuid.New().String()
}
