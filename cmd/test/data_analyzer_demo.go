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
	fmt.Println("========== 测试 data_analyzer 工具 ==========")

	// 初始化工具系统
	ctx := context.Background()
	knowledgeBase := brain.NewKnowledgeBase()

	err := tools.RegisterAllTools(knowledgeBase)
	if err != nil {
		log.Fatal("工具注册失败:", err)
	}

	fmt.Println("✅ 工具系统初始化成功")

	// 获取数据分析工具
	dataTool, exists := tools.Get("data_analyzer")
	if !exists {
		log.Fatal("数据分析工具不存在")
	}

	fmt.Println("✅ 数据分析工具获取成功\n")

	// 测试数据
	testData := []float64{10, 15, 12, 18, 20, 25, 22, 30, 28, 35, 100} // 包含一个异常值

	// 测试1: 统计分析
	fmt.Println("📊 测试1: 统计分析（全部）")
	statsArgs := map[string]interface{}{
		"action":        "statistics",
		"data":          testData,
		"analysis_type": "all",
	}

	argsJSON, _ := json.Marshal(statsArgs)
	result, err := dataTool.InvokableRun(ctx, string(argsJSON))
	if err != nil {
		log.Fatal("统计分析失败:", err)
	}

	fmt.Printf("✅ 统计分析结果: %s\n\n", result)

	// 测试2: 均值计算
	fmt.Println("📊 测试2: 均值计算")
	meanArgs := map[string]interface{}{
		"action":        "statistics",
		"data":          testData,
		"analysis_type": "mean",
	}

	argsJSON, _ = json.Marshal(meanArgs)
	result, err = dataTool.InvokableRun(ctx, string(argsJSON))
	if err != nil {
		log.Fatal("均值计算失败:", err)
	}

	fmt.Printf("✅ 均值结果: %s\n\n", result)

	// 测试3: 聚合分析
	fmt.Println("📊 测试3: 聚合分析")
	aggregateArgs := map[string]interface{}{
		"action":             "aggregate",
		"data":               testData,
		"aggregate_function": "sum",
	}

	argsJSON, _ = json.Marshal(aggregateArgs)
	result, err = dataTool.InvokableRun(ctx, string(argsJSON))
	if err != nil {
		log.Fatal("聚合分析失败:", err)
	}

	fmt.Printf("✅ 聚合分析结果: %s\n\n", result)

	// 测试4: 趋势分析
	fmt.Println("📊 测试4: 趋势分析")
	trendArgs := map[string]interface{}{
		"action": "trend",
		"data":   testData,
	}

	argsJSON, _ = json.Marshal(trendArgs)
	result, err = dataTool.InvokableRun(ctx, string(argsJSON))
	if err != nil {
		log.Fatal("趋势分析失败:", err)
	}

	fmt.Printf("✅ 趋势分析结果: %s\n\n", result)

	// 测试5: 异常检测（Z-score方法）
	fmt.Println("📊 测试5: 异常检测（Z-score方法）")
	anomalyArgs := map[string]interface{}{
		"action":    "anomaly",
		"data":      testData,
		"method":    "zscore",
		"threshold": 2.0,
	}

	argsJSON, _ = json.Marshal(anomalyArgs)
	result, err = dataTool.InvokableRun(ctx, string(argsJSON))
	if err != nil {
		log.Fatal("异常检测失败:", err)
	}

	fmt.Printf("✅ 异常检测结果: %s\n\n", result)

	// 测试6: 异常检测（IQR方法）
	fmt.Println("📊 测试6: 异常检测（IQR方法）")
	iqrArgs := map[string]interface{}{
		"action":    "anomaly",
		"data":      testData,
		"method":    "iqr",
		"threshold": 1.5,
	}

	argsJSON, _ = json.Marshal(iqrArgs)
	result, err = dataTool.InvokableRun(ctx, string(argsJSON))
	if err != nil {
		log.Fatal("异常检测失败:", err)
	}

	fmt.Printf("✅ IQR异常检测结果: %s\n\n", result)

	// 测试7: 相关性分析
	fmt.Println("📊 测试7: 相关性分析")
	correlationArgs := map[string]interface{}{
		"action": "correlation",
		"data":   testData,
	}

	argsJSON, _ = json.Marshal(correlationArgs)
	result, err = dataTool.InvokableRun(ctx, string(argsJSON))
	if err != nil {
		log.Fatal("相关性分析失败:", err)
	}

	fmt.Printf("✅ 相关性分析结果: %s\n\n", result)

	// 测试8: 使用正常数据（无异常值）
	fmt.Println("📊 测试8: 正常数据统计分析")
	normalData := []float64{10, 12, 15, 18, 20, 22, 25, 28, 30}
	normalArgs := map[string]interface{}{
		"action":        "statistics",
		"data":          normalData,
		"analysis_type": "all",
	}

	argsJSON, _ = json.Marshal(normalArgs)
	result, err = dataTool.InvokableRun(ctx, string(argsJSON))
	if err != nil {
		log.Fatal("正常数据统计分析失败:", err)
	}

	fmt.Printf("✅ 正常数据统计: %s\n\n", result)

	fmt.Println("========== 测试完成 ==========")
}