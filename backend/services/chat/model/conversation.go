package model

import (
	"time"
)

// Conversation 对话模型
type Conversation struct {
	ID          int64     `gorm:"primaryKey;column:id" json:"id"`
	UserID      *int64    `gorm:"column:user_id" json:"user_id"`
	CharacterID int64     `gorm:"column:character_id" json:"character_id"`
	Title       string    `gorm:"column:title" json:"title"`
	Status      int32     `gorm:"column:status" json:"status"`
	CreatedAt   time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt   time.Time `gorm:"column:updated_at" json:"updated_at"`
}

// TableName 指定表名
func (Conversation) TableName() string {
	return "conversations"
}
