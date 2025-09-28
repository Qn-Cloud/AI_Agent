package handler

import (
	"ai-roleplay/services/speech/api/internal/logic"
	"ai-roleplay/services/speech/api/internal/svc"
	"net/http"

	xhttp "github.com/zeromicro/x/http"
)

// 健康检查
func healthCheckHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := logic.NewHealthCheckLogic(r.Context(), svcCtx)
		resp, err := l.HealthCheck()
		if err != nil {
			xhttp.JsonBaseResponseCtx(r.Context(), w, err)
		} else {
			xhttp.JsonBaseResponseCtx(r.Context(), w, resp)
		}
	}
}
