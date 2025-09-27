package model

import (
	"time"
)

// Message 消息模型
type Message struct {
	ID             int64     `gorm:"primaryKey;column:id" json:"id"`
	ConversationID int64     `gorm:"column:conversation_id" json:"conversation_id"`
	UserID         *int64    `gorm:"column:user_id" json:"user_id"`
	CharacterID    *int64    `gorm:"column:character_id" json:"character_id"`
	Type           string    `gorm:"column:type" json:"type"` // 'user', 'ai'
	Content        string    `gorm:"column:content" json:"content"`
	AudioFileID    *int64    `gorm:"column:audio_file_id" json:"audio_file_id"`
	Status         int32     `gorm:"column:status" json:"status"`
	CreatedAt      time.Time `gorm:"column:created_at" json:"created_at"`
}

// TableName 指定表名
func (Message) TableName() string {
	return "messages"
}
