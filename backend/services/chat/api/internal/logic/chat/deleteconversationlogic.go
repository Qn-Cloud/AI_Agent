package chat

import (
	"context"

	"ai-roleplay/services/chat/api/internal/repo"
	"ai-roleplay/services/chat/api/internal/svc"
	"ai-roleplay/services/chat/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteConversationLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 删除对话
func NewDeleteConversationLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteConversationLogic {
	return &DeleteConversationLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteConversationLogic) DeleteConversation(req *types.ConversationRequest) (resp *types.BaseResponse, err error) {
	// 参数验证
	if req.ID <= 0 {
		return &types.BaseResponse{
			Code: 400,
			Msg:  "ID不能为空",
		}, nil
	}

	// 创建repo实例
	chatRepo := repo.NewChatServiceRepo(l.ctx, l.svcCtx)

	// 删除对话
	if err := chatRepo.DeleteConversation(req.ID); err != nil {
		l.Logger.Error("DeleteConversation failed: ", err)
		return &types.BaseResponse{
			Code: 500,
			Msg:  "删除对话失败",
		}, nil
	}

	return &types.BaseResponse{
		Code: 200,
		Msg:  "删除对话成功",
	}, nil

	return
}
