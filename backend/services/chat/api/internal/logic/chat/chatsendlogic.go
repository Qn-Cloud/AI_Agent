package chat

import (
	"context"
	"fmt"
	"io"
	"strings"
	"time"

	common "ai-roleplay/common/utils"
	"ai-roleplay/services/chat/api/internal/converter"
	"ai-roleplay/services/chat/api/internal/prompt"
	"ai-roleplay/services/chat/api/internal/repo"
	"ai-roleplay/services/chat/api/internal/svc"
	"ai-roleplay/services/chat/api/internal/types"
	"ai-roleplay/services/chat/model"

	llm_model "ai-roleplay/services/chat/api/internal/model"

	"github.com/cloudwego/eino/schema"
	"github.com/zeromicro/go-zero/core/logx"
)

type ChatSendLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 发送消息并获取SSE流式响应
func NewChatSendLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ChatSendLogic {
	return &ChatSendLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ChatSendLogic) ChatSend(req *types.ChatSendRequest) error {
	// todo: add your logic here and delete this line

	return nil
}

func (l *ChatSendLogic) Sse(req *types.ChatSendRequest, client chan<- *types.ChatSSEEvent) error {

	userId := int64(1)

	// 记录请求日志
	l.Infof("User %d starting SSE chat - ConversationId: %d, ContentLength: %d",
		userId, req.ConversationId, len(req.Content))

	chatRepo := repo.NewChatServiceRepo(l.ctx, l.svcCtx)

	// 1、处理会话ID
	conversationId := req.ConversationId
	if req.ConversationId == 0 {
		var err error
		converter := converter.NewChatConverter()
		conversation := converter.FromCreateConversationRequest(&types.CreateConversationRequest{
			CharacterID: req.CharacterID,
			Title:       "新对话",
		})

		err = chatRepo.CreateConversation(conversation)
		if err != nil {
			l.sendError(client, fmt.Sprintf("创建对话失败: %v", err))
			return err
		}
	}

	// 2、保存用户消息
	_, err := chatRepo.AddMessage(&model.Message{
		ConversationID: conversationId,
		Content:        req.Content,
		Type:           common.AI_Role_User,
	})
	if err != nil {
		l.sendError(client, fmt.Sprintf("保存用户消息失败: %v", err))
		return err
	}

	// 3、发送思考状态
	l.sendEvent(client, &types.ChatSSEEvent{
		Type:           common.AI_SSE_Event_Thinking,
		Content:        req.Content,
		ConversationID: conversationId,
	})

	// 4、获取对话历史
	chatHistory, err := l.getChatHistory(conversationId)
	if err != nil {
		l.sendError(client, fmt.Sprintf("获取对话历史失败: %v", err))
		return err
	}

	// 5、调用LLM流式生成
	return l.streamCallModelWithChannel(client, req, chatHistory, conversationId, userId)
}

func (l *ChatSendLogic) getChatHistory(conversation_id int64) ([]*schema.Message, error) {
	chatRepo := repo.NewChatServiceRepo(l.ctx, l.svcCtx)
	chatHistory, err := chatRepo.GetConversationMessages(conversation_id)
	if err != nil {
		return nil, err
	}
	var messages []*schema.Message
	for _, message := range chatHistory {
		role := convertRoleForLLM(message.Type)
		msg := &schema.Message{
			Role:    schema.RoleType(role),
			Content: message.Content,
		}
		messages = append(messages, msg)
	}
	return messages, nil
}
func convertRoleForLLM(role string) string {
	switch role {
	case "ai":
		return "assistant" // 将 "ai" 转换为 "assistant"
	case "user":
		return "user" // 保持不变
	case "system":
		return "system" // 保持不变
	default:
		return "assistant" // 默认为 assistant
	}
}

func (l *ChatSendLogic) sendEvent(client chan<- *types.ChatSSEEvent, resp *types.ChatSSEEvent) {
	select {
	case client <- resp:
	case <-l.ctx.Done():
		l.Info("Context cancelled, stopping event send")
	default:
		l.Error("Channel full, dropping event")
	}
}

// 通过Channel发送错误
func (l *ChatSendLogic) sendError(client chan<- *types.ChatSSEEvent, errMsg string) {
	l.sendEvent(client, &types.ChatSSEEvent{
		Type:  common.AI_SSE_Event_Error,
		Error: errMsg,
	})
}

func (l *ChatSendLogic) streamCallModelWithChannel(client chan<- *types.ChatSSEEvent, req *types.ChatSendRequest, chatHistory []*schema.Message, conversationId int64, userId int64) error {
	// 设置超时
	ctx, cancel := context.WithTimeout(l.ctx, 60*time.Second)
	defer cancel()

	// 创建模型

	chatModel := llm_model.CreateDeepSeekChatModel(ctx)
	promptMsg := prompt.CreateMessageFromTemplate(req.Content, chatHistory)

	// 开始流式生成
	l.Info("Starting LLM stream generation")
	streamReader := prompt.GenerateStream(ctx, chatModel, promptMsg)
	defer streamReader.Close()

	var fullContent strings.Builder

	for {
		select {
		case <-ctx.Done():
			l.sendError(client, "请求超时")
			return ctx.Err()
		default:
		}

		recv, err := streamReader.Recv()
		if err == io.EOF {
			l.Info("LLM stream completed")
			break
		}
		if err != nil {
			l.Errorf("LLM stream error: %v", err)
			l.sendError(client, fmt.Sprintf("接收流式数据失败: %v", err))
			return err
		}

		if recv.Content != "" {
			fullContent.WriteString(recv.Content)

			// 发送增量内容
			l.sendEvent(client, &types.ChatSSEEvent{
				Type:           common.AI_SSE_Event_Message,
				Delta:          recv.Content,
				Content:        fullContent.String(),
				ConversationID: conversationId,
			})
		}
	}

	finalContent := fullContent.String()
	l.Infof("Final content length: %d", len(finalContent))

	// 保存AI回复到数据库
	msgService := repo.NewChatServiceRepo(l.ctx, l.svcCtx)
	msgId, err := msgService.AddMessage(&model.Message{
		ConversationID: conversationId,
		Type:           common.AI_Role_Assistant,
		Content:        finalContent,
	})
	if err != nil {
		l.sendError(client, fmt.Sprintf("保存AI消息失败: %v", err))
		return err
	}

	// 发送完成事件
	l.sendEvent(client, &types.ChatSSEEvent{
		Type:           common.AI_SSE_Event_Done,
		Done:           true,
		MessageId:      msgId,
		Content:        finalContent,
		ConversationID: conversationId,
	})

	l.Info("SSE stream completed successfully")
	return nil
}
