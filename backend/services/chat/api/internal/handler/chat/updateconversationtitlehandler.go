package chat

import (
	"net/http"

	"ai-roleplay/services/chat/api/internal/logic/chat"
	"ai-roleplay/services/chat/api/internal/svc"
	"ai-roleplay/services/chat/api/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// 更新对话标题
func UpdateConversationTitleHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.UpdateTitleRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := chat.NewUpdateConversationTitleLogic(r.Context(), svcCtx)
		resp, err := l.UpdateConversationTitle(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
