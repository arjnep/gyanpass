package entity

import (
	"time"

	"github.com/google/uuid"
)

type Notification struct {
	ID        uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	UserID    uuid.UUID `gorm:"not null" json:"user_id" binding:"required"`
	Message   string    `gorm:"not null" json:"message" binding:"required"`
	Type      string    `json:"type"` // "info", "warning", "error"
	CreatedAt time.Time `gorm:"not null" json:"created_at"`
	IsRead    bool      `gorm:"default:false" json:"is_read"`
}
