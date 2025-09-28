package public

import (
	"ai-roleplay/services/character/api/internal/converter"
	"ai-roleplay/services/character/api/internal/repo"
	"ai-roleplay/services/character/api/internal/svc"
	"ai-roleplay/services/character/api/internal/types"
	"context"
	"encoding/json"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateCharacterLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateCharacterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateCharacterLogic {
	return &UpdateCharacterLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateCharacterLogic) UpdateCharacter(req *types.UpdateCharacterRequest) (resp *types.UpdateCharacterResponse, err error) {
	// 参数验证
	if req.ID <= 0 {
		return &types.UpdateCharacterResponse{
			Code: 400,
			Msg:  "角色ID无效",
		}, nil
	}

	if req.Name == "" {
		return &types.UpdateCharacterResponse{
			Code: 400,
			Msg:  "角色名称不能为空",
		}, nil
	}

	// 创建repo实例
	characterRepo := repo.NewCharacterServiceRepo(l.ctx, l.svcCtx)

	// 获取现有角色
	existingCharacter, err := characterRepo.GetCharacterByID(req.ID)
	if err != nil {
		l.Logger.Error("GetCharacterByID failed: ", err)
		return &types.UpdateCharacterResponse{
			Code: 500,
			Msg:  "获取角色信息失败",
		}, nil
	}

	if existingCharacter == nil {
		return &types.UpdateCharacterResponse{
			Code: 404,
			Msg:  "角色不存在",
		}, nil
	}

	// 权限检查：只能更新自己创建的角色
	// TODO: 从JWT token中获取真实的用户ID
	currentUserID := int64(1) // 临时硬编码
	if existingCharacter.CreatorID == nil || *existingCharacter.CreatorID != currentUserID {
		return &types.UpdateCharacterResponse{
			Code: 403,
			Msg:  "无权限更新此角色",
		}, nil
	}

	// 更新基础字段
	existingCharacter.Name = req.Name
	existingCharacter.Description = &req.Description
	existingCharacter.ShortDesc = &req.ShortDesc
	existingCharacter.CategoryID = &req.CategoryID
	existingCharacter.Prompt = &req.Prompt
	existingCharacter.Status = req.Status
	existingCharacter.IsPublic = int32(boolToInt(req.IsPublic))
	existingCharacter.UpdatedAt = time.Now()

	// 设置头像
	if req.Avatar != "" {
		existingCharacter.Avatar = &req.Avatar
	}

	// 处理标签
	if len(req.Tags) > 0 {
		tagsJSON, err := json.Marshal(req.Tags)
		if err != nil {
			l.Logger.Error("Marshal tags failed: ", err)
			return &types.UpdateCharacterResponse{
				Code: 500,
				Msg:  "标签格式错误",
			}, nil
		}
		tagsStr := string(tagsJSON)
		existingCharacter.Tags = &tagsStr
	}

	// 处理性格设置
	personalityJSON, err := json.Marshal(req.Personality)
	if err != nil {
		l.Logger.Error("Marshal personality failed: ", err)
		return &types.UpdateCharacterResponse{
			Code: 500,
			Msg:  "性格设置格式错误",
		}, nil
	}
	personalityStr := string(personalityJSON)
	existingCharacter.Personality = &personalityStr

	// 处理语音设置
	voiceSettingsJSON, err := json.Marshal(req.VoiceSettings)
	if err != nil {
		l.Logger.Error("Marshal voice settings failed: ", err)
		return &types.UpdateCharacterResponse{
			Code: 500,
			Msg:  "语音设置格式错误",
		}, nil
	}
	voiceSettingsStr := string(voiceSettingsJSON)
	existingCharacter.VoiceSettings = &voiceSettingsStr

	// 更新角色
	if err := characterRepo.UpdateCharacter(existingCharacter); err != nil {
		l.Logger.Error("UpdateCharacter failed: ", err)
		return &types.UpdateCharacterResponse{
			Code: 500,
			Msg:  "更新角色失败",
		}, nil
	}

	// 转换数据格式
	converter := converter.NewCharacterConverter()
	characterItem := converter.ToCharacterItem(existingCharacter)

	return &types.UpdateCharacterResponse{
		Code:      0,
		Msg:       "更新成功",
		Character: *characterItem,
	}, nil
}
