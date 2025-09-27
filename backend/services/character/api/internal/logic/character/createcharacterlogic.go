package character

import (
	"context"

	"ai-roleplay/services/character/api/internal/svc"
	"ai-roleplay/services/character/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateCharacterLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 创建自定义角色
func NewCreateCharacterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateCharacterLogic {
	return &CreateCharacterLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateCharacterLogic) CreateCharacter(req *types.CreateCharacterRequest) (resp *types.CreateCharacterResponse, err error) {
	// todo: add your logic here and delete this line

	return
}
