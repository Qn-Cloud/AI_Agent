package public

import (
	"ai-roleplay/services/character/api/internal/repo"
	"ai-roleplay/services/character/api/internal/svc"
	"ai-roleplay/services/character/api/internal/types"
	"context"
	"encoding/json"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdatePersonalityLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdatePersonalityLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdatePersonalityLogic {
	return &UpdatePersonalityLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdatePersonalityLogic) UpdatePersonality(req *types.UpdatePersonalityRequest) (resp *types.UpdatePersonalityResponse, err error) {
	// 参数验证
	if req.ID <= 0 {
		return &types.UpdatePersonalityResponse{
			Code: 400,
			Msg:  "角色ID无效",
		}, nil
	}

	// 获取当前用户ID
	// TODO: 从JWT token中获取真实的用户ID
	currentUserID := int64(1) // 临时硬编码

	// 创建repo实例
	characterRepo := repo.NewCharacterServiceRepo(l.ctx, l.svcCtx)

	// 检查角色是否存在且有权限
	existingCharacter, err := characterRepo.GetCharacterByID(req.ID)
	if err != nil {
		l.Logger.Error("GetCharacterByID failed: ", err)
		return &types.UpdatePersonalityResponse{
			Code: 500,
			Msg:  "获取角色信息失败",
		}, nil
	}

	if existingCharacter == nil {
		return &types.UpdatePersonalityResponse{
			Code: 404,
			Msg:  "角色不存在",
		}, nil
	}

	// 权限检查：只能更新自己创建的角色
	if existingCharacter.CreatorID == nil || *existingCharacter.CreatorID != currentUserID {
		return &types.UpdatePersonalityResponse{
			Code: 403,
			Msg:  "无权限更新此角色",
		}, nil
	}

	// 将性格设置转换为JSON
	personalityJSON, err := json.Marshal(req.Personality)
	if err != nil {
		l.Logger.Error("Marshal personality failed: ", err)
		return &types.UpdatePersonalityResponse{
			Code: 500,
			Msg:  "性格设置格式错误",
		}, nil
	}

	// 更新性格设置
	if err := characterRepo.UpdatePersonality(req.ID, currentUserID, string(personalityJSON)); err != nil {
		l.Logger.Error("UpdatePersonality failed: ", err)
		return &types.UpdatePersonalityResponse{
			Code: 500,
			Msg:  "更新性格设置失败",
		}, nil
	}

	return &types.UpdatePersonalityResponse{
		Code: 0,
		Msg:  "更新成功",
	}, nil
}
