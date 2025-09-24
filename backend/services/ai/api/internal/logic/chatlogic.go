package logic

import (
	"context"
	"fmt"
	"strings"
	"time"

	"ai-roleplay/common/response"
	"ai-roleplay/services/ai/api/internal/svc"
	"ai-roleplay/services/ai/api/internal/types"

	"github.com/sashabaranov/go-openai"
	"github.com/zeromicro/go-zero/core/logx"
)

type ChatLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewChatLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ChatLogic {
	return &ChatLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ChatLogic) Chat(req *types.ChatRequest) (resp *types.ChatResponse, err error) {
	// 1. 参数验证
	if err := l.validateChatRequest(req); err != nil {
		return &types.ChatResponse{
			Code: response.INVALID_PARAMS,
			Msg:  err.Error(),
		}, nil
	}

	// 2. 内容安全检查
	if err := l.contentSafetyCheck(req.UserInput); err != nil {
		return &types.ChatResponse{
			Code: response.CONTENT_UNSAFE,
			Msg:  "输入内容包含不当信息",
		}, nil
	}

	// 3. 构建对话上下文
	messages := l.buildChatContext(req.CharacterPrompt, req.Context, req.UserInput)

	// 4. 检查请求频率限制
	if err := l.checkRateLimit(req.UserId); err != nil {
		return &types.ChatResponse{
			Code: response.RATE_LIMIT_EXCEEDED,
			Msg:  "请求过于频繁，请稍后再试",
		}, nil
	}

	// 5. 调用大语言模型
	aiResponse, err := l.callLLM(messages, req.MaxTokens, req.Temperature)
	if err != nil {
		logx.Errorf("调用LLM失败: %v", err)
		return &types.ChatResponse{
			Code: response.AI_SERVICE_ERROR,
			Msg:  "AI服务暂时不可用",
		}, nil
	}

	// 6. 后处理AI回复
	processedResponse := l.postProcessResponse(aiResponse)

	// 7. 再次进行内容安全检查
	if err := l.contentSafetyCheck(processedResponse); err != nil {
		logx.Warnf("AI生成内容触发安全检查: %s", processedResponse)
		processedResponse = "抱歉，我无法回答这个问题。让我们聊聊其他话题吧。"
	}

	// 8. 记录对话统计
	go l.recordChatMetrics(req.ConversationId, req.CharacterId, req.UserInput, processedResponse)

	// 9. 更新缓存
	go l.updateConversationCache(req.ConversationId, req.UserInput, processedResponse)

	return &types.ChatResponse{
		Code: response.SUCCESS,
		Msg:  "生成成功",
		Data: &types.ChatData{
			Content:      processedResponse,
			TokenUsed:    l.estimateTokenUsage(req.UserInput, processedResponse),
			ResponseTime: time.Since(time.Now()).Milliseconds(),
			MessageId:    l.generateMessageId(),
		},
	}, nil
}

// 验证聊天请求
func (l *ChatLogic) validateChatRequest(req *types.ChatRequest) error {
	if req.ConversationId <= 0 {
		return fmt.Errorf("对话ID无效")
	}

	if req.CharacterId <= 0 {
		return fmt.Errorf("角色ID无效")
	}

	if len(strings.TrimSpace(req.UserInput)) == 0 {
		return fmt.Errorf("用户输入不能为空")
	}

	if len(req.UserInput) > 2000 {
		return fmt.Errorf("输入内容过长，最多2000个字符")
	}

	if req.MaxTokens <= 0 || req.MaxTokens > 4000 {
		req.MaxTokens = 1000 // 默认值
	}

	if req.Temperature < 0 || req.Temperature > 2 {
		req.Temperature = 0.7 // 默认值
	}

	return nil
}

// 内容安全检查
func (l *ChatLogic) contentSafetyCheck(content string) error {
	// 1. 基础敏感词过滤
	sensitiveWords := l.svcCtx.Config.ContentSafety.SensitiveWords
	lowerContent := strings.ToLower(content)

	for _, word := range sensitiveWords {
		if strings.Contains(lowerContent, strings.ToLower(word)) {
			return fmt.Errorf("内容包含敏感词汇")
		}
	}

	// 2. 调用第三方内容审核服务（如阿里云、腾讯云）
	if l.svcCtx.Config.ContentSafety.Enable {
		passed, err := l.callContentModerationAPI(content)
		if err != nil {
			logx.Errorf("内容审核API调用失败: %v", err)
			// API失败时使用基础过滤，不阻断服务
		} else if !passed {
			return fmt.Errorf("内容未通过安全检查")
		}
	}

	return nil
}

// 构建对话上下文
func (l *ChatLogic) buildChatContext(characterPrompt string, context []*types.Message, userInput string) []openai.ChatCompletionMessage {
	messages := []openai.ChatCompletionMessage{
		{
			Role:    openai.ChatMessageRoleSystem,
			Content: characterPrompt,
		},
	}

	// 添加历史对话
	for _, msg := range context {
		var role string
		if msg.Type == "user" {
			role = openai.ChatMessageRoleUser
		} else {
			role = openai.ChatMessageRoleAssistant
		}

		messages = append(messages, openai.ChatCompletionMessage{
			Role:    role,
			Content: msg.Content,
		})
	}

	// 添加当前用户输入
	messages = append(messages, openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleUser,
		Content: userInput,
	})

	return messages
}

// 检查频率限制
func (l *ChatLogic) checkRateLimit(userId int64) error {
	// 检查用户请求频率
	key := fmt.Sprintf("rate_limit:user:%d", userId)

	// 使用Redis滑动窗口算法
	count, err := l.svcCtx.Redis.Incr(l.ctx, key).Result()
	if err != nil {
		logx.Errorf("Redis操作失败: %v", err)
		return nil // Redis失败时不限制
	}

	if count == 1 {
		// 设置过期时间（1分钟）
		l.svcCtx.Redis.Expire(l.ctx, key, time.Minute)
	}

	// 每分钟最多20次请求
	if count > 20 {
		return fmt.Errorf("请求频率超限")
	}

	return nil
}

// 调用大语言模型
func (l *ChatLogic) callLLM(messages []openai.ChatCompletionMessage, maxTokens int, temperature float64) (string, error) {
	client := openai.NewClient(l.svcCtx.Config.OpenAI.ApiKey)

	req := openai.ChatCompletionRequest{
		Model:       l.svcCtx.Config.OpenAI.Model,
		Messages:    messages,
		MaxTokens:   maxTokens,
		Temperature: float32(temperature),
		TopP:        1,
		Stream:      false,
	}

	// 设置超时
	ctx, cancel := context.WithTimeout(l.ctx, 30*time.Second)
	defer cancel()

	resp, err := client.CreateChatCompletion(ctx, req)
	if err != nil {
		// 处理不同类型的错误
		if strings.Contains(err.Error(), "rate limit") {
			return "", fmt.Errorf("AI服务请求频率超限")
		}
		if strings.Contains(err.Error(), "invalid_request_error") {
			return "", fmt.Errorf("请求参数错误")
		}
		if strings.Contains(err.Error(), "context_length_exceeded") {
			return "", fmt.Errorf("对话内容过长")
		}

		return "", fmt.Errorf("AI服务调用失败: %w", err)
	}

	if len(resp.Choices) == 0 {
		return "", fmt.Errorf("AI服务返回空响应")
	}

	return resp.Choices[0].Message.Content, nil
}

// 后处理AI回复
func (l *ChatLogic) postProcessResponse(response string) string {
	// 1. 去除首尾空白
	response = strings.TrimSpace(response)

	// 2. 处理特殊格式
	// 移除可能的markdown格式标记
	response = strings.ReplaceAll(response, "```", "")
	response = strings.ReplaceAll(response, "**", "")

	// 3. 长度限制
	maxLength := 1500
	if len([]rune(response)) > maxLength {
		runes := []rune(response)
		response = string(runes[:maxLength]) + "..."
	}

	// 4. 确保回复符合角色设定
	// 这里可以添加更复杂的后处理逻辑

	return response
}

// 调用内容审核API
func (l *ChatLogic) callContentModerationAPI(content string) (bool, error) {
	// 这里可以集成阿里云、腾讯云等内容审核服务
	// 示例代码，实际需要根据具体服务实现

	// 模拟API调用
	time.Sleep(100 * time.Millisecond)

	// 简单规则：包含特定关键词返回false
	forbiddenKeywords := []string{"违法", "暴力", "政治敏感"}
	lowerContent := strings.ToLower(content)

	for _, keyword := range forbiddenKeywords {
		if strings.Contains(lowerContent, keyword) {
			return false, nil
		}
	}

	return true, nil
}

// 记录对话统计
func (l *ChatLogic) recordChatMetrics(conversationId, characterId int64, userInput, aiResponse string) {
	metrics := &model.ChatMetrics{
		ConversationId:   conversationId,
		CharacterId:      characterId,
		UserInputLength:  len(userInput),
		AiResponseLength: len(aiResponse),
		TokenUsed:        l.estimateTokenUsage(userInput, aiResponse),
		ResponseTime:     0, // 实际需要测量
		CreatedAt:        time.Now(),
	}

	if err := l.svcCtx.ChatMetricsModel.Insert(context.Background(), metrics); err != nil {
		logx.Errorf("记录对话统计失败: %v", err)
	}
}

// 更新对话缓存
func (l *ChatLogic) updateConversationCache(conversationId int64, userInput, aiResponse string) {
	// 更新对话上下文缓存
	cacheKey := fmt.Sprintf("conversation:context:%d", conversationId)

	// 获取现有缓存
	cachedMessages, err := l.svcCtx.CacheManager.GetConversationContext(context.Background(), conversationId)
	if err != nil {
		logx.Errorf("获取对话缓存失败: %v", err)
		return
	}

	// 添加新消息
	newMessages := append(cachedMessages,
		&model.Message{Type: "user", Content: userInput, CreatedAt: time.Now()},
		&model.Message{Type: "ai", Content: aiResponse, CreatedAt: time.Now()},
	)

	// 保持最近20条消息
	if len(newMessages) > 20 {
		newMessages = newMessages[len(newMessages)-20:]
	}

	// 更新缓存
	if err := l.svcCtx.CacheManager.SetConversationContext(context.Background(), conversationId, newMessages); err != nil {
		logx.Errorf("更新对话缓存失败: %v", err)
	}
}

// 估算Token使用量
func (l *ChatLogic) estimateTokenUsage(userInput, aiResponse string) int {
	// 简单估算：中文字符约等于1个token，英文单词约等于1个token
	totalChars := len([]rune(userInput + aiResponse))
	return int(float64(totalChars) * 1.2) // 加上一些开销
}

// 生成消息ID
func (l *ChatLogic) generateMessageId() string {
	return fmt.Sprintf("msg_%d_%d", time.Now().UnixNano(), l.ctx.Value("userId"))
}
