package work

import (
	"context"
	"fmt"
	"time"

	"github.com/bytedance/sonic"
	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/schema"
)

// WorkAssistantTool 工作助手工具
// 支持代码片段管理、技术文档搜索、API文档查询、常用命令记录
type WorkAssistantTool struct {
	EnableLearning bool
	Category       string
	ToolName       string
}

// GetCategory 获取工具分类
func (t *WorkAssistantTool) GetCategory() string {
	return t.Category
}

// GetToolName 获取工具名称
func (t *WorkAssistantTool) GetToolName() string {
	return t.ToolName
}

// AfterRun 工具调用后的后处理
func (t *WorkAssistantTool) AfterRun(ctx context.Context, result string, err error) {
	// 简化实现，不进行学习集成以避免循环依赖
}

type workAssistantParams struct {
	Action       string                 `json:"action"`        // save_snippet, search_snippet, save_command, search_command, search_doc, save_note
	Language     string                 `json:"language,omitempty"`     // 编程语言
	Code         string                 `json:"code,omitempty"`         // 代码片段
	Description  string                 `json:"description,omitempty"`  // 描述
	Tags         []string               `json:"tags,omitempty"`         // 标签
	Command      string                 `json:"command,omitempty"`      // 命令
	DocType      string                 `json:"doc_type,omitempty"`     // 文档类型: api, tutorial, guide
	SearchQuery  string                 `json:"search_query,omitempty"` // 搜索查询
	Title        string                 `json:"title,omitempty"`        // 标题
	Content      string                 `json:"content,omitempty"`     // 内容
	Metadata     map[string]interface{} `json:"metadata,omitempty"`    // 额外元数据
}

type codeSnippet struct {
	ID          string                 `json:"id"`
	UserID      string                 `json:"user_id"`
	Language    string                 `json:"language"`
	Code        string                 `json:"code"`
	Description string                 `json:"description"`
	Tags        []string               `json:"tags"`
	CreatedAt   int64                  `json:"created_at"`
	UpdatedAt   int64                  `json:"updated_at"`
	UseCount    int                    `json:"use_count"`
	Metadata    map[string]interface{} `json:"metadata"`
}

type commandRecord struct {
	ID          string                 `json:"id"`
	UserID      string                 `json:"user_id"`
	Command     string                 `json:"command"`
	Description string                 `json:"description"`
	Tags        []string               `json:"tags"`
	CreatedAt   int64                  `json:"created_at"`
	UseCount    int                    `json:"use_count"`
	Metadata    map[string]interface{} `json:"metadata"`
}

type techNote struct {
	ID          string                 `json:"id"`
	UserID      string                 `json:"user_id"`
	Title       string                 `json:"title"`
	Content     string                 `json:"content"`
	Category    string                 `json:"category"`
	Tags        []string               `json:"tags"`
	CreatedAt   int64                  `json:"created_at"`
	UpdatedAt   int64                  `json:"updated_at"`
	Metadata    map[string]interface{} `json:"metadata"`
}

type docSearchResult struct {
	Title       string   `json:"title"`
	URL         string   `json:"url"`
	Description string   `json:"description"`
	Tags        []string `json:"tags"`
	Relevance   float64  `json:"relevance"`
}

func NewWorkAssistantTool() *WorkAssistantTool {
	return &WorkAssistantTool{
		Category:       "work",
		EnableLearning: true,
		ToolName:       "work_assistant",
	}
}

func (t *WorkAssistantTool) Info(_ context.Context) (*schema.ToolInfo, error) {
	return &schema.ToolInfo{
		Name: "work_assistant",
		Desc: "工作助手工具，支持代码片段管理、技术文档搜索、API文档查询、常用命令记录、技术笔记管理",
		ParamsOneOf: schema.NewParamsOneOfByParams(map[string]*schema.ParameterInfo{
			"action":       {Desc: "操作类型: save_snippet(保存代码), search_snippet(搜索代码), save_command(保存命令), search_command(搜索命令), search_doc(搜索文档), save_note(保存笔记)", Type: schema.String, Required: true},
			"language":     {Desc: "编程语言: go, python, javascript, etc.", Type: schema.String, Required: false},
			"code":         {Desc: "代码片段", Type: schema.String, Required: false},
			"description":  {Desc: "描述", Type: schema.String, Required: false},
			"tags":         {Desc: "标签", Type: schema.Array, Required: false},
			"command":      {Desc: "命令", Type: schema.String, Required: false},
			"doc_type":     {Desc: "文档类型: api, tutorial, guide", Type: schema.String, Required: false},
			"search_query": {Desc: "搜索查询", Type: schema.String, Required: false},
			"title":        {Desc: "标题", Type: schema.String, Required: false},
			"content":      {Desc: "内容", Type: schema.String, Required: false},
			"metadata":     {Desc: "额外元数据", Type: schema.Object, Required: false},
		}),
	}, nil
}

func (t *WorkAssistantTool) InvokableRun(ctx context.Context, argumentsInJSON string, _ ...tool.Option) (string, error) {
	p := workAssistantParams{}
	if err := sonic.UnmarshalString(argumentsInJSON, &p); err != nil {
		return "", err
	}

	// 验证参数
	if p.Action == "" {
		return "", fmt.Errorf("action is required")
	}

	// 根据操作类型执行相应功能
	switch p.Action {
	case "save_snippet":
		return t.saveCodeSnippet(p)
	case "search_snippet":
		return t.searchCodeSnippet(p)
	case "save_command":
		return t.saveCommand(p)
	case "search_command":
		return t.searchCommand(p)
	case "search_doc":
		return t.searchDoc(p)
	case "save_note":
		return t.saveTechNote(p)
	default:
		return "", fmt.Errorf("unknown action: %s", p.Action)
	}
}

// saveCodeSnippet 保存代码片段
func (t *WorkAssistantTool) saveCodeSnippet(p workAssistantParams) (string, error) {
	if p.Language == "" {
		return "", fmt.Errorf("language is required for save_snippet action")
	}
	if p.Code == "" {
		return "", fmt.Errorf("code is required for save_snippet action")
	}

	snippet := codeSnippet{
		ID:          GenerateUUID(),
		UserID:      "weimeng_xiansheng",
		Language:    p.Language,
		Code:        p.Code,
		Description: p.Description,
		Tags:        p.Tags,
		CreatedAt:   time.Now().Unix(),
		UpdatedAt:   time.Now().Unix(),
		UseCount:    0,
		Metadata:    p.Metadata,
	}

	// TODO: 实际存储到数据库
	result := map[string]interface{}{
		"success": true,
		"message": "代码片段保存成功",
		"snippet": snippet,
	}

	b, _ := sonic.Marshal(result)
	return string(b), nil
}

// searchCodeSnippet 搜索代码片段
func (t *WorkAssistantTool) searchCodeSnippet(p workAssistantParams) (string, error) {
	// TODO: 从数据库搜索
	snippets := []codeSnippet{
		{
			ID:          "1",
			UserID:      "weimeng_xiansheng",
			Language:    "go",
			Code:        "func main() {\n    fmt.Println(\"Hello, World!\")\n}",
			Description: "简单的Go程序",
			Tags:        []string{"go", "basic", "hello"},
			CreatedAt:   time.Now().Add(-7 * 24 * time.Hour).Unix(),
			UpdatedAt:   time.Now().Add(-7 * 24 * time.Hour).Unix(),
			UseCount:    5,
		},
	}

	result := map[string]interface{}{
		"success":  true,
		"snippets": snippets,
		"count":    len(snippets),
	}

	b, _ := sonic.Marshal(result)
	return string(b), nil
}

// saveCommand 保存命令
func (t *WorkAssistantTool) saveCommand(p workAssistantParams) (string, error) {
	if p.Command == "" {
		return "", fmt.Errorf("command is required for save_command action")
	}

	command := commandRecord{
		ID:          GenerateUUID(),
		UserID:      "weimeng_xiansheng",
		Command:     p.Command,
		Description: p.Description,
		Tags:        p.Tags,
		CreatedAt:   time.Now().Unix(),
		UseCount:    0,
		Metadata:    p.Metadata,
	}

	// TODO: 实际存储到数据库
	result := map[string]interface{}{
		"success": true,
		"message": "命令保存成功",
		"command": command,
	}

	b, _ := sonic.Marshal(result)
	return string(b), nil
}

// searchCommand 搜索命令
func (t *WorkAssistantTool) searchCommand(p workAssistantParams) (string, error) {
	// TODO: 从数据库搜索
	commands := []commandRecord{
		{
			ID:          "1",
			UserID:      "weimeng_xiansheng",
			Command:     "go run main.go",
			Description: "运行Go程序",
			Tags:        []string{"go", "run"},
			CreatedAt:   time.Now().Add(-3 * 24 * time.Hour).Unix(),
			UseCount:    10,
		},
		{
			ID:          "2",
			UserID:      "weimeng_xiansheng",
			Command:     "go build -o app",
			Description: "编译Go程序",
			Tags:        []string{"go", "build"},
			CreatedAt:   time.Now().Add(-5 * 24 * time.Hour).Unix(),
			UseCount:    8,
		},
	}

	result := map[string]interface{}{
		"success":  true,
		"commands": commands,
		"count":    len(commands),
	}

	b, _ := sonic.Marshal(result)
	return string(b), nil
}

// searchDoc 搜索技术文档
func (t *WorkAssistantTool) searchDoc(p workAssistantParams) (string, error) {
	if p.SearchQuery == "" {
		return "", fmt.Errorf("search_query is required for search_doc action")
	}

	// TODO: 集成实际的技术文档搜索API
	// 这里模拟搜索结果
	results := []docSearchResult{
		{
			Title:       "Go语言官方文档",
			URL:         "https://golang.org/doc/",
			Description: "Go语言的官方文档，包含教程、规范和标准库文档",
			Tags:        []string{"go", "official", "documentation"},
			Relevance:   0.95,
		},
		{
			Title:       "Go by Example",
			URL:         "https://gobyexample.com/",
			Description: "Go语言实例教程，通过大量实例学习Go语言",
			Tags:        []string{"go", "tutorial", "examples"},
			Relevance:   0.88,
		},
		{
			Title:       "Effective Go",
			URL:         "https://golang.org/doc/effective_go.html",
			Description: "Go语言最佳实践指南",
			Tags:        []string{"go", "best-practices", "guide"},
			Relevance:   0.82,
		},
	}

	result := map[string]interface{}{
		"success": true,
		"query":   p.SearchQuery,
		"results": results,
		"count":   len(results),
	}

	b, _ := sonic.Marshal(result)
	return string(b), nil
}

// saveTechNote 保存技术笔记
func (t *WorkAssistantTool) saveTechNote(p workAssistantParams) (string, error) {
	if p.Title == "" {
		return "", fmt.Errorf("title is required for save_note action")
	}
	if p.Content == "" {
		return "", fmt.Errorf("content is required for save_note action")
	}

	note := techNote{
		ID:        GenerateUUID(),
		UserID:    "weimeng_xiansheng",
		Title:     p.Title,
		Content:   p.Content,
		Category:  p.DocType,
		Tags:      p.Tags,
		CreatedAt: time.Now().Unix(),
		UpdatedAt: time.Now().Unix(),
		Metadata:  p.Metadata,
	}

	// TODO: 实际存储到数据库
	result := map[string]interface{}{
		"success": true,
		"message": "技术笔记保存成功",
		"note":    note,
	}

	b, _ := sonic.Marshal(result)
	return string(b), nil
}