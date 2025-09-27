package character

import (
	"context"

	"ai-roleplay/services/character/api/internal/svc"
	"ai-roleplay/services/character/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetMyCharactersLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取我创建的角色
func NewGetMyCharactersLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetMyCharactersLogic {
	return &GetMyCharactersLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetMyCharactersLogic) GetMyCharacters(req *types.MyCharacterRequest) (resp *types.MyCharacterResponse, err error) {
	// todo: add your logic here and delete this line

	return
}
