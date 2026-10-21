/*
 * Copyright 2024 CloudWeGo Authors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package main

import (
	"context"
	"github.com/cloudwego/eino-ext/components/model/openai"
	"github.com/cloudwego/eino/components/prompt"
	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/schema"
	"log/slog"
)

func main() {
	openAIAPIKey := "YOUR_API_KEY"
	openAIModelName := "qwen-plus"
	openAIBaseURL := "https://dashscope.aliyuncs.com/compatible-mode/v1"

	ctx := context.Background()

	// 初始化 tools
	//todoTools := []tool.BaseTool{
	//	&httpSearch.HttpWeatherTool{},
	//	&httpSearch.HttpSearchTool{},
	//	&httpSearch.HttpSearchEngineTool{},
	//}

	// 创建并配置 ChatModel
	chatModel, err := openai.NewChatModel(ctx, &openai.ChatModelConfig{
		BaseURL: openAIBaseURL,
		Model:   openAIModelName,
		APIKey:  openAIAPIKey,
	})
	if err != nil {
		slog.Error("NewChatModel failed, err=%v", err)
		return
	}
	//创建模版
	template := prompt.FromMessages(schema.FString,
		schema.SystemMessage("你是一个{role}"),
		schema.MessagesPlaceholder("history_msg", false),
		&schema.Message{
			Role:    schema.User,
			Content: "请帮我{task}",
		},
	)

	//准备变量
	variables := map[string]any{
		"role": "助手",
		"task": "你还知道我吗?",
		"history_msg": []*schema.Message{
			{
				Role:    schema.User,
				Content: "我叫超级无敌宇宙爆龙战士",
			},
			{
				Role:    schema.Assistant,
				Content: "好的,超级无敌宇宙爆龙战士!",
			},
			{
				Role:    schema.User,
				Content: "我喜欢打lol",
			},
			{
				Role:    schema.Assistant,
				Content: "超级无敌宇宙爆龙战士喜欢玩lol游戏!",
			},
			{
				Role:    schema.User,
				Content: "我今年27岁",
			},
			{
				Role:    schema.Assistant,
				Content: "超级无敌宇宙爆龙战士今年27岁了!",
			},
		},
	}

	//格式化模版
	result, err := template.Format(context.Background(), variables)
	if err != nil {
		slog.Error("Format failed, err=%v", err)
		return
	}
	msg, err := chatModel.Generate(context.Background(), result)

	if err != nil {
		slog.Error("Generate failed, err=%v", err)
		return
	}
	slog.Info("Generate result: %+v", msg)
	slog.Info("over")

	//// 获取工具信息, 用于绑定到 ChatModel
	//toolInfos := make([]*schema.ToolInfo, 0, len(todoTools))
	//var info *schema.ToolInfo
	//for _, todoTool := range todoTools {
	//	info, err = todoTool.Info(ctx)
	//	if err != nil {
	//		slog.Info("get ToolInfo failed, err=%v", err)
	//		return
	//	}
	//	toolInfos = append(toolInfos, info)
	//}
	//
	//// 将 tools 绑定到 ChatModel
	//err = chatModel.BindTools(toolInfos)
	//if err != nil {
	//	slog.Error("BindTools failed, err=%v", err)
	//	return
	//}
	//
	//// 创建 tools 节点
	//todoToolsNode, err := compose.NewToolNode(ctx, &compose.ToolsNodeConfig{
	//	Tools: todoTools,
	//})
	//if err != nil {
	//	slog.Error("NewToolNode failed, err=%v", err)
	//	return
	//}
	//
	//// 构建完整的处理链
	//chain := compose.NewChain[[]*schema.Message, []*schema.Message]()
	//chain.
	//	AppendChatModel(chatModel, compose.WithNodeName("chat_model")).
	//	AppendToolsNode(todoToolsNode, compose.WithNodeName("tools"))
	//
	//// 编译并运行 chain
	//agent, err := chain.Compile(ctx)
	//if err != nil {
	//	slog.Error("chain.Compile failed, err=%v", err)
	//	return
	//}
	//// 运行示例
	//resp, err := agent.Invoke(ctx, []*schema.Message{
	//	{
	//		Role:    schema.System,
	//		Content: "你是一个助手，可以回答用户的一般问题。当用户询问天气相关的问题时，请使用天气工具获取数据。当用户需要获取互联网上的信息、新闻、知识或事实时，请使用搜索引擎工具。对于其他类型的问题、闲聊或情感表达，请使用通用搜索工具。根据问题的性质选择最合适的工具。",
	//	},
	//	{
	//		Role:    schema.User,
	//		Content: fmt.Sprintf("我想了解一下国内的新闻!", str),
	//	},
	//})
	//if err != nil {
	//	slog.Error("agent.Invoke failed, err=%v", err)
	//	return
	//}
	//
	//// 输出结果
	//for idx, msg := range resp {
	//	slog.Info("\n")
	//	slog.Info("message %d: %s: %s", idx, msg.Role, msg.Content)
	//}
}

func composeChain() {
	compose.NewChain[map[string]any, []*schema.Message]()
}
