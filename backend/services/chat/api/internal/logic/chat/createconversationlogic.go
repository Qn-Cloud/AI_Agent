package chat

import (
	"ai-roleplay/services/chat/api/internal/converter"
	"ai-roleplay/services/chat/api/internal/repo"
	"ai-roleplay/services/chat/api/internal/svc"
	"ai-roleplay/services/chat/api/internal/types"
	"context"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateConversationLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateConversationLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateConversationLogic {
	return &CreateConversationLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateConversationLogic) CreateConversation(req *types.CreateConversationRequest) (resp *types.CreateConversationResponse, err error) {
	// 参数验证
	if req.CharacterID <= 0 {
		return &types.CreateConversationResponse{
			Code: 400,
			Msg:  "角色ID无效",
		}, nil
	}

	if req.Title == "" {
		req.Title = "新对话" // 设置默认标题
	}

	// 创建repo实例
	chatRepo := repo.NewChatServiceRepo(l.ctx, l.svcCtx)

	// 转换请求为数据模型
	converter := converter.NewChatConverter()
	conversation := converter.FromCreateConversationRequest(req)

	// 创建对话
	if err := chatRepo.CreateConversation(conversation); err != nil {
		l.Logger.Error("CreateConversation failed: ", err)
		return &types.CreateConversationResponse{
			Code: 500,
			Msg:  "创建对话失败",
		}, nil
	}

	// 构建响应
	resp = converter.BuildCreateConversationResponse(conversation)
	return resp, nil
}
