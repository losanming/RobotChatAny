package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"new-robot-self/config"
	"new-robot-self/pkg/rag"
)

func main() {
	fmt.Println("=== 添加威猛先生健康知识到知识库 ===")
	
	// 加载配置
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("加载配置失败: %v", err)
	}

	// 初始化RAG服务
	ragService := rag.NewService(cfg.RAG)
	
	// 健康检查
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	
	if err := ragService.HealthCheck(ctx); err != nil {
		log.Fatalf("RAG服务健康检查失败: %v", err)
	}
	fmt.Println("✅ RAG服务健康检查通过")

	// 威猛先生的健康知识文档
	healthDocs := []rag.Document{
		{
			ID:      "weimeng_basic_info",
			Content: "威猛先生是一名程序员，男性，身高165cm，体重76.6kg，BMI指数约为28.2，属于超重范围。他有轻微的糖尿病，主要原因是肥胖导致的。他的目标是减肥和改善糖尿病状况。",
			Metadata: map[string]interface{}{
				"category": "basic_info",
				"author":   "system",
				"priority": "high",
			},
		},
		{
			ID:      "weimeng_exercise_plan",
			Content: "威猛先生采用每三天一次的运动循环：第一天晚上室内慢跑1小时，第二天室内慢走1小时，第三天有氧全身运动20分钟。这个运动计划旨在帮助他减肥和改善糖尿病状况，建议长期坚持。",
			Metadata: map[string]interface{}{
				"category": "exercise",
				"author":   "system",
				"priority": "high",
			},
		},
		{
			ID:      "weimeng_diet_habits",
			Content: "威猛先生的饮食习惯：每天早上一瓶200ml的纯牛奶；中午100g米饭一荤一素；晚上偶尔一根黄瓜一杯黑咖啡。这个饮食计划相对清淡，有助于控制血糖和体重，建议继续保持并适当调整。",
			Metadata: map[string]interface{}{
				"category": "diet",
				"author":   "system",
				"priority": "high",
			},
		},
		{
			ID:      "weimeng_diabetes_management",
			Content: "威猛先生有轻微的糖尿病，主要是由肥胖引起的。管理建议包括：1. 控制体重，目标BMI控制在24以下；2. 规律运动，每周至少3-4次；3. 饮食控制，减少高糖高脂食物；4. 定期监测血糖；5. 保持良好的作息习惯。",
			Metadata: map[string]interface{}{
				"category": "health",
				"author":   "system",
				"priority": "high",
			},
		},
		{
			ID:      "weimeng_weight_loss_goals",
			Content: "威猛先生的减肥目标：从目前的76.6kg减到理想体重范围。根据身高165cm，理想体重应该在55-68kg之间。建议分阶段进行，每月减重2-3kg，避免快速减肥对身体造成伤害。同时需要配合饮食控制和规律运动。",
			Metadata: map[string]interface{}{
				"category": "weight_loss",
				"author":   "system",
				"priority": "medium",
			},
		},
		{
			ID:      "weimeng_programmer_lifestyle",
			Content: "作为程序员，威猛先生的工作特点是长时间坐着。建议：1. 每工作1小时起身活动5-10分钟；2. 保持正确的坐姿，避免颈椎和腰椎问题；3. 控制屏幕时间，保护视力；4. 合理安排工作休息时间；5. 保持规律作息，避免熬夜。",
			Metadata: map[string]interface{}{
				"category": "lifestyle",
				"author":   "system",
				"priority": "medium",
			},
		},
		{
			ID:      "weimeng_diet_recommendations",
			Content: "针对威猛先生的糖尿病和减肥需求，饮食建议包括：1. 增加蔬菜摄入，每天至少300-500g；2. 选择低GI值的碳水化合物，如糙米、燕麦等；3. 适量优质蛋白质，如鱼、鸡肉、豆制品；4. 控制油脂摄入，选择健康油脂如橄榄油；5. 少食多餐，避免暴饮暴食；6. 充足饮水，每天1.5-2L。",
			Metadata: map[string]interface{}{
				"category": "diet_recommendation",
				"author":   "system",
				"priority": "high",
			},
		},
		{
			ID:      "weimeng_exercise_benefits",
			Content: "运动对威猛先生的好处：1. 帮助控制血糖，提高胰岛素敏感性；2. 促进脂肪燃烧，帮助减重；3. 改善心血管健康；4. 缓解工作压力，改善睡眠；5. 增强体质，提高免疫力。建议选择适合的运动方式，如快走、慢跑、游泳等有氧运动。",
			Metadata: map[string]interface{}{
				"category": "exercise_benefits",
				"author":   "system",
				"priority": "medium",
			},
		},
	}

	// 批量添加文档
	fmt.Printf("准备添加 %d 个健康知识文档...\n", len(healthDocs))
	
	if err := ragService.AddDocuments(ctx, healthDocs); err != nil {
		log.Fatalf("添加健康知识文档失败: %v", err)
	}
	
	fmt.Printf("✅ 成功添加 %d 个健康知识文档到知识库\n", len(healthDocs))
	
	// 等待索引更新
	time.Sleep(2 * time.Second)
	
	// 测试检索
	fmt.Println("\n=== 测试知识库检索 ===")
	testQueries := []string{
		"威猛先生的运动计划",
		"他的饮食习惯",
		"糖尿病管理建议",
		"减肥目标",
	}
	
	for i, query := range testQueries {
		fmt.Printf("\n--- 测试查询 %d: %s ---\n", i+1, query)
		
		results, err := ragService.Search(ctx, rag.SearchQuery{
			Query: query,
			TopK:  3,
		})
		if err != nil {
			log.Printf("检索失败: %v", err)
			continue
		}
		
		if len(results) == 0 {
			fmt.Println("没有找到相关结果")
			continue
		}
		
		for j, result := range results {
			fmt.Printf("  结果 %d (相似度: %.4f):\n", j+1, result.Score)
			fmt.Printf("    内容: %s\n", result.Content)
		}
	}
	
	fmt.Println("\n=== 健康知识添加完成 ===")
}
