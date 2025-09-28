package chat

import (
	"context"
	"time"

	"ai-roleplay/services/chat/api/internal/repo"
	"ai-roleplay/services/chat/api/internal/svc"
	"ai-roleplay/services/chat/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ChatHistoryBeforeLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 侧边栏历史
func NewChatHistoryBeforeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ChatHistoryBeforeLogic {
	return &ChatHistoryBeforeLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ChatHistoryBeforeLogic) ChatHistoryBefore(req *types.ChatBeforeRequest) (resp *types.ChatBeforeResponse, err error) {

	resp = &types.ChatBeforeResponse{
		Todays:     make([]types.HistoryItem, 0),
		Yesterdays: make([]types.HistoryItem, 0),
		Befores:    make([]types.HistoryItem, 0),
	}

	chatRepo := repo.NewChatServiceRepo(l.ctx, l.svcCtx)
	conversation, _, err := chatRepo.GetConversationList(&types.ConversationListRequest{
		UserID:   req.UserId,
		Page:     1,
		PageSize: 10,
	})

	for _, conversation := range conversation {

		res := types.HistoryItem{
			ConversationID: conversation.ID,
			CharacterID:    conversation.CharacterID,
			CreatedAt:      conversation.CreatedAt.Format("2006-01-02 15:04:05"),
			CharacterName:  conversation.Title,
		}

		now := time.Now().Local()
		today := now.Format("2006-01-02")
		yesterday := now.AddDate(0, 0, -1).Format("2006-01-02")
		// 也将数据库中的时间转换为本地时区并格式化为日期字符串
		conversationDate := conversation.CreatedAt.Local().Format("2006-01-02")
		// 然后比较日期字符串
		switch conversationDate {
		case today:
			resp.Todays = append(resp.Todays, res)
		case yesterday:
			resp.Yesterdays = append(resp.Yesterdays, res)
		default:
			resp.Befores = append(resp.Befores, res)
		}
	}

	if err != nil {
		return nil, err
	}

	return
}
