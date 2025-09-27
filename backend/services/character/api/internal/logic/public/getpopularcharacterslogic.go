package public

import (
	"ai-roleplay/services/character/api/internal/converter"
	"ai-roleplay/services/character/api/internal/repo"
	"ai-roleplay/services/character/api/internal/svc"
	"ai-roleplay/services/character/api/internal/types"
	"context"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetPopularCharactersLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetPopularCharactersLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetPopularCharactersLogic {
	return &GetPopularCharactersLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetPopularCharactersLogic) GetPopularCharacters(req *types.PopularCharacterRequest) (resp *types.PopularCharacterResponse, err error) {
	// 设置默认值
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 20
	}

	// 创建repo实例
	characterRepo := repo.NewCharacterServiceRepo(l.ctx, l.svcCtx)

	// 获取热门角色
	characters, total, err := characterRepo.GetPopularCharacters(req)
	if err != nil {
		l.Logger.Error("GetPopularCharacters failed: ", err)
		return &types.PopularCharacterResponse{
			Code: 500,
			Msg:  "获取热门角色失败",
		}, nil
	}

	// 转换数据格式
	converter := converter.NewCharacterConverter()
	briefList := converter.ToCharacterBriefList(characters)
	pagination := converter.BuildPagination(req.Page, req.PageSize, total)

	resp = &types.PopularCharacterResponse{
		Code:  0,
		Msg:   "获取成功",
		Total: total,
		Page:  pagination,
		List:  briefList,
	}

	return resp, nil
}
