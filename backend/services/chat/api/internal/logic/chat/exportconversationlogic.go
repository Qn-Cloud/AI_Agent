package chat

import (
	"context"

	"ai-roleplay/services/chat/api/internal/repo"
	"ai-roleplay/services/chat/api/internal/svc"
	"ai-roleplay/services/chat/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ExportConversationLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 导出对话记录
func NewExportConversationLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ExportConversationLogic {
	return &ExportConversationLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ExportConversationLogic) ExportConversation(req *types.ConversationRequest) (resp *types.ExportResponse, err error) {
	// 参数验证
	if req.ID <= 0 {
		return &types.ExportResponse{
			Code: 400,
			Msg:  "ID不能为空",
		}, nil
	}

	// 创建repo实例
	chatRepo := repo.NewChatServiceRepo(l.ctx, l.svcCtx)

	// 导出对话
	if err := chatRepo.ExportConversation(req.ID); err != nil {
		l.Logger.Error("ExportConversation failed: ", err)
		return &types.ExportResponse{
			Code: 500,
			Msg:  "导出对话失败",
		}, nil
	}

	return &types.ExportResponse{
		Code: 200,
		Msg:  "导出对话成功",
	}, nil
}
