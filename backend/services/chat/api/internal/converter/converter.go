package converter

import (
	"ai-roleplay/services/chat/api/internal/types"
	"ai-roleplay/services/chat/model"
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
	content.WriteString(fmt.Sprintf("对话标题: %s\n", conversation.Title))
	content.WriteString(fmt.Sprintf("创建时间: %s\n", conversation.CreatedAt.Format("2006-01-02 15:04:05")))
	content.WriteString(fmt.Sprintf("更新时间: %s\n", conversation.UpdatedAt.Format("2006-01-02 15:04:05")))
	content.WriteString("=" + strings.Repeat("=", 50) + "\n\n")

	for _, message := range messages {
		var sender string
		if message.Type == "user" {
			sender = "用户"
		} else {
			sender = "AI"
		}

		content.WriteString(fmt.Sprintf("[%s] %s:\n%s\n\n",
			message.CreatedAt.Format("15:04:05"),
			sender,
			message.Content))
	}

	return &types.ExportResponse{
		Code:     0,
		Msg:      "导出成功",
		Data:     content.String(),
		Format:   "txt",
		Filename: fmt.Sprintf("conversation_%d_%s.txt", conversation.ID, time.Now().Format("20060102_150405")),
	}
}

// FromSendMessageRequest 从发送消息请求创建消息模型
func (c *ChatConverter) FromSendMessageRequest(req *types.SendMessageRequest) *model.Message {
	message := &model.Message{
		ConversationID: req.ConversationID,
		Type:           "user", // 发送的消息都是用户消息
		Content:        req.Content,
		Status:         1, // 正常状态
		CreatedAt:      time.Now(),
	}

	if req.CharacterID > 0 {
		message.CharacterID = &req.CharacterID
	}

	return message
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
