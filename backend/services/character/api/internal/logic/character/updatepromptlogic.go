package character

import (
	"context"

	"ai-roleplay/services/character/api/internal/svc"
	"ai-roleplay/services/character/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdatePromptLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 更新角色提示词
func NewUpdatePromptLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdatePromptLogic {
	return &UpdatePromptLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdatePromptLogic) UpdatePrompt(req *types.UpdatePromptRequest) (resp *types.UpdatePromptResponse, err error) {
	// todo: add your logic here and delete this line

	return
}
