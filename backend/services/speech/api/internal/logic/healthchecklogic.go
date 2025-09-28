package logic

import (
	"context"

	"ai-roleplay/services/speech/api/internal/svc"
	"ai-roleplay/services/speech/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type HealthCheckLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 健康检查
func NewHealthCheckLogic(ctx context.Context, svcCtx *svc.ServiceContext) *HealthCheckLogic {
	return &HealthCheckLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *HealthCheckLogic) HealthCheck() (resp *types.BaseResponse, err error) {
	l.Logger.Info("语音服务健康检查")

	// 这里可以添加各种健康检查逻辑
	// 例如：检查数据库连接、检查外部服务状态等

	return &types.BaseResponse{
		Code: 200,
		Msg:  "语音服务运行正常",
	}, nil
}
