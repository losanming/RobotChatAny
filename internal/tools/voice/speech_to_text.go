package voice

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/bytedance/sonic"
	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/schema"
)

// SpeechToTextTool 阿里云语音转文字工具
// 支持实时语音识别和录音文件识别
type SpeechToTextTool struct {
	EnableLearning bool
	Category       string
	ToolName       string
	client         *http.Client
	apiKey         string
	appKey         string
}

type speechToTextParams struct {
	Action       string `json:"action"`        // transcribe_realtime, transcribe_file
	AudioData    string `json:"audio_data,omitempty"`    // Base64编码的音频数据
	AudioFile    string `json:"audio_file,omitempty"`    // 音频文件路径
	Format       string `json:"format,omitempty"`       // 音频格式: wav, mp3, aac, m4a
	SampleRate   int    `json:"sample_rate,omitempty"`  // 采样率: 16000, 8000
	Language     string `json:"language,omitempty"`     // 语言: zh-CN, en-US
	EnablePunctuation bool `json:"enable_punctuation,omitempty"` // 是否添加标点
	EnableITN    bool   `json:"enable_itn,omitempty"`    // 是否启用ITN（逆文本标准化）
}

type transcriptionResult struct {
	Success     bool   `json:"success"`
	Text        string `json:"text"`
	Duration    float64 `json:"duration"`    // 音频时长（秒）
	WordCount   int    `json:"word_count"`   // 词数
	Sentences   []sentence `json:"sentences"` // 分句结果
	Confidence  float64 `json:"confidence"`  // 置信度
	ProcessedAt int64  `json:"processed_at"` // 处理时间
}

type sentence struct {
	Text       string  `json:"text"`
	StartTime  float64 `json:"start_time"` // 开始时间（秒）
	EndTime    float64 `json:"end_time"`   // 结束时间（秒）
	Confidence float64 `json:"confidence"` // 置信度
}

// GetCategory 获取工具分类
func (t *SpeechToTextTool) GetCategory() string {
	return t.Category
}

// GetToolName 获取工具名称
func (t *SpeechToTextTool) GetToolName() string {
	return t.ToolName
}

// AfterRun 工具调用后的后处理
func (t *SpeechToTextTool) AfterRun(ctx context.Context, result string, err error) {
	// 简化实现，不进行学习集成以避免循环依赖
}

func NewSpeechToTextTool() *SpeechToTextTool {
	return &SpeechToTextTool{
		Category:       "voice",
		EnableLearning: true,
		ToolName:       "speech_to_text",
		client:         &http.Client{Timeout: 30 * time.Second},
		apiKey:         os.Getenv("ALIYUN_ACCESS_KEY_ID"),
		appKey:         os.Getenv("ALIYUN_ACCESS_KEY_SECRET"),
	}
}

func (t *SpeechToTextTool) Info(_ context.Context) (*schema.ToolInfo, error) {
	return &schema.ToolInfo{
		Name: "speech_to_text",
		Desc: "阿里云语音转文字工具，支持实时语音识别和录音文件识别，将语音转换为文字",
		ParamsOneOf: schema.NewParamsOneOfByParams(map[string]*schema.ParameterInfo{
			"action":             {Desc: "操作类型: transcribe_realtime(实时识别), transcribe_file(文件识别)", Type: schema.String, Required: true},
			"audio_data":         {Desc: "Base64编码的音频数据（实时识别使用）", Type: schema.String, Required: false},
			"audio_file":         {Desc: "音频文件路径（文件识别使用）", Type: schema.String, Required: false},
			"format":             {Desc: "音频格式: wav, mp3, aac, m4a", Type: schema.String, Required: false},
			"sample_rate":        {Desc: "采样率: 16000, 8000", Type: schema.Integer, Required: false},
			"language":           {Desc: "语言: zh-CN, en-US", Type: schema.String, Required: false},
			"enable_punctuation": {Desc: "是否添加标点符号", Type: schema.Boolean, Required: false},
			"enable_itn":         {Desc: "是否启用ITN（逆文本标准化）", Type: schema.Boolean, Required: false},
		}),
	}, nil
}

func (t *SpeechToTextTool) InvokableRun(ctx context.Context, argumentsInJSON string, _ ...tool.Option) (string, error) {
	p := speechToTextParams{}
	if err := sonic.UnmarshalString(argumentsInJSON, &p); err != nil {
		return "", err
	}

	// 验证参数
	if p.Action == "" {
		return "", fmt.Errorf("action is required")
	}

	// 设置默认值
	if p.Format == "" {
		p.Format = "wav"
	}
	if p.SampleRate == 0 {
		p.SampleRate = 16000
	}
	if p.Language == "" {
		p.Language = "zh-CN"
	}
	if !p.EnablePunctuation {
		p.EnablePunctuation = true
	}

	// 根据操作类型执行相应功能
	switch p.Action {
	case "transcribe_realtime":
		return t.transcribeRealtime(ctx, p)
	case "transcribe_file":
		return t.transcribeFile(ctx, p)
	default:
		return "", fmt.Errorf("unknown action: %s", p.Action)
	}
}

// transcribeRealtime 实时语音识别
func (t *SpeechToTextTool) transcribeRealtime(ctx context.Context, p speechToTextParams) (string, error) {
	if p.AudioData == "" {
		return "", fmt.Errorf("audio_data is required for transcribe_realtime action")
	}

	// TODO: 调用阿里云实时语音识别API
	// 这里先返回模拟结果
	result := transcriptionResult{
		Success:     true,
		Text:        "这是一个模拟的语音识别结果",
		Duration:    5.2,
		WordCount:   10,
		Confidence:  0.95,
		ProcessedAt: time.Now().Unix(),
		Sentences: []sentence{
			{
				Text:       "这是一个模拟的语音识别结果",
				StartTime:  0.0,
				EndTime:    5.2,
				Confidence: 0.95,
			},
		},
	}

	b, _ := sonic.Marshal(result)
	return string(b), nil
}

// transcribeFile 录音文件识别
func (t *SpeechToTextTool) transcribeFile(ctx context.Context, p speechToTextParams) (string, error) {
	if p.AudioFile == "" {
		return "", fmt.Errorf("audio_file is required for transcribe_file action")
	}

	// 检查文件是否存在
	if _, err := os.Stat(p.AudioFile); os.IsNotExist(err) {
		return "", fmt.Errorf("audio file not found: %s", p.AudioFile)
	}

	// TODO: 调用阿里云录音文件识别API
	// 这里先返回模拟结果
	result := transcriptionResult{
		Success:     true,
		Text:        "这是从音频文件中识别出的文字内容",
		Duration:    10.5,
		WordCount:   15,
		Confidence:  0.92,
		ProcessedAt: time.Now().Unix(),
		Sentences: []sentence{
			{
				Text:       "这是从音频文件中识别出的文字内容",
				StartTime:  0.0,
				EndTime:    10.5,
				Confidence: 0.92,
			},
		},
	}

	b, _ := sonic.Marshal(result)
	return string(b), nil
}

// callAliyunASR 调用阿里云语音识别API
func (t *SpeechToTextTool) callAliyunASR(ctx context.Context, audioData []byte, format string, sampleRate int) (*transcriptionResult, error) {
	// 阿里云语音识别API端点
	apiURL := "https://nls-meta.cn-shanghai.aliyuncs.com/upload"

	// 构建请求参数
	params := url.Values{}
	params.Set("format", format)
	params.Set("sample_rate", fmt.Sprintf("%d", sampleRate))
	params.Set("language", "zh-CN")
	params.Set("enable_punctuation", "true")

	// TODO: 实现阿里云API签名和请求
	// 这里需要实现阿里云的签名算法
	
	req, err := http.NewRequestWithContext(ctx, "POST", apiURL+"?"+params.Encode(), nil)
	if err != nil {
		return nil, err
	}

	resp, err := t.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned status: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// 解析响应
	var result transcriptionResult
	if err := sonic.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	return &result, nil
}