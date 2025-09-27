package public

import (
	"ai-roleplay/services/character/api/internal/converter"
	"ai-roleplay/services/character/api/internal/repo"
	"ai-roleplay/services/character/api/internal/svc"
	"ai-roleplay/services/character/api/internal/types"
	"context"

	"github.com/zeromicro/go-zero/core/logx"
)

type SearchCharactersLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSearchCharactersLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SearchCharactersLogic {
	return &SearchCharactersLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SearchCharactersLogic) SearchCharacters(req *types.SearchCharacterRequest) (resp *types.SearchCharacterResponse, err error) {
	// 参数验证
	if req.Keyword == "" {
		return &types.SearchCharacterResponse{
			Code: 400,
			Msg:  "搜索关键词不能为空",
		}, nil
	}

	// 创建repo实例
	characterRepo := repo.NewCharacterServiceRepo(l.ctx, l.svcCtx)

	// 搜索角色
	characters, total, err := characterRepo.SearchCharacters(req)
	if err != nil {
		l.Logger.Error("SearchCharacters failed: ", err)
		return &types.SearchCharacterResponse{
			Code: 500,
			Msg:  "搜索角色失败",
		}, nil
	}

	// 转换数据格式
	converter := converter.NewCharacterConverter()
	briefList := converter.ToCharacterBriefList(characters)
	pagination := converter.BuildPagination(req.Page, req.PageSize, total)

	resp = &types.SearchCharacterResponse{
		Code:  0,
		Msg:   "搜索成功",
		Total: total,
		Page:  pagination,
		List:  briefList,
	}

	return resp, nil
}
