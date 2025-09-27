package public

import (
	"ai-roleplay/services/character/api/internal/repo"
	"ai-roleplay/services/character/api/internal/svc"
	"ai-roleplay/services/character/api/internal/types"
	"context"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetCharacterCategoriesLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetCharacterCategoriesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetCharacterCategoriesLogic {
	return &GetCharacterCategoriesLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetCharacterCategoriesLogic) GetCharacterCategories() (resp *types.CharacterCategoriesResponse, err error) {
	// 创建repo实例
	characterRepo := repo.NewCharacterServiceRepo(l.ctx, l.svcCtx)

	// 获取角色分类
	categories, err := characterRepo.GetCharacterCategories()
	if err != nil {
		l.Logger.Error("GetCharacterCategories failed: ", err)
		return &types.CharacterCategoriesResponse{
			Code: 500,
			Msg:  "获取角色分类失败",
		}, nil
	}

	resp = &types.CharacterCategoriesResponse{
		Code:       0,
		Msg:        "获取成功",
		Categories: categories,
	}

	return resp, nil
}
