package chat

import (
	"context"

	"ai-roleplay/services/chat/api/internal/svc"
	"ai-roleplay/services/chat/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type SearchConversationsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 搜索对话
func NewSearchConversationsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SearchConversationsLogic {
	return &SearchConversationsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SearchConversationsLogic) SearchConversations(req *types.SearchConversationRequest) (resp *types.ConversationListResponse, err error) {
	// todo: add your logic here and delete this line

	return
}
