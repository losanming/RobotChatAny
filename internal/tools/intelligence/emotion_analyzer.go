package intelligence

import (
	"context"
	"encoding/json"

	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/schema"
)

// EmotionAnalyzer 情感分析工具
type EmotionAnalyzer struct {
	enableLearning bool
	category       string
	toolName       string
}

func NewEmotionAnalyzer() *EmotionAnalyzer {
	return &EmotionAnalyzer{
		enableLearning: true,
		category:       "intelligence",
		toolName:       "emotion_analyzer",
	}
}

func (e *EmotionAnalyzer) Info(ctx context.Context) (*schema.ToolInfo, error) {
	return &schema.ToolInfo{
		Name: "emotion_analyzer",
		Desc: "情感分析工具，支持情感识别、情绪分析、情感趋势",
		ParamsOneOf: schema.NewParamsOneOfByParams(map[string]*schema.ParameterInfo{
			"action": {Desc: "操作类型", Type: schema.String, Required: true, Enum: []string{"analyze", "detect", "trend"}},
			"text": {Desc: "文本内容", Type: schema.String, Required: false},
		}),
	}, nil
}

func (e *EmotionAnalyzer) InvokableRun(ctx context.Context, args string, opts ...tool.Option) (string, error) {
	var params map[string]interface{}
	json.Unmarshal([]byte(args), &params)
	result := map[string]interface{}{
		"success": true,
		"message": "情感分析成功（演示版本）",
	}
	return toJSON(result), nil
}

func (e *EmotionAnalyzer) AfterRun(ctx context.Context, result string, err error) {}