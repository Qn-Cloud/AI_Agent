package public

import (
	"ai-roleplay/services/character/api/internal/svc"
	"ai-roleplay/services/character/api/internal/types"
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

// BatchLogicManager 批量逻辑管理器
type BatchLogicManager struct {
	logx.Logger
	ctx       context.Context
	svcCtx    *svc.ServiceContext
	logics    map[string]interface{}
	factories map[string]LogicFactory
	mutex     sync.RWMutex
}

// NewBatchLogicManager 创建批量逻辑管理器
func NewBatchLogicManager(ctx context.Context, svcCtx *svc.ServiceContext) *BatchLogicManager {
	manager := &BatchLogicManager{
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
func (m *BatchLogicManager) registerFactories() {
	// 注册各种Logic工厂
	m.RegisterFactory(&GetCharacterListLogicFactory{})
	m.RegisterFactory(&SearchCharactersLogicFactory{})
	m.RegisterFactory(&GetCharacterDetailLogicFactory{})
	m.RegisterFactory(&GetRecommendedCharactersLogicFactory{})
	m.RegisterFactory(&GetCharacterCategoriesLogicFactory{})
	m.RegisterFactory(&GetPopularCharactersLogicFactory{})
	m.RegisterFactory(&GetCharacterTagsLogicFactory{})
}

// RegisterFactory 注册Logic工厂
func (m *BatchLogicManager) RegisterFactory(factory LogicFactory) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.factories[factory.GetLogicName()] = factory
}

// BatchCreateLogics 批量创建Logic实例
func (m *BatchLogicManager) BatchCreateLogics(logicNames []string) error {
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
func (m *BatchLogicManager) GetLogic(logicName string) (interface{}, error) {
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
func (m *BatchLogicManager) BatchCallFunctions(requests []BatchCallRequest) []BatchCallResponse {
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
func (m *BatchLogicManager) callFunction(req BatchCallRequest) BatchCallResponse {
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

// GetCharacterListLogicFactory 角色列表Logic工厂
type GetCharacterListLogicFactory struct{}

func (f *GetCharacterListLogicFactory) CreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) interface{} {
	return NewGetCharacterListLogic(ctx, svcCtx)
}

func (f *GetCharacterListLogicFactory) GetLogicName() string {
	return "GetCharacterListLogic"
}

// SearchCharactersLogicFactory 搜索角色Logic工厂
type SearchCharactersLogicFactory struct{}

func (f *SearchCharactersLogicFactory) CreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) interface{} {
	return NewSearchCharactersLogic(ctx, svcCtx)
}

func (f *SearchCharactersLogicFactory) GetLogicName() string {
	return "SearchCharactersLogic"
}

// GetCharacterDetailLogicFactory 角色详情Logic工厂
type GetCharacterDetailLogicFactory struct{}

func (f *GetCharacterDetailLogicFactory) CreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) interface{} {
	return NewGetCharacterDetailLogic(ctx, svcCtx)
}

func (f *GetCharacterDetailLogicFactory) GetLogicName() string {
	return "GetCharacterDetailLogic"
}

// GetRecommendedCharactersLogicFactory 推荐角色Logic工厂
type GetRecommendedCharactersLogicFactory struct{}

func (f *GetRecommendedCharactersLogicFactory) CreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) interface{} {
	return NewGetRecommendedCharactersLogic(ctx, svcCtx)
}

func (f *GetRecommendedCharactersLogicFactory) GetLogicName() string {
	return "GetRecommendedCharactersLogic"
}

// GetCharacterCategoriesLogicFactory 角色分类Logic工厂
type GetCharacterCategoriesLogicFactory struct{}

func (f *GetCharacterCategoriesLogicFactory) CreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) interface{} {
	return NewGetCharacterCategoriesLogic(ctx, svcCtx)
}

func (f *GetCharacterCategoriesLogicFactory) GetLogicName() string {
	return "GetCharacterCategoriesLogic"
}

// GetPopularCharactersLogicFactory 热门角色Logic工厂
type GetPopularCharactersLogicFactory struct{}

func (f *GetPopularCharactersLogicFactory) CreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) interface{} {
	return NewGetPopularCharactersLogic(ctx, svcCtx)
}

func (f *GetPopularCharactersLogicFactory) GetLogicName() string {
	return "GetPopularCharactersLogic"
}

// GetCharacterTagsLogicFactory 角色标签Logic工厂
type GetCharacterTagsLogicFactory struct{}

func (f *GetCharacterTagsLogicFactory) CreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) interface{} {
	return NewGetCharacterTagsLogic(ctx, svcCtx)
}

func (f *GetCharacterTagsLogicFactory) GetLogicName() string {
	return "GetCharacterTagsLogic"
}

// =============================================================================
// 便捷方法
// =============================================================================

// BatchLogicHelper 批量Logic助手
type BatchLogicHelper struct {
	manager *BatchLogicManager
}

// NewBatchLogicHelper 创建批量Logic助手
func NewBatchLogicHelper(ctx context.Context, svcCtx *svc.ServiceContext) *BatchLogicHelper {
	return &BatchLogicHelper{
		manager: NewBatchLogicManager(ctx, svcCtx),
	}
}

// InitAllLogics 初始化所有Logic
func (h *BatchLogicHelper) InitAllLogics() error {
	logicNames := []string{
		"GetCharacterListLogic",
		"SearchCharactersLogic",
		"GetCharacterDetailLogic",
		"GetRecommendedCharactersLogic",
		"GetCharacterCategoriesLogic",
		"GetPopularCharactersLogic",
		"GetCharacterTagsLogic",
	}

	return h.manager.BatchCreateLogics(logicNames)
}

// GetCharacterListLogic 获取角色列表Logic
func (h *BatchLogicHelper) GetCharacterListLogic() (*GetCharacterListLogic, error) {
	logic, err := h.manager.GetLogic("GetCharacterListLogic")
	if err != nil {
		return nil, err
	}
	return logic.(*GetCharacterListLogic), nil
}

// SearchCharactersLogic 获取搜索角色Logic
func (h *BatchLogicHelper) SearchCharactersLogic() (*SearchCharactersLogic, error) {
	logic, err := h.manager.GetLogic("SearchCharactersLogic")
	if err != nil {
		return nil, err
	}
	return logic.(*SearchCharactersLogic), nil
}

// GetCharacterDetailLogic 获取角色详情Logic
func (h *BatchLogicHelper) GetCharacterDetailLogic() (*GetCharacterDetailLogic, error) {
	logic, err := h.manager.GetLogic("GetCharacterDetailLogic")
	if err != nil {
		return nil, err
	}
	return logic.(*GetCharacterDetailLogic), nil
}

// GetRecommendedCharactersLogic 获取推荐角色Logic
func (h *BatchLogicHelper) GetRecommendedCharactersLogic() (*GetRecommendedCharactersLogic, error) {
	logic, err := h.manager.GetLogic("GetRecommendedCharactersLogic")
	if err != nil {
		return nil, err
	}
	return logic.(*GetRecommendedCharactersLogic), nil
}

// GetCharacterCategoriesLogic 获取角色分类Logic
func (h *BatchLogicHelper) GetCharacterCategoriesLogic() (*GetCharacterCategoriesLogic, error) {
	logic, err := h.manager.GetLogic("GetCharacterCategoriesLogic")
	if err != nil {
		return nil, err
	}
	return logic.(*GetCharacterCategoriesLogic), nil
}

// GetPopularCharactersLogic 获取热门角色Logic
func (h *BatchLogicHelper) GetPopularCharactersLogic() (*GetPopularCharactersLogic, error) {
	logic, err := h.manager.GetLogic("GetPopularCharactersLogic")
	if err != nil {
		return nil, err
	}
	return logic.(*GetPopularCharactersLogic), nil
}

// GetCharacterTagsLogic 获取角色标签Logic
func (h *BatchLogicHelper) GetCharacterTagsLogic() (*GetCharacterTagsLogic, error) {
	logic, err := h.manager.GetLogic("GetCharacterTagsLogic")
	if err != nil {
		return nil, err
	}
	return logic.(*GetCharacterTagsLogic), nil
}

// BatchExecuteCharacterOperations 批量执行角色相关操作
func (h *BatchLogicHelper) BatchExecuteCharacterOperations(operations []CharacterOperation) ([]CharacterOperationResult, error) {
	var results []CharacterOperationResult
	var wg sync.WaitGroup
	resultChan := make(chan CharacterOperationResult, len(operations))

	for _, op := range operations {
		wg.Add(1)
		go func(operation CharacterOperation) {
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

// CharacterOperation 角色操作定义
type CharacterOperation struct {
	Type string      `json:"type"` // "list", "search", "detail", "recommended", "categories", "popular", "tags"
	Data interface{} `json:"data"` // 请求参数
}

// CharacterOperationResult 角色操作结果
type CharacterOperationResult struct {
	Type   string      `json:"type"`
	Result interface{} `json:"result"`
	Error  error       `json:"error,omitempty"`
}

// executeOperation 执行单个操作
func (h *BatchLogicHelper) executeOperation(op CharacterOperation) CharacterOperationResult {
	result := CharacterOperationResult{Type: op.Type}

	switch op.Type {
	case "list":
		if req, ok := op.Data.(*types.CharacterListRequest); ok {
			logic, err := h.GetCharacterListLogic()
			if err != nil {
				result.Error = err
				return result
			}
			resp, err := logic.GetCharacterList(req)
			result.Result = resp
			result.Error = err
		} else {
			result.Error = fmt.Errorf("invalid request type for list operation")
		}

	case "search":
		if req, ok := op.Data.(*types.SearchCharacterRequest); ok {
			logic, err := h.SearchCharactersLogic()
			if err != nil {
				result.Error = err
				return result
			}
			resp, err := logic.SearchCharacters(req)
			result.Result = resp
			result.Error = err
		} else {
			result.Error = fmt.Errorf("invalid request type for search operation")
		}

	case "detail":
		if req, ok := op.Data.(*types.CharacterDetailRequest); ok {
			logic, err := h.GetCharacterDetailLogic()
			if err != nil {
				result.Error = err
				return result
			}
			resp, err := logic.GetCharacterDetail(req)
			result.Result = resp
			result.Error = err
		} else {
			result.Error = fmt.Errorf("invalid request type for detail operation")
		}

	case "recommended":
		if req, ok := op.Data.(*types.RecommendedCharacterRequest); ok {
			logic, err := h.GetRecommendedCharactersLogic()
			if err != nil {
				result.Error = err
				return result
			}
			resp, err := logic.GetRecommendedCharacters(req)
			result.Result = resp
			result.Error = err
		} else {
			result.Error = fmt.Errorf("invalid request type for recommended operation")
		}

	case "categories":
		logic, err := h.GetCharacterCategoriesLogic()
		if err != nil {
			result.Error = err
			return result
		}
		resp, err := logic.GetCharacterCategories()
		result.Result = resp
		result.Error = err

	case "popular":
		if req, ok := op.Data.(*types.PopularCharacterRequest); ok {
			logic, err := h.GetPopularCharactersLogic()
			if err != nil {
				result.Error = err
				return result
			}
			resp, err := logic.GetPopularCharacters(req)
			result.Result = resp
			result.Error = err
		} else {
			result.Error = fmt.Errorf("invalid request type for popular operation")
		}

	case "tags":
		logic, err := h.GetCharacterTagsLogic()
		if err != nil {
			result.Error = err
			return result
		}
		resp, err := logic.GetCharacterTags()
		result.Result = resp
		result.Error = err

	default:
		result.Error = fmt.Errorf("unknown operation type: %s", op.Type)
	}

	return result
}
