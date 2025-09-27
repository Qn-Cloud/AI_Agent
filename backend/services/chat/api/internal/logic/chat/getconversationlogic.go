package chat

import (
	"ai-roleplay/services/chat/api/internal/converter"
	"ai-roleplay/services/chat/api/internal/repo"
	"ai-roleplay/services/chat/api/internal/svc"
	"ai-roleplay/services/chat/api/internal/types"
	"context"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetConversationLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetConversationLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetConversationLogic {
	return &GetConversationLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetConversationLogic) GetConversation(req *types.ConversationRequest) (resp *types.ConversationResponse, err error) {
	// 参数验证
	if req.ID <= 0 {
		return &types.ConversationResponse{
			Code: 400,
			Msg:  "对话ID无效",
		}, nil
	}

	// 创建repo实例
	chatRepo := repo.NewChatServiceRepo(l.ctx, l.svcCtx)

	// 获取对话详情
	conversation, err := chatRepo.GetConversationByID(req.ID)
	if err != nil {
		l.Logger.Error("GetConversationByID failed: ", err)
		return &types.ConversationResponse{
			Code: 500,
			Msg:  "获取对话详情失败",
		}, nil
	}

	if conversation == nil {
		return &types.ConversationResponse{
			Code: 404,
			Msg:  "对话不存在",
		}, nil
	}

	// 转换数据格式
	converter := converter.NewChatConverter()
	resp = converter.BuildConversationResponse(conversation)

	return resp, nil
}
