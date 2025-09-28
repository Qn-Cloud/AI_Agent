package handler

import (
	"net/http"

	"ai-roleplay/services/speech/api/internal/logic"
	"ai-roleplay/services/speech/api/internal/svc"
	"ai-roleplay/services/speech/api/internal/types"

	"github.com/zeromicro/go-zero/rest/httpx"
	xhttp "github.com/zeromicro/x/http"
)

// 语音转文字
func speechToTextHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.SpeechToTextRequest

		// 对于multipart form data，只解析非文件字段
		if err := httpx.Parse(r, &req); err != nil {
			// 如果解析失败但是multipart数据，继续处理
			if r.Header.Get("Content-Type") == "" || r.Header.Get("Content-Type")[:19] != "multipart/form-data" {
				xhttp.JsonBaseResponseCtx(r.Context(), w, err)
				return
			}
		}

		l := logic.NewSpeechToTextLogic(r.Context(), svcCtx, r)
		resp, err := l.SpeechToText(&req)
		if err != nil {
			xhttp.JsonBaseResponseCtx(r.Context(), w, err)
		} else {
			xhttp.JsonBaseResponseCtx(r.Context(), w, resp)
		}
	}
}
