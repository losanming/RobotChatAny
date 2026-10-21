package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"new-robot-self/internal/brain"
	"new-robot-self/internal/tools"
)

func main() {
	fmt.Println("========== 测试 file_processor 工具 ==========")

	// 初始化工具系统
	ctx := context.Background()
	knowledgeBase := brain.NewKnowledgeBase()

	err := tools.RegisterAllTools(knowledgeBase)
	if err != nil {
		log.Fatal("工具注册失败:", err)
	}

	fmt.Println("✅ 工具系统初始化成功")

	// 获取文件处理工具
	fileTool, exists := tools.Get("file_processor")
	if !exists {
		log.Fatal("文件处理工具不存在")
	}

	fmt.Println("✅ 文件处理工具获取成功\n")

	// 创建测试文件
	testContent := `这是一个测试文件。
用于测试文件处理工具的功能。
包含中文和English混合内容。
Hello World!
测试数据：12345
特殊字符：!@#$%^&*()`

	testFilePath := "/tmp/test_file_processor.txt"
	err = ioutil.WriteFile(testFilePath, []byte(testContent), 0644)
	if err != nil {
		log.Fatal("创建测试文件失败:", err)
	}

	fmt.Printf("✅ 测试文件创建成功: %s\n\n", testFilePath)

	// 测试1: 获取文件信息
	fmt.Println("📄 测试1: 获取文件信息")
	infoArgs := map[string]interface{}{
		"action":    "get_info",
		"file_path": testFilePath,
	}

	argsJSON, _ := json.Marshal(infoArgs)
	result, err := fileTool.InvokableRun(ctx, string(argsJSON))
	if err != nil {
		log.Fatal("获取文件信息失败:", err)
	}

	fmt.Printf("✅ 文件信息: %s\n\n", result)

	// 测试2: 读取文本内容
	fmt.Println("📄 测试2: 读取文本内容")
	readArgs := map[string]interface{}{
		"action":    "read_text",
		"file_path": testFilePath,
		"encoding":  "utf-8",
	}

	argsJSON, _ = json.Marshal(readArgs)
	result, err = fileTool.InvokableRun(ctx, string(argsJSON))
	if err != nil {
		log.Fatal("读取文本失败:", err)
	}

	fmt.Printf("✅ 文件内容: %s\n\n", result)

	// 测试3: 分析内容
	fmt.Println("📄 测试3: 分析内容")
	analyzeArgs := map[string]interface{}{
		"action":    "analyze_content",
		"file_path": testFilePath,
	}

	argsJSON, _ = json.Marshal(analyzeArgs)
	result, err = fileTool.InvokableRun(ctx, string(argsJSON))
	if err != nil {
		log.Fatal("分析内容失败:", err)
	}

	fmt.Printf("✅ 内容分析: %s\n\n", result)

	// 测试4: 格式转换
	fmt.Println("📄 测试4: 格式转换为JSON")
	convertArgs := map[string]interface{}{
		"action":        "convert_format",
		"file_path":      testFilePath,
		"output_format": "json",
	}

	argsJSON, _ = json.Marshal(convertArgs)
	result, err = fileTool.InvokableRun(ctx, string(argsJSON))
	if err != nil {
		log.Fatal("格式转换失败:", err)
	}

	fmt.Printf("✅ 格式转换结果: %s\n\n", result)

	// 测试5: 格式转换为Markdown
	fmt.Println("📄 测试5: 格式转换为Markdown")
	convertMarkdownArgs := map[string]interface{}{
		"action":        "convert_format",
		"file_path":      testFilePath,
		"output_format": "markdown",
	}

	argsJSON, _ = json.Marshal(convertMarkdownArgs)
	result, err = fileTool.InvokableRun(ctx, string(argsJSON))
	if err != nil {
		log.Fatal("格式转换失败:", err)
	}

	fmt.Printf("✅ Markdown转换结果: %s\n\n", result)

	// 清理测试文件
	err = ioutil.WriteFile(testFilePath, []byte(""), 0644)
	if err != nil {
		log.Printf("清理测试文件失败: %v\n", err)
	}

	fmt.Println("========== 测试完成 ==========")
}