package chat

import (
	"net/http"

	"ai-roleplay/services/chat/api/internal/logic/chat"
	"ai-roleplay/services/chat/api/internal/svc"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// 导出对话记录
func ExportConversationHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := chat.NewExportConversationLogic(r.Context(), svcCtx)
		resp, err := l.ExportConversation()
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
