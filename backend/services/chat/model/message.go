package model

import (
	"time"
)

// Message 消息模型
type Message struct {
	ID             int64     `gorm:"primaryKey;column:id" json:"id"`
	ConversationID int64     `gorm:"column:conversation_id" json:"conversation_id"`
	Type           string    `gorm:"column:type" json:"type"` // 'user', 'ai'
	Content        string    `gorm:"column:content" json:"content"`
	AudioID        *int64    `gorm:"column:audio_id" json:"audio_id"`
	Metadata       *string   `gorm:"column:metadata" json:"metadata"` // JSON字符串
	TokenUsed      int32     `gorm:"column:token_used;default:0" json:"token_used"`
	ProcessingTime int32     `gorm:"column:processing_time;default:0" json:"processing_time"` // 毫秒
	CreatedAt      time.Time `gorm:"column:created_at" json:"created_at"`
}

// TableName 指定表名
func (Message) TableName() string {
	return "messages"
}
