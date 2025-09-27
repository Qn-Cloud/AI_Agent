package public

import (
	"net/http"

	"ai-roleplay/services/character/api/internal/logic/public"
	"ai-roleplay/services/character/api/internal/svc"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// 获取角色标签
func GetCharacterTagsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := public.NewGetCharacterTagsLogic(r.Context(), svcCtx)
		resp, err := l.GetCharacterTags()
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
