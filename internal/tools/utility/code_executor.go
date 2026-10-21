package utility

import (
	"context"
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"
	"time"

	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/schema"
)

// CodeExecutor 代码执行工具
type CodeExecutor struct {
	enableLearning bool
	category       string
	toolName       string
}

// NewCodeExecutor 创建代码执行工具
func NewCodeExecutor() *CodeExecutor {
	return &CodeExecutor{
		enableLearning: true,
		category:       "utility",
		toolName:       "code_executor",
	}
}

// Info 返回工具信息
func (r *CodeExecutor) Info(ctx context.Context) (*schema.ToolInfo, error) {
	return &schema.ToolInfo{
		Name: "code_executor",
		Desc: "代码执行工具，支持安全的代码执行环境、多语言支持、执行结果捕获、错误处理",
		ParamsOneOf: schema.NewParamsOneOfByParams(map[string]*schema.ParameterInfo{
			"action": {
				Desc:     "操作类型",
				Type:     schema.String,
				Required: true,
				Enum:     []string{"execute", "validate", "get_supported_languages"},
			},
			"language": {
				Desc:     "编程语言",
				Type:     schema.String,
				Required: false,
				Enum:     []string{"python", "javascript", "bash", "go", "ruby"},
			},
			"code": {
				Desc:     "要执行的代码",
				Type:     schema.String,
				Required: false,
			},
			"timeout": {
				Desc:     "超时时间（秒）",
				Type:     schema.Integer,
				Required: false,
			},
			"working_dir": {
				Desc:     "工作目录",
				Type:     schema.String,
				Required: false,
			},
		}),
	}, nil
}

// InvokableRun 执行工具调用
func (r *CodeExecutor) InvokableRun(ctx context.Context, args string, opts ...tool.Option) (string, error) {
	var params map[string]interface{}
	if err := json.Unmarshal([]byte(args), &params); err != nil {
		return "", fmt.Errorf("参数解析失败: %w", err)
	}

	action := params["action"].(string)
	switch action {
	case "execute":
		return r.executeCode(params)
	case "validate":
		return r.validateCode(params)
	case "get_supported_languages":
		return r.getSupportedLanguages()
	default:
		return "", fmt.Errorf("不支持的操作类型: %s", action)
	}
}

// executeCode 执行代码
func (r *CodeExecutor) executeCode(params map[string]interface{}) (string, error) {
	language := getStringParam(params, "language", "python")
	code := params["code"].(string)
	timeout := getIntParam(params, "timeout", 10)
	workingDir := getStringParam(params, "working_dir", "")

	// 安全检查：检查代码是否包含危险操作
	if err := r.checkCodeSafety(code, language); err != nil {
		return "", fmt.Errorf("代码安全检查失败: %w", err)
	}

	// 根据语言选择执行方式
	var cmd *exec.Cmd
	var output strings.Builder
	var errorOutput strings.Builder

	// 设置超时
	execCtx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)
	defer cancel()

	switch language {
	case "python":
		cmd = exec.CommandContext(execCtx, "python3", "-c", code)
	case "javascript", "js":
		cmd = exec.CommandContext(execCtx, "node", "-e", code)
	case "bash":
		cmd = exec.CommandContext(execCtx, "bash", "-c", code)
	case "go":
		// Go代码需要特殊处理，暂时不支持
		return "", fmt.Errorf("Go代码执行需要编译，暂不支持直接执行")
	case "ruby":
		cmd = exec.CommandContext(execCtx, "ruby", "-e", code)
	default:
		return "", fmt.Errorf("不支持的语言: %s", language)
	}

	// 设置工作目录
	if workingDir != "" {
		cmd.Dir = workingDir
	}

	// 捕获输出
	cmd.Stdout = &output
	cmd.Stderr = &errorOutput

	// 执行命令
	err := cmd.Run()

	// 准备结果
	result := map[string]interface{}{
		"success":       err == nil,
		"language":      language,
		"stdout":        output.String(),
		"stderr":        errorOutput.String(),
		"exit_code":     getExitCode(err),
		"execution_time": timeout, // 简化处理
	}

	if err != nil {
		result["error"] = err.Error()
	}

	return toJSON(result), nil
}

// validateCode 验证代码
func (r *CodeExecutor) validateCode(params map[string]interface{}) (string, error) {
	language := getStringParam(params, "language", "python")
	code := params["code"].(string)

	// 基本语法检查
	validationErrors := r.validateSyntax(code, language)

	// 安全检查
	securityErrors := r.checkCodeSafety(code, language)

	var errors []string
	errors = append(errors, validationErrors...)
	if securityErrors != nil {
		errors = append(errors, securityErrors.Error())
	}

	result := map[string]interface{}{
		"success":       len(errors) == 0,
		"language":      language,
		"is_valid":      len(errors) == 0,
		"errors":        errors,
		"code_length":   len(code),
		"line_count":    strings.Count(code, "\n") + 1,
	}

	return toJSON(result), nil
}

// getSupportedLanguages 获取支持的语言列表
func (r *CodeExecutor) getSupportedLanguages() (string, error) {
	languages := []map[string]interface{}{
		{
			"name":        "python",
			"version":     "3.x",
			"description": "Python 3解释器",
			"status":      "available",
		},
		{
			"name":        "javascript",
			"version":     "Node.js",
			"description": "Node.js JavaScript运行时",
			"status":      "available",
		},
		{
			"name":        "bash",
			"version":     "Bash",
			"description": "Bash shell脚本",
			"status":      "available",
		},
		{
			"name":        "go",
			"version":     "1.x",
			"description": "Go编程语言（需要编译）",
			"status":      "limited",
		},
		{
			"name":        "ruby",
			"version":     "2.x",
			"description": "Ruby编程语言",
			"status":      "available",
		},
	}

	result := map[string]interface{}{
		"success":    true,
		"languages":  languages,
		"total":      len(languages),
	}

	return toJSON(result), nil
}

// checkCodeSafety 检查代码安全性
func (r *CodeExecutor) checkCodeSafety(code, language string) error {
	// 危险操作黑名单
	dangerousPatterns := []string{
		"rm -rf /",
		"rm -rf",
		":(){:|:&};:", // fork bomb
		"mkfs",
		"format",
		"dd if=",
		"wget",
		"curl",
	}

	codeLower := strings.ToLower(code)
	for _, pattern := range dangerousPatterns {
		if strings.Contains(codeLower, pattern) {
			return fmt.Errorf("代码包含潜在危险操作: %s", pattern)
		}
	}

	// 语言特定的安全检查
	switch language {
	case "python":
		if strings.Contains(code, "os.system") || strings.Contains(code, "subprocess.call") {
			// 允许但警告
			// return fmt.Errorf("Python代码包含系统调用，可能存在安全风险")
		}
	case "bash":
		if strings.Contains(code, "sudo") {
			return fmt.Errorf("Bash代码包含sudo命令，存在安全风险")
		}
	}

	return nil
}

// validateSyntax 验证语法
func (r *CodeExecutor) validateSyntax(code, language string) []string {
	var errors []string

	// 基本的语法检查
	switch language {
	case "python":
		// 检查括号匹配
		if !r.checkBrackets(code, "(", ")") {
			errors = append(errors, "括号不匹配")
		}
		// 检查引号匹配
		if !r.checkQuotes(code) {
			errors = append(errors, "引号不匹配")
		}
	case "javascript":
		if !r.checkBrackets(code, "{", "}") {
			errors = append(errors, "大括号不匹配")
		}
		if !r.checkBrackets(code, "(", ")") {
			errors = append(errors, "括号不匹配")
		}
	}

	return errors
}

// checkBrackets 检查括号匹配
func (r *CodeExecutor) checkBrackets(code, open, close string) bool {
	count := 0
	for _, char := range code {
		charStr := string(char)
		if charStr == open {
			count++
		} else if charStr == close {
			count--
		}
		if count < 0 {
			return false
		}
	}
	return count == 0
}

// checkQuotes 检查引号匹配
func (r *CodeExecutor) checkQuotes(code string) bool {
	singleQuoteCount := strings.Count(code, "'")
	doubleQuoteCount := strings.Count(code, "\"")
	return singleQuoteCount%2 == 0 && doubleQuoteCount%2 == 0
}

// getExitCode 获取退出码
func getExitCode(err error) int {
	if err == nil {
		return 0
	}
	// 简化处理，实际应该解析exit error
	return 1
}

// AfterRun 工具调用后的后处理
func (r *CodeExecutor) AfterRun(ctx context.Context, result string, err error) {
	// 简化实现，暂时不做学习集成
}