package initialize

import (
	"context"
	"fmt"
	"time"

	"new-robot-self/config"
	"new-robot-self/global"
	"new-robot-self/pkg/rag"
)

// InitRAG 初始化RAG服务
func InitRAG() error {
	cfg := config.GetConfig()
	
	// 检查RAG配置是否存在
	if cfg.RAG.Embedding.BaseURL == "" || cfg.RAG.Qdrant.BaseURL == "" {
		return fmt.Errorf("RAG config is incomplete")
	}

	// 创建RAG服务
	global.RAGService = rag.NewService(cfg.RAG)

	// 健康检查
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := global.RAGService.HealthCheck(ctx); err != nil {
		return fmt.Errorf("RAG service health check failed: %w", err)
	}

	fmt.Println("RAG service initialized successfully")
	return nil
}
