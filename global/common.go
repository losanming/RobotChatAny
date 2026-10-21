package global

import (
	"new-robot-self/internal/ai"
	"new-robot-self/pkg/rag"
)

// DefaultConfig 默认配置
// @@TODO 需要从环境配置获取敏感数据
// 使用阿里云Deepseek模型的推荐配置
var DefaultConfig = ai.Config{
	// 阿里云API密钥
	APIKey: "YOUR_API_KEY",

	// 阿里云Deepseek模型的API端点
	Endpoint: "https://dashscope.aliyuncs.com/compatible-mode/v1",

	// 使用Qwen-Turbo模型
	Model: "qwen-turbo",

	// 温度0.7，提供创造性但不会太发散的回答
	Temperature: 0.7,

	// TopP 0.95，允许一定的多样性
	TopP: 0.95,

	// 限制回复最多800个token
	MaxTokens: 800,
}

// 高德地图APIkey
var AmapKey = "YOUR_AMAP_KEY"

// 视觉模型
var VlConfig = ai.Config{
	// 阿里云API密钥
	APIKey: "YOUR_API_KEY",

	// 阿里云Deepseek模型的API端点
	Endpoint: "https://dashscope.aliyuncs.com/compatible-mode/v1",

	// 使用Qwen2.5-VL Max模型
	Model: "qwen-vl-max",

	// 温度0.7，提供创造性但不会太发散的回答
	Temperature: 0.7,

	// TopP 0.95，允许一定的多样性
	TopP: 0.95,

	// 限制回复最多150个token
	MaxTokens: 150,
}

const TriggerPath = "asset/triggers.json"
const FetcherPath = "asset/resources.json"
const HostName = "Dangwing'ai"

// RAGService RAG服务全局实例
var RAGService *rag.Service
