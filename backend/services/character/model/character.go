package model

import (
	"encoding/json"
	"time"
)

type Character struct {
	ID            int64     `gorm:"primaryKey;column:id" json:"id"`
	Name          string    `gorm:"column:name" json:"name"`
	Avatar        *string   `gorm:"column:avatar" json:"avatar"`
	Description   *string   `gorm:"column:description" json:"description"`
	ShortDesc     *string   `gorm:"column:short_desc" json:"short_desc"`
	CategoryID    *int64    `gorm:"column:category_id" json:"category_id"`
	Tags          *string   `gorm:"column:tags" json:"tags"`
	Prompt        *string   `gorm:"column:prompt" json:"prompt"`
	Personality   *string   `gorm:"column:personality" json:"personality"`
	VoiceSettings *string   `gorm:"column:voice_settings" json:"voice_settings"`
	Status        int32     `gorm:"column:status" json:"status"`
	IsPublic      int32     `gorm:"column:is_public" json:"is_public"`
	CreatorID     *int64    `gorm:"column:creator_id" json:"creator_id"`
	Rating        float64   `gorm:"column:rating" json:"rating"`
	RatingCount   int32     `gorm:"column:rating_count" json:"rating_count"`
	FavoriteCount int32     `gorm:"column:favorite_count" json:"favorite_count"`
	ChatCount     int32     `gorm:"column:chat_count" json:"chat_count"`
	CreatedAt     time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt     time.Time `gorm:"column:updated_at" json:"updated_at"`
}

// TableName 指定表名
func (Character) TableName() string {
	return "characters"
}

func (c *Character) GetTags() []string {
	if c.Tags == nil {
		return []string{}
	}

	var tags []string
	json.Unmarshal([]byte(*c.Tags), &tags)
	return tags
}
