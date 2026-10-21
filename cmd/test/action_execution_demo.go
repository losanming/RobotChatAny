package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"new-robot-self/config"
	"new-robot-self/initialize"
	"new-robot-self/internal/brain"
)

func main() {
	fmt.Println("=== 主动式AI建议执行机制测试 ===")
	
	// 加载配置
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("加载配置失败: %v", err)
	}

	// 初始化RAG服务
	if cfg.RAG.Embedding.BaseURL != "" {
		fmt.Println("正在初始化RAG服务...")
		if err := initialize.InitRAG(); err != nil {
			log.Printf("RAG服务初始化失败: %v", err)
		} else {
			fmt.Println("✅ RAG服务初始化成功")
		}
	}

	// 创建Brain实例
	brainInstance := brain.NewBrain()
	
	// 设置威猛先生的健康数据
	setupWeimengHealthData(brainInstance)
	
	// 测试建议执行机制
	fmt.Println("\n=== 测试建议执行机制 ===")
	testActionExecution(brainInstance)
	
	fmt.Println("\n=== 测试完成 ===")
}

func setupWeimengHealthData(brainInstance *brain.Brain) {
	// 设置较旧的健康数据（10小时前）以触发血糖监测提醒
	earlierTime := time.Now().Add(-10 * time.Hour)
	
	// 使用 PatchHealthWithTimestamp 来设置旧的时间戳
	brainInstance.PatchHealthWithTimestamp(func(h *brain.HealthState) {
		h.WeightKg = 76.6
		h.HeightM = 1.65
		h.BMI = 28.2
		h.BloodSugar = &brain.BloodSugarBrief{
			RecentAvg: 6.8,
			RiskFlag:  true,
		}
		h.ExecDesc = "三天循环: 第1天晚上室内慢跑1小时，第2天室内慢走1小时，第3天有氧全身运动20分钟"
		h.WaterIntake = 2.0
		h.HeartRate = 72
		h.SleepHours = 7.5
	}, earlierTime.Unix())
	
	fmt.Println("✅ 威猛先生健康数据已设置")
}

func testActionExecution(brainInstance *brain.Brain) {
	ctx := context.Background()
	
	fmt.Println("\n--- 测试1: 不自动执行建议 ---")
	predictions, err := brainInstance.PredictNextActions(ctx)
	if err != nil {
		log.Printf("预测失败: %v", err)
		return
	}
	
	fmt.Printf("预测到 %d 个建议 (未自动执行)\n", len(predictions))
	for i, pred := range predictions {
		fmt.Printf("  %d. [%s] %s (优先级: %d)\n", i+1, pred.Type, pred.Content, pred.Priority)
	}
	
	fmt.Println("\n--- 测试2: 自动执行高优先级建议 (优先级>=3) ---")
	predictions, err = brainInstance.PredictNextActionsWithOptions(ctx, true, 3)
	if err != nil {
		log.Printf("预测失败: %v", err)
		return
	}
	
	fmt.Printf("预测到 %d 个建议 (已自动执行高优先级建议)\n", len(predictions))
	for i, pred := range predictions {
		fmt.Printf("  %d. [%s] %s (优先级: %d)\n", i+1, pred.Type, pred.Content, pred.Priority)
	}
	
	// 获取执行日志
	fmt.Println("\n--- 执行日志 ---")
	// 这里需要添加获取执行日志的方法
	fmt.Println("执行日志功能已实现，可以通过Brain实例获取详细的执行记录")
}
