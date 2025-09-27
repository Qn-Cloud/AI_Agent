package chat

import (
	"ai-roleplay/services/chat/api/internal/config"
	"ai-roleplay/services/chat/api/internal/svc"
	"ai-roleplay/services/chat/api/internal/types"
	"context"
	"testing"

	"github.com/zeromicro/go-zero/core/conf"
)

var (
	svcCtx *svc.ServiceContext
	ctx    context.Context
	// logic
	batchDeleteConversationsLogic *BatchDeleteConversationsLogic
	createConversationLogic       *CreateConversationLogic
	deleteConversationLogic       *DeleteConversationLogic
	exportConversationLogic       *ExportConversationLogic
	getConversationLogic          *GetConversationLogic
	getConversationListLogic      *GetConversationListLogic
	getMessagesLogic              *GetMessagesLogic
	searchConversationsLogic      *SearchConversationsLogic
	sendMessageLogic              *SendMessageLogic
	updateConversationTitleLogic  *UpdateConversationTitleLogic
	getConversationHistoryLogic   *GetConversationHistoryLogic
)

func init() {
	var c config.Config
	conf.MustLoad("../../../etc/chat-api.yaml", &c)
	svcCtx = svc.NewServiceContext(c)
	ctx = context.Background()
}

func TestBatchDeleteConversationsLogic(t *testing.T) {
	batchDeleteConversationsLogic = NewBatchDeleteConversationsLogic(ctx, svcCtx)
	resp, err := batchDeleteConversationsLogic.BatchDeleteConversations(&types.BatchDeleteRequest{
		ConversationIDs: []int64{1, 2, 3},
	})
	if err != nil {
		t.Fatalf("BatchDeleteConversations failed: %v", err)
	}
	t.Logf("BatchDeleteConversations resp: %v\n", resp)
}

func TestCreateConversationLogic(t *testing.T) {
	createConversationLogic = NewCreateConversationLogic(ctx, svcCtx)
	resp, err := createConversationLogic.CreateConversation(&types.CreateConversationRequest{
		Title: "test",
	})
	if err != nil {
		t.Fatalf("CreateConversation failed: %v", err)
	}
	t.Logf("CreateConversation resp: %v\n", resp)
}

func TestDeleteConversationLogic(t *testing.T) {
	deleteConversationLogic = NewDeleteConversationLogic(ctx, svcCtx)
	resp, err := deleteConversationLogic.DeleteConversation(&types.ConversationRequest{
		ID: 1,
	})
	if err != nil {
		t.Fatalf("DeleteConversation failed: %v", err)
	}
	t.Logf("DeleteConversation resp: %v\n", resp)
}

func TestExportConversationLogic(t *testing.T) {
	exportConversationLogic = NewExportConversationLogic(ctx, svcCtx)
	resp, err := exportConversationLogic.ExportConversation(&types.ConversationRequest{
		ID: 1,
	})
	if err != nil {
		t.Fatalf("ExportConversation failed: %v", err)
	}
	t.Logf("ExportConversation resp: %v\n", resp)
}

func TestGetConversationHistoryLogic(t *testing.T) {
	getConversationHistoryLogic = NewGetConversationHistoryLogic(ctx, svcCtx)
	resp, err := getConversationHistoryLogic.GetConversationHistory(&types.GetConversationHistoryRequest{
		UserID: 1,
	})
	if err != nil {
		t.Fatalf("GetConversationHistory failed: %v", err)
	}
	t.Logf("GetConversationHistory resp: %v\n", resp)
}
