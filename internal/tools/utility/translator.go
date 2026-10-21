package utility

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/schema"
)

// Translator 翻译工具
type Translator struct {
	enableLearning bool
	category       string
	toolName       string
}

// NewTranslator 创建翻译工具
func NewTranslator() *Translator {
	return &Translator{
		enableLearning: true,
		category:       "utility",
		toolName:       "translator",
	}
}

// Info 返回工具信息
func (t *Translator) Info(ctx context.Context) (*schema.ToolInfo, error) {
	return &schema.ToolInfo{
		Name: "translator",
		Desc: "翻译工具，支持多语言翻译、批量翻译、翻译质量评估（演示版本）",
		ParamsOneOf: schema.NewParamsOneOfByParams(map[string]*schema.ParameterInfo{
			"action": {
				Desc:     "操作类型",
				Type:     schema.String,
				Required: true,
				Enum:     []string{"translate", "detect_language", "get_supported_languages"},
			},
			"text": {
				Desc:     "要翻译的文本",
				Type:     schema.String,
				Required: false,
			},
			"source_lang": {
				Desc:     "源语言",
				Type:     schema.String,
				Required: false,
			},
			"target_lang": {
				Desc:     "目标语言",
				Type:     schema.String,
				Required: false,
			},
			"engine": {
				Desc:     "翻译引擎",
				Type:     schema.String,
				Required: false,
				Enum:     []string{"demo", "mock"},
			},
		}),
	}, nil
}

// InvokableRun 执行工具调用
func (t *Translator) InvokableRun(ctx context.Context, args string, opts ...tool.Option) (string, error) {
	var params map[string]interface{}
	if err := json.Unmarshal([]byte(args), &params); err != nil {
		return "", fmt.Errorf("参数解析失败: %w", err)
	}

	action := params["action"].(string)
	switch action {
	case "translate":
		return t.translate(params)
	case "detect_language":
		return t.detectLanguage(params)
	case "get_supported_languages":
		return t.getSupportedLanguages()
	default:
		return "", fmt.Errorf("不支持的操作类型: %s", action)
	}
}

// translate 翻译文本
func (t *Translator) translate(params map[string]interface{}) (string, error) {
	text := params["text"].(string)
	sourceLang := getStringParam(params, "source_lang", "auto")
	targetLang := getStringParam(params, "target_lang", "en")
	engine := getStringParam(params, "engine", "demo")

	// 演示版本：模拟翻译
	translatedText, confidence := t.mockTranslate(text, sourceLang, targetLang)

	result := map[string]interface{}{
		"success":       true,
		"original_text": text,
		"translated_text": translatedText,
		"source_lang":   sourceLang,
		"target_lang":   targetLang,
		"engine":        engine,
		"confidence":    confidence,
		"char_count":    len(text),
	}

	return toJSON(result), nil
}

// detectLanguage 检测语言
func (t *Translator) detectLanguage(params map[string]interface{}) (string, error) {
	text := params["text"].(string)

	// 简单的语言检测
	detectedLang := t.simpleLanguageDetection(text)

	result := map[string]interface{}{
		"success":        true,
		"text":           text,
		"detected_lang":  detectedLang,
		"confidence":     0.8,
	}

	return toJSON(result), nil
}

// getSupportedLanguages 获取支持的语言列表
func (t *Translator) getSupportedLanguages() (string, error) {
	languages := []map[string]interface{}{
		{"code": "zh", "name": "中文", "native_name": "中文"},
		{"code": "en", "name": "English", "native_name": "English"},
		{"code": "ja", "name": "Japanese", "native_name": "日本語"},
		{"code": "ko", "name": "Korean", "native_name": "한국어"},
		{"code": "fr", "name": "French", "native_name": "Français"},
		{"code": "de", "name": "German", "native_name": "Deutsch"},
		{"code": "es", "name": "Spanish", "native_name": "Español"},
		{"code": "ru", "name": "Russian", "native_name": "Русский"},
	}

	result := map[string]interface{}{
		"success":   true,
		"languages": languages,
		"total":     len(languages),
	}

	return toJSON(result), nil
}

// mockTranslate 模拟翻译（演示版本）
func (t *Translator) mockTranslate(text, sourceLang, targetLang string) (string, float64) {
	// 演示版本：简单模拟翻译
	// 实际应用中应该调用真实的翻译API

	if sourceLang == targetLang {
		return text, 1.0
	}

	// 简单的演示翻译逻辑
	translatedText := fmt.Sprintf("[翻译] %s (%s -> %s)", text, sourceLang, targetLang)
	confidence := 0.85

	return translatedText, confidence
}

// simpleLanguageDetection 简单语言检测
func (t *Translator) simpleLanguageDetection(text string) string {
	// 检测中文字符
	chineseCount := 0
	for _, r := range text {
		if r >= 0x4e00 && r <= 0x9fff {
			chineseCount++
		}
	}

	if chineseCount > len(text)/2 {
		return "zh"
	}

	// 检测日文字符
	japaneseCount := 0
	for _, r := range text {
		if (r >= 0x3040 && r <= 0x309f) || (r >= 0x30a0 && r <= 0x30ff) {
			japaneseCount++
		}
	}

	if japaneseCount > len(text)/3 {
		return "ja"
	}

	// 默认返回英语
	return "en"
}

// AfterRun 工具调用后的后处理
func (t *Translator) AfterRun(ctx context.Context, result string, err error) {
	// 简化实现，暂时不做学习集成
}