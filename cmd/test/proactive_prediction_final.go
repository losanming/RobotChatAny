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
	fmt.Println("=== 主动式AI个人全能助手 - 智能预测引擎测试 ===")
	
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
	
	// 测试预测功能
	fmt.Println("\n=== 测试智能预测功能 ===")
	testPredictionEngine(brainInstance)
	
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
	fmt.Printf("   体重: 76.6kg, BMI: 28.2\n")
	fmt.Printf("   血糖: 6.8 mmol/L, 心率: 72\n")
	fmt.Printf("   数据更新时间: 10小时前 (用于测试血糖监测提醒)\n")
}

func testPredictionEngine(brainInstance *brain.Brain) {
	fmt.Println("开始预测威猛先生的需求...")
	
	ctx := context.Background()
	
	// 预测下一步行动
	predictions, err := brainInstance.PredictNextActions(ctx)
	if err != nil {
		log.Printf("预测失败: %v", err)
		return
	}
	
	if len(predictions) == 0 {
		fmt.Println("当前没有预测到需要主动干预的需求")
		return
	}
	
	fmt.Printf("\n🔮 预测到 %d 个主动建议:\n", len(predictions))
	
	for i, prediction := range predictions {
		fmt.Printf("\n--- 建议 %d ---\n", i+1)
		fmt.Printf("类型: %s\n", prediction.Type)
		fmt.Printf("置信度: %.1f%%\n", prediction.Confidence*100)
		fmt.Printf("优先级: %d\n", prediction.Priority)
		fmt.Printf("内容: %s\n", prediction.Content)
		fmt.Printf("理由: %s\n", prediction.Reason)
		
		// 执行高优先级的建议
		if prediction.Priority >= 3 {
			fmt.Printf("🎯 执行建议...\n")
			fmt.Printf("   %s\n", prediction.Content)
		}
	}
}
