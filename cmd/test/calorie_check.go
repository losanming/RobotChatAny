package main

import (
	"context"
	"encoding/json"
	"fmt"

	"new-robot-self/internal/tools/health"
)

func main() {
	tool := health.NewDietManagerTool()
	ctx := context.Background()

	// 调用工具计算100g大米的卡路里
	params := map[string]interface{}{
		"action":       "calculate",
		"food_name":    "米饭",
		"amount":       100,
		"amount_unit":  "g",
	}

	paramsJSON, _ := json.Marshal(params)
	result, err := tool.InvokableRun(ctx, string(paramsJSON))
	if err != nil {
		fmt.Printf("❌ 调用失败: %v\n", err)
		return
	}

	fmt.Println("=== 100g大米营养信息 ===")
	fmt.Println(result)
	
	// 解析结果以便更好展示
	var data map[string]interface{}
	json.Unmarshal([]byte(result), &data)
	
	if nutrition, ok := data["nutrition"].(map[string]interface{}); ok {
		fmt.Println("\n详细营养成分：")
		if calories, ok := nutrition["calories"].(float64); ok {
			fmt.Printf("🔥 卡路里: %.1f 千卡\n", calories)
		}
		if protein, ok := nutrition["protein"].(float64); ok {
			fmt.Printf("🥩 蛋白质: %.1f 克\n", protein)
		}
		if carbs, ok := nutrition["carbs"].(float64); ok {
			fmt.Printf("🍚 碳水化合物: %.1f 克\n", carbs)
		}
		if fat, ok := nutrition["fat"].(float64); ok {
			fmt.Printf("🥑 脂肪: %.1f 克\n", fat)
		}
		if gi, ok := nutrition["gi"].(float64); ok {
			fmt.Printf("📊 血糖生成指数(GI): %.0f\n", gi)
		}
	}
	
	if friendly, ok := data["diabetes_friendly"].(bool); ok {
		if friendly {
			fmt.Println("\n✅ 糖尿病友好")
		} else {
			fmt.Println("\n⚠️  对糖尿病患者不友好（高GI食物）")
		}
	}
}