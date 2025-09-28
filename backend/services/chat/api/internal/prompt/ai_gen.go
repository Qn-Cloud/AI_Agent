package prompt

import (
	"context"
	"fmt"
	"io"
	"log"
	"strings"

	"github.com/cloudwego/eino/components/model"
	"github.com/cloudwego/eino/components/prompt"
	"github.com/cloudwego/eino/schema"
)

func Generate(ctx context.Context, llm model.ToolCallingChatModel, in []*schema.Message) *schema.Message {
	result, err := llm.Generate(ctx, in)
	if err != nil {
		log.Fatalf("llm generate failed: %v", err)
	}
	return result
}

func GenerateStream(ctx context.Context, llm model.ToolCallingChatModel, in []*schema.Message) *schema.StreamReader[*schema.Message] {
	result, err := llm.Stream(ctx, in)
	if err != nil {
		log.Fatalf("llm generate failed: %v", err)
	}
	return result
}

func ReportStream(sr *schema.StreamReader[*schema.Message]) {
	defer sr.Close()
	var contentBuffer strings.Builder
	for {
		recv, err := sr.Recv()
		if err == io.EOF { // 流式输出结束
			return
		}
		if err != nil {
			log.Fatalf("recv failed : %v", err)
		}
		if recv.Content != "" {
			contentBuffer.WriteString(recv.Content)
		}
	}

}

func createTemplate() prompt.ChatTemplate {
	// 创建模板，使用 FString 格式
	return prompt.FromMessages(schema.FString,
		// 系统消息模板
		schema.SystemMessage("你是一个{role}。你需要用{style}的语气回答问题。你的目标是帮助程序员保持积极乐观的心态，提供技术建议的同时也要关注他们的心理健康。"),

		// 插入需要的对话历史（新对话的话这里不填）
		schema.MessagesPlaceholder("chat_history", true),

		// 用户消息模板
		schema.UserMessage("问题: {question}"),
	)
}

func CreateMessageFromTemplate(content string, chatHistory []*schema.Message) []*schema.Message {
	template := createTemplate()
	// 使用模板生成消息
	fmt.Println("history", chatHistory)
	// messages, err := template.Format(context.Background(), map[string]any{
	// 	"role":     "程序员鼓励师",
	// 	"style":    "积极、温暖且专业",
	// 	"question": "我的代码一直报错，感觉好沮丧，该怎么办？",
	// 	// "content":  content,
	// 	// 对话历史（这个例子里模拟两轮对话历史）
	// 	"chat_history": []*schema.Message{
	// 		schema.UserMessage("你好"),
	// 		schema.AssistantMessage("嘿！我是你的程序员鼓励师！记住，每个优秀的程序员都是从 Debug 中成长起来的。有什么我可以帮你的吗？", nil),
	// 		schema.UserMessage("我觉得自己写的代码太烂了"),
	// 		schema.AssistantMessage("每个程序员都经历过这个阶段！重要的是你在不断学习和进步。让我们一起看看代码，我相信通过重构和优化，它会变得更好。记住，Rome wasn't built in a day，代码质量是通过持续改进来提升的。", nil),
	// 	},
	// 	// "chat_history": chatHistory,
	// })
	messages, err := template.Format(context.Background(), map[string]any{
		"role":         "程序员鼓励师",
		"style":        "积极、温暖且专业",
		"question":     content,     // 使用用户输入内容
		"chat_history": chatHistory, // 使用实际对话历史
	})
	if err != nil {
		log.Fatalf("format template failed: %v\n", err)
	}
	return messages
}
