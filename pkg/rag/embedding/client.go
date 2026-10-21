package embedding

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"new-robot-self/config"
)

// Client Embedding客户端
type Client struct {
	config     config.EmbeddingConfig
	httpClient *http.Client
}

// NewClient 创建Embedding客户端
func NewClient(cfg config.EmbeddingConfig) *Client {
	return &Client{
		config: cfg,
		httpClient: &http.Client{
			Timeout: time.Duration(cfg.TimeoutSeconds) * time.Second,
		},
	}
}

// EmbeddingRequest Embedding请求
type EmbeddingRequest struct {
	Model    string   `json:"model"`
	Input    []string `json:"input"`
	Encode   bool     `json:"encode,omitempty"`
}

// EmbeddingResponse Embedding响应
type EmbeddingResponse struct {
	Object string      `json:"object"`
	Data   []Embedding `json:"data"`
	Model  string      `json:"model"`
	Usage  Usage       `json:"usage"`
}

// Embedding 单个embedding结果
type Embedding struct {
	Object    string    `json:"object"`
	Embedding []float64 `json:"embedding"`
	Index     int       `json:"index"`
}

// Usage 使用情况
type Usage struct {
	PromptTokens int `json:"prompt_tokens"`
	TotalTokens  int `json:"total_tokens"`
}

// EmbeddingResult embedding结果
type EmbeddingResult struct {
	Vector     []float64
	Dimensions int
	TokenUsage int
}

// Embed 文本向量化
func (c *Client) Embed(ctx context.Context, texts []string) (*EmbeddingResult, error) {
	if len(texts) == 0 {
		return nil, fmt.Errorf("texts cannot be empty")
	}

	// 构建请求URL
	url := c.config.BaseURL + c.config.EndpointPath

	// 构建请求体
	reqBody := EmbeddingRequest{
		Model:  c.config.Model,
		Input:  texts,
		Encode: true,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("marshal request failed: %w", err)
	}

	// 发送HTTP请求
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("create request failed: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	if c.config.APIKey != "" {
		req.Header.Set("Authorization", "Bearer "+c.config.APIKey)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("send request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("request failed with status %d: %s", resp.StatusCode, string(body))
	}

	// 解析响应
	var response EmbeddingResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("decode response failed: %w", err)
	}

	if len(response.Data) == 0 {
		return nil, fmt.Errorf("no embedding data in response")
	}

	// 返回第一个文本的embedding结果（简化处理）
	result := &EmbeddingResult{
		Vector:     response.Data[0].Embedding,
		Dimensions: len(response.Data[0].Embedding),
		TokenUsage: response.Usage.TotalTokens,
	}

	return result, nil
}

// EmbedSingle 单个文本向量化
func (c *Client) EmbedSingle(ctx context.Context, text string) (*EmbeddingResult, error) {
	return c.Embed(ctx, []string{text})
}
