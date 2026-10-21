package awareness

import (
	"fmt"
	"strings"
)

// BuildPrompt 根据上下文生成适合交给大模型的环境描述
// 约定：所有时间字段均为毫秒时间戳（int64）
func BuildPrompt(c Context) string {
	var b strings.Builder

	// 基本头
	b.WriteString("# 环境感知\n")
	b.WriteString(fmt.Sprintf("ts_ms: %d\n", c.TS))
	if c.DeviceID != "" {
		b.WriteString(fmt.Sprintf("device_id: %s\n", c.DeviceID))
	}
	if c.Source != "" {
		b.WriteString(fmt.Sprintf("source: %s\n", c.Source))
	}

	// 位置与地址
	b.WriteString("\n## 位置\n")
	b.WriteString(fmt.Sprintf("gcj02: lat=%.6f, lon=%.6f; in_china=%v\n", c.LocTrans.LatGCJ02, c.LocTrans.LonGCJ02, c.LocTrans.InChina))
	b.WriteString(fmt.Sprintf("wgs84: lat=%.6f, lon=%.6f\n", c.LocTrans.LatWGS84, c.LocTrans.LonWGS84))
	if c.Regeo != nil {
		addr := c.Regeo.Address
		if addr == "" {
			addr = "(无格式化地址)"
		}
		b.WriteString(fmt.Sprintf("address: %s\n", addr))
		b.WriteString(fmt.Sprintf("adcode: %s, city: %s, district: %s, township: %s\n",
			c.Regeo.Adcode, strings.Join(c.Regeo.City, "/"), c.Regeo.District, c.Regeo.Township))
		if len(c.Regeo.POIs) > 0 {
			b.WriteString("nearby_poi_topk:\n")
			for i, p := range c.Regeo.POIs {
				if i >= 3 { // Top-3 展示
					break
				}
				b.WriteString(fmt.Sprintf("- %s (%s) type=%s\n", p.Name, p.ID, p.Type))
			}
		}
	}

	// 周边（可选）
	if c.Around != nil && len(c.Around.POIs) > 0 {
		b.WriteString("\n## 周边\n")
		b.WriteString(fmt.Sprintf("keywords=%s types=%s radius=%dm count=%d\n", c.Around.Keywords, c.Around.Types, c.Around.Radius, c.Around.Count))
		for i, p := range c.Around.POIs {
			if i >= 5 { // Top-5 展示
				break
			}
			b.WriteString(fmt.Sprintf("- %s type=%s addr=%s loc=%s\n", p.Name, p.Type, p.Address, p.Location))
		}
	}

	// 天气（可选）
	if c.Weather != nil {
		b.WriteString("\n## 天气\n")
		b.WriteString(fmt.Sprintf("query_city_adcode=%s, extensions=%s, count=%d\n", c.Weather.CityAdcode, c.Weather.Extensions, c.Weather.Count))
		if len(c.Weather.Lives) > 0 {
			lv := c.Weather.Lives[0]
			b.WriteString(fmt.Sprintf("live: %s %s°C 湿度%s%% 风向%s 风力%s 报告时间ms=%d\n",
				lv.Weather, lv.TemperatureC, lv.Humidity, lv.WindDirection, lv.WindPower, lv.ReportTimeMS))
		}
		if len(c.Weather.Forecasts) > 0 {
			fc := c.Weather.Forecasts[0]
			b.WriteString(fmt.Sprintf("forecast(city=%s province=%s) 报告时间ms=%d\n", fc.City, fc.Province, fc.ReportTimeMS))
			for i, d := range fc.Casts {
				if i >= 3 { // 最多展示3天
					break
				}
				b.WriteString(fmt.Sprintf("- %s 周%s: 白天%s/%s°C 夜间%s/%s°C 风向%s/%s 级\n",
					d.Date, d.Week, d.DayWeather, d.DayTempC, d.NightWeather, d.NightTempC, d.DayWind, d.NightWind))
			}
		}
	}

	// 视觉（占位）
	if c.Vision != nil {
		b.WriteString("\n## 视觉(占位)\n")
		if c.Vision.Summary != "" {
			b.WriteString("summary: " + c.Vision.Summary + "\n")
		}
		if len(c.Vision.Labels) > 0 {
			b.WriteString("labels: ")
			b.WriteString(strings.Join(c.Vision.Labels, ", ") + "\n")
		}
	}

	// 指令
	b.WriteString("\n# 你的任务\n")
	b.WriteString("你是设备侧的主动式助手，请理解当前用户所处的地理与环境，结合天气和周边信息，推断用户可能的意图，给出下一步的有帮助的建议或问题。回答务必简洁、具体、中文。\n")
	return b.String()
}
