package public

import (
	"ai-roleplay/services/character/api/internal/config"
	"ai-roleplay/services/character/api/internal/svc"
	"ai-roleplay/services/character/api/internal/types"
	"context"
	"testing"

	"github.com/zeromicro/go-zero/core/conf"
)

var (
	svcCtx *svc.ServiceContext
	ctx    context.Context
	// logic
	getCharacterCategoriesLogic   *GetCharacterCategoriesLogic
	getCharacterDetailLogic       *GetCharacterDetailLogic
	getCharacterListLogic         *GetCharacterListLogic
	getCharacterTagsLogic         *GetCharacterTagsLogic
	getPopularCharactersLogic     *GetPopularCharactersLogic
	getRecommendedCharactersLogic *GetRecommendedCharactersLogic
	searchCharacterLogic          *SearchCharactersLogic
)

func init() {
	var c config.Config
	conf.MustLoad("../../../etc/character-api.yaml", &c)
	svcCtx = svc.NewServiceContext(c)
	ctx = context.Background()
}

func TestGetCharacterCategoriesLogic(t *testing.T) {
	getCharacterCategoriesLogic = NewGetCharacterCategoriesLogic(ctx, svcCtx)
	resp, err := getCharacterCategoriesLogic.GetCharacterCategories()
	if err != nil {
		t.Fatalf("GetCharacterCategories failed: %v", err)
	}
	t.Logf("GetCharacterCategories resp: %v\n", resp)
}

func TestGetCharacterDetailLogic(t *testing.T) {
	getCharacterDetailLogic = NewGetCharacterDetailLogic(ctx, svcCtx)
	resp, err := getCharacterDetailLogic.GetCharacterDetail(&types.CharacterDetailRequest{
		ID: 1,
	})
	if err != nil {
		t.Fatalf("GetCharacterDetail failed: %v", err)
	}
	t.Logf("GetCharacterDetail resp: %v\n", resp)
}

func TestGetCharacterListLogic(t *testing.T) {
	getCharacterListLogic = NewGetCharacterListLogic(ctx, svcCtx)
	resp, err := getCharacterListLogic.GetCharacterList(&types.CharacterListRequest{
		Page: 1,
	})
	if err != nil {
		t.Fatalf("GetCharacterList failed: %v", err)
	}
	t.Logf("GetCharacterList resp: %v\n", resp)
}

func TestGetCharacterTagsLogic(t *testing.T) {
	getCharacterTagsLogic = NewGetCharacterTagsLogic(ctx, svcCtx)
	resp, err := getCharacterTagsLogic.GetCharacterTags()
	if err != nil {
		t.Fatalf("GetCharacterTags failed: %v", err)
	}
	t.Logf("GetCharacterTags resp: %v\n", resp)
}

func TestGetPopularCharactersLogic(t *testing.T) {
	getPopularCharactersLogic = NewGetPopularCharactersLogic(ctx, svcCtx)
	resp, err := getPopularCharactersLogic.GetPopularCharacters(&types.PopularCharacterRequest{
		Page: 1,
	})
	if err != nil {
		t.Fatalf("GetPopularCharacters failed: %v", err)
	}
	t.Logf("GetPopularCharacters resp: %v\n", resp)
}

func TestGetRecommendedCharactersLogic(t *testing.T) {
	getRecommendedCharactersLogic = NewGetRecommendedCharactersLogic(ctx, svcCtx)
	resp, err := getRecommendedCharactersLogic.GetRecommendedCharacters(&types.RecommendedCharacterRequest{
		Count: 10,
	})
	if err != nil {
		t.Fatalf("GetRecommendedCharacters failed: %v", err)
	}
	t.Logf("GetRecommendedCharacters resp: %v\n", resp)
}

func TestSearchCharactersLogic(t *testing.T) {
	searchCharacterLogic = NewSearchCharactersLogic(ctx, svcCtx)
	resp, err := searchCharacterLogic.SearchCharacters(&types.SearchCharacterRequest{
		Page:     1,
		PageSize: 10,
		Keyword:  "莎士比亚",
	})
	if err != nil {
		t.Fatalf("SearchCharacters failed: %v", err)
	}
	t.Logf("SearchCharacters resp: %v\n", resp)
}
