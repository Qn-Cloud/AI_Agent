package repo

import (
	common "ai-roleplay/common/utils"
	"ai-roleplay/services/chat/api/internal/svc"
	"ai-roleplay/services/chat/api/internal/types"
	"ai-roleplay/services/chat/model"
	"context"
	"fmt"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

type ChatServiceRepo struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewChatServiceRepo(ctx context.Context, svcCtx *svc.ServiceContext) *ChatServiceRepo {
	return &ChatServiceRepo{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// SendMessage 发送消息
func (r *ChatServiceRepo) SendMessage(message *model.Message) error {
	db := r.svcCtx.Db.WithContext(r.ctx)

	if err := db.Create(message).Error; err != nil {
		r.Logger.Error("SendMessage failed: ", err)
		return err
	}

	return nil
}

// CreateConversation 创建对话
func (r *ChatServiceRepo) CreateConversation(conversation *model.Conversation) error {
	db := r.svcCtx.Db.WithContext(r.ctx)

	if err := db.Create(conversation).Error; err != nil {
		r.Logger.Error("CreateConversation failed: ", err)
		return err
	}

	return nil
}

// GetConversationByID 根据ID获取对话
func (r *ChatServiceRepo) GetConversationByID(id int64) (*model.Conversation, error) {
	db := r.svcCtx.Db.WithContext(r.ctx)

	var conversation model.Conversation
	if err := db.Where("id = ? AND status != ?", id, common.Deleted).First(&conversation).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		r.Logger.Error("GetConversationByID failed: ", err)
		return nil, err
	}

	return &conversation, nil
}

// GetConversationList 获取对话列表
func (r *ChatServiceRepo) GetConversationList(req *types.ConversationListRequest) ([]model.Conversation, int64, error) {
	db := r.svcCtx.Db.WithContext(r.ctx)

	// 设置默认值
	page := req.Page
	if page <= 0 {
		page = 1
	}
	pageSize := req.PageSize
	if pageSize <= 0 {
		pageSize = 20
	}
	offset := (page - 1) * pageSize

	// 构建查询条件
	query := db.Model(&model.Conversation{})

	// 过滤已删除的记录
	query = query.Where("status != ?", common.Deleted)

	// 用户筛选（如果有用户ID参数的话）
	if req.UserID > 0 {
		query = query.Where("user_id = ?", req.UserID)
	}

	// 角色筛选
	if req.CharacterID > 0 {
		query = query.Where("character_id = ?", req.CharacterID)
	}

	// 统计总数
	var total int64
	if err := query.Count(&total).Error; err != nil {
		r.Logger.Error("GetConversationList count failed: ", err)
		return nil, 0, err
	}

	// 查询数据
	var conversations []model.Conversation
	if err := query.Order("updated_at DESC").Offset(offset).Limit(pageSize).Find(&conversations).Error; err != nil {
		r.Logger.Error("GetConversationList find failed: ", err)
		return nil, 0, err
	}

	return conversations, total, nil
}

// GetMessages 获取对话消息列表
func (r *ChatServiceRepo) GetMessages(req *types.MessageListRequest) ([]model.Message, int64, error) {
	db := r.svcCtx.Db.WithContext(r.ctx)

	// 设置默认值
	page := req.Page
	if page <= 0 {
		page = 1
	}
	pageSize := req.PageSize
	if pageSize <= 0 {
		pageSize = 50
	}
	offset := (page - 1) * pageSize

	// 构建查询条件
	query := db.Model(&model.Message{}).Where("conversation_id = ?", req.ConversationID)

	// 统计总数
	var total int64
	if err := query.Count(&total).Error; err != nil {
		r.Logger.Error("GetMessages count failed: ", err)
		return nil, 0, err
	}

	// 查询数据，按时间正序
	var messages []model.Message
	if err := query.Order("created_at ASC").Offset(offset).Limit(pageSize).Find(&messages).Error; err != nil {
		r.Logger.Error("GetMessages find failed: ", err)
		return nil, 0, err
	}

	return messages, total, nil
}

// DeleteConversation 删除对话（软删除）
func (r *ChatServiceRepo) DeleteConversation(id int64) error {
	db := r.svcCtx.Db.WithContext(r.ctx)

	// 软删除：只修改状态为已删除
	if err := db.Model(&model.Conversation{}).Where("id = ?", id).
		Update("status", common.Deleted).Error; err != nil {
		r.Logger.Error("DeleteConversation failed: ", err)
		return err
	}

	return nil
}

// ClearMessages 清空对话消息
func (r *ChatServiceRepo) ClearMessages(conversationID int64) error {
	db := r.svcCtx.Db.WithContext(r.ctx)

	if err := db.Where("conversation_id = ?", conversationID).Delete(&model.Message{}).Error; err != nil {
		r.Logger.Error("ClearMessages failed: ", err)
		return err
	}

	return nil
}

// UpdateConversationTitle 更新对话标题
func (r *ChatServiceRepo) UpdateConversationTitle(id int64, title string) error {
	db := r.svcCtx.Db.WithContext(r.ctx)

	if err := db.Model(&model.Conversation{}).Where("id = ?", id).
		Update("title", title).Error; err != nil {
		r.Logger.Error("UpdateConversationTitle failed: ", err)
		return err
	}

	return nil
}

// SearchConversations 搜索对话
func (r *ChatServiceRepo) SearchConversations(req *types.SearchConversationRequest) ([]model.Conversation, int64, error) {
	db := r.svcCtx.Db.WithContext(r.ctx)

	page := req.Page
	if page <= 0 {
		page = 1
	}
	pageSize := req.PageSize
	if pageSize <= 0 {
		pageSize = 20
	}
	offset := (page - 1) * pageSize

	// 构建搜索查询
	query := db.Model(&model.Conversation{})

	// 过滤已删除的记录
	query = query.Where("status != ?", common.Deleted)

	if req.Keyword != "" {
		keyword := "%" + req.Keyword + "%"
		query = query.Where("title LIKE ?", keyword)
	}

	// 用户筛选
	if req.UserID > 0 {
		query = query.Where("user_id = ?", req.UserID)
	}

	// 统计总数
	var total int64
	if err := query.Count(&total).Error; err != nil {
		r.Logger.Error("SearchConversations count failed: ", err)
		return nil, 0, err
	}

	// 查询数据
	var conversations []model.Conversation
	if err := query.Order("updated_at DESC").Offset(offset).Limit(pageSize).Find(&conversations).Error; err != nil {
		r.Logger.Error("SearchConversations find failed: ", err)
		return nil, 0, err
	}

	return conversations, total, nil
}

// ExportConversation 导出对话记录
func (r *ChatServiceRepo) ExportConversation(conversationID int64) (*model.Conversation, []model.Message, error) {
	db := r.svcCtx.Db.WithContext(r.ctx)

	// 获取对话信息
	var conversation model.Conversation
	if err := db.Where("id = ?", conversationID).First(&conversation).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil, fmt.Errorf("conversation not found")
		}
		r.Logger.Error("ExportConversation get conversation failed: ", err)
		return nil, nil, err
	}

	// 获取所有消息
	var messages []model.Message
	if err := db.Where("conversation_id = ?", conversationID).
		Order("created_at ASC").Find(&messages).Error; err != nil {
		r.Logger.Error("ExportConversation get messages failed: ", err)
		return nil, nil, err
	}

	return &conversation, messages, nil
}

// BatchDeleteConversations 批量删除对话（软删除）
func (r *ChatServiceRepo) BatchDeleteConversations(conversationIDs []int64) error {
	db := r.svcCtx.Db.WithContext(r.ctx)

	// 批量软删除：只修改状态为已删除
	if err := db.Model(&model.Conversation{}).Where("id IN ?", conversationIDs).
		Update("status", common.Deleted).Error; err != nil {
		r.Logger.Error("BatchDeleteConversations failed: ", err)
		return err
	}

	return nil
}

// UpdateConversationUpdatedAt 更新对话的最后更新时间
func (r *ChatServiceRepo) UpdateConversationUpdatedAt(conversationID int64) error {
	db := r.svcCtx.Db.WithContext(r.ctx)

	if err := db.Model(&model.Conversation{}).Where("id = ?", conversationID).
		Update("updated_at", time.Now()).Error; err != nil {
		r.Logger.Error("UpdateConversationUpdatedAt failed: ", err)
		return err
	}

	return nil
}

// GetMessageByID 根据ID获取消息
func (r *ChatServiceRepo) GetMessageByID(id int64) (*model.Message, error) {
	db := r.svcCtx.Db.WithContext(r.ctx)

	var message model.Message
	if err := db.Where("id = ?", id).First(&message).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		r.Logger.Error("GetMessageByID failed: ", err)
		return nil, err
	}

	return &message, nil
}

func (r *ChatServiceRepo) GetMessageByConversationID(conversationID int64) ([]model.Message, error) {
	db := r.svcCtx.Db.WithContext(r.ctx)

	var messages []model.Message
	if err := db.Where("conversation_id = ?", conversationID).Find(&messages).Error; err != nil {
		r.Logger.Error("GetMessageByConversationID failed: ", err)
		return nil, err
	}

	return messages, nil
}

// UpdateMessage 更新消息
func (r *ChatServiceRepo) UpdateMessage(message *model.Message) error {
	db := r.svcCtx.Db.WithContext(r.ctx)

	if err := db.Save(message).Error; err != nil {
		r.Logger.Error("UpdateMessage failed: ", err)
		return err
	}

	return nil
}

// DeleteMessage 删除消息
func (r *ChatServiceRepo) DeleteMessage(id int64) error {
	db := r.svcCtx.Db.WithContext(r.ctx)

	if err := db.Where("id = ?", id).Delete(&model.Message{}).Error; err != nil {
		r.Logger.Error("DeleteMessage failed: ", err)
		return err
	}

	return nil
}

// GetConversationsByUserID 根据用户ID获取对话列表
func (r *ChatServiceRepo) GetConversationsByUserID(userID int64, page, pageSize int, characterID int64) ([]model.Conversation, int64, error) {
	db := r.svcCtx.Db.WithContext(r.ctx)

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 20
	}
	offset := (page - 1) * pageSize

	query := db.Model(&model.Conversation{}).Where("user_id = ?", userID)
	if characterID > 0 {
		query = query.Where("character_id = ?", characterID)
	}
	query = query.Where("status != ?", common.Deleted)
	// 统计总数
	var total int64
	if err := query.Count(&total).Error; err != nil {
		r.Logger.Error("GetConversationsByUserID count failed: ", err)
		return nil, 0, err
	}

	// 查询数据
	var conversations []model.Conversation
	if err := query.Order("updated_at DESC").Offset(offset).Limit(pageSize).Find(&conversations).Error; err != nil {
		r.Logger.Error("GetConversationsByUserID find failed: ", err)
		return nil, 0, err
	}

	return conversations, total, nil
}

func (r *ChatServiceRepo) GetConversationHistory(req *types.GetConversationHistoryRequest) ([]model.Conversation, int64, error) {
	db := r.svcCtx.Db.WithContext(r.ctx)

	page := req.Page
	if page <= 0 {
		page = 1
	}
	pageSize := req.PageSize
	if pageSize <= 0 {
		pageSize = 20
	}
	offset := (page - 1) * pageSize

	query := db.Model(&model.Conversation{})

	// 过滤已删除的记录
	query = query.Where("status != ?", common.Deleted)

	if req.Keyword != "" {
		keyword := "%" + req.Keyword + "%"
		query = query.Where("title LIKE ?", keyword)
	}

	if req.UserID > 0 {
		query = query.Where("user_id = ?", req.UserID)
	}

	if req.StartTime != "" {
		query = query.Where("created_at >= ?", req.StartTime)
	}

	if req.EndTime != "" {
		query = query.Where("created_at <= ?", req.EndTime)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		r.Logger.Error("GetConversationHistory count failed: ", err)
		return nil, 0, err
	}

	var conversations []model.Conversation
	if err := query.Order("updated_at DESC").Offset(offset).Limit(pageSize).Find(&conversations).Error; err != nil {
		r.Logger.Error("GetConversationHistory find failed: ", err)
		return nil, 0, err
	}

	return conversations, total, nil
}

func (r *ChatServiceRepo) GetConversationDuration(conversationID int64) (int64, error) {
	db := r.svcCtx.Db.WithContext(r.ctx)

	var conversation model.Conversation
	if err := db.Where("id = ?", conversationID).First(&conversation).Error; err != nil {
		r.Logger.Error("GetConversationDuration failed: ", err)
		return 0, err
	}

	return int64(conversation.UpdatedAt.Sub(conversation.CreatedAt).Seconds()), nil
}

func (r *ChatServiceRepo) GetConversationLastMessage(conversationID int64) (*model.Message, error) {
	db := r.svcCtx.Db.WithContext(r.ctx)

	var message model.Message
	if err := db.Where("conversation_id = ?", conversationID).Order("created_at DESC").First(&message).Error; err != nil {
		r.Logger.Error("GetConversationLastMessage failed: ", err)
		return nil, err
	}

	return &message, nil
}

func (r *ChatServiceRepo) GetConversationActiveDays(conversationID int64) (int64, error) {
	db := r.svcCtx.Db.WithContext(r.ctx)

	var conversation model.Conversation
	if err := db.Where("id = ?", conversationID).First(&conversation).Error; err != nil {
		r.Logger.Error("GetConversationActiveDays failed: ", err)
		return 0, err
	}

	return int64(conversation.UpdatedAt.Sub(conversation.CreatedAt).Hours() / 24), nil
}

func (r *ChatServiceRepo) GetConversationMessages(conversationID int64) ([]model.Message, error) {
	db := r.svcCtx.Db.WithContext(r.ctx)

	var messages []model.Message
	if err := db.Where("conversation_id = ?", conversationID).Find(&messages).Error; err != nil {
		r.Logger.Error("GetConversationMessages failed: ", err)
		return nil, err
	}

	return messages, nil
}

func (r *ChatServiceRepo) AddMessage(message *model.Message) (int64, error) {
	db := r.svcCtx.Db.WithContext(r.ctx)

	if err := db.Create(message).Error; err != nil {
		r.Logger.Error("AddMessage failed: ", err)
		return 0, err
	}

	return message.ID, nil
}
