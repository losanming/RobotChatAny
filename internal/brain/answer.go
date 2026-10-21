package brain

import (
	"encoding/json"
	"fmt"
)

// PipelineOutput 工具执行管道输出结构
type PipelineOutput struct {
	Address string            `json:"address,omitempty"`
	Adcode  string            `json:"adcode,omitempty"`
	Lat     float64           `json:"lat,omitempty"`
	Lng     float64           `json:"lng,omitempty"`
	Around  map[string]any    `json:"around,omitempty"`
	Weather map[string]any    `json:"weather,omitempty"`
	Steps   []StepRaw         `json:"steps"`
}

// StepRaw 工具执行步骤原始输出
type StepRaw struct {
	Name string `json:"name"`
	Out  string `json:"out_json"`
}

// RenderAnswer 将 PipelineOutput 渲染为中文自然语言回答（无 LLM 版本）。
func RenderAnswer(p PipelineOutput) string {
	addr := p.Address
	if addr == "" {
		addr = "未知位置"
	}
	base := fmt.Sprintf("你目前在：%s", addr)
	if p.Adcode != "" {
		base += fmt.Sprintf("（adcode=%s）", p.Adcode)
	}

	// 周边概览
	around := ""
	if p.Around != nil {
		// 取出 pois 数量
		if v, ok := p.Around["count"]; ok {
			around += fmt.Sprintf("，附近共检索到 %v 个地点", v)
		}
		// 尝试从 pois 中展示前几个名称
		if pois, ok := p.Around["pois"]; ok {
			bs, _ := json.Marshal(pois)
			var items []map[string]any
			_ = json.Unmarshal(bs, &items)
			n := 0
			names := make([]string, 0, 3)
			for _, it := range items {
				if name, ok := it["name"].(string); ok && name != "" {
					names = append(names, name)
					n++
					if n >= 3 {
						break
					}
				}
			}
			if len(names) > 0 {
				around += fmt.Sprintf("，例如：%s", joinWithComma(names))
			}
		}
	}

	// 天气
	weather := ""
	if p.Weather != nil {
		if lives, ok := p.Weather["lives"]; ok {
			bs, _ := json.Marshal(lives)
			var arr []map[string]any
			_ = json.Unmarshal(bs, &arr)
			if len(arr) > 0 {
				w := arr[0]
				if wth, ok2 := w["weather"].(string); ok2 && wth != "" {
					weather = fmt.Sprintf("，当前天气：%s", wth)
					if t, ok3 := w["temperature"].(string); ok3 && t != "" {
						weather += fmt.Sprintf("，气温：%s℃", t)
					}
				}
			}
		}
	}

	return base + around + weather
}

func joinWithComma(xs []string) string {
	out := ""
	for i, s := range xs {
		if i > 0 {
			out += "、"
		}
		out += s
	}
	return out
}
