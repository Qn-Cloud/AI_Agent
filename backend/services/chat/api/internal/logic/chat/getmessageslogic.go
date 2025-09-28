package chat

import (
	"context"

	"ai-roleplay/services/chat/api/internal/converter"
	"ai-roleplay/services/chat/api/internal/repo"
	"ai-roleplay/services/chat/api/internal/svc"
	"ai-roleplay/services/chat/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetMessagesLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取对话消息历史
func NewGetMessagesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetMessagesLogic {
	return &GetMessagesLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetMessagesLogic) GetMessages(req *types.MessageListRequest) (resp *types.MessageListResponse, err error) {
	l.Logger.Infof("GetMessages请求参数: %+v", req)

	// 参数验证
	if req.ConversationID <= 0 {
		l.Logger.Errorf("对话ID无效: %d", req.ConversationID)
		return &types.MessageListResponse{
			Code: 400,
			Msg:  "对话ID无效",
		}, nil
	}

	// 创建repo实例
	chatRepo := repo.NewChatServiceRepo(l.ctx, l.svcCtx)

	// 获取消息列表
	messages, total, err := chatRepo.GetMessages(req)
	if err != nil {
		l.Logger.Error("GetMessages failed: ", err)
		return &types.MessageListResponse{
			Code: 500,
			Msg:  "获取消息列表失败",
		}, nil
	}

	l.Logger.Infof("查询到 %d 条消息，总数: %d", len(messages), total)

	// 转换数据格式
	converter := converter.NewChatConverter()
	messageList := converter.ToMessageList(messages)

	// 计算分页信息
	page := req.Page
	if page <= 0 {
		page = 1
	}
	pageSize := req.PageSize
	if pageSize <= 0 {
		pageSize = 50
	}
	hasMore := int64(page*pageSize) < total

	return &types.MessageListResponse{
		Code:     0,
		Msg:      "获取成功",
		Messages: messageList,
		Total:    total,
		Page:     page,
		HasMore:  hasMore,
	}, nil
}
