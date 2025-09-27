package character

import (
	"context"

	"ai-roleplay/services/character/api/internal/svc"
	"ai-roleplay/services/character/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateVoiceSettingsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 更新语音设置
func NewUpdateVoiceSettingsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateVoiceSettingsLogic {
	return &UpdateVoiceSettingsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateVoiceSettingsLogic) UpdateVoiceSettings(req *types.UpdateVoiceSettingsRequest) (resp *types.UpdateVoiceSettingsResponse, err error) {
	// todo: add your logic here and delete this line

	return
}
