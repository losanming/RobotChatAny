package initialize

import (
	"errors"
	"fmt"
	"new-robot-self/internal/prompt"
	"new-robot-self/public"
	public_contextManager "new-robot-self/public/contextManager"
	public_context_prompt "new-robot-self/public/contextPromptGenerator"
)

func InitPromptGenerator() error {
	Generator := prompt.NewGenerator(public.Fetcher)
	if Generator == nil {
		return errors.New("prompt生成器初始化失败")
	}
	public_context_prompt.ContextPrompt = prompt.NewContextPromptGenerator(public_contextManager.ContextManager)
	public.Generator = Generator
	fmt.Printf("prompt 初始化配置成功 \n")
	return nil
}
