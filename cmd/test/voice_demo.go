package main

import (
	"context"
	"encoding/json"
	"fmt"

	"new-robot-self/internal/tools/voice"
)

func main() {
	fmt.Println("=== 语音转文字功能测试 ===")
	fmt.Println()

	tool := voice.NewSpeechToTextTool()
	ctx := context.Background()

	// 测试实时语音识别
	testRealtimeTranscription(ctx, tool)

	// 测试文件语音识别
	testFileTranscription(ctx, tool)

	fmt.Println("=== 测试完成 ===")
}

func testRealtimeTranscription(ctx context.Context, tool *voice.SpeechToTextTool) {
	fmt.Println("【1】测试实时语音识别")

	params := map[string]interface{}{
		"action":             "transcribe_realtime",
		"audio_data":         "模拟的Base64音频数据", // 实际使用时需要真实的Base64编码音频
		"format":             "wav",
		"sample_rate":        16000,
		"language":           "zh-CN",
		"enable_punctuation": true,
	}

	paramsJSON, _ := json.Marshal(params)
	result, err := tool.InvokableRun(ctx, string(paramsJSON))
	if err != nil {
		fmt.Printf("❌ 实时语音识别失败: %v\n", err)
	} else {
		fmt.Printf("✅ 实时语音识别成功: %s\n", result)
		
		// 解析结果
		var transcription TranscriptionResult
		json.Unmarshal([]byte(result), &transcription)
		
		if transcription.Success {
			fmt.Printf("📝 识别文字: %s\n", transcription.Text)
			fmt.Printf("⏱️  时长: %.2f秒\n", transcription.Duration)
			fmt.Printf("🎯 置信度: %.2f\n", transcription.Confidence)
			fmt.Printf("📊 词数: %d\n", transcription.WordCount)
		}
	}
	
	fmt.Println()
}

func testFileTranscription(ctx context.Context, tool *voice.SpeechToTextTool) {
	fmt.Println("【2】测试文件语音识别")

	params := map[string]interface{}{
		"action":             "transcribe_file",
		"audio_file":         "/tmp/test_audio.wav", // 需要提供实际存在的音频文件
		"format":             "wav",
		"sample_rate":        16000,
		"language":           "zh-CN",
		"enable_punctuation": true,
	}

	paramsJSON, _ := json.Marshal(params)
	result, err := tool.InvokableRun(ctx, string(paramsJSON))
	if err != nil {
		fmt.Printf("❌ 文件语音识别失败: %v\n", err)
		fmt.Printf("💡 提示: 需要提供实际存在的音频文件路径\n")
	} else {
		fmt.Printf("✅ 文件语音识别成功: %s\n", result)
		
		// 解析结果
		var transcription TranscriptionResult
		json.Unmarshal([]byte(result), &transcription)
		
		if transcription.Success {
			fmt.Printf("📝 识别文字: %s\n", transcription.Text)
			fmt.Printf("⏱️  时长: %.2f秒\n", transcription.Duration)
			fmt.Printf("🎯 置信度: %.2f\n", transcription.Confidence)
			fmt.Printf("📊 词数: %d\n", transcription.WordCount)
			
			// 显示分句结果
			if len(transcription.Sentences) > 0 {
				fmt.Println("📋 分句结果:")
				for i, sentence := range transcription.Sentences {
					fmt.Printf("   %d. %s (%.1f-%.1fs, 置信度: %.2f)\n", 
						i+1, sentence.Text, sentence.StartTime, sentence.EndTime, sentence.Confidence)
				}
			}
		}
	}
	
	fmt.Println()
}

// 为了测试方便，添加导出的结构体
type TranscriptionResult struct {
	Success     bool      `json:"success"`
	Text        string    `json:"text"`
	Duration    float64   `json:"duration"`
	WordCount   int       `json:"word_count"`
	Sentences   []Sentence `json:"sentences"`
	Confidence  float64   `json:"confidence"`
	ProcessedAt int64     `json:"processed_at"`
}

type Sentence struct {
	Text       string  `json:"text"`
	StartTime  float64 `json:"start_time"`
	EndTime    float64 `json:"end_time"`
	Confidence float64 `json:"confidence"`
}