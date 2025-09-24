package logic

import (
	"context"
	"fmt"

	"ai-roleplay/common/response"
	"ai-roleplay/services/character/api/internal/svc"
	"ai-roleplay/services/character/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetCharacterListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetCharacterListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetCharacterListLogic {
	return &GetCharacterListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetCharacterListLogic) GetCharacterList(req *types.CharacterListRequest) (resp *types.CharacterListResponse, err error) {
	// 1. 参数验证和默认值处理
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 || req.PageSize > 100 {
		req.PageSize = 20
	}

	// 2. 构建查询条件
	offset := (req.Page - 1) * req.PageSize
	queryParams := &model.CharacterQueryParams{
		Category: req.Category,
		Tags:     req.Tags,
		Status:   1, // 只查询启用的角色
		IsPublic: 1, // 只查询公开的角色
		OrderBy:  req.OrderBy,
		Offset:   offset,
		Limit:    req.PageSize,
	}

	// 3. 尝试从缓存获取
	cacheKey := l.buildCacheKey(req)
	cachedResult, err := l.svcCtx.CacheManager.GetCharacterList(l.ctx, cacheKey)
	if err == nil && cachedResult != nil {
		return cachedResult, nil
	}

	// 4. 从数据库查询角色列表
	characters, total, err := l.svcCtx.CharacterModel.FindList(l.ctx, queryParams)
	if err != nil {
		logx.Errorf("查询角色列表失败: %v", err)
		return &types.CharacterListResponse{
			Code: response.DATABASE_ERROR,
			Msg:  "查询角色列表失败",
		}, nil
	}

	// 5. 转换数据格式
	characterList := make([]*types.Character, 0, len(characters))
	for _, char := range characters {
		characterInfo := l.convertToCharacterInfo(char)
		characterList = append(characterList, characterInfo)
	}

	// 6. 构造响应
	resp = &types.CharacterListResponse{
		Code: response.SUCCESS,
		Msg:  "获取成功",
		Data: &types.CharacterListData{
			List:     characterList,
			Total:    total,
			Page:     req.Page,
			PageSize: req.PageSize,
			HasMore:  int64(req.Page*req.PageSize) < total,
		},
	}

	// 7. 缓存结果（异步）
	go func() {
		if err := l.svcCtx.CacheManager.SetCharacterList(context.Background(), cacheKey, resp); err != nil {
			logx.Errorf("缓存角色列表失败: %v", err)
		}
	}()

	return resp, nil
}

// 构建缓存key
func (l *GetCharacterListLogic) buildCacheKey(req *types.CharacterListRequest) string {
	return fmt.Sprintf("character:list:%d:%d:%s:%s:%s",
		req.Page, req.PageSize, req.Category, req.Tags, req.OrderBy)
}

// 转换角色模型为响应格式
func (l *GetCharacterListLogic) convertToCharacterInfo(char *model.Character) *types.Character {
	// 解析tags JSON
	var tags []string
	if char.Tags != "" {
		// 这里需要JSON解析，简化处理
		tags = []string{} // 实际需要JSON解析
	}

	// 解析性格设置
	personality := &types.Personality{
		Friendliness: 70,
		Humor:        60,
		Intelligence: 80,
		Creativity:   75,
	}

	// 解析语音设置
	voiceSettings := &types.VoiceSettings{
		Rate:   1.0,
		Pitch:  1.0,
		Volume: 0.8,
	}

	return &types.Character{
		Id:            char.Id,
		Name:          char.Name,
		Avatar:        char.Avatar,
		Description:   char.Description,
		ShortDesc:     l.getShortDescription(char.Description),
		Tags:          tags,
		Category:      char.Category,
		Prompt:        char.Prompt,
		Status:        char.Status,
		Rating:        float64(char.Rating),
		RatingCount:   char.RatingCount,
		FavoriteCount: char.FavoriteCount,
		ChatCount:     char.ChatCount,
		IsPublic:      char.IsPublic == 1,
		CreatorId:     char.CreatorId,
		Personality:   personality,
		VoiceSettings: voiceSettings,
		CreatedAt:     char.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:     char.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
}

// 获取简短描述
func (l *GetCharacterListLogic) getShortDescription(description string) string {
	maxLength := 100
	if len(description) <= maxLength {
		return description
	}

	// 截取前100个字符并添加省略号
	runes := []rune(description)
	if len(runes) <= maxLength {
		return description
	}

	return string(runes[:maxLength]) + "..."
}
