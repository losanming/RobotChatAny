package main

import (
	"context"
	"fmt"
	"github.com/cloudwego/eino-ext/components/document/loader/file"
	"github.com/cloudwego/eino/components/document"
)

func main() {
	loader, _ := file.NewFileLoader(context.Background(), &file.FileLoaderConfig{
		UseNameAsID: true,
	})
	filePath := "cmd/test/需求.md"
	docs, _ := loader.Load(context.Background(), document.Source{URI: filePath})
	fmt.Println(docs[0].Content)
}
