package chat

import (
	"ai-roleplay/services/chat/api/internal/converter"
	"ai-roleplay/services/chat/api/internal/repo"
	"ai-roleplay/services/chat/api/internal/svc"
	"ai-roleplay/services/chat/api/internal/types"
	"context"

	"github.com/zeromicro/go-zero/core/logx"
)

type BatchDeleteConversationsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewBatchDeleteConversationsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BatchDeleteConversationsLogic {
	return &BatchDeleteConversationsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *BatchDeleteConversationsLogic) BatchDeleteConversations(req *types.BatchDeleteRequest) (resp *types.BaseResponse, err error) {
	// 参数验证
	if len(req.ConversationIDs) == 0 {
		return &types.BaseResponse{
			Code: 400,
			Msg:  "请选择要删除的对话",
		}, nil
	}

	// 验证ID有效性
	for _, id := range req.ConversationIDs {
		if id <= 0 {
			return &types.BaseResponse{
				Code: 400,
				Msg:  "对话ID无效",
			}, nil
		}
	}

	// 创建repo实例
	chatRepo := repo.NewChatServiceRepo(l.ctx, l.svcCtx)

	// 批量删除对话
	if err := chatRepo.BatchDeleteConversations(req.ConversationIDs); err != nil {
		l.Logger.Error("BatchDeleteConversations failed: ", err)
		return &types.BaseResponse{
			Code: 500,
			Msg:  "批量删除对话失败",
		}, nil
	}

	// 构建响应
	converter := converter.NewChatConverter()
	resp = converter.BuildBaseResponse(0, "批量删除成功")

	return resp, nil
}
