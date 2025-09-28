package public

import (
	"ai-roleplay/services/character/api/internal/repo"
	"ai-roleplay/services/character/api/internal/svc"
	"ai-roleplay/services/character/api/internal/types"
	"context"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdatePromptLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdatePromptLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdatePromptLogic {
	return &UpdatePromptLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdatePromptLogic) UpdatePrompt(req *types.UpdatePromptRequest) (resp *types.UpdatePromptResponse, err error) {
	// 参数验证
	if req.ID <= 0 {
		return &types.UpdatePromptResponse{
			Code: 400,
			Msg:  "角色ID无效",
		}, nil
	}

	if req.Prompt == "" {
		return &types.UpdatePromptResponse{
			Code: 400,
			Msg:  "提示词不能为空",
		}, nil
	}

	currentUserID := int64(1)

	// 创建repo实例
	characterRepo := repo.NewCharacterServiceRepo(l.ctx, l.svcCtx)

	// 检查角色是否存在且有权限
	existingCharacter, err := characterRepo.GetCharacterByID(req.ID)
	if err != nil {
		l.Logger.Error("GetCharacterByID failed: ", err)
		return &types.UpdatePromptResponse{
			Code: 500,
			Msg:  "获取角色信息失败",
		}, nil
	}

	if existingCharacter == nil {
		return &types.UpdatePromptResponse{
			Code: 404,
			Msg:  "角色不存在",
		}, nil
	}

	// 权限检查：只能更新自己创建的角色
	if existingCharacter.CreatorID == nil || *existingCharacter.CreatorID != currentUserID {
		return &types.UpdatePromptResponse{
			Code: 403,
			Msg:  "无权限更新此角色",
		}, nil
	}

	// 更新提示词
	if err := characterRepo.UpdatePrompt(req.ID, currentUserID, req.Prompt); err != nil {
		l.Logger.Error("UpdatePrompt failed: ", err)
		return &types.UpdatePromptResponse{
			Code: 500,
			Msg:  "更新提示词失败",
		}, nil
	}

	return &types.UpdatePromptResponse{
		Code: 0,
		Msg:  "更新成功",
	}, nil
}
