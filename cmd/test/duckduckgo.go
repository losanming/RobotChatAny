package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	duck "github.com/cloudwego/eino-ext/components/tool/duckduckgo"
	"github.com/cloudwego/eino-ext/components/tool/duckduckgo/ddgsearch"
)

func main() {
	ctx := context.Background()
	// Create configuration
	config := &duck.Config{
		MaxResults: 3, // Limit to return 3 results
		Region:     ddgsearch.RegionJP,
		DDGConfig: &ddgsearch.Config{
			Timeout:    10 * time.Second,
			Cache:      true,
			MaxRetries: 5,
		},
	}

	// Create search client
	tool, err := duck.NewTool(ctx, config)
	if err != nil {
		log.Fatalf("NewTool of duckduckgo failed, err=%v", err)
	}

	// Create search request
	searchReq := &duck.SearchRequest{
		Query: "Golang programming development",
		Page:  1,
	}

	jsonReq, err := json.Marshal(searchReq)
	if err != nil {
		log.Fatalf("Marshal of search request failed, err=%v", err)
	}

	// Execute search
	resp, err := tool.InvokableRun(ctx, string(jsonReq))
	if err != nil {
		log.Fatalf("Search of duckduckgo failed, err=%v", err)
	}

	var searchResp duck.SearchResponse
	if err := json.Unmarshal([]byte(resp), &searchResp); err != nil {
		log.Fatalf("Unmarshal of search response failed, err=%v", err)
	}

	// Print results
	fmt.Println("Search Results:")
	fmt.Println("==============")
	for i, result := range searchResp.Results {
		fmt.Printf("\n%d. Title: %s\n", i+1, result.Title)
		fmt.Printf("   Link: %s\n", result.Link)
		fmt.Printf("   Description: %s\n", result.Description)
	}
	fmt.Println("")
	fmt.Println("==============")
}
