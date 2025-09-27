package chat

import (
	"context"

	"ai-roleplay/services/chat/api/internal/svc"
	"ai-roleplay/services/chat/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ClearMessagesLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 清空对话消息
func NewClearMessagesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ClearMessagesLogic {
	return &ClearMessagesLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ClearMessagesLogic) ClearMessages() (resp *types.BaseResponse, err error) {
	// todo: add your logic here and delete this line

	return
}
