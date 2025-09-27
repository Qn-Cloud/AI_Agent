package chat

import (
	"ai-roleplay/services/chat/api/internal/converter"
	"ai-roleplay/services/chat/api/internal/repo"
	"ai-roleplay/services/chat/api/internal/svc"
	"ai-roleplay/services/chat/api/internal/types"
	"context"

	"github.com/zeromicro/go-zero/core/logx"
)

type ExportConversationLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

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
			Msg:  "对话ID无效",
		}, nil
	}

	// 创建repo实例
	chatRepo := repo.NewChatServiceRepo(l.ctx, l.svcCtx)

	// 获取对话信息和消息记录
	conversation, messages, err := chatRepo.ExportConversation(req.ID)
	if err != nil {
		l.Logger.Error("ExportConversation failed: ", err)
		return &types.ExportResponse{
			Code: 500,
			Msg:  "获取对话数据失败",
		}, nil
	}

	if conversation == nil {
		return &types.ExportResponse{
			Code: 404,
			Msg:  "对话不存在",
		}, nil
	}

	// 使用转换器生成导出内容
	converter := converter.NewChatConverter()
	resp = converter.BuildExportResponse(conversation, messages)

	return resp, nil
}
