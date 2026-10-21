package brain

import (
	"context"
	"fmt"
	"strings"
	"time"

	"new-robot-self/global"
	"new-robot-self/pkg/rag"
)

const (
	// RAG触发关键词
	RAGTriggerKeywords = "健康|运动|饮食|糖尿病|体重|血糖|作息|计划|建议"
)

// GenerateRAGEnhancedResponse 生成RAG增强的回答
func (b *Brain) GenerateRAGEnhancedResponse(ctx context.Context, query string) (string, error) {
	b.mu.RLock()
	defer b.mu.RUnlock()
	
	startTime := time.Now()
	
	// 1. 判断是否需要使用RAG
	if !b.shouldUseRAG(query) {
		return b.GenerateStandardResponse(ctx, query)
	}

	// 2. RAG检索相关文档
	searchCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	
	searchResults, err := global.RAGService.Search(searchCtx, rag.SearchQuery{
		Query: query,
		TopK:  3,
	})
	if err != nil {
		// RAG失败时降级到标准回答
		fmt.Printf("RAG检索失败，降级到标准回答: %v\n", err)
		return b.GenerateStandardResponse(ctx, query)
	}

	if len(searchResults) == 0 {
		fmt.Println("RAG检索无结果，降级到标准回答")
		return b.GenerateStandardResponse(ctx, query)
	}

	fmt.Printf("RAG检索到 %d 个相关文档，耗时: %v\n", len(searchResults), time.Since(startTime))

	// 3. 构建增强prompt
	enhancedPrompt := b.buildRAGPrompt(query, searchResults)

	// 4. 调用LLM生成回答
	response, err := b.callLLM(ctx, enhancedPrompt)
	if err != nil {
		fmt.Printf("LLM调用失败，降级到标准回答: %v\n", err)
		return b.GenerateStandardResponse(ctx, query)
	}

	fmt.Printf("RAG增强回答完成，总耗时: %v\n", time.Since(startTime))
	return response, nil
}

// shouldUseRAG 判断是否需要使用RAG
func (b *Brain) shouldUseRAG(query string) bool {
	if query == "" {
		return false
	}
	
	// 检查是否包含触发关键词
	for _, keyword := range strings.Split(RAGTriggerKeywords, "|") {
		if strings.Contains(query, keyword) {
			return true
		}
	}
	return false
}

// buildRAGPrompt 构建RAG增强的prompt
func (b *Brain) buildRAGPrompt(query string, results []rag.SearchResult) string {
	var contextBuilder strings.Builder

	contextBuilder.WriteString("基于以下知识库内容回答问题：\n\n")
	
	for i, result := range results {
		contextBuilder.WriteString(fmt.Sprintf("【参考资料%d】%s\n", i+1, result.Content))
		if result.Score > 0.7 {
			contextBuilder.WriteString(fmt.Sprintf("(相关度: %.1f%%)\n", result.Score*100))
		}
		contextBuilder.WriteString("\n")
	}

	contextBuilder.WriteString(fmt.Sprintf("问题: %s\n", query))
	contextBuilder.WriteString("\n请基于上述参考资料回答问题，如果参考资料中没有相关信息，请基于你的通用知识回答。")

	return contextBuilder.String()
}

// GenerateStandardResponse 生成标准回答（非RAG）
func (b *Brain) GenerateStandardResponse(ctx context.Context, query string) (string, error) {
	// 调用现有的对话逻辑
	prompt := b.GetPrompt(ctx, 5)
	
	// 这里需要调用实际的LLM服务
	// 暂时返回一个占位回答
	return fmt.Sprintf("标准回答: %s (基于当前上下文: %s)", query, prompt), nil
}

// callLLM 调用LLM生成回答
func (b *Brain) callLLM(ctx context.Context, prompt string) (string, error) {
	// 这里调用实际的LLM服务
	// 暂时返回占位回答，后续需要集成实际的AI服务
	
	// 模拟LLM调用
	time.Sleep(100 * time.Millisecond)
	
	return fmt.Sprintf("RAG增强回答 (基于知识库): %s", prompt), nil
}

// AddKnowledgeToBrain 添加知识到Brain的知识库
func (b *Brain) AddKnowledgeToBrain(ctx context.Context, doc rag.Document) error {
	return global.RAGService.AddDocument(ctx, doc)
}

// SearchKnowledge 搜索Brain知识库
func (b *Brain) SearchKnowledge(ctx context.Context, query string) ([]rag.SearchResult, error) {
	return global.RAGService.Search(ctx, rag.SearchQuery{
		Query: query,
		TopK:  5,
	})
}
