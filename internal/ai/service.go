// Package ai 提供AI服务相关功能
package ai

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"sync"

	"github.com/cloudwego/eino-ext/components/model/ark"
	"github.com/cloudwego/eino/components/prompt"
	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/schema"
)

// Config AI服务配置结构体
type Config struct {
	// APIKey 阿里云API密钥
	APIKey string `yaml:"api_key"`

	// Endpoint API服务端点URL
	Endpoint string `yaml:"endpoint"`

	// Model 使用的模型名称
	Model string `yaml:"model"`

	// Temperature 控制输出的随机性(0.0-1.0)
	// 值越高，回答越有创意；值越低，回答越确定
	Temperature float64 `yaml:"temperature"`

	// TopP 控制输出的多样性(0.0-1.0)
	// 值越高，选择的词越多样；值越低，选择越保守
	TopP float64 `yaml:"top_p"`

	// MaxTokens 限制回复的最大token数
	// 防止回复过长，通常设置在100-500之间
	MaxTokens int `yaml:"max_tokens"`
}

// GetChatModel 返回已初始化的 ChatModel 实例。
// 注意：调用方需确保在 initialize.InitAiService 或显式 InitModel 之后再调用。
func (s *Service) GetChatModel() *ark.ChatModel {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.model
}

// InvokeWithTools 使用当前模型和提供的工具集合执行一轮对话。
// 内部构建 ChatModel->ToolsNode 的 Chain，直接返回 LLM 的消息输出。
func (s *Service) InvokeWithTools(ctx context.Context, reg map[string]tool.InvokableTool, msgs []*schema.Message) ([]*schema.Message, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if s.model == nil {
		return nil, errors.New("ai.Service model is not initialized")
	}

	// 收集 ToolInfo 并绑定到 ChatModel
	infos := make([]*schema.ToolInfo, 0, len(reg))
	tools := make([]tool.BaseTool, 0, len(reg))
	for name, t := range reg {
		info, err := t.Info(ctx)
		if err != nil {
			slog.Error("get ToolInfo failed", "tool", name, "err", err)
			return nil, err
		}
		infos = append(infos, info)
		tools = append(tools, t)
	}
	if err := s.model.BindTools(infos); err != nil {
		return nil, err
	}

	// 构建 ToolsNode 并连接成链
	toolsNode, err := compose.NewToolNode(ctx, &compose.ToolsNodeConfig{Tools: tools})
	if err != nil {
		return nil, err
	}
	chain := compose.NewChain[[]*schema.Message, []*schema.Message]()
	chain.
		AppendChatModel(s.model, compose.WithNodeName("chat_model")).
		AppendToolsNode(toolsNode, compose.WithNodeName("tools"))
	agent, err := chain.Compile(ctx)
	if err != nil {
		return nil, err
	}
	return agent.Invoke(ctx, msgs)
}

// Service AI服务结构体
// 管理与阿里云AI模型的通信
type Service struct {
	// config 服务配置
	config Config
	model  *ark.ChatModel
	// client HTTP客户端
	client *http.Client
	ctx    context.Context
	// mu 读写锁，保证并发安全
	mu sync.RWMutex
}

// Message 单条消息结构
type Message struct {
	// Role 消息发送者的角色(system/user/assistant)
	Role string `json:"role"`
	// Content 消息内容
	Content string `json:"content"`
	// 时间戳
	Timestamp int64 `json:"timestamp"`
}

type VlMessage struct {
	Message []VlBaseMessage
}

type VlBaseMessage struct {
	Role    string        `json:"role"`
	Content []interface{} `json:"content"`
}

type VlTextData struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

type VlImageData struct {
	Type     string `json:"type"`
	ImageUrl struct {
		Url string `json:"url"`
	} `json:"image_url"`
}

func NewService(cfg Config) *Service {
	return &Service{
		ctx:    context.Background(),
		config: cfg,
		client: &http.Client{},
	}
}

func (s *Service) InitModel() error {
	if s.config.MaxTokens == 0 || s.config.MaxTokens > 500 {
		s.config.MaxTokens = 500
	}
	if s.config.Model == "" {
		s.config.Model = "qwen-turbo"
	}
	m, err := ark.NewChatModel(s.ctx, &ark.ChatModelConfig{
		APIKey:    s.config.APIKey,
		Model:     s.config.Model,
		MaxTokens: &s.config.MaxTokens,
		BaseURL:   s.config.Endpoint,
	})
	if err != nil {
		return err
	}
	s.model = m
	return nil
}

// Chat 发送聊天请求到AI模型并返回完整响应
func (s *Service) Chat(ctx context.Context, messages Message) (string, error) {
	// 加读锁保证并发安全
	s.mu.RLock()
	defer s.mu.RUnlock()

	temp := prompt.FromMessages(schema.FString,
		&schema.Message{
			Role:    schema.System,
			Content: "你是一个{role}",
		},
		&schema.Message{
			Role:    schema.User,
			Content: "{task}。",
		},
	)

	variables := map[string]any{
		"role": messages.Role,
		"task": messages.Content,
	}

	// 格式化模板
	result, err := temp.Format(ctx, variables)
	if err != nil {
		slog.Error("格式化模板失败", "err", err)
		return "", err
	}
	// 获取回复
	reader, err1 := s.model.Generate(ctx, result)
	if err1 != nil {
		slog.Error("获取ai数据流失败", "err", err1)
		return "", err1
	}
	return reader.Content, nil
}
