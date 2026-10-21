package main

import (
	"fmt"
	"log"
	"new-robot-self/internal/resource"
	"time"
)

func main() {
	// 创建资源获取器
	fetcher, err := resource.NewFetcher("asset/resources.json")
	if err != nil {
		log.Fatalf("创建资源获取器失败: %v", err)
	}

	// 获取游戏和社会新闻
	resources, err := fetcher.GetRandomResources([]resource.ResourceType{
		resource.Tech,
		resource.Social,
	}, 5)
	if err != nil {
		log.Fatalf("获取资源失败: %v", err)
	}
	fmt.Printf("获取资源成功: %d\n", len(resources))
	t := time.Now()
	resources_1, err := fetcher.GetRandomResources([]resource.ResourceType{
		resource.Tech,
		resource.Social,
	}, 5)
	fmt.Printf("第二次获取资源成功: %d\n", len(resources_1))
	fmt.Println("第二次获取耗时", time.Since(t))
}
