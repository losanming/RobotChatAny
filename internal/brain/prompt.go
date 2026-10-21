package brain

import (
	"fmt"
	"strings"
)

// BuildBrainPrompt 将 WorldState 压缩为给大模型的决策 Prompt
// 内容：Top-K 目的假设 + 环境要点（来自 awareness）+ 近对话线索
func BuildBrainPrompt(s *WorldState, topK int) string {
	var b strings.Builder
	b.WriteString("# 大脑决策输入\n")
	b.WriteString(fmt.Sprintf("state_version: %s\nupdated_ts_ms: %d\n", s.Version, s.UpdatedTS))

	// 目的假设 Top-K
	b.WriteString("\n## 目的假设(Top-K)\n")
	k := topK
	if k > len(s.GoalHypotheses) {
		k = len(s.GoalHypotheses)
	}
	for i := 0; i < k; i++ {
		g := s.GoalHypotheses[i]
		b.WriteString(fmt.Sprintf("- [%d] %s (conf=%.2f, last=%d)\n", i+1, g.Title, g.Confidence, g.LastUpdated))
	}
	if k == 0 {
		b.WriteString("(无当前目的假设)\n")
	}

	// 环境要点
	if s.Env != nil {
		a := s.Env.Awareness
		b.WriteString("\n## 环境要点\n")
		b.WriteString(fmt.Sprintf("gcj02=(%.6f,%.6f), in_china=%v\n", a.LocTrans.LatGCJ02, a.LocTrans.LonGCJ02, a.LocTrans.InChina))
		if a.Regeo != nil {
			addr := a.Regeo.Address
			if addr == "" {
				addr = "(无格式化地址)"
			}
			b.WriteString(fmt.Sprintf("addr=%s; adcode=%s; city=%s\n", addr, a.Regeo.Adcode, strings.Join(a.Regeo.City, "/")))
		}
		if a.Weather != nil && len(a.Weather.Lives) > 0 {
			lv := a.Weather.Lives[0]
			b.WriteString(fmt.Sprintf("weather=%s %s°C 湿度%s%% 风向%s 风力%s\n", lv.Weather, lv.TemperatureC, lv.Humidity, lv.WindDirection, lv.WindPower))
		}
		if a.Around != nil && len(a.Around.POIs) > 0 {
			b.WriteString("nearby_top3:\n")
			for i := 0; i < len(a.Around.POIs) && i < 3; i++ {
				p := a.Around.POIs[i]
				b.WriteString(fmt.Sprintf("- %s type=%s addr=%s\n", p.Name, p.Type, p.Address))
			}
		}
	}

	// 对话线索
	if s.Dialog != nil && len(s.Dialog.History) > 0 {
		b.WriteString("\n## 近对话(末2条)\n")
		start := len(s.Dialog.History) - 2
		if start < 0 { start = 0 }
		for i := start; i < len(s.Dialog.History); i++ {
			t := s.Dialog.History[i]
			b.WriteString(fmt.Sprintf("- [%s] %s\n", t.Role, t.Content))
		}
		if s.Dialog.LastUserIntent != "" {
			b.WriteString(fmt.Sprintf("last_user_intent: %s\n", s.Dialog.LastUserIntent))
		}
	}

	// 指令
	b.WriteString("\n# 你的任务\n")
	b.WriteString("请在不确定初始目的的前提下，结合“目的假设与环境要点”，先进行目的确认或补充关键缺失信息，然后给出1-2条下一步建议或要问的问题（中文，简洁具体）。\n")
	return b.String()
}
