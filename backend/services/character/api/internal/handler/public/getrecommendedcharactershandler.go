package public

import (
	"net/http"

	"ai-roleplay/services/character/api/internal/logic/public"
	"ai-roleplay/services/character/api/internal/svc"
	"ai-roleplay/services/character/api/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// 获取推荐角色
func GetRecommendedCharactersHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.RecommendedCharacterRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := public.NewGetRecommendedCharactersLogic(r.Context(), svcCtx)
		resp, err := l.GetRecommendedCharacters(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
