package character

import (
	"context"

	"ai-roleplay/services/character/api/internal/svc"
	"ai-roleplay/services/character/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateCharacterLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 更新角色信息
func NewUpdateCharacterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateCharacterLogic {
	return &UpdateCharacterLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateCharacterLogic) UpdateCharacter(req *types.UpdateCharacterRequest) (resp *types.UpdateCharacterResponse, err error) {
	// todo: add your logic here and delete this line

	return
}
