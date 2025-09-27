package public

import (
	"ai-roleplay/services/character/api/internal/converter"
	"ai-roleplay/services/character/api/internal/repo"
	"ai-roleplay/services/character/api/internal/svc"
	"ai-roleplay/services/character/api/internal/types"
	"context"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetRecommendedCharactersLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetRecommendedCharactersLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetRecommendedCharactersLogic {
	return &GetRecommendedCharactersLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetRecommendedCharactersLogic) GetRecommendedCharacters(req *types.RecommendedCharacterRequest) (resp *types.RecommendedCharacterResponse, err error) {
	// 设置默认值
	count := req.Count
	if count <= 0 {
		count = 10
	}

	// 创建repo实例
	characterRepo := repo.NewCharacterServiceRepo(l.ctx, l.svcCtx)

	// 获取推荐角色
	characters, err := characterRepo.GetRecommendedCharacters(count)
	if err != nil {
		l.Logger.Error("GetRecommendedCharacters failed: ", err)
		return &types.RecommendedCharacterResponse{
			Code: 500,
			Msg:  "获取推荐角色失败",
		}, nil
	}

	// 转换数据格式
	converter := converter.NewCharacterConverter()
	briefList := converter.ToCharacterBriefList(characters)
	pagination := converter.BuildPagination(1, count, int64(len(characters)))

	resp = &types.RecommendedCharacterResponse{
		Code:  0,
		Msg:   "获取成功",
		Total: int64(len(characters)),
		Page:  pagination,
		List:  briefList,
	}

	return resp, nil
}
