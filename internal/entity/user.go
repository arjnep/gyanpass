package entity

import (
	"github.com/google/uuid"
)

type User struct {
	UID       uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"uid"`
	FirstName string    `gorm:"not null" json:"first_name" binding:"required"`
	LastName  string    `gorm:"not null" json:"last_name" binding:"required"`
	Email     string    `gorm:"unique;not null" json:"email,omitempty" binding:"required,email"`
	Phone     string    `gorm:"unique;not null" json:"phone,omitempty" binding:"required"`
	Password  string    `gorm:"not null" json:"-" binding:"required,min=8"`
	Role      string    `gorm:"default:user" json:"role,omitempty"` // "admin", "user"
}
