package chat

import (
	"net/http"

	"ai-roleplay/services/chat/api/internal/logic/chat"
	"ai-roleplay/services/chat/api/internal/svc"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// 清空对话消息
func ClearMessagesHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := chat.NewClearMessagesLogic(r.Context(), svcCtx)
		resp, err := l.ClearMessages()
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
