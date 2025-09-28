package public

import (
	"ai-roleplay/services/character/api/internal/repo"
	"ai-roleplay/services/character/api/internal/svc"
	"ai-roleplay/services/character/api/internal/types"
	"context"
	"encoding/json"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateVoiceSettingsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateVoiceSettingsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateVoiceSettingsLogic {
	return &UpdateVoiceSettingsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateVoiceSettingsLogic) UpdateVoiceSettings(req *types.UpdateVoiceSettingsRequest) (resp *types.UpdateVoiceSettingsResponse, err error) {
	// 参数验证
	if req.ID <= 0 {
		return &types.UpdateVoiceSettingsResponse{
			Code: 400,
			Msg:  "角色ID无效",
		}, nil
	}

	currentUserID := int64(1)

	// 创建repo实例
	characterRepo := repo.NewCharacterServiceRepo(l.ctx, l.svcCtx)

	// 检查角色是否存在且有权限
	existingCharacter, err := characterRepo.GetCharacterByID(req.ID)
	if err != nil {
		l.Logger.Error("GetCharacterByID failed: ", err)
		return &types.UpdateVoiceSettingsResponse{
			Code: 500,
			Msg:  "获取角色信息失败",
		}, nil
	}

	if existingCharacter == nil {
		return &types.UpdateVoiceSettingsResponse{
			Code: 404,
			Msg:  "角色不存在",
		}, nil
	}

	// 权限检查：只能更新自己创建的角色
	if existingCharacter.CreatorID == nil || *existingCharacter.CreatorID != currentUserID {
		return &types.UpdateVoiceSettingsResponse{
			Code: 403,
			Msg:  "无权限更新此角色",
		}, nil
	}

	// 将语音设置转换为JSON
	voiceSettingsJSON, err := json.Marshal(req.VoiceSettings)
	if err != nil {
		l.Logger.Error("Marshal voice settings failed: ", err)
		return &types.UpdateVoiceSettingsResponse{
			Code: 500,
			Msg:  "语音设置格式错误",
		}, nil
	}

	// 更新语音设置
	if err := characterRepo.UpdateVoiceSettings(req.ID, currentUserID, string(voiceSettingsJSON)); err != nil {
		l.Logger.Error("UpdateVoiceSettings failed: ", err)
		return &types.UpdateVoiceSettingsResponse{
			Code: 500,
			Msg:  "更新语音设置失败",
		}, nil
	}

	return &types.UpdateVoiceSettingsResponse{
		Code: 0,
		Msg:  "更新成功",
	}, nil
}
