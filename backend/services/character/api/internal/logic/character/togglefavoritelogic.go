package character

import (
	"context"

	"ai-roleplay/services/character/api/internal/svc"
	"ai-roleplay/services/character/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ToggleFavoriteLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 收藏/取消收藏角色
func NewToggleFavoriteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ToggleFavoriteLogic {
	return &ToggleFavoriteLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ToggleFavoriteLogic) ToggleFavorite(req *types.ToggleFavoriteRequest) (resp *types.ToggleFavoriteResponse, err error) {
	// todo: add your logic here and delete this line

	return
}
