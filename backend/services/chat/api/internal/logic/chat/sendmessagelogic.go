package chat

import (
	"ai-roleplay/services/chat/api/internal/converter"
	"ai-roleplay/services/chat/api/internal/repo"
	"ai-roleplay/services/chat/api/internal/svc"
	"ai-roleplay/services/chat/api/internal/types"
	"context"

	"github.com/zeromicro/go-zero/core/logx"
)

type SendMessageLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSendMessageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SendMessageLogic {
	return &SendMessageLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SendMessageLogic) SendMessage(req *types.SendMessageRequest) (resp *types.SendMessageResponse, err error) {
	// 参数验证
	if req.ConversationID <= 0 {
		return &types.SendMessageResponse{
			Code: 400,
			Msg:  "对话ID无效",
		}, nil
	}

	if req.Content == "" {
		return &types.SendMessageResponse{
			Code: 400,
			Msg:  "消息内容不能为空",
		}, nil
	}

	if req.Type != "user" && req.Type != "ai" {
		return &types.SendMessageResponse{
			Code: 400,
			Msg:  "消息类型无效",
		}, nil
	}

	// 创建repo实例
	chatRepo := repo.NewChatServiceRepo(l.ctx, l.svcCtx)

	// 验证对话是否存在
	conversation, err := chatRepo.GetConversationByID(req.ConversationID)
	if err != nil {
		l.Logger.Error("GetConversationByID failed: ", err)
		return &types.SendMessageResponse{
			Code: 500,
			Msg:  "获取对话信息失败",
		}, nil
	}

	if conversation == nil {
		return &types.SendMessageResponse{
			Code: 404,
			Msg:  "对话不存在",
		}, nil
	}

	// 转换请求为数据模型
	converter := converter.NewChatConverter()
	message := converter.FromSendMessageRequest(req)

	// 保存消息
	if err := chatRepo.SendMessage(message); err != nil {
		l.Logger.Error("SendMessage failed: ", err)
		return &types.SendMessageResponse{
			Code: 500,
			Msg:  "发送消息失败",
		}, nil
	}

	// 更新对话的最后更新时间
	if err := chatRepo.UpdateConversationUpdatedAt(req.ConversationID); err != nil {
		l.Logger.Error("UpdateConversationUpdatedAt failed: ", err)
		// 不影响主流程，只记录日志
	}

	// 构建响应
	resp = converter.BuildSendMessageResponse(message)
	return resp, nil
}
