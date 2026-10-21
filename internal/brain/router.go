package brain

import "strings"

// Intent 表示用户意图的最小集合
// where_am_i: 我现在在哪里/在哪里
// what_is_around: 这里有什么/附近/周边
// weather_now: 天气/温度/下雨
// screen_context: 屏幕/显示什么

type Intent string

const (
	IntentWhereAmI     Intent = "where_am_i"
	IntentWhatIsAround Intent = "what_is_around"
	IntentWeatherNow   Intent = "weather_now"
	IntentScreenCtx    Intent = "screen_context"
)

// DetectIntents 基于关键词的轻量意图识别（低成本、可控）
func DetectIntents(text string) []Intent {
	lt := strings.ToLower(strings.TrimSpace(text))
	cn := strings.TrimSpace(text)
	var out []Intent

	// where am i
	if strings.Contains(cn, "我在哪") || strings.Contains(cn, "在哪里") || strings.Contains(cn, "我现在在哪") {
		out = append(out, IntentWhereAmI)
	}
	// what is around
	if strings.Contains(cn, "这里有什么") || strings.Contains(cn, "周边") || strings.Contains(cn, "附近") || strings.Contains(cn, "周围") {
		out = append(out, IntentWhatIsAround)
	}
	// weather
	if strings.Contains(cn, "天气") || strings.Contains(cn, "温度") || strings.Contains(cn, "下雨") || strings.Contains(lt, "rain") {
		out = append(out, IntentWeatherNow)
	}
	// screen
	if strings.Contains(cn, "屏幕") || strings.Contains(cn, "显示什么") || strings.Contains(cn, "窗口") {
		out = append(out, IntentScreenCtx)
	}

	// 如果都没命中，根据开放提问的泛化，常见的是位置/周边问法，此处可按需默认加 where_am_i
	if len(out) == 0 {
		// 保守返回空，由上层决定是否使用 LLM 做兜底分类
	}
	return out
}
