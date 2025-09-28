package public

import (
	"ai-roleplay/services/character/api/internal/repo"
	"ai-roleplay/services/character/api/internal/svc"
	"ai-roleplay/services/character/api/internal/types"
	"context"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteCharacterLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteCharacterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteCharacterLogic {
	return &DeleteCharacterLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteCharacterLogic) DeleteCharacter(req *types.DeleteCharacterRequest) (resp *types.DeleteCharacterResponse, err error) {
	// 参数验证
	if req.ID <= 0 {
		return &types.DeleteCharacterResponse{
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
		return &types.DeleteCharacterResponse{
			Code: 500,
			Msg:  "获取角色信息失败",
		}, nil
	}

	if existingCharacter == nil {
		return &types.DeleteCharacterResponse{
			Code: 404,
			Msg:  "角色不存在",
		}, nil
	}

	// 权限检查：只能删除自己创建的角色
	// TODO: 从JWT token中获取真实的用户ID
	currentUserID := int64(1) // 临时硬编码
	if existingCharacter.CreatorID == nil || *existingCharacter.CreatorID != currentUserID {
		return &types.DeleteCharacterResponse{
			Code: 403,
			Msg:  "无权限删除此角色",
		}, nil
	}

	// 检查角色是否被使用（有对话记录）
	if existingCharacter.ChatCount > 0 {
		// 软删除：只更新状态，不真正删除
		// 这里可以根据业务需求决定是否允许删除有对话记录的角色
		l.Logger.Infof("尝试删除有对话记录的角色: character_id: %d, chat_count: %d", req.ID, existingCharacter.ChatCount)
	}

	// 删除角色
	if err := characterRepo.DeleteCharacter(req.ID, currentUserID); err != nil {
		l.Logger.Error("DeleteCharacter failed: ", err)
		return &types.DeleteCharacterResponse{
			Code: 500,
			Msg:  "删除角色失败",
		}, nil
	}

	return &types.DeleteCharacterResponse{
		Code: 0,
		Msg:  "删除成功",
	}, nil
}
