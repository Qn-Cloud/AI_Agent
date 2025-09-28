package public

import (
	"ai-roleplay/services/character/api/internal/converter"
	"ai-roleplay/services/character/api/internal/repo"
	"ai-roleplay/services/character/api/internal/svc"
	"ai-roleplay/services/character/api/internal/types"
	"ai-roleplay/services/character/model"
	"context"
	"encoding/json"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateCharacterLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateCharacterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateCharacterLogic {
	return &CreateCharacterLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateCharacterLogic) CreateCharacter(req *types.CreateCharacterRequest) (resp *types.CreateCharacterResponse, err error) {
	// 参数验证
	if req.Name == "" {
		return &types.CreateCharacterResponse{
			Code: 400,
			Msg:  "角色名称不能为空",
		}, nil
	}

	// 创建repo实例
	characterRepo := repo.NewCharacterServiceRepo(l.ctx, l.svcCtx)

	// 构建角色模型
	character := &model.Character{
		Name:          req.Name,
		Description:   &req.Description,
		ShortDesc:     &req.ShortDesc,
		CategoryID:    &req.CategoryID,
		Prompt:        &req.Prompt,
		Status:        1, // 正常状态
		IsPublic:      int32(boolToInt(req.IsPublic)),
		Rating:        0.0,
		RatingCount:   0,
		FavoriteCount: 0,
		ChatCount:     0,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	// 设置头像
	if req.Avatar != "" {
		character.Avatar = &req.Avatar
	}

	// 设置创建者ID（从JWT中获取，这里暂时模拟）
	// TODO: 从JWT token中获取真实的用户ID
	creatorID := int64(1) // 临时硬编码
	character.CreatorID = &creatorID

	// 处理标签
	if len(req.Tags) > 0 {
		tagsJSON, err := json.Marshal(req.Tags)
		if err != nil {
			l.Logger.Error("Marshal tags failed: ", err)
			return &types.CreateCharacterResponse{
				Code: 500,
				Msg:  "标签格式错误",
			}, nil
		}
		tagsStr := string(tagsJSON)
		character.Tags = &tagsStr
	}

	// 处理性格设置
	personalityJSON, err := json.Marshal(req.Personality)
	if err != nil {
		l.Logger.Error("Marshal personality failed: ", err)
		return &types.CreateCharacterResponse{
			Code: 500,
			Msg:  "性格设置格式错误",
		}, nil
	}
	personalityStr := string(personalityJSON)
	character.Personality = &personalityStr

	// 处理语音设置
	voiceSettingsJSON, err := json.Marshal(req.VoiceSettings)
	if err != nil {
		l.Logger.Error("Marshal voice settings failed: ", err)
		return &types.CreateCharacterResponse{
			Code: 500,
			Msg:  "语音设置格式错误",
		}, nil
	}
	voiceSettingsStr := string(voiceSettingsJSON)
	character.VoiceSettings = &voiceSettingsStr

	// 创建角色
	if err := characterRepo.CreateCharacter(character); err != nil {
		l.Logger.Error("CreateCharacter failed: ", err)
		return &types.CreateCharacterResponse{
			Code: 500,
			Msg:  "创建角色失败",
		}, nil
	}

	// 转换数据格式
	converter := converter.NewCharacterConverter()
	characterItem := converter.ToCharacterItem(character)

	return &types.CreateCharacterResponse{
		Code:      0,
		Msg:       "创建成功",
		Character: *characterItem,
	}, nil
}

// 辅助函数：bool转int
func boolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}
