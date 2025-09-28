package chat

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"ai-roleplay/services/chat/api/internal/logic/chat"
	"ai-roleplay/services/chat/api/internal/svc"
	"ai-roleplay/services/chat/api/internal/types"

	"github.com/zeromicro/go-zero/core/logc"
	"github.com/zeromicro/go-zero/core/threading"
	"github.com/zeromicro/go-zero/rest/httpx"
	xhttp "github.com/zeromicro/x/http"
)

// 发送消息并获取SSE流式响应
func ChatSendHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.ChatSendRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := chat.NewChatSendLogic(r.Context(), svcCtx)
		err := l.ChatSend(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.Ok(w)
		}
	}
}

func ChatSendHandler2(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.ChatSendRequest
		if err := httpx.Parse(r, &req); err != nil {
			xhttp.JsonBaseResponseCtx(r.Context(), w, err)
			return
		}

		client := make(chan *types.ChatSSEEvent, 16)
		defer func() {
			close(client)
		}()
		l := chat.NewChatSendLogic(r.Context(), svcCtx)
		ctx, cancel := context.WithCancel(r.Context())
		defer cancel() // 确保后台任务取消
		threading.GoSafeCtx(ctx, func() {
			err := l.Sse(&req, client)
			if err != nil {
				logc.Errorf(r.Context(), "sseHandler", logc.Field("error", err))
				return
			}
		})
		for {
			select {
			case data, ok := <-client:
				if !ok { // 通道已关闭
					return
				}
				if err := writeSSE(w, data); err != nil {
					return
				}
			case <-ctx.Done():
				return
			}
		}
	}
}
func writeSSE(w http.ResponseWriter, event *types.ChatSSEEvent) error {
	data, err := json.Marshal(event)
	if err != nil {
		return err
	}
	if _, err := fmt.Fprintf(w, "data: %s\n\n", data); err != nil {
		return err // 写入失败（连接已关闭）
	}
	if flusher, ok := w.(http.Flusher); ok {
		flusher.Flush()
	}
	return nil
}
