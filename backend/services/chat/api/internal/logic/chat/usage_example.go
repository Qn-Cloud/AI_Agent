package chat

import (
	"ai-roleplay/services/chat/api/internal/svc"
	"ai-roleplay/services/chat/api/internal/types"
	"context"
	"fmt"
	"log"
	"time"
)

// UsageExample 展示如何使用批量处理模板
func UsageExample(ctx context.Context, svcCtx *svc.ServiceContext) {
	// 1. 创建批量助手
	helper := NewBatchChatHelper(ctx, svcCtx)

	// 2. 初始化所有Logic
	if err := helper.InitAllLogics(); err != nil {
		log.Printf("初始化Logic失败: %v", err)
		return
	}

	// 3. 准备批量操作
	operations := []ChatOperation{
		// 创建对话
		{
			Type: "create_conversation",
			Data: &types.CreateConversationRequest{
				CharacterID: 1,
				Title:       "与哈利波特的对话",
			},
		},
		// 发送消息
		{
			Type: "send_message",
			Data: &types.SendMessageRequest{
				ConversationID: 1,
				CharacterID:    1,
				Content:        "你好，哈利波特！",
			},
		},
		// 获取对话详情
		{
			Type: "get_conversation",
			Data: &types.ConversationRequest{
				ID: 1,
			},
		},
		// 批量删除对话
		{
			Type: "batch_delete",
			Data: &types.BatchDeleteRequest{
				ConversationIDs: []int64{2, 3, 4},
			},
		},
	}

	// 4. 执行批量操作
	results, err := helper.BatchExecuteChatOperations(operations)
	if err != nil {
		log.Printf("批量执行失败: %v", err)
		return
	}

	// 5. 处理结果
	for i, result := range results {
		fmt.Printf("操作 %d (%s):\n", i+1, result.Type)
		if result.Error != nil {
			fmt.Printf("  错误: %v\n", result.Error)
		} else {
			fmt.Printf("  成功: %+v\n", result.Result)
		}
		fmt.Println()
	}
}

// SingleOperationExample 展示单独操作的用法
func SingleOperationExample(ctx context.Context, svcCtx *svc.ServiceContext) {
	// 1. 直接使用Logic
	sendLogic := NewSendMessageLogic(ctx, svcCtx)

	resp, err := sendLogic.SendMessage(&types.SendMessageRequest{
		ConversationID: 1,
		CharacterID:    1,
		Content:        "Hello World!",
	})

	if err != nil {
		log.Printf("发送消息失败: %v", err)
		return
	}

	fmt.Printf("发送消息成功: %+v\n", resp)

	// 2. 使用批量管理器获取Logic
	manager := NewBatchChatLogicManager(ctx, svcCtx)

	logic, err := manager.GetLogic("CreateConversationLogic")
	if err != nil {
		log.Printf("获取Logic失败: %v", err)
		return
	}

	createLogic := logic.(*CreateConversationLogic)
	createResp, err := createLogic.CreateConversation(&types.CreateConversationRequest{
		CharacterID: 2,
		Title:       "新的对话",
	})

	if err != nil {
		log.Printf("创建对话失败: %v", err)
		return
	}

	fmt.Printf("创建对话成功: %+v\n", createResp)
}

// ParallelOperationExample 展示并发操作的用法
func ParallelOperationExample(ctx context.Context, svcCtx *svc.ServiceContext) {
	// 使用BatchCallFunctions进行并发调用
	manager := NewBatchChatLogicManager(ctx, svcCtx)

	requests := []BatchCallRequest{
		{
			LogicName:  "SendMessageLogic",
			MethodName: "SendMessage",
			Args: []interface{}{
				&types.SendMessageRequest{
					ConversationID: 1,
					CharacterID:    1,
					Content:        "消息1",
				},
			},
		},
		{
			LogicName:  "SendMessageLogic",
			MethodName: "SendMessage",
			Args: []interface{}{
				&types.SendMessageRequest{
					ConversationID: 1,
					CharacterID:    1,
					Content:        "消息2",
				},
			},
		},
		{
			LogicName:  "GetConversationLogic",
			MethodName: "GetConversation",
			Args: []interface{}{
				&types.ConversationRequest{ID: 1},
			},
		},
	}

	responses := manager.BatchCallFunctions(requests)

	for i, response := range responses {
		fmt.Printf("并发调用 %d:\n", i+1)
		fmt.Printf("  Logic: %s.%s\n", response.LogicName, response.MethodName)
		if response.Error != nil {
			fmt.Printf("  错误: %v\n", response.Error)
		} else {
			fmt.Printf("  结果: %+v\n", response.Result)
		}
		fmt.Println()
	}
}

// PerformanceComparisonExample 展示性能对比
func PerformanceComparisonExample(ctx context.Context, svcCtx *svc.ServiceContext) {

	helper := NewBatchChatHelper(ctx, svcCtx)
	helper.InitAllLogics()

	// 准备测试数据
	operations := make([]ChatOperation, 100)
	for i := 0; i < 100; i++ {
		operations[i] = ChatOperation{
			Type: "send_message",
			Data: &types.SendMessageRequest{
				ConversationID: 1,
				CharacterID:    1,
				Content:        fmt.Sprintf("测试消息 %d", i),
			},
		}
	}

	// 1. 顺序执行
	start := time.Now()
	sendLogic := NewSendMessageLogic(ctx, svcCtx)
	for i := 0; i < 100; i++ {
		sendLogic.SendMessage(&types.SendMessageRequest{
			ConversationID: 1,
			CharacterID:    1,
			Content:        fmt.Sprintf("顺序消息 %d", i),
		})
	}
	sequentialTime := time.Since(start)

	// 2. 批量并发执行
	start = time.Now()
	helper.BatchExecuteChatOperations(operations)
	parallelTime := time.Since(start)

	fmt.Printf("性能对比:\n")
	fmt.Printf("  顺序执行: %v\n", sequentialTime)
	fmt.Printf("  并发执行: %v\n", parallelTime)
	fmt.Printf("  性能提升: %.2fx\n", float64(sequentialTime)/float64(parallelTime))
}
