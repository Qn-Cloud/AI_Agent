package public

import (
	"ai-roleplay/services/character/api/internal/repo"
	"ai-roleplay/services/character/api/internal/svc"
	"ai-roleplay/services/character/api/internal/types"
	"context"

	"github.com/zeromicro/go-zero/core/logx"
)

type ToggleFavoriteLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewToggleFavoriteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ToggleFavoriteLogic {
	return &ToggleFavoriteLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ToggleFavoriteLogic) ToggleFavorite(req *types.ToggleFavoriteRequest) (resp *types.ToggleFavoriteResponse, err error) {
	// 参数验证
	if req.ID <= 0 {
		return &types.ToggleFavoriteResponse{
			Code: 400,
			Msg:  "角色ID无效",
		}, nil
	}

	// 创建repo实例
	characterRepo := repo.NewCharacterServiceRepo(l.ctx, l.svcCtx)

	// 检查角色是否存在
	existingCharacter, err := characterRepo.GetCharacterByID(req.ID)
	if err != nil {
		l.Logger.Error("GetCharacterByID failed: ", err)
		return &types.ToggleFavoriteResponse{
			Code: 500,
			Msg:  "获取角色信息失败",
		}, nil
	}

	if existingCharacter == nil {
		return &types.ToggleFavoriteResponse{
			Code: 404,
			Msg:  "角色不存在",
		}, nil
	}

	// 获取当前用户ID
	// TODO: 从JWT token中获取真实的用户ID
	currentUserID := int64(1) // 临时硬编码

	// 切换收藏状态
	isFavorite, err := characterRepo.ToggleFavorite(currentUserID, req.ID)
	if err != nil {
		l.Logger.Error("ToggleFavorite failed: ", err)
		return &types.ToggleFavoriteResponse{
			Code: 500,
			Msg:  "操作失败",
		}, nil
	}

	var msg string
	if isFavorite {
		msg = "收藏成功"
	} else {
		msg = "取消收藏成功"
	}

	return &types.ToggleFavoriteResponse{
		Code:       0,
		Msg:        msg,
		IsFavorite: isFavorite,
	}, nil
}
