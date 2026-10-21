package initialize

import (
	"fmt"
	"new-robot-self/global"
	"new-robot-self/internal/resource"
	"new-robot-self/public"
)

func InitFetcher() error {
	fetcher, err := resource.NewFetcher(global.FetcherPath)
	if err != nil {
		return err
	}
	public.Fetcher = fetcher
	fmt.Printf("fetcher 初始化配置成功 \n")
	return nil
}
