package public

import (
	"ai-roleplay/services/character/api/internal/converter"
	"ai-roleplay/services/character/api/internal/repo"
	"ai-roleplay/services/character/api/internal/svc"
	"ai-roleplay/services/character/api/internal/types"
	"context"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetMyFavoritesLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetMyFavoritesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetMyFavoritesLogic {
	return &GetMyFavoritesLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetMyFavoritesLogic) GetMyFavorites(req *types.FavoriteCharacterRequest) (resp *types.FavoriteCharacterResponse, err error) {
	// 获取当前用户ID

	currentUserID := int64(1) // 临时硬编码

	// 创建repo实例
	characterRepo := repo.NewCharacterServiceRepo(l.ctx, l.svcCtx)

	// 获取我的收藏角色
	characters, total, err := characterRepo.GetMyFavorites(currentUserID, req)
	if err != nil {
		l.Logger.Error("GetMyFavorites failed: ", err)
		return &types.FavoriteCharacterResponse{
			Code: 500,
			Msg:  "获取收藏列表失败",
		}, nil
	}

	// 转换数据格式
	converter := converter.NewCharacterConverter()
	listResp := converter.BuildCharacterListResponse(characters, total, req.Page, req.PageSize)

	// 构建收藏响应
	resp = &types.FavoriteCharacterResponse{
		Code:  0,
		Msg:   "获取成功",
		Total: listResp.Total,
		Page:  listResp.Page,
		List:  listResp.List,
	}

	return resp, nil
}
