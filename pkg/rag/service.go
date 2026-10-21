package rag

import (
	"context"
	"fmt"
	"hash/fnv"
	"strings"
	"time"

	"new-robot-self/config"
	"new-robot-self/pkg/rag/embedding"
	"new-robot-self/pkg/rag/qdrant"
)

// Service RAG服务
type Service struct {
	embeddingClient *embedding.Client
	qdrantClient   *qdrant.Client
	config         config.RAGConfig
}

// NewService 创建RAG服务
func NewService(cfg config.RAGConfig) *Service {
	return &Service{
		embeddingClient: embedding.NewClient(cfg.Embedding),
		qdrantClient:   qdrant.NewClient(cfg.Qdrant),
		config:         cfg,
	}
}

// Document 文档结构
type Document struct {
	ID       string
	Content  string
	Metadata map[string]interface{}
}

// SearchQuery 搜索查询
type SearchQuery struct {
	Query string
	TopK  int
}

// SearchResult 搜索结果
type SearchResult struct {
	Content  string
	Score    float64
	Metadata map[string]interface{}
}

// AddDocument 添加文档到向量数据库
func (s *Service) AddDocument(ctx context.Context, doc Document) error {
	// 1. 生成向量
	embedResult, err := s.embeddingClient.EmbedSingle(ctx, doc.Content)
	if err != nil {
		return fmt.Errorf("embed document failed: %w", err)
	}

	// 2. 构建向量点（使用UUID作为ID）
	point := qdrant.Point{
		ID:     generateUUID(doc.ID), // 转换为UUID格式
		Vector: embedResult.Vector,
		Payload: map[string]interface{}{
			"content":  doc.Content,
			"metadata": doc.Metadata,
			"doc_id":   doc.ID, // 保存原始ID
		},
	}

	// 3. 插入向量数据库
	if err := s.qdrantClient.UpsertSingle(ctx, point); err != nil {
		return fmt.Errorf("upsert to qdrant failed: %w", err)
	}

	return nil
}

// AddDocuments 批量添加文档
func (s *Service) AddDocuments(ctx context.Context, docs []Document) error {
	if len(docs) == 0 {
		return fmt.Errorf("documents cannot be empty")
	}

	points := make([]qdrant.Point, 0, len(docs))

	for _, doc := range docs {
		// 生成向量
		embedResult, err := s.embeddingClient.EmbedSingle(ctx, doc.Content)
		if err != nil {
			return fmt.Errorf("embed document %s failed: %w", doc.ID, err)
		}

		// 构建向量点（使用UUID作为ID）
		point := qdrant.Point{
			ID:     generateUUID(doc.ID), // 转换为UUID格式
			Vector: embedResult.Vector,
			Payload: map[string]interface{}{
				"content":  doc.Content,
				"metadata": doc.Metadata,
				"doc_id":   doc.ID, // 保存原始ID
			},
		}

		points = append(points, point)
	}

	// 批量插入向量数据库
	if err := s.qdrantClient.Upsert(ctx, points); err != nil {
		return fmt.Errorf("batch upsert to qdrant failed: %w", err)
	}

	return nil
}

// Search 搜索相关文档
func (s *Service) Search(ctx context.Context, query SearchQuery) ([]SearchResult, error) {
	if query.Query == "" {
		return nil, fmt.Errorf("query cannot be empty")
	}

	// 1. 生成查询向量
	embedResult, err := s.embeddingClient.EmbedSingle(ctx, query.Query)
	if err != nil {
		return nil, fmt.Errorf("embed query failed: %w", err)
	}

	// 2. 向量搜索
	topK := query.TopK
	if topK <= 0 {
		topK = s.config.Qdrant.TopK
	}

	// 临时设置topK（这里需要修改qdrant client支持动态topK）
	// 为了简单起见，我们使用配置中的topK
	qdrantResults, err := s.qdrantClient.Search(ctx, embedResult.Vector)
	if err != nil {
		return nil, fmt.Errorf("search in qdrant failed: %w", err)
	}

	// 3. 转换结果
	results := make([]SearchResult, 0, len(qdrantResults))
	for _, result := range qdrantResults {
		content, _ := result.Payload["content"].(string)
		metadata, _ := result.Payload["metadata"].(map[string]interface{})

		results = append(results, SearchResult{
			Content:  content,
			Score:    result.Score,
			Metadata: metadata,
		})
	}

	return results, nil
}

// HealthCheck 健康检查
func (s *Service) HealthCheck(ctx context.Context) error {
	// 检查Qdrant连接
	if err := s.qdrantClient.GetCollectionInfo(ctx); err != nil {
		// 如果集合不存在，尝试创建
		if strings.Contains(err.Error(), "doesn't exist") {
			fmt.Printf("集合不存在，尝试创建集合: %s\n", s.config.Qdrant.Collection)
			if createErr := s.qdrantClient.CreateCollection(ctx, s.config.Embedding.Dimensions); createErr != nil {
				return fmt.Errorf("创建集合失败: %w", createErr)
			}
			fmt.Println("集合创建成功")
		} else {
			return fmt.Errorf("qdrant health check failed: %w", err)
		}
	}

	// 检查Embedding服务（简单测试）
	testCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	
	_, err := s.embeddingClient.EmbedSingle(testCtx, "health check")
	if err != nil {
		return fmt.Errorf("embedding service health check failed: %w", err)
	}

	return nil
}

// GetEmbeddingClient 获取embedding客户端（供其他模块直接使用）
func (s *Service) GetEmbeddingClient() *embedding.Client {
	return s.embeddingClient
}

// GetQdrantClient 获取qdrant客户端（供其他模块直接使用）
func (s *Service) GetQdrantClient() *qdrant.Client {
	return s.qdrantClient
}

// generateUUID 从字符串生成UUID格式的ID
func generateUUID(input string) string {
	h := fnv.New32a()
	h.Write([]byte(input))
	hash := h.Sum32()
	
	// 简单的UUID格式生成
	return fmt.Sprintf("%08x-%04x-%04x-%04x-%012x",
		hash,
		hash>>16,
		hash>>24,
		hash>>20,
		hash)
}
