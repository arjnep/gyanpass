package entity

import "github.com/google/uuid"

type Book struct {
	ID             uint        `gorm:"not null;primaryKey" json:"id"`
	Title          string      `gorm:"not null" json:"title" binding:"required"`
	Author         string      `gorm:"not null" json:"author" binding:"required"`
	Genre          string      `json:"genre" binding:"omitempty"`
	Description    Description `gorm:"embedded" json:"description" binding:"required"`
	ImageUrl       string      `gorm:"not null" json:"image_url" binding:"required"`
	UserID         uuid.UUID   `gorm:"not null" json:"user_id,omitempty"`
	Owner          User        `gorm:"foreignKey:UserID" json:"owner"`
	PickupLocation Location    `gorm:"embedded" json:"location,omitempty" binding:"required"`
	IsActive       bool        `json:"is_active"`
}

type Description struct {
	Message           string `json:"message" binding:"omitempty"`
	Condition         string `gorm:"not null" json:"condition" binding:"required"`
	PreferredExchange string `gorm:"not null" json:"preferred_exchange" binding:"required"`
}

type Location struct {
	Address   string  `json:"address" binding:"omitempty"`
	Latitude  float64 `gorm:"not null" json:"latitude,omitempty" binding:"required"`
	Longitude float64 `gorm:"not null" json:"longitude,omitempty" binding:"required"`
}
