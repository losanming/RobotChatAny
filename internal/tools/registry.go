package tools

import (
	"context"
	"fmt"
	"sync"

	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/schema"
)

// ToolRegistry 统一工具注册中心
// 管理所有工具的注册、查询和调用
type ToolRegistry struct {
	mu     sync.RWMutex
	tools  map[string]BaseTool
	config map[string]ToolConfig
}

// NewToolRegistry 创建新的工具注册中心
func NewToolRegistry() *ToolRegistry {
	return &ToolRegistry{
		tools:  make(map[string]BaseTool),
		config: make(map[string]ToolConfig),
	}
}

// Register 注册工具
// name: 工具名称（必须唯一）
// tool: 工具实例
// config: 工具配置
func (r *ToolRegistry) Register(name string, tool BaseTool, config ToolConfig) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.tools[name]; exists {
		return fmt.Errorf("tool %s already registered", name)
	}

	r.tools[name] = tool
	r.config[name] = config
	return nil
}

// Get 获取工具
func (r *ToolRegistry) Get(name string) (BaseTool, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	tool, exists := r.tools[name]
	return tool, exists
}

// GetConfig 获取工具配置
func (r *ToolRegistry) GetConfig(name string) (ToolConfig, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	config, exists := r.config[name]
	return config, exists
}

// List 列出所有工具名称
func (r *ToolRegistry) List() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()

	names := make([]string, 0, len(r.tools))
	for name := range r.tools {
		names = append(names, name)
	}
	return names
}

// ListByCategory 按分类列出工具
func (r *ToolRegistry) ListByCategory(category string) []string {
	r.mu.RLock()
	defer r.mu.RUnlock()

	names := make([]string, 0)
	for name := range r.tools {
		// 使用Info方法获取工具信息，然后检查描述中的分类信息
		// 这里简化处理，实际应该在ToolInfo中添加Category字段
		// 暂时返回所有工具
		names = append(names, name)
	}
	return names
}

// ToEinoTools 转换为Eino工具集合
// 用于注册到大模型系统中
func (r *ToolRegistry) ToEinoTools() map[string]tool.InvokableTool {
	r.mu.RLock()
	defer r.mu.RUnlock()

	einoTools := make(map[string]tool.InvokableTool)
	for name, baseTool := range r.tools {
		einoTools[name] = &einoToolWrapper{
			BaseTool: baseTool,
			name:     name,
			registry: r,
		}
	}
	return einoTools
}

// GetToolInfo 获取工具信息
func (r *ToolRegistry) GetToolInfo(ctx context.Context, name string) (*schema.ToolInfo, error) {
	tool, exists := r.Get(name)
	if !exists {
		return nil, fmt.Errorf("tool %s not found", name)
	}
	return tool.Info(ctx)
}

// GetAllToolInfo 获取所有工具信息
func (r *ToolRegistry) GetAllToolInfo(ctx context.Context) ([]*schema.ToolInfo, error) {
	infos := make([]*schema.ToolInfo, 0)
	for name := range r.tools {
		info, err := r.GetToolInfo(ctx, name)
		if err != nil {
			return nil, fmt.Errorf("failed to get tool info for %s: %w", name, err)
		}
		infos = append(infos, info)
	}
	return infos, nil
}

// GetGlobalRegistry 获取全局注册中心实例
func GetGlobalRegistry() *ToolRegistry {
	return globalRegistry
}

// einoToolWrapper 将BaseTool包装为Eino的InvokableTool
type einoToolWrapper struct {
	BaseTool
	name     string
	registry *ToolRegistry
}

// Info 实现Eino的Info接口
func (w *einoToolWrapper) Info(ctx context.Context) (*schema.ToolInfo, error) {
	return w.BaseTool.Info(ctx)
}

// InvokableRun 实现Eino的InvokableRun接口
func (w *einoToolWrapper) InvokableRun(ctx context.Context, args string, opts ...tool.Option) (string, error) {
	result, err := w.BaseTool.InvokableRun(ctx, args, opts...)
	// 执行后处理（学习集成）
	w.BaseTool.AfterRun(ctx, result, err)
	return result, err
}

// globalRegistry 全局工具注册中心
var globalRegistry = NewToolRegistry()

// Register 注册工具到全局注册中心
func Register(name string, tool BaseTool, config ToolConfig) error {
	return globalRegistry.Register(name, tool, config)
}

// Get 从全局注册中心获取工具
func Get(name string) (BaseTool, bool) {
	return globalRegistry.Get(name)
}

// GetConfig 从全局注册中心获取工具配置
func GetConfig(name string) (ToolConfig, bool) {
	return globalRegistry.GetConfig(name)
}

// List 列出全局注册中心的所有工具
func List() []string {
	return globalRegistry.List()
}

// ListByCategory 按分类列出工具
func ListByCategory(category string) []string {
	return globalRegistry.ListByCategory(category)
}

// ToEinoTools 转换为Eino工具集合
func ToEinoTools() map[string]tool.InvokableTool {
	return globalRegistry.ToEinoTools()
}

// GetToolInfo 获取工具信息
func GetToolInfo(ctx context.Context, name string) (*schema.ToolInfo, error) {
	return globalRegistry.GetToolInfo(ctx, name)
}

// GetAllToolInfo 获取所有工具信息
func GetAllToolInfo(ctx context.Context) ([]*schema.ToolInfo, error) {
	return globalRegistry.GetAllToolInfo(ctx)
}
