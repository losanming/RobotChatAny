package initialize

import (
	"fmt"
	"new-robot-self/pkg/screenshot/quality"
	public_analyzer "new-robot-self/public/analyzer"
)

func InitAnalyzer() error {
	a := quality.NewAnalyzer()
	if a != nil {
		public_analyzer.Analyzer = a
	}
	fmt.Printf("analyzer 图片筛选 初始化配置成功 \n")
	return nil
}
