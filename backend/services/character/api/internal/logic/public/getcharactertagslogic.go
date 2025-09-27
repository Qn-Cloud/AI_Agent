package public

import (
	"ai-roleplay/services/character/api/internal/repo"
	"ai-roleplay/services/character/api/internal/svc"
	"ai-roleplay/services/character/api/internal/types"
	"context"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetCharacterTagsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetCharacterTagsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetCharacterTagsLogic {
	return &GetCharacterTagsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetCharacterTagsLogic) GetCharacterTags() (resp *types.CharacterTagsResponse, err error) {
	// 创建repo实例
	characterRepo := repo.NewCharacterServiceRepo(l.ctx, l.svcCtx)

	// 获取角色标签
	tags, err := characterRepo.GetCharacterTags()
	if err != nil {
		l.Logger.Error("GetCharacterTags failed: ", err)
		return &types.CharacterTagsResponse{
			Code: 500,
			Msg:  "获取角色标签失败",
		}, nil
	}

	resp = &types.CharacterTagsResponse{
		Code: 0,
		Msg:  "获取成功",
		Tags: tags,
	}

	return resp, nil
}
