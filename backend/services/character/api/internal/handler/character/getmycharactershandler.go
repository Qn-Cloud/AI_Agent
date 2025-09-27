package character

import (
	"net/http"

	"ai-roleplay/services/character/api/internal/logic/character"
	"ai-roleplay/services/character/api/internal/svc"
	"ai-roleplay/services/character/api/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// 获取我创建的角色
func GetMyCharactersHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.MyCharacterRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := character.NewGetMyCharactersLogic(r.Context(), svcCtx)
		resp, err := l.GetMyCharacters(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
