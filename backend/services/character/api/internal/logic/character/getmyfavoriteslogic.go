package character

import (
	"context"

	"ai-roleplay/services/character/api/internal/svc"
	"ai-roleplay/services/character/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetMyFavoritesLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取我的收藏角色
func NewGetMyFavoritesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetMyFavoritesLogic {
	return &GetMyFavoritesLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetMyFavoritesLogic) GetMyFavorites(req *types.FavoriteCharacterRequest) (resp *types.FavoriteCharacterResponse, err error) {
	// todo: add your logic here and delete this line

	return
}
