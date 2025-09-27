package public

import (
	"net/http"

	"ai-roleplay/services/character/api/internal/logic/public"
	"ai-roleplay/services/character/api/internal/svc"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// 获取角色分类
func GetCharacterCategoriesHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := public.NewGetCharacterCategoriesLogic(r.Context(), svcCtx)
		resp, err := l.GetCharacterCategories()
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
