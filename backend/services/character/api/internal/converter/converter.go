package converter

import (
	"ai-roleplay/services/character/api/internal/types"
	"ai-roleplay/services/character/model"
	"encoding/json"
)

// CharacterConverter 角色转换器
type CharacterConverter struct{}

// NewCharacterConverter 创建角色转换器
func NewCharacterConverter() *CharacterConverter {
	return &CharacterConverter{}
}

// ToCharacterBrief 将数据库模型转换为简要信息类型
func (c *CharacterConverter) ToCharacterBrief(character *model.Character) *types.CharacterBrief {
	if character == nil {
		return nil
	}

	// 解析标签
	tags := character.GetTags()

	// 处理可选字段
	avatar := ""
	if character.Avatar != nil {
		avatar = *character.Avatar
	}

	shortDesc := ""
	if character.ShortDesc != nil {
		shortDesc = *character.ShortDesc
	}

	return &types.CharacterBrief{
		ID:            character.ID,
		Name:          character.Name,
		Avatar:        avatar,
		ShortDesc:     shortDesc,
		CategoryName:  "", // 这里需要从关联查询中获取
		Tags:          tags,
		Rating:        character.Rating,
		RatingCount:   character.RatingCount,
		FavoriteCount: character.FavoriteCount,
		ChatCount:     character.ChatCount,
		IsPublic:      character.IsPublic == 1,
	}
}

// ToCharacterBriefList 将数据库模型列表转换为简要信息列表
func (c *CharacterConverter) ToCharacterBriefList(characters []model.Character) []types.CharacterBrief {
	result := make([]types.CharacterBrief, 0, len(characters))
	for _, character := range characters {
		if brief := c.ToCharacterBrief(&character); brief != nil {
			result = append(result, *brief)
		}
	}
	return result
}

// ToCharacterItem 将数据库模型转换为详细信息类型
func (c *CharacterConverter) ToCharacterItem(character *model.Character) *types.CharacterItem {
	if character == nil {
		return nil
	}

	// 解析JSON字段
	var personality types.CharacterPersonality
	if character.Personality != nil {
		json.Unmarshal([]byte(*character.Personality), &personality)
	}

	var voiceSettings types.CharacterVoiceSettings
	if character.VoiceSettings != nil {
		json.Unmarshal([]byte(*character.VoiceSettings), &voiceSettings)
	}

	// 解析标签
	tags := character.GetTags()

	// 处理可选字段
	avatar := ""
	if character.Avatar != nil {
		avatar = *character.Avatar
	}

	description := ""
	if character.Description != nil {
		description = *character.Description
	}

	shortDesc := ""
	if character.ShortDesc != nil {
		shortDesc = *character.ShortDesc
	}

	prompt := ""
	if character.Prompt != nil {
		prompt = *character.Prompt
	}

	categoryID := int64(0)
	if character.CategoryID != nil {
		categoryID = *character.CategoryID
	}

	creatorID := int64(0)
	if character.CreatorID != nil {
		creatorID = *character.CreatorID
	}

	return &types.CharacterItem{
		ID:            character.ID,
		Name:          character.Name,
		Avatar:        avatar,
		Description:   description,
		ShortDesc:     shortDesc,
		CategoryID:    categoryID,
		CategoryName:  "",
		Tags:          tags,
		Prompt:        prompt,
		Personality:   personality,
		VoiceSettings: voiceSettings,
		Status:        character.Status,
		IsPublic:      character.IsPublic == 1,
		CreatorID:     creatorID,
		CreatorName:   "",
		Rating:        character.Rating,
		RatingCount:   character.RatingCount,
		FavoriteCount: character.FavoriteCount,
		ChatCount:     character.ChatCount,
	}
}

// BuildPagination 构建分页信息
func (c *CharacterConverter) BuildPagination(page, pageSize int, total int64) *types.Pagination {

	totalPage := int(total) / pageSize
	if int(total)%pageSize > 0 {
		totalPage++
	}

	return &types.Pagination{
		TotalCount: int(total),
		TotalPage:  totalPage,
		Page:       page,
		PageSize:   pageSize,
	}
}

// BuildCharacterListResponse 构建角色列表响应
func (c *CharacterConverter) BuildCharacterListResponse(
	characters []model.Character,
	total int64,
	page, pageSize int,
) *types.CharacterListResponse {
	briefList := c.ToCharacterBriefList(characters)
	pagination := c.BuildPagination(page, pageSize, total)

	return &types.CharacterListResponse{
		Code:  0,
		Msg:   "success",
		Total: total,
		Page:  pagination,
		List:  briefList,
	}
}

// BuildCharacterDetailResponse 构建角色详情响应
func (c *CharacterConverter) BuildCharacterDetailResponse(character *model.Character) *types.CharacterDetailResponse {
	if character == nil {
		return &types.CharacterDetailResponse{
			Code: 404,
			Msg:  "角色不存在",
		}
	}

	item := c.ToCharacterItem(character)
	return &types.CharacterDetailResponse{
		Code:      0,
		Msg:       "success",
		Character: *item,
	}
}
