package public_brain

import (
	"log"

	"new-robot-self/internal/brain"
	"new-robot-self/internal/tools"
)

var BrainService *brain.Brain

func InitTools() {
	// 创建知识库
	knowledgeBase := brain.NewKnowledgeBase()

	// 注册所有工具
	if err := tools.RegisterAllTools(knowledgeBase); err != nil {
		log.Printf("注册工具失败: %v", err)
	} else {
		log.Println("工具注册成功")
	}

	// 设置health_manager_tool的brainService
	if healthTool, exists := tools.Get("health_manager"); exists {
		if hm, ok := healthTool.(*tools.HealthManagerTool); ok {
			hm.SetBrainService(BrainService)
			log.Println("health_manager_tool的brainService已设置")
		}
	}
}

