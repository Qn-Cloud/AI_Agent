package repo

import (
	"ai-roleplay/services/character/api/internal/svc"
	"ai-roleplay/services/character/api/internal/types"
	"ai-roleplay/services/character/model"
	"context"
	"fmt"
	"strings"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

type CharacterServiceRepo struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCharacterServiceRepo(ctx context.Context, svcCtx *svc.ServiceContext) *CharacterServiceRepo {
	return &CharacterServiceRepo{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// GetCharacterList 获取角色列表
func (r *CharacterServiceRepo) GetCharacterList(req *types.CharacterListRequest) ([]model.Character, int64, error) {
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
	query := db.Model(&model.Character{}).Where("status = ? AND is_public = ?", 1, 1)

	// 分类筛选
	if req.CategoryID > 0 {
		query = query.Where("category_id = ?", req.CategoryID)
	}

	// 标签筛选
	if req.Tags != "" {
		tags := strings.Split(req.Tags, ",")
		for _, tag := range tags {
			tag = strings.TrimSpace(tag)
			if tag != "" {
				query = query.Where("JSON_CONTAINS(tags, ?)", fmt.Sprintf(`"%s"`, tag))
			}
		}
	}

	// 关键词搜索
	if req.Keyword != "" {
		keyword := "%" + req.Keyword + "%"
		query = query.Where("name LIKE ? OR description LIKE ? OR short_desc LIKE ?", keyword, keyword, keyword)
	}

	// 统计总数
	var total int64
	if err := query.Count(&total).Error; err != nil {
		r.Logger.Error("GetCharacterList count failed: ", err)
		return nil, 0, err
	}

	// 排序
	orderBy := req.OrderBy
	if orderBy == "" {
		orderBy = "created_at"
	}
	if req.OrderDesc {
		orderBy += " DESC"
	} else {
		orderBy += " ASC"
	}

	// 查询数据
	var characters []model.Character
	if err := query.Order(orderBy).Offset(offset).Limit(pageSize).Find(&characters).Error; err != nil {
		r.Logger.Error("GetCharacterList find failed: ", err)
		return nil, 0, err
	}

	return characters, total, nil
}

// SearchCharacters 搜索角色
func (r *CharacterServiceRepo) SearchCharacters(req *types.SearchCharacterRequest) ([]model.Character, int64, error) {
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
	query := db.Model(&model.Character{}).Where("status = ? AND is_public = ?", 1, 1)

	if req.Keyword != "" {
		keyword := "%" + req.Keyword + "%"
		query = query.Where("name LIKE ? OR description LIKE ? OR short_desc LIKE ?", keyword, keyword, keyword)
	}

	// 统计总数
	var total int64
	if err := query.Count(&total).Error; err != nil {
		r.Logger.Error("SearchCharacters count failed: ", err)
		return nil, 0, err
	}

	// 查询数据，按相关度排序
	var characters []model.Character
	if err := query.Order("chat_count DESC, favorite_count DESC, created_at DESC").
		Offset(offset).Limit(pageSize).Find(&characters).Error; err != nil {
		r.Logger.Error("SearchCharacters find failed: ", err)
		return nil, 0, err
	}

	return characters, total, nil
}

// GetCharacterByID 根据ID获取角色详情
func (r *CharacterServiceRepo) GetCharacterByID(id int64) (*model.Character, error) {
	db := r.svcCtx.Db.WithContext(r.ctx)

	var character model.Character
	if err := db.Where("id = ? AND status = ?", id, 1).First(&character).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		r.Logger.Error("GetCharacterByID failed: ", err)
		return nil, err
	}

	return &character, nil
}

// GetRecommendedCharacters 获取推荐角色
func (r *CharacterServiceRepo) GetRecommendedCharacters(count int) ([]model.Character, error) {
	db := r.svcCtx.Db.WithContext(r.ctx)

	if count <= 0 {
		count = 10
	}

	var characters []model.Character
	// 推荐逻辑：按评分、收藏数、对话数综合排序
	if err := db.Where("status = ? AND is_public = ?", 1, 1).
		Order("(rating * 0.4 + favorite_count * 0.3 + chat_count * 0.3) DESC, created_at DESC").
		Limit(count).Find(&characters).Error; err != nil {
		r.Logger.Error("GetRecommendedCharacters failed: ", err)
		return nil, err
	}

	return characters, nil
}

// GetCharacterCategories 获取角色分类列表
func (r *CharacterServiceRepo) GetCharacterCategories() ([]string, error) {
	db := r.svcCtx.Db.WithContext(r.ctx)

	var categories []string
	if err := db.Model(&model.Character{}).
		Select("DISTINCT(SELECT name FROM character_categories WHERE id = characters.category_id) as category_name").
		Where("status = ? AND is_public = ? AND category_id IS NOT NULL", 1, 1).
		Pluck("category_name", &categories).Error; err != nil {
		r.Logger.Error("GetCharacterCategories failed: ", err)
		return nil, err
	}

	// 过滤空值
	var result []string
	for _, cat := range categories {
		if cat != "" {
			result = append(result, cat)
		}
	}

	return result, nil
}

// GetPopularCharacters 获取热门角色
func (r *CharacterServiceRepo) GetPopularCharacters(req *types.PopularCharacterRequest) ([]model.Character, int64, error) {
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

	query := db.Model(&model.Character{}).Where("status = ? AND is_public = ?", 1, 1)

	// 统计总数
	var total int64
	if err := query.Count(&total).Error; err != nil {
		r.Logger.Error("GetPopularCharacters count failed: ", err)
		return nil, 0, err
	}

	// 按热度排序：对话数 > 收藏数 > 评分
	var characters []model.Character
	if err := query.Order("chat_count DESC, favorite_count DESC, rating DESC").
		Offset(offset).Limit(pageSize).Find(&characters).Error; err != nil {
		r.Logger.Error("GetPopularCharacters find failed: ", err)
		return nil, 0, err
	}

	return characters, total, nil
}

// GetCharacterTags 获取所有角色标签
func (r *CharacterServiceRepo) GetCharacterTags() ([]string, error) {
	db := r.svcCtx.Db.WithContext(r.ctx)

	var characters []model.Character
	if err := db.Select("tags").Where("status = ? AND is_public = ? AND tags IS NOT NULL", 1, 1).
		Find(&characters).Error; err != nil {
		r.Logger.Error("GetCharacterTags failed: ", err)
		return nil, err
	}

	// 收集所有标签
	tagSet := make(map[string]bool)
	for _, char := range characters {
		tags := char.GetTags()
		for _, tag := range tags {
			if tag != "" {
				tagSet[tag] = true
			}
		}
	}

	// 转换为切片
	var result []string
	for tag := range tagSet {
		result = append(result, tag)
	}

	return result, nil
}

// CreateCharacter 创建角色
func (r *CharacterServiceRepo) CreateCharacter(character *model.Character) error {
	db := r.svcCtx.Db.WithContext(r.ctx)

	if err := db.Create(character).Error; err != nil {
		r.Logger.Error("CreateCharacter failed: ", err)
		return err
	}

	return nil
}

// UpdateCharacter 更新角色
func (r *CharacterServiceRepo) UpdateCharacter(character *model.Character) error {
	db := r.svcCtx.Db.WithContext(r.ctx)

	if err := db.Save(character).Error; err != nil {
		r.Logger.Error("UpdateCharacter failed: ", err)
		return err
	}

	return nil
}

// DeleteCharacter 删除角色（软删除）
func (r *CharacterServiceRepo) DeleteCharacter(id, creatorID int64) error {
	db := r.svcCtx.Db.WithContext(r.ctx)

	// 软删除：更新状态为2（禁用）
	if err := db.Model(&model.Character{}).
		Where("id = ? AND creator_id = ?", id, creatorID).
		Update("status", 2).Error; err != nil {
		r.Logger.Error("DeleteCharacter failed: ", err)
		return err
	}

	return nil
}

// ToggleFavorite 切换收藏状态
func (r *CharacterServiceRepo) ToggleFavorite(userID, characterID int64) (bool, error) {
	db := r.svcCtx.Db.WithContext(r.ctx)

	// 检查是否已收藏
	var count int64
	if err := db.Table("character_favorites").
		Where("user_id = ? AND character_id = ?", userID, characterID).
		Count(&count).Error; err != nil {
		r.Logger.Error("ToggleFavorite check failed: ", err)
		return false, err
	}

	if count > 0 {
		// 已收藏，取消收藏
		if err := db.Exec("DELETE FROM character_favorites WHERE user_id = ? AND character_id = ?",
			userID, characterID).Error; err != nil {
			r.Logger.Error("ToggleFavorite delete failed: ", err)
			return false, err
		}

		// 更新收藏数
		if err := db.Model(&model.Character{}).
			Where("id = ?", characterID).
			UpdateColumn("favorite_count", gorm.Expr("favorite_count - 1")).Error; err != nil {
			r.Logger.Error("ToggleFavorite update count failed: ", err)
		}

		return false, nil
	} else {
		// 未收藏，添加收藏
		if err := db.Exec("INSERT INTO character_favorites (user_id, character_id, created_at) VALUES (?, ?, NOW())",
			userID, characterID).Error; err != nil {
			r.Logger.Error("ToggleFavorite insert failed: ", err)
			return false, err
		}

		// 更新收藏数
		if err := db.Model(&model.Character{}).
			Where("id = ?", characterID).
			UpdateColumn("favorite_count", gorm.Expr("favorite_count + 1")).Error; err != nil {
			r.Logger.Error("ToggleFavorite update count failed: ", err)
		}

		return true, nil
	}
}

// GetMyFavorites 获取我的收藏角色
func (r *CharacterServiceRepo) GetMyFavorites(userID int64, req *types.FavoriteCharacterRequest) ([]model.Character, int64, error) {
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

	// 构建查询：通过favorites表关联
	query := db.Model(&model.Character{}).
		Joins("INNER JOIN character_favorites ON characters.id = character_favorites.character_id").
		Where("character_favorites.user_id = ? AND characters.status = ?", userID, 1)

	// 获取总数
	var total int64
	if err := query.Count(&total).Error; err != nil {
		r.Logger.Error("GetMyFavorites count failed: ", err)
		return nil, 0, err
	}

	// 获取数据
	var characters []model.Character
	if err := query.Offset(offset).Limit(pageSize).
		Order("character_favorites.created_at DESC").
		Find(&characters).Error; err != nil {
		r.Logger.Error("GetMyFavorites find failed: ", err)
		return nil, 0, err
	}

	return characters, total, nil
}

// GetMyCharacters 获取我创建的角色
func (r *CharacterServiceRepo) GetMyCharacters(userID int64, req *types.MyCharacterRequest) ([]model.Character, int64, error) {
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

	// 构建查询
	query := db.Model(&model.Character{}).
		Where("creator_id = ? AND status != ?", userID, 2) // 排除已删除的

	// 状态筛选
	if req.Status > 0 {
		query = query.Where("status = ?", req.Status)
	}

	// 公开性筛选
	if req.IsPublic != nil {
		isPublicValue := 0
		if *req.IsPublic {
			isPublicValue = 1
		}
		query = query.Where("is_public = ?", isPublicValue)
	}

	// 获取总数
	var total int64
	if err := query.Count(&total).Error; err != nil {
		r.Logger.Error("GetMyCharacters count failed: ", err)
		return nil, 0, err
	}

	// 获取数据
	var characters []model.Character
	if err := query.Offset(offset).Limit(pageSize).
		Order("created_at DESC").
		Find(&characters).Error; err != nil {
		r.Logger.Error("GetMyCharacters find failed: ", err)
		return nil, 0, err
	}

	return characters, total, nil
}

// UpdatePrompt 更新提示词
func (r *CharacterServiceRepo) UpdatePrompt(id, creatorID int64, prompt string) error {
	db := r.svcCtx.Db.WithContext(r.ctx)

	if err := db.Model(&model.Character{}).
		Where("id = ? AND creator_id = ?", id, creatorID).
		Update("prompt", prompt).Error; err != nil {
		r.Logger.Error("UpdatePrompt failed: ", err)
		return err
	}

	return nil
}

// UpdatePersonality 更新性格设置
func (r *CharacterServiceRepo) UpdatePersonality(id, creatorID int64, personality string) error {
	db := r.svcCtx.Db.WithContext(r.ctx)

	if err := db.Model(&model.Character{}).
		Where("id = ? AND creator_id = ?", id, creatorID).
		Update("personality", personality).Error; err != nil {
		r.Logger.Error("UpdatePersonality failed: ", err)
		return err
	}

	return nil
}

// UpdateVoiceSettings 更新语音设置
func (r *CharacterServiceRepo) UpdateVoiceSettings(id, creatorID int64, voiceSettings string) error {
	db := r.svcCtx.Db.WithContext(r.ctx)

	if err := db.Model(&model.Character{}).
		Where("id = ? AND creator_id = ?", id, creatorID).
		Update("voice_settings", voiceSettings).Error; err != nil {
		r.Logger.Error("UpdateVoiceSettings failed: ", err)
		return err
	}

	return nil
}
