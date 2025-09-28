package handler

import (
	"net/http"

	"ai-roleplay/services/speech/api/internal/logic"
	"ai-roleplay/services/speech/api/internal/svc"
	"ai-roleplay/services/speech/api/internal/types"

	"github.com/zeromicro/go-zero/rest/httpx"
	xhttp "github.com/zeromicro/x/http"
)

// 文字转语音
func textToSpeechHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.TextToSpeechRequest
		if err := httpx.Parse(r, &req); err != nil {
			xhttp.JsonBaseResponseCtx(r.Context(), w, err)
			return
		}

		l := logic.NewTextToSpeechLogic(r.Context(), svcCtx)
		resp, err := l.TextToSpeech(&req)
		if err != nil {
			xhttp.JsonBaseResponseCtx(r.Context(), w, err)
		} else {
			xhttp.JsonBaseResponseCtx(r.Context(), w, resp)
		}
	}
}
