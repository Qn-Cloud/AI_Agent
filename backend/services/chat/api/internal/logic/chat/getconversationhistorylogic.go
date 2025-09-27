package chat

import (
	"context"

	"ai-roleplay/services/chat/api/internal/repo"
	"ai-roleplay/services/chat/api/internal/svc"
	"ai-roleplay/services/chat/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetConversationHistoryLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取对话历史
func NewGetConversationHistoryLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetConversationHistoryLogic {
	return &GetConversationHistoryLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetConversationHistoryLogic) GetConversationHistory(req *types.GetConversationHistoryRequest) (resp *types.GetConversationHistoryResponse, err error) {
	// 创建repo实例
	chatRepo := repo.NewChatServiceRepo(l.ctx, l.svcCtx)

	// 获取对话历史
	conversations, total, err := chatRepo.GetConversationsByUserID(req.UserID, req.Page, req.PageSize)
	if err != nil {
		l.Logger.Error("GetConversationHistory failed: ", err)
		return nil, err
	}
	result := []types.ConversationHistoryItem{}
	var messageCount int
	characterSet := map[int64]bool{}
	for _, conversation := range conversations {
		item := types.ConversationHistoryItem{
			ConversationID: conversation.ID,
		}
		messages, err := chatRepo.GetMessageByConversationID(conversation.ID)
		if err != nil {
			l.Logger.Error("GetMessageByConversationID failed: ", err)
			return nil, err
		}
		item.MessageCount = int64(len(messages))
		item.LastMessageTime = messages[len(messages)-1].CreatedAt.Format("2006-01-02 15:04:05")
		// lastContent
		lastMessage, err := chatRepo.GetConversationLastMessage(conversation.ID)
		if err != nil {
			l.Logger.Error("GetConversationLastMessage failed: ", err)
			return nil, err
		}
		item.LastMessageContent = lastMessage.Content
		// ConversationDuration
		duration, err := chatRepo.GetConversationDuration(conversation.ID)
		item.ConversationDuration = duration

		// messageCount
		messageCount += int(item.MessageCount)

		characterSet[conversation.CharacterID] = true
		result = append(result, item)
	}
	var activeDays int
	if len(conversations) > 0 {
		firstConversationTime := conversations[0].UpdatedAt
		lastConversationTime := conversations[len(conversations)-1].UpdatedAt
		activeDays = int(lastConversationTime.Sub(firstConversationTime).Hours() / 24)
	}

	return &types.GetConversationHistoryResponse{
		List:              result,
		ConversationTotal: int(total),
		MessageCount:      messageCount,
		ActiveDays:        activeDays,
		CharacterCount:    int(len(characterSet)),
	}, nil
}
