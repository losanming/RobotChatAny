package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"new-robot-self/config"
	"new-robot-self/internal/brain"
)

func main() {
	fmt.Println("=== 主动式AI预测功能准确性测试 ===")
	
	// 加载配置
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("加载配置失败: %v", err)
	}

	if cfg.RAG.Embedding.BaseURL != "" {
		fmt.Println("✅ RAG服务配置已加载")
	}

	// 创建Brain实例
	brainInstance := brain.NewBrain()
	
	// 测试场景1：健康数据较旧，应该触发血糖监测提醒
	fmt.Println("\n=== 测试场景1：健康数据较旧 ===")
	testOldHealthData(brainInstance)
	
	// 测试场景2：运动时间
	fmt.Println("\n=== 测试场景2：模拟运动时间 ===")
	testExerciseTime(brainInstance)
	
	// 测试场景3：饭点时间
	fmt.Println("\n=== 测试场景3：模拟饭点时间 ===")
	testMealTime(brainInstance)
	
	// 测试场景4：工作时间
	fmt.Println("\n=== 测试场景4：模拟工作时间 ===")
	testWorkTime(brainInstance)
	
	// 测试场景5：周末时间
	fmt.Println("\n=== 测试场景5：模拟周末时间 ===")
	testWeekendTime(brainInstance)
	
	fmt.Println("\n=== 所有测试完成 ===")
}

func testOldHealthData(brainInstance *brain.Brain) {
	// 设置较旧的健康数据（10小时前）
	earlierTime := time.Now().Add(-10 * time.Hour)
	
	fmt.Printf("设置健康数据，时间戳: %d (10小时前)\n", earlierTime.Unix())
	fmt.Printf("当前时间: %d\n", time.Now().Unix())
	fmt.Printf("时间差: %.1f 小时\n", time.Since(earlierTime).Hours())
	
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
	
	predictions, err := brainInstance.PredictNextActions(context.Background())
	if err != nil {
		log.Printf("预测失败: %v", err)
		return
	}
	
	fmt.Printf("预测到 %d 个建议\n", len(predictions))
	for _, pred := range predictions {
		fmt.Printf("- [%s] %s (优先级: %d, 置信度: %.1f%%)\n", 
			pred.Type, pred.Content, pred.Priority, pred.Confidence*100)
	}
	
	// 验证：应该有血糖监测提醒
	hasBloodSugarReminder := false
	for _, pred := range predictions {
		if pred.Action == "blood_sugar_reminder" {
			hasBloodSugarReminder = true
			break
		}
	}
	
	if hasBloodSugarReminder {
		fmt.Println("✅ 测试通过：成功触发血糖监测提醒")
	} else {
		fmt.Println("❌ 测试失败：未触发血糖监测提醒")
	}
}

func testExerciseTime(brainInstance *brain.Brain) {
	// 由于无法修改系统时间，这里只是说明性测试
	fmt.Println("说明：运动时间测试需要在晚上19-21点进行")
	fmt.Println("预期结果：应该触发运动提醒，优先级4，置信度80%")
}

func testMealTime(brainInstance *brain.Brain) {
	// 由于无法修改系统时间，这里只是说明性测试
	fmt.Println("说明：饭点时间测试需要在7点、12点或18点进行")
	fmt.Println("预期结果：应该触发饮食建议，优先级3，置信度60%")
}

func testWorkTime(brainInstance *brain.Brain) {
	// 由于无法修改系统时间，这里只是说明性测试
	fmt.Println("说明：工作时间测试需要在工作日9-18点进行")
	fmt.Println("预期结果：应该触发工作专注时间建议和休息提醒")
}

func testWeekendTime(brainInstance *brain.Brain) {
	// 由于无法修改系统时间，这里只是说明性测试
	fmt.Println("说明：周末时间测试需要在周六或周日进行")
	fmt.Println("预期结果：应该触发周末活动建议")
}
