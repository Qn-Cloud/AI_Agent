package public

import (
	"ai-roleplay/services/character/api/internal/converter"
	"ai-roleplay/services/character/api/internal/repo"
	"ai-roleplay/services/character/api/internal/svc"
	"ai-roleplay/services/character/api/internal/types"
	"context"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetCharacterListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetCharacterListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetCharacterListLogic {
	return &GetCharacterListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetCharacterListLogic) GetCharacterList(req *types.CharacterListRequest) (resp *types.CharacterListResponse, err error) {
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 10
	}

	// 创建repo实例
	characterRepo := repo.NewCharacterServiceRepo(l.ctx, l.svcCtx)

	// 获取角色列表
	characters, total, err := characterRepo.GetCharacterList(req)
	if err != nil {
		l.Logger.Error("GetCharacterList failed: ", err)
		return &types.CharacterListResponse{
			Code: 500,
			Msg:  "获取角色列表失败",
		}, nil
	}

	// 转换数据格式
	converter := converter.NewCharacterConverter()
	resp = converter.BuildCharacterListResponse(characters, total, req.Page, req.PageSize)

	return resp, nil
}
