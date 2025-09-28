package model

import (
	"context"
	"log"

	"github.com/cloudwego/eino-ext/components/model/openai"
	"github.com/cloudwego/eino/components/model"
)

var (
	DeepSeek_Api_Key = "sk-636bd4e9653f45ba852ff77d8d3502f7"
	DeepSeek_Url     = "https://api.deepseek.com"
	DeepSeek_Model   = "deepseek-chat"
	//DeepSeekBearApiKey = "Bearer sk-0c9c8a4a2cc14304b1fc47c63301101c"
)

func CreateDeepSeekChatModel(ctx context.Context) model.ToolCallingChatModel {
	chatModel, err := openai.NewChatModel(ctx, &openai.ChatModelConfig{
		BaseURL: DeepSeek_Url,
		Model:   DeepSeek_Model,
		APIKey:  DeepSeek_Api_Key,
	})
	if err != nil {
		log.Fatalf("create openai chat model failed, err=%v", err)
	}
	return chatModel
}
