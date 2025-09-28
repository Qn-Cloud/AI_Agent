package public

import (
	"ai-roleplay/services/character/api/internal/converter"
	"ai-roleplay/services/character/api/internal/repo"
	"ai-roleplay/services/character/api/internal/svc"
	"ai-roleplay/services/character/api/internal/types"
	"context"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetMyCharactersLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetMyCharactersLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetMyCharactersLogic {
	return &GetMyCharactersLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetMyCharactersLogic) GetMyCharacters(req *types.MyCharacterRequest) (resp *types.MyCharacterResponse, err error) {
	// 获取当前用户ID
	currentUserID := int64(1) // 临时硬编码

	// 创建repo实例
	characterRepo := repo.NewCharacterServiceRepo(l.ctx, l.svcCtx)

	// 获取我创建的角色
	characters, total, err := characterRepo.GetMyCharacters(currentUserID, req)
	if err != nil {
		l.Logger.Error("GetMyCharacters failed: ", err)
		return &types.MyCharacterResponse{
			Code: 500,
			Msg:  "获取角色列表失败",
		}, nil
	}

	// 转换数据格式
	converter := converter.NewCharacterConverter()
	listResp := converter.BuildCharacterListResponse(characters, total, req.Page, req.PageSize)

	// 构建我的角色响应
	resp = &types.MyCharacterResponse{
		Code:  0,
		Msg:   "获取成功",
		Total: listResp.Total,
		Page:  listResp.Page,
		List:  listResp.List,
	}

	return resp, nil
}
