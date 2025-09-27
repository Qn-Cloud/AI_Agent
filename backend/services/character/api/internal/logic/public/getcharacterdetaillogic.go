package public

import (
	"ai-roleplay/services/character/api/internal/converter"
	"ai-roleplay/services/character/api/internal/repo"
	"ai-roleplay/services/character/api/internal/svc"
	"ai-roleplay/services/character/api/internal/types"
	"context"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetCharacterDetailLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetCharacterDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetCharacterDetailLogic {
	return &GetCharacterDetailLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetCharacterDetailLogic) GetCharacterDetail(req *types.CharacterDetailRequest) (resp *types.CharacterDetailResponse, err error) {
	// 参数验证
	if req.ID <= 0 {
		return &types.CharacterDetailResponse{
			Code: 400,
			Msg:  "角色ID无效",
		}, nil
	}

	// 创建repo实例
	characterRepo := repo.NewCharacterServiceRepo(l.ctx, l.svcCtx)

	// 获取角色详情
	character, err := characterRepo.GetCharacterByID(req.ID)
	if err != nil {
		l.Logger.Error("GetCharacterByID failed: ", err)
		return &types.CharacterDetailResponse{
			Code: 500,
			Msg:  "获取角色详情失败",
		}, nil
	}

	if character == nil {
		return &types.CharacterDetailResponse{
			Code: 404,
			Msg:  "角色不存在",
		}, nil
	}

	// 转换数据格式
	converter := converter.NewCharacterConverter()
	resp = converter.BuildCharacterDetailResponse(character)

	return resp, nil
}
