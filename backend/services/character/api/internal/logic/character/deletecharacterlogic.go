package character

import (
	"context"

	"ai-roleplay/services/character/api/internal/svc"
	"ai-roleplay/services/character/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteCharacterLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 删除角色
func NewDeleteCharacterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteCharacterLogic {
	return &DeleteCharacterLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteCharacterLogic) DeleteCharacter(req *types.DeleteCharacterRequest) (resp *types.DeleteCharacterResponse, err error) {
	// todo: add your logic here and delete this line

	return
}
