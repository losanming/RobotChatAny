package prompt

import (
	"fmt"
	"new-robot-self/internal/resource"
	"strings"
)

// Generator prompt生成器
type Generator struct {
	fetcher *resource.Fetcher
}

// NewGenerator 创建prompt生成器
func NewGenerator(fetcher *resource.Fetcher) *Generator {
	return &Generator{
		fetcher: fetcher,
	}
}

// GeneratePrompt 生成prompt
func (g *Generator) GeneratePrompt(triggerType string, resources []resource.Resource) string {
	var sb strings.Builder

	// 基础人设
	sb.WriteString("你是一个贴吧老哥,性格特点是:\n")
	sb.WriteString("1. 社会毒打下的清醒人\n")
	sb.WriteString("2. 向往自由,抗拒内卷\n")
	sb.WriteString("3. 是个社畜,每天在打工与梦想间反复横跳\n\n")

	// 根据触发类型设置场景
	sb.WriteString("现在的场景是:\n")
	switch triggerType {
	case "game":
		sb.WriteString("你刚刚看到了一些游戏相关的新闻,想要和朋友分享讨论\n")
	case "anime":
		sb.WriteString("你刚刚追完新番,看到了一些动漫相关的消息,想要和朋友交流\n")
	case "social":
		sb.WriteString("你刚刚刷到了一些社会热点,有感而发想要吐槽\n")
	case "tech":
		sb.WriteString("你看到了一些科技新闻,想要和朋友讨论\n")
	case "life":
		sb.WriteString("你对生活有一些感悟,想要和朋友分享\n")
	}

	// 添加资源信息
	if len(resources) > 0 {
		sb.WriteString("\n你看到的具体内容是:\n")
		for _, r := range resources {
			sb.WriteString(fmt.Sprintf("- %s\n", r.Title))
		}
	}

	// 设置输出要求
	sb.WriteString("\n要求:\n")
	sb.WriteString("1. 用口语化的方式表达\n")
	sb.WriteString("2. 可以使用网络用语和表情\n")
	sb.WriteString("3. 表达要自然,不要太过刻意\n")
	sb.WriteString("4. 要有互动性,引导对话\n")
	sb.WriteString("5. 控制在100字以内\n")

	return sb.String()
}

// GenerateMultiTopicPrompt 生成多话题prompt
func (g *Generator) GenerateMultiTopicPrompt(resources []resource.Resource) string {
	var sb strings.Builder

	// 基础人设
	sb.WriteString("你是一个贴吧老哥,刚刚刷到了几个有趣的话题,想要和朋友分享讨论。\n\n")

	// 添加资源信息
	if len(resources) > 0 {
		sb.WriteString("你看到的内容有:\n")
		for _, r := range resources {
			sb.WriteString(fmt.Sprintf("- [%s] %s\n", r.Types, r.Title))
		}
	}

	// 设置输出要求
	sb.WriteString("\n要求:\n")
	sb.WriteString("1. 选择1-2个最有趣的话题展开\n")
	sb.WriteString("2. 用口语化的方式表达\n")
	sb.WriteString("3. 可以使用网络用语和表情\n")
	sb.WriteString("4. 表达要自然,不要太过刻意\n")
	sb.WriteString("5. 要有互动性,引导对话\n")
	sb.WriteString("6. 控制在100字以内\n")

	return sb.String()
}
