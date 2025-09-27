package chat

import (
	"context"

	"ai-roleplay/services/chat/api/internal/svc"
	"ai-roleplay/services/chat/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateConversationTitleLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 更新对话标题
func NewUpdateConversationTitleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateConversationTitleLogic {
	return &UpdateConversationTitleLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateConversationTitleLogic) UpdateConversationTitle(req *types.UpdateTitleRequest) (resp *types.BaseResponse, err error) {
	// todo: add your logic here and delete this line

	return
}
