package character

import (
	"net/http"

	"ai-roleplay/services/character/api/internal/logic/character"
	"ai-roleplay/services/character/api/internal/svc"
	"ai-roleplay/services/character/api/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// 更新角色提示词
func UpdatePromptHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.UpdatePromptRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := character.NewUpdatePromptLogic(r.Context(), svcCtx)
		resp, err := l.UpdatePrompt(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
