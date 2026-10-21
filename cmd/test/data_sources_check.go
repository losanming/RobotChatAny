package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

// 测试不同的数据源
func main() {
	fmt.Println("=== 数据源测试 ===")
	
	// 测试RSS源
	testRSS("https://36kr.com/feed", "36氪RSS")
	testRSS("https://sspai.com/feed", "少数派RSS")
	testRSS("https://www.ithome.com/rss/", "IT之家RSS")
	
	// 测试聚合API
	// testAPI("https://api.injahow.cn/", "聚合API")
	
	fmt.Println("\n=== 测试完成 ===")
}

func testRSS(url, name string) {
	fmt.Printf("\n测试 %s: %s\n", name, url)
	
	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	
	resp, err := client.Get(url)
	if err != nil {
		fmt.Printf("❌ 请求失败: %v\n", err)
		return
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		fmt.Printf("❌ HTTP状态: %d\n", resp.StatusCode)
		return
	}
	
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("❌ 读取失败: %v\n", err)
		return
	}
	
	fmt.Printf("✅ 成功获取数据，长度: %d 字节\n", len(body))
	
	// 检查是否是RSS格式
	if len(body) > 100 {
		content := string(body[:100])
		if contains(content, "<rss") || contains(content, "<feed") {
			fmt.Printf("✅ 确认为RSS/Atom格式\n")
		} else {
			fmt.Printf("⚠️  可能不是RSS格式，前100字符: %s\n", content)
		}
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) && findSubstring(s, substr))
}

func findSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
