package qdrant

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

// Client Qdrant客户端
type Client struct {
	config     config.QdrantConfig
	httpClient *http.Client
}

// NewClient 创建Qdrant客户端
func NewClient(cfg config.QdrantConfig) *Client {
	return &Client{
		config: cfg,
		httpClient: &http.Client{
			Timeout: time.Duration(cfg.TimeoutSeconds) * time.Second,
		},
	}
}

// Point 向量点
type Point struct {
	ID      string                 `json:"id"`
	Vector  []float64              `json:"vector"`
	Payload map[string]interface{} `json:"payload"`
}

// SearchRequest 搜索请求
type SearchRequest struct {
	Vector        []float64 `json:"vector"`
	Limit         int       `json:"limit"`
	ScoreThreshold float64  `json:"score_threshold,omitempty"`
	WithPayload   bool      `json:"with_payload,omitempty"`
}

// SearchResponse 搜索响应
type SearchResponse struct {
	Result []SearchResult `json:"result"`
	Status string         `json:"status"`
	Time   float64        `json:"time"`
}

// SearchResult 单个搜索结果
type SearchResult struct {
	ID     interface{}            `json:"id"` // 可以是string或uint64
	Score  float64                `json:"score"`
	Payload map[string]interface{} `json:"payload,omitempty"`
}

// UpsertRequest 批量插入请求
type UpsertRequest struct {
	Points []Point `json:"points"`
}

// UpsertResponse 批量插入响应
type UpsertResponse struct {
	Status    string `json:"status"`
	Time      float64 `json:"time"`
	Operation string `json:"operation"`
}

// Upsert 批量插入向量点
func (c *Client) Upsert(ctx context.Context, points []Point) error {
	if len(points) == 0 {
		return fmt.Errorf("points cannot be empty")
	}

	url := fmt.Sprintf("%s/collections/%s/points", c.config.BaseURL, c.config.Collection)

	reqBody := UpsertRequest{
		Points: points,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return fmt.Errorf("marshal request failed: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "PUT", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("create request failed: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	if c.config.APIKey != "" {
		req.Header.Set("api-key", c.config.APIKey)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("send request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("request failed with status %d: %s", resp.StatusCode, string(body))
	}

	var response UpsertResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return fmt.Errorf("decode response failed: %w", err)
	}

	if response.Status != "ok" {
		return fmt.Errorf("upsert failed with status: %s", response.Status)
	}

	return nil
}

// UpsertSingle 插入单个向量点
func (c *Client) UpsertSingle(ctx context.Context, point Point) error {
	return c.Upsert(ctx, []Point{point})
}

// Search 向量搜索
func (c *Client) Search(ctx context.Context, vector []float64) ([]SearchResult, error) {
	if len(vector) == 0 {
		return nil, fmt.Errorf("vector cannot be empty")
	}

	url := fmt.Sprintf("%s/collections/%s/points/search", c.config.BaseURL, c.config.Collection)

	reqBody := SearchRequest{
		Vector:        vector,
		Limit:         c.config.TopK,
		ScoreThreshold: c.config.ScoreThreshold,
		WithPayload:   true,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("marshal request failed: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("create request failed: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	if c.config.APIKey != "" {
		req.Header.Set("api-key", c.config.APIKey)
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

	var response SearchResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("decode response failed: %w", err)
	}

	return response.Result, nil
}

// GetCollectionInfo 获取集合信息（用于健康检查）
func (c *Client) GetCollectionInfo(ctx context.Context) error {
	url := fmt.Sprintf("%s/collections/%s", c.config.BaseURL, c.config.Collection)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return fmt.Errorf("create request failed: %w", err)
	}

	if c.config.APIKey != "" {
		req.Header.Set("api-key", c.config.APIKey)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("send request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("request failed with status %d: %s", resp.StatusCode, string(body))
	}

	return nil
}

// CreateCollection 创建集合
func (c *Client) CreateCollection(ctx context.Context, vectorSize int) error {
	url := fmt.Sprintf("%s/collections/%s", c.config.BaseURL, c.config.Collection)

	// 构建创建集合的请求体
	requestBody := map[string]interface{}{
		"vectors": map[string]interface{}{
			"size":     vectorSize,
			"distance": "Cosine",
		},
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return fmt.Errorf("marshal request failed: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "PUT", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("create request failed: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	if c.config.APIKey != "" {
		req.Header.Set("api-key", c.config.APIKey)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("send request failed: %w", err)
	}
	defer resp.Body.Close()

	// 如果集合已存在，返回成功
	if resp.StatusCode == http.StatusOK {
		return nil
	}

	// 如果是其他错误，返回错误
	body, _ := io.ReadAll(resp.Body)
	return fmt.Errorf("request failed with status %d: %s", resp.StatusCode, string(body))
}
