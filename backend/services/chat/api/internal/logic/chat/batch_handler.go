package chat

import (
	"ai-roleplay/services/chat/api/internal/svc"
	"ai-roleplay/services/chat/api/internal/types"
	"context"
	"fmt"
	"reflect"
	"sync"

	"github.com/zeromicro/go-zero/core/logx"
)

// LogicFactory 逻辑工厂接口
type LogicFactory interface {
	CreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) interface{}
	GetLogicName() string
}

// BatchChatLogicManager 批量聊天逻辑管理器
type BatchChatLogicManager struct {
	logx.Logger
	ctx       context.Context
	svcCtx    *svc.ServiceContext
	logics    map[string]interface{}
	factories map[string]LogicFactory
	mutex     sync.RWMutex
}

// NewBatchChatLogicManager 创建批量聊天逻辑管理器
func NewBatchChatLogicManager(ctx context.Context, svcCtx *svc.ServiceContext) *BatchChatLogicManager {
	manager := &BatchChatLogicManager{
		Logger:    logx.WithContext(ctx),
		ctx:       ctx,
		svcCtx:    svcCtx,
		logics:    make(map[string]interface{}),
		factories: make(map[string]LogicFactory),
	}

	// 注册所有Logic工厂
	manager.registerFactories()

	return manager
}

// registerFactories 注册所有Logic工厂
func (m *BatchChatLogicManager) registerFactories() {
	// 注册各种Logic工厂
	m.RegisterFactory(&SendMessageLogicFactory{})
	m.RegisterFactory(&CreateConversationLogicFactory{})
	m.RegisterFactory(&GetConversationLogicFactory{})
	m.RegisterFactory(&GetConversationListLogicFactory{})
	m.RegisterFactory(&GetMessagesLogicFactory{})
	m.RegisterFactory(&DeleteConversationLogicFactory{})
	m.RegisterFactory(&ClearMessagesLogicFactory{})
	m.RegisterFactory(&UpdateConversationTitleLogicFactory{})
	m.RegisterFactory(&SearchConversationsLogicFactory{})
	m.RegisterFactory(&ExportConversationLogicFactory{})
	m.RegisterFactory(&BatchDeleteConversationsLogicFactory{})
}

// RegisterFactory 注册Logic工厂
func (m *BatchChatLogicManager) RegisterFactory(factory LogicFactory) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.factories[factory.GetLogicName()] = factory
}

// BatchCreateLogics 批量创建Logic实例
func (m *BatchChatLogicManager) BatchCreateLogics(logicNames []string) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	for _, name := range logicNames {
		if factory, exists := m.factories[name]; exists {
			m.logics[name] = factory.CreateLogic(m.ctx, m.svcCtx)
			m.Logger.Infof("Created logic: %s", name)
		} else {
			return fmt.Errorf("logic factory not found: %s", name)
		}
	}

	return nil
}

// GetLogic 获取Logic实例
func (m *BatchChatLogicManager) GetLogic(logicName string) (interface{}, error) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	if logic, exists := m.logics[logicName]; exists {
		return logic, nil
	}

	// 如果不存在，尝试创建
	if factory, exists := m.factories[logicName]; exists {
		m.mutex.RUnlock()
		m.mutex.Lock()
		logic := factory.CreateLogic(m.ctx, m.svcCtx)
		m.logics[logicName] = logic
		m.mutex.Unlock()
		m.mutex.RLock()
		return logic, nil
	}

	return nil, fmt.Errorf("logic not found: %s", logicName)
}

// BatchCall 批量调用函数
type BatchCallRequest struct {
	LogicName  string
	MethodName string
	Args       []interface{}
}

type BatchCallResponse struct {
	LogicName  string
	MethodName string
	Result     []interface{}
	Error      error
}

// BatchCallFunctions 批量调用函数
func (m *BatchChatLogicManager) BatchCallFunctions(requests []BatchCallRequest) []BatchCallResponse {
	var wg sync.WaitGroup
	responses := make([]BatchCallResponse, len(requests))

	for i, req := range requests {
		wg.Add(1)
		go func(index int, request BatchCallRequest) {
			defer wg.Done()
			responses[index] = m.callFunction(request)
		}(i, req)
	}

	wg.Wait()
	return responses
}

// callFunction 调用单个函数
func (m *BatchChatLogicManager) callFunction(req BatchCallRequest) BatchCallResponse {
	response := BatchCallResponse{
		LogicName:  req.LogicName,
		MethodName: req.MethodName,
	}

	// 获取Logic实例
	logic, err := m.GetLogic(req.LogicName)
	if err != nil {
		response.Error = err
		return response
	}

	// 使用反射调用方法
	logicValue := reflect.ValueOf(logic)
	method := logicValue.MethodByName(req.MethodName)
	if !method.IsValid() {
		response.Error = fmt.Errorf("method not found: %s.%s", req.LogicName, req.MethodName)
		return response
	}

	// 准备参数
	args := make([]reflect.Value, len(req.Args))
	for i, arg := range req.Args {
		args[i] = reflect.ValueOf(arg)
	}

	// 调用方法
	results := method.Call(args)

	// 处理结果
	response.Result = make([]interface{}, len(results))
	for i, result := range results {
		response.Result[i] = result.Interface()
	}

	return response
}

// =============================================================================
// Logic工厂实现
// =============================================================================

// SendMessageLogicFactory 发送消息Logic工厂
type SendMessageLogicFactory struct{}

func (f *SendMessageLogicFactory) CreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) interface{} {
	return NewSendMessageLogic(ctx, svcCtx)
}

func (f *SendMessageLogicFactory) GetLogicName() string {
	return "SendMessageLogic"
}

// CreateConversationLogicFactory 创建对话Logic工厂
type CreateConversationLogicFactory struct{}

func (f *CreateConversationLogicFactory) CreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) interface{} {
	return NewCreateConversationLogic(ctx, svcCtx)
}

func (f *CreateConversationLogicFactory) GetLogicName() string {
	return "CreateConversationLogic"
}

// GetConversationLogicFactory 获取对话Logic工厂
type GetConversationLogicFactory struct{}

func (f *GetConversationLogicFactory) CreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) interface{} {
	return NewGetConversationLogic(ctx, svcCtx)
}

func (f *GetConversationLogicFactory) GetLogicName() string {
	return "GetConversationLogic"
}

// GetConversationListLogicFactory 获取对话列表Logic工厂
type GetConversationListLogicFactory struct{}

func (f *GetConversationListLogicFactory) CreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) interface{} {
	// return NewGetConversationListLogic(ctx, svcCtx) // 需要实现
	return nil
}

func (f *GetConversationListLogicFactory) GetLogicName() string {
	return "GetConversationListLogic"
}

// GetMessagesLogicFactory 获取消息列表Logic工厂
type GetMessagesLogicFactory struct{}

func (f *GetMessagesLogicFactory) CreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) interface{} {
	// return NewGetMessagesLogic(ctx, svcCtx) // 需要实现
	return nil
}

func (f *GetMessagesLogicFactory) GetLogicName() string {
	return "GetMessagesLogic"
}

// DeleteConversationLogicFactory 删除对话Logic工厂
type DeleteConversationLogicFactory struct{}

func (f *DeleteConversationLogicFactory) CreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) interface{} {
	// return NewDeleteConversationLogic(ctx, svcCtx) // 需要实现
	return nil
}

func (f *DeleteConversationLogicFactory) GetLogicName() string {
	return "DeleteConversationLogic"
}

// ClearMessagesLogicFactory 清空消息Logic工厂
type ClearMessagesLogicFactory struct{}

func (f *ClearMessagesLogicFactory) CreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) interface{} {
	// return NewClearMessagesLogic(ctx, svcCtx) // 需要实现
	return nil
}

func (f *ClearMessagesLogicFactory) GetLogicName() string {
	return "ClearMessagesLogic"
}

// UpdateConversationTitleLogicFactory 更新对话标题Logic工厂
type UpdateConversationTitleLogicFactory struct{}

func (f *UpdateConversationTitleLogicFactory) CreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) interface{} {
	// return NewUpdateConversationTitleLogic(ctx, svcCtx) // 需要实现
	return nil
}

func (f *UpdateConversationTitleLogicFactory) GetLogicName() string {
	return "UpdateConversationTitleLogic"
}

// SearchConversationsLogicFactory 搜索对话Logic工厂
type SearchConversationsLogicFactory struct{}

func (f *SearchConversationsLogicFactory) CreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) interface{} {
	// return NewSearchConversationsLogic(ctx, svcCtx) // 需要实现
	return nil
}

func (f *SearchConversationsLogicFactory) GetLogicName() string {
	return "SearchConversationsLogic"
}

// ExportConversationLogicFactory 导出对话Logic工厂
type ExportConversationLogicFactory struct{}

func (f *ExportConversationLogicFactory) CreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) interface{} {
	// return NewExportConversationLogic(ctx, svcCtx) // 需要实现
	return nil
}

func (f *ExportConversationLogicFactory) GetLogicName() string {
	return "ExportConversationLogic"
}

// BatchDeleteConversationsLogicFactory 批量删除对话Logic工厂
type BatchDeleteConversationsLogicFactory struct{}

func (f *BatchDeleteConversationsLogicFactory) CreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) interface{} {
	return NewBatchDeleteConversationsLogic(ctx, svcCtx)
}

func (f *BatchDeleteConversationsLogicFactory) GetLogicName() string {
	return "BatchDeleteConversationsLogic"
}

// =============================================================================
// 便捷方法
// =============================================================================

// BatchChatHelper 批量聊天助手
type BatchChatHelper struct {
	manager *BatchChatLogicManager
}

// NewBatchChatHelper 创建批量聊天助手
func NewBatchChatHelper(ctx context.Context, svcCtx *svc.ServiceContext) *BatchChatHelper {
	return &BatchChatHelper{
		manager: NewBatchChatLogicManager(ctx, svcCtx),
	}
}

// InitAllLogics 初始化所有Logic
func (h *BatchChatHelper) InitAllLogics() error {
	logicNames := []string{
		"SendMessageLogic",
		"CreateConversationLogic",
		"GetConversationLogic",
		"GetConversationListLogic",
		"GetMessagesLogic",
		"DeleteConversationLogic",
		"ClearMessagesLogic",
		"UpdateConversationTitleLogic",
		"SearchConversationsLogic",
		"ExportConversationLogic",
		"BatchDeleteConversationsLogic",
	}

	return h.manager.BatchCreateLogics(logicNames)
}

// ChatOperation 聊天操作定义
type ChatOperation struct {
	Type string      `json:"type"` // "send_message", "create_conversation", "get_conversation", etc.
	Data interface{} `json:"data"` // 请求参数
}

// ChatOperationResult 聊天操作结果
type ChatOperationResult struct {
	Type   string      `json:"type"`
	Result interface{} `json:"result"`
	Error  error       `json:"error,omitempty"`
}

// BatchExecuteChatOperations 批量执行聊天相关操作
func (h *BatchChatHelper) BatchExecuteChatOperations(operations []ChatOperation) ([]ChatOperationResult, error) {
	var results []ChatOperationResult
	var wg sync.WaitGroup
	resultChan := make(chan ChatOperationResult, len(operations))

	for _, op := range operations {
		wg.Add(1)
		go func(operation ChatOperation) {
			defer wg.Done()
			result := h.executeOperation(operation)
			resultChan <- result
		}(op)
	}

	go func() {
		wg.Wait()
		close(resultChan)
	}()

	for result := range resultChan {
		results = append(results, result)
	}

	return results, nil
}

// executeOperation 执行单个操作
func (h *BatchChatHelper) executeOperation(op ChatOperation) ChatOperationResult {
	result := ChatOperationResult{Type: op.Type}

	switch op.Type {
	case "send_message":
		if req, ok := op.Data.(*types.SendMessageRequest); ok {
			logic, err := h.manager.GetLogic("SendMessageLogic")
			if err != nil {
				result.Error = err
				return result
			}
			resp, err := logic.(*SendMessageLogic).SendMessage(req)
			result.Result = resp
			result.Error = err
		} else {
			result.Error = fmt.Errorf("invalid request type for send_message operation")
		}

	case "create_conversation":
		if req, ok := op.Data.(*types.CreateConversationRequest); ok {
			logic, err := h.manager.GetLogic("CreateConversationLogic")
			if err != nil {
				result.Error = err
				return result
			}
			resp, err := logic.(*CreateConversationLogic).CreateConversation(req)
			result.Result = resp
			result.Error = err
		} else {
			result.Error = fmt.Errorf("invalid request type for create_conversation operation")
		}

	case "get_conversation":
		if req, ok := op.Data.(*types.ConversationRequest); ok {
			logic, err := h.manager.GetLogic("GetConversationLogic")
			if err != nil {
				result.Error = err
				return result
			}
			resp, err := logic.(*GetConversationLogic).GetConversation(req)
			result.Result = resp
			result.Error = err
		} else {
			result.Error = fmt.Errorf("invalid request type for get_conversation operation")
		}

	case "batch_delete":
		if req, ok := op.Data.(*types.BatchDeleteRequest); ok {
			logic, err := h.manager.GetLogic("BatchDeleteConversationsLogic")
			if err != nil {
				result.Error = err
				return result
			}
			resp, err := logic.(*BatchDeleteConversationsLogic).BatchDeleteConversations(req)
			result.Result = resp
			result.Error = err
		} else {
			result.Error = fmt.Errorf("invalid request type for batch_delete operation")
		}

	default:
		result.Error = fmt.Errorf("unknown operation type: %s", op.Type)
	}

	return result
}
