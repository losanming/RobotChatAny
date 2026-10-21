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
	// 加载配置
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("加载配置失败: %v", err)
	}

	// 创建RAG服务
	ragService := rag.NewService(cfg.RAG)

	// 健康检查
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	fmt.Println("=== RAG服务健康检查 ===")
	if err := ragService.HealthCheck(ctx); err != nil {
		log.Fatalf("RAG服务健康检查失败: %v", err)
	}
	fmt.Println("✅ RAG服务健康检查通过")

	// 测试添加文档
	fmt.Println("\n=== 测试添加文档 ===")
	testDocs := []rag.Document{
		{
			ID:      "doc_001",
			Content: "威猛先生是一个程序员，身高165cm，体重76.6kg，有轻微的糖尿病。他每三天进行一次运动循环：第一天晚上室内慢跑1小时，第二天室内慢走1小时，第三天有氧全身运动20分钟。",
			Metadata: map[string]interface{}{
				"category": "health",
				"author":   "system",
			},
		},
		{
			ID:      "doc_002",
			Content: "威猛先生的饮食习惯：每天早上一瓶200ml的纯牛奶；中午100g米饭一荤一素；晚上偶尔一根黄瓜一杯黑咖啡。他的目标是减肥和改善糖尿病状况。",
			Metadata: map[string]interface{}{
				"category": "diet",
				"author":   "system",
			},
		},
		{
			ID:      "doc_003",
			Content: "RobotChatAny是一个主动式多模态AI助手项目，具备视觉、听觉和会说话的能力。它使用Golang开发，采用xorm作为ORM框架，支持PostgreSQL、Redis、NSQ等技术栈。",
			Metadata: map[string]interface{}{
				"category": "project",
				"author":   "system",
			},
		},
	}

	if err := ragService.AddDocuments(ctx, testDocs); err != nil {
		log.Fatalf("添加文档失败: %v", err)
	}
	fmt.Printf("✅ 成功添加 %d 个文档到向量数据库\n", len(testDocs))

	// 等待一下确保数据已索引
	time.Sleep(2 * time.Second)

	// 测试搜索
	fmt.Println("\n=== 测试向量搜索 ===")
	searchQueries := []rag.SearchQuery{
		{
			Query: "威猛先生的运动计划是什么？",
			TopK:  3,
		},
		{
			Query: "这个项目使用什么技术栈？",
			TopK:  3,
		},
		{
			Query: "糖尿病患者的饮食建议",
			TopK:  3,
		},
	}

	for i, query := range searchQueries {
		fmt.Printf("\n--- 查询 %d: %s ---\n", i+1, query.Query)
		
		results, err := ragService.Search(ctx, query)
		if err != nil {
			log.Printf("搜索失败: %v", err)
			continue
		}

		if len(results) == 0 {
			fmt.Println("没有找到相关结果")
			continue
		}

		for j, result := range results {
			fmt.Printf("  结果 %d (相似度: %.4f):\n", j+1, result.Score)
			fmt.Printf("    内容: %s\n", result.Content)
			if len(result.Metadata) > 0 {
				fmt.Printf("    元数据: %v\n", result.Metadata)
			}
		}
	}

	// 测试embedding单独功能
	fmt.Println("\n=== 测试Embedding功能 ===")
	embeddingClient := ragService.GetEmbeddingClient()
	testText := "这是一个测试文本"
	embedResult, err := embeddingClient.EmbedSingle(ctx, testText)
	if err != nil {
		log.Printf("Embedding测试失败: %v", err)
	} else {
		fmt.Printf("✅ 文本 '%s' 的向量维度: %d\n", testText, embedResult.Dimensions)
		fmt.Printf("   Token使用: %d\n", embedResult.TokenUsage)
		fmt.Printf("   向量前5个值: %v\n", embedResult.Vector[:5])
	}

	fmt.Println("\n=== RAG功能测试完成 ===")
}
