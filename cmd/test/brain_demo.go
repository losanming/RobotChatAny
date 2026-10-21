package main

import (
	"context"
	"fmt"
	"time"

	"new-robot-self/internal/awareness"
	"new-robot-self/internal/brain"
	loc "new-robot-self/pkg/location"
)

func main() {
	// 构造一个极简的 awareness.Context（实际项目中可对接 awareness_demo 的 JSON 输出）
	aCtx := awareness.Context{
		TS:       time.Now().Unix(),
		DeviceID: "dev-001",
		Source:   "gps",
		LocTrans: loc.TransformedLocation{
			TS:       time.Now().Unix(),
			LatGCJ02: 39.908823,
			LonGCJ02: 116.397470,
			LatWGS84: 39.909,
			LonWGS84: 116.397,
			InChina:  true,
			Quality:  10,
		},
		Regeo: &loc.ReverseGeocodeResult{
			TS:       0,
			LatGCJ02: 0,
			LonGCJ02: 0,
			Adcode:   "",
			Address:  "",
			Country:  "",
			Province: "",
			City:     nil,
			CityCode: "",
			District: "",
			Township: "",
			TownCode: "",
			POIs:     nil,
			Provider: "",
		},
		Around: &loc.PlaceAroundResult{
			TS:        0,
			LatGCJ02:  0,
			LonGCJ02:  0,
			Radius:    0,
			Keywords:  "",
			Types:     "",
			SortRule:  "",
			Region:    "",
			CityLimit: false,
			PageSize:  0,
			PageNum:   0,
			Count:     0,
			POIs:      nil,
			Provider:  "",
		},
		Weather: &loc.WeatherResult{
			TS:         0,
			Provider:   "",
			Extensions: "",
			CityAdcode: "",
			Count:      0,
			Lives:      nil,
			Forecasts:  nil,
		},
		Vision: &awareness.ImageInsight{
			Labels:  nil,
			Objects: nil,
			Summary: "",
			Extra:   nil,
		},
	}

	b := brain.NewBrain()
	b.UpdateFromAwareness(context.Background(), aCtx)
	// 模拟一条对话输入
	b.UpdateFromDialog(context.Background(), "user", "我想下午找个咖啡店写代码", time.Now().Unix())

	fmt.Println("===== BRAIN PROMPT =====")
	fmt.Println(b.GetPrompt(context.Background(), 3))
}
