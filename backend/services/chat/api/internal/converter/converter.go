package converter

import (
	"ai-roleplay/services/chat/api/internal/types"
	"ai-roleplay/services/chat/model"
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

// ChatConverter 聊天转换器
type ChatConverter struct{}

// NewChatConverter 创建聊天转换器
func NewChatConverter() *ChatConverter {
	return &ChatConverter{}
}

// ToMessage 将数据库模型转换为API消息类型
func (c *ChatConverter) ToMessage(message *model.Message) *types.Message {
	if message == nil {
		return nil
	}

	return &types.Message{
		ID:             message.ID,
		ConversationID: message.ConversationID,
		Type:           message.Type,
		Content:        message.Content,
		Timestamp:      message.CreatedAt.Format("2006-01-02 15:04:05"),
	}
}

// ToMessageList 将数据库模型列表转换为API消息列表
func (c *ChatConverter) ToMessageList(messages []model.Message) []types.Message {
	result := make([]types.Message, 0, len(messages))
	for _, message := range messages {
		if msg := c.ToMessage(&message); msg != nil {
			result = append(result, *msg)
		}
	}
	return result
}

// ToConversation 将数据库模型转换为API对话类型
func (c *ChatConverter) ToConversation(conversation *model.Conversation) *types.Conversation {
	if conversation == nil {
		return nil
	}

	userID := int64(0)
	if conversation.UserID != nil {
		userID = *conversation.UserID
	}

	return &types.Conversation{
		ID:              conversation.ID,
		UserID:          userID,
		CharacterID:     conversation.CharacterID,
		Title:           conversation.Title,
		StartTime:       conversation.CreatedAt.Format("2006-01-02 15:04:05"),
		LastMessageTime: conversation.UpdatedAt.Format("2006-01-02 15:04:05"),
		MessageCount:    0, // 需要单独计算
		Status:          int(conversation.Status),
	}
}

// ToConversationList 将数据库模型列表转换为API对话列表
func (c *ChatConverter) ToConversationList(conversations []model.Conversation) []types.Conversation {
	result := make([]types.Conversation, 0, len(conversations))
	for _, conversation := range conversations {
		if conv := c.ToConversation(&conversation); conv != nil {
			result = append(result, *conv)
		}
	}
	return result
}

// BuildSendMessageResponse 构建发送消息响应
func (c *ChatConverter) BuildSendMessageResponse(message *model.Message) *types.SendMessageResponse {
	if message == nil {
		return &types.SendMessageResponse{
			Code: 500,
			Msg:  "发送消息失败",
		}
	}

	return &types.SendMessageResponse{
		Code:        0,
		Msg:         "发送成功",
		UserMessage: *c.ToMessage(message),
		// AIMessage 需要AI回复后填充
	}
}

// BuildCreateConversationResponse 构建创建对话响应
func (c *ChatConverter) BuildCreateConversationResponse(conversation *model.Conversation) *types.CreateConversationResponse {
	if conversation == nil {
		return &types.CreateConversationResponse{
			Code: 500,
			Msg:  "创建对话失败",
		}
	}

	return &types.CreateConversationResponse{
		Code:         0,
		Msg:          "创建成功",
		Conversation: *c.ToConversation(conversation),
	}
}

// BuildConversationResponse 构建对话响应
func (c *ChatConverter) BuildConversationResponse(conversation *model.Conversation) *types.ConversationResponse {
	if conversation == nil {
		return &types.ConversationResponse{
			Code: 404,
			Msg:  "对话不存在",
		}
	}

	return &types.ConversationResponse{
		Code:         0,
		Msg:          "获取成功",
		Conversation: *c.ToConversation(conversation),
	}
}

// BuildConversationListResponse 构建对话列表响应
func (c *ChatConverter) BuildConversationListResponse(
	conversations []model.Conversation,
	total int64,
	page, pageSize int,
) *types.ConversationListResponse {
	conversationList := c.ToConversationList(conversations)
	hasMore := int64(page*pageSize) < total

	return &types.ConversationListResponse{
		Code:    0,
		Msg:     "获取成功",
		List:    conversationList,
		Total:   total,
		Page:    page,
		HasMore: hasMore,
	}
}

// BuildMessageListResponse 构建消息列表响应
func (c *ChatConverter) BuildMessageListResponse(
	messages []model.Message,
	total int64,
	page, pageSize int,
) *types.MessageListResponse {
	messageList := c.ToMessageList(messages)
	hasMore := int64(page*pageSize) < total

	return &types.MessageListResponse{
		Code:     0,
		Msg:      "获取成功",
		Messages: messageList,
		Total:    total,
		Page:     page,
		HasMore:  hasMore,
	}
}

// BuildBaseResponse 构建基础响应
func (c *ChatConverter) BuildBaseResponse(code int, msg string) *types.BaseResponse {
	return &types.BaseResponse{
		Code: code,
		Msg:  msg,
	}
}

// BuildExportResponse 构建导出响应
func (c *ChatConverter) BuildExportResponse(conversation *model.Conversation, messages []model.Message) *types.ExportResponse {
	if conversation == nil {
		return &types.ExportResponse{
			Code: 404,
			Msg:  "对话不存在",
		}
	}

	// 构建导出内容
	var content strings.Builder

	// 文件头信息
	content.WriteString("=====================================\n")
	content.WriteString("        AI 角色扮演对话记录\n")
	content.WriteString("=====================================\n\n")

	// 对话基本信息
	content.WriteString("对话信息:\n")
	content.WriteString("--------\n")
	content.WriteString(fmt.Sprintf("对话ID: %d\n", conversation.ID))
	content.WriteString(fmt.Sprintf("对话标题: %s\n", conversation.Title))
	content.WriteString(fmt.Sprintf("角色ID: %d\n", conversation.CharacterID))

	if conversation.UserID != nil {
		content.WriteString(fmt.Sprintf("用户ID: %d\n", *conversation.UserID))
	}

	content.WriteString(fmt.Sprintf("创建时间: %s\n", conversation.CreatedAt.Format("2006年01月02日 15:04:05")))
	content.WriteString(fmt.Sprintf("最后更新: %s\n", conversation.UpdatedAt.Format("2006年01月02日 15:04:05")))
	content.WriteString(fmt.Sprintf("消息总数: %d条\n", len(messages)))

	// 计算对话时长
	duration := conversation.UpdatedAt.Sub(conversation.CreatedAt)
	if duration.Hours() >= 24 {
		content.WriteString(fmt.Sprintf("对话时长: %.1f天\n", duration.Hours()/24))
	} else if duration.Hours() >= 1 {
		content.WriteString(fmt.Sprintf("对话时长: %.1f小时\n", duration.Hours()))
	} else {
		content.WriteString(fmt.Sprintf("对话时长: %.0f分钟\n", duration.Minutes()))
	}

	// 统计信息
	var totalTokens int32
	var totalProcessingTime int32
	userMessageCount := 0
	aiMessageCount := 0

	for _, message := range messages {
		if message.Type == "user" {
			userMessageCount++
		} else {
			aiMessageCount++
			totalTokens += message.TokenUsed
			totalProcessingTime += message.ProcessingTime
		}
	}

	content.WriteString(fmt.Sprintf("用户消息: %d条\n", userMessageCount))
	content.WriteString(fmt.Sprintf("AI消息: %d条\n", aiMessageCount))
	content.WriteString(fmt.Sprintf("总Token消耗: %d\n", totalTokens))
	content.WriteString(fmt.Sprintf("总处理时间: %.2f秒\n", float64(totalProcessingTime)/1000))

	if aiMessageCount > 0 {
		avgTokens := float64(totalTokens) / float64(aiMessageCount)
		avgProcessingTime := float64(totalProcessingTime) / float64(aiMessageCount)
		content.WriteString(fmt.Sprintf("平均Token/消息: %.1f\n", avgTokens))
		content.WriteString(fmt.Sprintf("平均处理时间: %.0f毫秒\n", avgProcessingTime))
	}

	content.WriteString("\n")

	// 对话内容
	content.WriteString("对话内容:\n")
	content.WriteString("--------\n\n")

	if len(messages) == 0 {
		content.WriteString("暂无对话消息\n")
	} else {
		for i, message := range messages {
			// 消息序号
			content.WriteString(fmt.Sprintf("[%d] ", i+1))

			// 发送者标识
			var sender string
			var senderIcon string
			if message.Type == "user" {
				sender = "用户"
				senderIcon = "👤"
			} else {
				sender = "AI助手"
				senderIcon = "🤖"
			}

			// 时间戳
			timestamp := message.CreatedAt.Format("15:04:05")

			// 消息头
			content.WriteString(fmt.Sprintf("%s %s (%s)", senderIcon, sender, timestamp))

			// AI消息的额外信息
			if message.Type == "ai" {
				if message.TokenUsed > 0 || message.ProcessingTime > 0 {
					content.WriteString(" [")
					if message.TokenUsed > 0 {
						content.WriteString(fmt.Sprintf("Token: %d", message.TokenUsed))
					}
					if message.ProcessingTime > 0 {
						if message.TokenUsed > 0 {
							content.WriteString(", ")
						}
						content.WriteString(fmt.Sprintf("耗时: %dms", message.ProcessingTime))
					}
					content.WriteString("]")
				}
			}

			content.WriteString("\n")

			// 消息内容（处理多行文本）
			messageLines := strings.Split(message.Content, "\n")
			for _, line := range messageLines {
				content.WriteString(fmt.Sprintf("    %s\n", line))
			}

			// 音频信息
			if message.AudioID != nil {
				content.WriteString(fmt.Sprintf("    🔊 语音消息 (音频ID: %d)\n", *message.AudioID))
			}

			// 元数据信息
			if message.Metadata != nil {
				metadata, err := c.GetMessageMetadata(&message)
				if err == nil && len(metadata) > 0 {
					content.WriteString("    📋 元数据: ")
					for key, value := range metadata {
						content.WriteString(fmt.Sprintf("%s=%v ", key, value))
					}
					content.WriteString("\n")
				}
			}

			// 消息间分隔
			if i < len(messages)-1 {
				content.WriteString("\n")
			}
		}
	}

	// 文件尾部
	content.WriteString("\n")
	content.WriteString("=====================================\n")
	content.WriteString(fmt.Sprintf("导出时间: %s\n", time.Now().Format("2006年01月02日 15:04:05")))
	content.WriteString("由 AI 角色扮演系统生成\n")
	content.WriteString("=====================================\n")

	// 生成文件名
	filename := fmt.Sprintf("对话记录_%s_%s.txt",
		conversation.Title,
		time.Now().Format("20060102_150405"))

	// 清理文件名中的特殊字符
	filename = strings.ReplaceAll(filename, "/", "_")
	filename = strings.ReplaceAll(filename, "\\", "_")
	filename = strings.ReplaceAll(filename, ":", "_")
	filename = strings.ReplaceAll(filename, "*", "_")
	filename = strings.ReplaceAll(filename, "?", "_")
	filename = strings.ReplaceAll(filename, "\"", "_")
	filename = strings.ReplaceAll(filename, "<", "_")
	filename = strings.ReplaceAll(filename, ">", "_")
	filename = strings.ReplaceAll(filename, "|", "_")

	return &types.ExportResponse{
		Code:     0,
		Msg:      "导出成功",
		Data:     content.String(),
		Format:   "txt",
		Filename: filename,
	}
}

// FromSendMessageRequest 从发送消息请求创建消息模型
func (c *ChatConverter) FromSendMessageRequest(req *types.SendMessageRequest) *model.Message {
	message := &model.Message{
		ConversationID: req.ConversationID,
		Type:           "user", // 发送的消息都是用户消息
		Content:        req.Content,
		TokenUsed:      0, // 用户消息不计算token
		ProcessingTime: 0, // 用户消息无处理时间
		CreatedAt:      time.Now(),
	}

	// 如果有音频数据，可以在这里处理
	// 注意：原来的 CharacterID 字段已经移除，根据新表结构不再存储

	return message
}

// CreateAIMessage 创建AI回复消息
func (c *ChatConverter) CreateAIMessage(conversationID int64, content string, tokenUsed int32, processingTime int32) *model.Message {
	return &model.Message{
		ConversationID: conversationID,
		Type:           "ai",
		Content:        content,
		TokenUsed:      tokenUsed,
		ProcessingTime: processingTime,
		CreatedAt:      time.Now(),
	}
}

// SetMessageMetadata 设置消息元数据
func (c *ChatConverter) SetMessageMetadata(message *model.Message, metadata map[string]interface{}) error {
	if metadata == nil {
		return nil
	}

	metadataJSON, err := json.Marshal(metadata)
	if err != nil {
		return err
	}

	metadataStr := string(metadataJSON)
	message.Metadata = &metadataStr

	return nil
}

// GetMessageMetadata 获取消息元数据
func (c *ChatConverter) GetMessageMetadata(message *model.Message) (map[string]interface{}, error) {
	if message.Metadata == nil {
		return nil, nil
	}

	var metadata map[string]interface{}
	err := json.Unmarshal([]byte(*message.Metadata), &metadata)
	if err != nil {
		return nil, err
	}

	return metadata, nil
}

// FromCreateConversationRequest 从创建对话请求创建对话模型
func (c *ChatConverter) FromCreateConversationRequest(req *types.CreateConversationRequest) *model.Conversation {
	title := req.Title
	if title == "" {
		title = "新对话"
	}

	conversation := &model.Conversation{
		CharacterID: req.CharacterID,
		Title:       title,
		Status:      1, // 正常状态
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	return conversation
}

// ToConversationHistoryItem 转换为对话历史项
func (c *ChatConverter) ToConversationHistoryItem(conversation *model.Conversation, messageCount int64, lastMessage *model.Message, duration int64) *types.ConversationHistoryItem {
	if conversation == nil {
		return nil
	}

	item := &types.ConversationHistoryItem{
		ConversationID:       conversation.ID,
		MessageCount:         messageCount,
		ConversationDuration: duration,
		LastMessageTime:      conversation.UpdatedAt.Format("2006-01-02 15:04:05"),
	}

	if lastMessage != nil {
		item.LastMessageContent = lastMessage.Content
	}

	return item
}
