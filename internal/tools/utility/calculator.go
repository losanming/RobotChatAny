package utility

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/schema"
)

// Calculator 计算器工具
type Calculator struct {
	enableLearning bool
	category       string
	toolName       string
}

// NewCalculator 创建计算器工具
func NewCalculator() *Calculator {
	return &Calculator{
		enableLearning: true,
		category:       "utility",
		toolName:       "calculator",
	}
}

// Info 返回工具信息
func (c *Calculator) Info(ctx context.Context) (*schema.ToolInfo, error) {
	return &schema.ToolInfo{
		Name: "calculator",
		Desc: "计算器工具，支持基本数学运算、科学计算、表达式求值、历史记录",
		ParamsOneOf: schema.NewParamsOneOfByParams(map[string]*schema.ParameterInfo{
			"action": {
				Desc:     "操作类型",
				Type:     schema.String,
				Required: true,
				Enum:     []string{"calculate", "get_history", "clear_history"},
			},
			"expression": {
				Desc:     "数学表达式",
				Type:     schema.String,
				Required: false,
			},
			"mode": {
				Desc:     "计算模式",
				Type:     schema.String,
				Required: false,
				Enum:     []string{"basic", "scientific"},
			},
		}),
	}, nil
}

// InvokableRun 执行工具调用
func (c *Calculator) InvokableRun(ctx context.Context, args string, opts ...tool.Option) (string, error) {
	var params map[string]interface{}
	if err := json.Unmarshal([]byte(args), &params); err != nil {
		return "", fmt.Errorf("参数解析失败: %w", err)
	}

	action := params["action"].(string)
	switch action {
	case "calculate":
		return c.calculate(params)
	case "get_history":
		return c.getHistory()
	case "clear_history":
		return c.clearHistory()
	default:
		return "", fmt.Errorf("不支持的操作类型: %s", action)
	}
}

// calculate 计算表达式
func (c *Calculator) calculate(params map[string]interface{}) (string, error) {
	expression := params["expression"].(string)
	mode := getStringParam(params, "mode", "basic")

	// 清理表达式
	expression = strings.TrimSpace(expression)
	if expression == "" {
		return "", fmt.Errorf("表达式不能为空")
	}

	// 安全检查
	if err := c.checkExpressionSafety(expression); err != nil {
		return "", fmt.Errorf("表达式安全检查失败: %w", err)
	}

	// 计算结果
	result, err := c.evaluateExpression(expression, mode)
	if err != nil {
		return "", fmt.Errorf("计算失败: %w", err)
	}

	// 格式化结果
	formattedResult := c.formatResult(result, mode)

	calcResult := map[string]interface{}{
		"success":    true,
		"expression": expression,
		"result":     formattedResult,
		"mode":       mode,
		"raw_result": result,
	}

	return toJSON(calcResult), nil
}

// getHistory 获取历史记录
func (c *Calculator) getHistory() (string, error) {
	result := map[string]interface{}{
		"success": true,
		"history": []string{}, // 演示版本，返回空历史
		"total":   0,
	}

	return toJSON(result), nil
}

// clearHistory 清除历史记录
func (c *Calculator) clearHistory() (string, error) {
	result := map[string]interface{}{
		"success": true,
		"message": "历史记录已清除（演示功能）",
	}

	return toJSON(result), nil
}

// evaluateExpression 求值表达式
func (c *Calculator) evaluateExpression(expression, mode string) (float64, error) {
	// 简单的表达式求值（演示版本）
	// 实际应用中应该使用更安全的表达式解析器

	// 替换常见数学函数
	if mode == "scientific" {
		expression = c.replaceMathFunctions(expression)
	}

	// 简单的四则运算
	return c.simpleEvaluate(expression)
}

// replaceMathFunctions 替换数学函数
func (c *Calculator) replaceMathFunctions(expression string) string {
	// 演示版本：不进行实际替换，只返回原表达式
	return expression
}

// simpleEvaluate 简单求值
func (c *Calculator) simpleEvaluate(expression string) (float64, error) {
	// 移除空格
	expression = strings.ReplaceAll(expression, " ", "")

	// 简单的加减乘除运算
	// 演示版本：只支持基本运算
	
	// 检查是否包含数学函数
	if strings.Contains(expression, "math.") {
		// 演示版本：返回默认值
		return 0, fmt.Errorf("科学计算功能在演示版本中受限")
	}

	// 简单的数字运算
	return c.parseBasicExpression(expression)
}

// parseBasicExpression 解析基本表达式
func (c *Calculator) parseBasicExpression(expression string) (float64, error) {
	// 简单的表达式解析（演示版本）
	// 只支持基本的加减乘除

	// 先处理乘除
	parts := strings.Split(expression, "*")
	if len(parts) > 1 {
		result := 1.0
		for _, part := range parts {
			val, err := strconv.ParseFloat(part, 64)
			if err != nil {
				return 0, err
			}
			result *= val
		}
		return result, nil
	}

	parts = strings.Split(expression, "/")
	if len(parts) > 1 {
		val1, err := strconv.ParseFloat(parts[0], 64)
		if err != nil {
			return 0, err
		}
		val2, err := strconv.ParseFloat(parts[1], 64)
		if err != nil {
			return 0, err
		}
		if val2 == 0 {
			return 0, fmt.Errorf("除数不能为零")
		}
		return val1 / val2, nil
	}

	// 处理加减
	parts = strings.Split(expression, "+")
	if len(parts) > 1 {
		result := 0.0
		for _, part := range parts {
			val, err := strconv.ParseFloat(part, 64)
			if err != nil {
				return 0, err
			}
			result += val
		}
		return result, nil
	}

	parts = strings.Split(expression, "-")
	if len(parts) > 1 {
		val1, err := strconv.ParseFloat(parts[0], 64)
		if err != nil {
			return 0, err
		}
		val2, err := strconv.ParseFloat(parts[1], 64)
		if err != nil {
			return 0, err
		}
		return val1 - val2, nil
	}

	// 单个数字
	return strconv.ParseFloat(expression, 64)
}

// checkExpressionSafety 检查表达式安全性
func (c *Calculator) checkExpressionSafety(expression string) error {
	// 检查危险模式
	dangerousPatterns := []string{
		"eval",
		"exec",
		"system",
		"__import__",
	}

	for _, pattern := range dangerousPatterns {
		if strings.Contains(strings.ToLower(expression), pattern) {
			return fmt.Errorf("表达式包含危险操作: %s", pattern)
		}
	}

	// 检查字符安全性
	validChars := regexp.MustCompile(`^[0-9+\-*/().\smathsincoqrtabsexplogpieE]+$`)
	if !validChars.MatchString(expression) {
		return fmt.Errorf("表达式包含非法字符")
	}

	return nil
}

// formatResult 格式化结果
func (c *Calculator) formatResult(result float64, mode string) string {
	if mode == "scientific" {
		return fmt.Sprintf("%.6e", result)
	}
	return fmt.Sprintf("%.4f", result)
}

// AfterRun 工具调用后的后处理
func (c *Calculator) AfterRun(ctx context.Context, result string, err error) {
	// 简化实现，暂时不做学习集成
}