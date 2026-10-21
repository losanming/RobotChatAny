package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"new-robot-self/config"
	"new-robot-self/global"
	"new-robot-self/internal/brain"
	"new-robot-self/pkg/rag"
)

func main() {
	fmt.Println("=== RAG集成功能测试 ===")
	
	// 加载配置
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("加载配置失败: %v", err)
	}

	// 初始化RAG服务
	ragService := rag.NewService(cfg.RAG)
	global.RAGService = ragService // 设置全局实例
	
	// 健康检查
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	
	if err := ragService.HealthCheck(ctx); err != nil {
		log.Fatalf("RAG服务健康检查失败: %v", err)
	}
	fmt.Println("✅ RAG服务健康检查通过")

	// 创建Brain实例
	brain := brain.NewBrain()
	
	// 测试RAG集成
	fmt.Println("\n=== 测试RAG集成对话 ===")
	testRAGIntegration(ctx, brain)
	
	fmt.Println("\n=== RAG集成测试完成 ===")
}

func testRAGIntegration(ctx context.Context, brain *brain.Brain) {
	// 测试问题列表
	testQueries := []string{
		"威猛先生的运动计划是什么？",
		"他的饮食习惯如何？",
		"糖尿病患者的饮食建议",
		"今天天气怎么样", // 这个不应该触发RAG
		"请介绍一下你自己", // 这个也不应该触发RAG
	}

	for i, query := range testQueries {
		fmt.Printf("\n--- 测试查询 %d: %s ---\n", i+1, query)
		
		startTime := time.Now()
		response, err := brain.GenerateRAGEnhancedResponse(ctx, query)
		if err != nil {
			log.Printf("查询失败: %v", err)
			continue
		}
		
		fmt.Printf("回答: %s\n", response)
		fmt.Printf("耗时: %v\n", time.Since(startTime))
		
		// 简单的延迟避免请求过快
		time.Sleep(500 * time.Millisecond)
	}
}
