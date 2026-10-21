package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"new-robot-self/internal/brain"
	"new-robot-self/internal/tools"
)

func main() {
	fmt.Println("========== 测试 code_executor 工具 ==========")

	// 初始化工具系统
	ctx := context.Background()
	knowledgeBase := brain.NewKnowledgeBase()

	err := tools.RegisterAllTools(knowledgeBase)
	if err != nil {
		log.Fatal("工具注册失败:", err)
	}

	fmt.Println("✅ 工具系统初始化成功")

	// 获取代码执行工具
	codeTool, exists := tools.Get("code_executor")
	if !exists {
		log.Fatal("代码执行工具不存在")
	}

	fmt.Println("✅ 代码执行工具获取成功\n")

	// 测试1: 获取支持的语言列表
	fmt.Println("💻 测试1: 获取支持的语言列表")
	langArgs := map[string]interface{}{
		"action": "get_supported_languages",
	}

	argsJSON, _ := json.Marshal(langArgs)
	result, err := codeTool.InvokableRun(ctx, string(argsJSON))
	if err != nil {
		log.Fatal("获取语言列表失败:", err)
	}

	fmt.Printf("✅ 语言列表: %s\n\n", result)

	// 测试2: 验证Python代码
	fmt.Println("💻 测试2: 验证Python代码")
	validatePythonArgs := map[string]interface{}{
		"action":   "validate",
		"language": "python",
		"code":     "print('Hello World')\nfor i in range(5):\n    print(i)",
	}

	argsJSON, _ = json.Marshal(validatePythonArgs)
	result, err = codeTool.InvokableRun(ctx, string(argsJSON))
	if err != nil {
		log.Fatal("验证Python代码失败:", err)
	}

	fmt.Printf("✅ Python代码验证: %s\n\n", result)

	// 测试3: 验证JavaScript代码
	fmt.Println("💻 测试3: 验证JavaScript代码")
	validateJSArgs := map[string]interface{}{
		"action":   "validate",
		"language": "javascript",
		"code":     "console.log('Hello World');\nfor(let i=0;i<5;i++){\n    console.log(i);\n}",
	}

	argsJSON, _ = json.Marshal(validateJSArgs)
	result, err = codeTool.InvokableRun(ctx, string(argsJSON))
	if err != nil {
		log.Fatal("验证JavaScript代码失败:", err)
	}

	fmt.Printf("✅ JavaScript代码验证: %s\n\n", result)

	// 测试4: 执行安全的Python代码
	fmt.Println("💻 测试4: 执行安全的Python代码")
	execPythonArgs := map[string]interface{}{
		"action":   "execute",
		"language": "python",
		"code":     "print('Hello from Python!')\nprint('计算结果:', 2+3)\nprint('列表:', [1,2,3,4,5])",
		"timeout":  5,
	}

	argsJSON, _ = json.Marshal(execPythonArgs)
	result, err = codeTool.InvokableRun(ctx, string(argsJSON))
	if err != nil {
		log.Fatal("执行Python代码失败:", err)
	}

	fmt.Printf("✅ Python代码执行结果: %s\n\n", result)

	// 测试5: 执行安全的JavaScript代码
	fmt.Println("💻 测试5: 执行安全的JavaScript代码")
	execJSArgs := map[string]interface{}{
		"action":   "execute",
		"language": "javascript",
		"code":     "console.log('Hello from JavaScript!');\nconsole.log('计算结果:', 2+3);\nconsole.log('数组:', [1,2,3,4,5]);",
		"timeout":  5,
	}

	argsJSON, _ = json.Marshal(execJSArgs)
	result, err = codeTool.InvokableRun(ctx, string(argsJSON))
	if err != nil {
		log.Fatal("执行JavaScript代码失败:", err)
	}

	fmt.Printf("✅ JavaScript代码执行结果: %s\n\n", result)

	// 测试6: 执行安全的Bash命令
	fmt.Println("💻 测试6: 执行安全的Bash命令")
	execBashArgs := map[string]interface{}{
		"action":   "execute",
		"language": "bash",
		"code":     "echo 'Hello from Bash!'\necho '当前目录: '\npwd\necho '文件列表:'\nls -la | head -5",
		"timeout":  5,
	}

	argsJSON, _ = json.Marshal(execBashArgs)
	result, err = codeTool.InvokableRun(ctx, string(argsJSON))
	if err != nil {
		log.Fatal("执行Bash命令失败:", err)
	}

	fmt.Printf("✅ Bash命令执行结果: %s\n\n", result)

	// 测试7: 检测危险代码
	fmt.Println("💻 测试7: 检测危险代码")
	dangerousArgs := map[string]interface{}{
		"action":   "validate",
		"language": "bash",
		"code":     "echo 'test' && rm -rf /tmp/test",
	}

	argsJSON, _ = json.Marshal(dangerousArgs)
	result, err = codeTool.InvokableRun(ctx, string(argsJSON))
	if err != nil {
		log.Fatal("检测危险代码失败:", err)
	}

	fmt.Printf("✅ 危险代码检测结果: %s\n\n", result)

	// 测试8: 执行数学计算的Python代码
	fmt.Println("💻 测试8: 执行数学计算的Python代码")
	mathArgs := map[string]interface{}{
		"action":   "execute",
		"language": "python",
		"code":     "import math\nprint('圆周率:', math.pi)\nprint('正弦值:', math.sin(math.pi/2))\nprint('平方根:', math.sqrt(16))",
		"timeout":  5,
	}

	argsJSON, _ = json.Marshal(mathArgs)
	result, err = codeTool.InvokableRun(ctx, string(argsJSON))
	if err != nil {
		log.Fatal("执行数学计算失败:", err)
	}

	fmt.Printf("✅ 数学计算结果: %s\n\n", result)

	fmt.Println("========== 测试完成 ==========")
	fmt.Println("⚠️ 注意: 代码执行工具具有安全限制，某些操作可能被阻止")
}