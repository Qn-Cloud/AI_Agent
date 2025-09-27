package character

import (
	"context"

	"ai-roleplay/services/character/api/internal/svc"
	"ai-roleplay/services/character/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdatePersonalityLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 更新角色性格设置
func NewUpdatePersonalityLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdatePersonalityLogic {
	return &UpdatePersonalityLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdatePersonalityLogic) UpdatePersonality(req *types.UpdatePersonalityRequest) (resp *types.UpdatePersonalityResponse, err error) {
	// todo: add your logic here and delete this line

	return
}
