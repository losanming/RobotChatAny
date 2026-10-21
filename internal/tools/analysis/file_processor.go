package analysis

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/schema"
)

// FileProcessor 文件处理工具
type FileProcessor struct {
	enableLearning bool
	category       string
	toolName       string
}

// NewFileProcessor 创建文件处理工具
func NewFileProcessor() *FileProcessor {
	return &FileProcessor{
		enableLearning: true,
		category:       "analysis",
		toolName:       "file_processor",
	}
}

// Info 返回工具信息
func (f *FileProcessor) Info(ctx context.Context) (*schema.ToolInfo, error) {
	return &schema.ToolInfo{
		Name: "file_processor",
		Desc: "文件处理工具，支持文本文件读取、基本信息提取、内容分析、格式转换",
		ParamsOneOf: schema.NewParamsOneOfByParams(map[string]*schema.ParameterInfo{
			"action": {
				Desc:     "操作类型",
				Type:     schema.String,
				Required: true,
				Enum:     []string{"read_text", "get_info", "analyze_content", "convert_format"},
			},
			"file_path": {
				Desc:     "文件路径",
				Type:     schema.String,
				Required: false,
			},
			"encoding": {
				Desc:     "文件编码",
				Type:     schema.String,
				Required: false,
			},
			"max_length": {
				Desc:     "最大读取长度",
				Type:     schema.Integer,
				Required: false,
			},
			"output_format": {
				Desc:     "输出格式",
				Type:     schema.String,
				Required: false,
				Enum:     []string{"text", "json", "markdown"},
			},
		}),
	}, nil
}

// InvokableRun 执行工具调用
func (f *FileProcessor) InvokableRun(ctx context.Context, args string, opts ...tool.Option) (string, error) {
	var params map[string]interface{}
	if err := json.Unmarshal([]byte(args), &params); err != nil {
		return "", fmt.Errorf("参数解析失败: %w", err)
	}

	action := params["action"].(string)
	switch action {
	case "read_text":
		return f.readText(params)
	case "get_info":
		return f.getFileInfo(params)
	case "analyze_content":
		return f.analyzeContent(params)
	case "convert_format":
		return f.convertFormat(params)
	default:
		return "", fmt.Errorf("不支持的操作类型: %s", action)
	}
}

// readText 读取文本文件内容
func (f *FileProcessor) readText(params map[string]interface{}) (string, error) {
	filePath := params["file_path"].(string)
	if filePath == "" {
		return "", fmt.Errorf("文件路径不能为空")
	}

	encoding := getStringParam(params, "encoding", "utf-8")
	maxLength := getIntParam(params, "max_length", 0)

	// 读取文件
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf("读取文件失败: %w", err)
	}

	// 转换为字符串
	text := string(content)

	// 限制长度
	if maxLength > 0 && len(text) > maxLength {
		text = text[:maxLength] + "... (内容已截断)"
	}

	result := map[string]interface{}{
		"success":  true,
		"file_path": filePath,
		"encoding": encoding,
		"content":  text,
		"length":   len(content),
		"truncated": maxLength > 0 && len(content) > maxLength,
	}

	return toJSON(result), nil
}

// getFileInfo 获取文件基本信息
func (f *FileProcessor) getFileInfo(params map[string]interface{}) (string, error) {
	filePath := params["file_path"].(string)
	if filePath == "" {
		return "", fmt.Errorf("文件路径不能为空")
	}

	// 获取文件信息
	fileInfo, err := ioutil.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf("获取文件信息失败: %w", err)
	}

	// 获取文件扩展名
	ext := strings.ToLower(filepath.Ext(filePath))

	// 分析文件类型
	fileType := f.detectFileType(ext)

	// 统计信息
	lineCount := strings.Count(string(fileInfo), "\n")
	wordCount := len(strings.Fields(string(fileInfo)))
	charCount := len(fileInfo)

	result := map[string]interface{}{
		"success":     true,
		"file_path":   filePath,
		"file_type":   fileType,
		"extension":   ext,
		"size_bytes":  len(fileInfo),
		"size_kb":     float64(len(fileInfo)) / 1024,
		"line_count":  lineCount,
		"word_count":  wordCount,
		"char_count": charCount,
	}

	return toJSON(result), nil
}

// analyzeContent 分析文件内容
func (f *FileProcessor) analyzeContent(params map[string]interface{}) (string, error) {
	filePath := params["file_path"].(string)
	if filePath == "" {
		return "", fmt.Errorf("文件路径不能为空")
	}

	// 读取文件
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf("读取文件失败: %w", err)
	}

	text := string(content)

	// 内容分析
	analysis := f.performContentAnalysis(text)

	result := map[string]interface{}{
		"success":  true,
		"file_path": filePath,
		"analysis": analysis,
	}

	return toJSON(result), nil
}

// convertFormat 转换文件格式
func (f *FileProcessor) convertFormat(params map[string]interface{}) (string, error) {
	filePath := params["file_path"].(string)
	outputFormat := getStringParam(params, "output_format", "text")

	if filePath == "" {
		return "", fmt.Errorf("文件路径不能为空")
	}

	// 读取文件
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf("读取文件失败: %w", err)
	}

	text := string(content)

	// 格式转换
	var convertedContent string
	switch outputFormat {
	case "json":
		convertedContent = fmt.Sprintf(`{"content": %q, "length": %d}`, text, len(text))
	case "markdown":
		convertedContent = fmt.Sprintf("```text\n%s\n```", text)
	default: // text
		convertedContent = text
	}

	result := map[string]interface{}{
		"success":          true,
		"file_path":        filePath,
		"output_format":    outputFormat,
		"converted_content": convertedContent,
		"original_length":  len(content),
		"converted_length": len(convertedContent),
	}

	return toJSON(result), nil
}

// detectFileType 检测文件类型
func (f *FileProcessor) detectFileType(ext string) string {
	fileTypes := map[string]string{
		".txt":  "text/plain",
		".md":   "text/markdown",
		".json": "application/json",
		".xml":  "application/xml",
		".html": "text/html",
		".csv":  "text/csv",
		".go":   "text/go",
		".py":   "text/python",
		".js":   "text/javascript",
		".ts":   "text/typescript",
		".java": "text/java",
		".c":    "text/c",
		".cpp":  "text/cpp",
		".h":    "text/c",
		".hpp":  "text/cpp",
	}

	if fileType, ok := fileTypes[ext]; ok {
		return fileType
	}

	return "application/octet-stream"
}

// performContentAnalysis 执行内容分析
func (f *FileProcessor) performContentAnalysis(text string) map[string]interface{} {
	// 基本统计
	lines := strings.Split(text, "\n")
	words := strings.Fields(text)
	
	// 语言检测（简单版）
	language := f.detectLanguage(text)
	
	// 内容特征
	hasNumbers := f.containsNumbers(text)
	hasSpecialChars := f.containsSpecialChars(text)
	avgLineLength := 0
	if len(lines) > 0 {
		totalLength := 0
		for _, line := range lines {
			totalLength += len(line)
		}
		avgLineLength = totalLength / len(lines)
	}

	// 关键词提取（简单版）
	keywords := f.extractKeywords(text, 10)

	return map[string]interface{}{
		"line_count":       len(lines),
		"word_count":       len(words),
		"char_count":       len(text),
		"avg_line_length":  avgLineLength,
		"detected_language": language,
		"has_numbers":      hasNumbers,
		"has_special_chars": hasSpecialChars,
		"keywords":         keywords,
		"readability_score": f.calculateReadability(text),
	}
}

// detectLanguage 检测语言（简单版）
func (f *FileProcessor) detectLanguage(text string) string {
	// 简单的语言检测逻辑
	chineseChars := 0
	englishChars := 0
	
	for _, r := range text {
		if r >= 0x4e00 && r <= 0x9fff {
			chineseChars++
		} else if (r >= 0x0041 && r <= 0x005a) || (r >= 0x0061 && r <= 0x007a) {
			englishChars++
		}
	}

	if chineseChars > englishChars {
		return "chinese"
	} else if englishChars > chineseChars {
		return "english"
	}
	
	return "unknown"
}

// containsNumbers 检查是否包含数字
func (f *FileProcessor) containsNumbers(text string) bool {
	for _, r := range text {
		if r >= '0' && r <= '9' {
			return true
		}
	}
	return false
}

// containsSpecialChars 检查是否包含特殊字符
func (f *FileProcessor) containsSpecialChars(text string) bool {
	specialChars := "!@#$%^&*()_+-=[]{}|;':\",./<>?"
	for _, r := range text {
		if strings.ContainsRune(specialChars, r) {
			return true
		}
	}
	return false
}

// extractKeywords 提取关键词（简单版）
func (f *FileProcessor) extractKeywords(text string, maxCount int) []string {
	// 简单的关键词提取：分词后统计频率
	words := strings.Fields(strings.ToLower(text))
	
	wordCount := make(map[string]int)
	for _, word := range words {
		if len(word) > 2 { // 忽略短词
			wordCount[word]++
		}
	}

	// 按频率排序
	type wordFreq struct {
		word  string
		count int
	}
	
	var wordFreqs []wordFreq
	for word, count := range wordCount {
		wordFreqs = append(wordFreqs, wordFreq{word, count})
	}

	// 简单排序
	for i := 0; i < len(wordFreqs); i++ {
		for j := i + 1; j < len(wordFreqs); j++ {
			if wordFreqs[j].count > wordFreqs[i].count {
				wordFreqs[i], wordFreqs[j] = wordFreqs[j], wordFreqs[i]
			}
		}
	}

	// 取前N个
	result := make([]string, 0, maxCount)
	for i := 0; i < min(maxCount, len(wordFreqs)); i++ {
		result = append(result, wordFreqs[i].word)
	}

	return result
}

// calculateReadability 计算可读性分数（简单版）
func (f *FileProcessor) calculateReadability(text string) float64 {
	if len(text) == 0 {
		return 0
	}

	sentences := strings.Count(text, ".") + strings.Count(text, "!") + strings.Count(text, "?")
	words := len(strings.Fields(text))
	
	if sentences == 0 {
		sentences = 1
	}

	// 简化的可读性计算
	avgSentenceLength := float64(words) / float64(sentences)
	readability := 100 - (avgSentenceLength * 1.5)
	
	if readability < 0 {
		readability = 0
	}
	if readability > 100 {
		readability = 100
	}

	return readability
}

// AfterRun 工具调用后的后处理
func (f *FileProcessor) AfterRun(ctx context.Context, result string, err error) {
	// 简化实现，暂时不做学习集成
}

// 辅助函数
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}