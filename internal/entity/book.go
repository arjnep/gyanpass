package entity

import "github.com/google/uuid"

type Book struct {
	ID             uint      `gorm:"not null;primaryKey" json:"id"`
	Title          string    `gorm:"not null" json:"title" binding:"required"`
	Author         string    `gorm:"not null" json:"author" binding:"required"`
	Genre          string    `json:"genre" binding:"omitempty"`
	Description    string    `json:"description" binding:"omitempty"`
	UserID         uuid.UUID `gorm:"not null" json:"user_id,omitempty"`
	Owner          User      `gorm:"foreignKey:UserID" json:"owner"`
	PickupLocation Location  `gorm:"embedded" json:"location,omitempty" binding:"required"`
}

type Location struct {
	Address   string  `json:"address" binding:"omitempty"`
	Latitude  float64 `gorm:"not null" json:"latitude" binding:"required"`
	Longitude float64 `gorm:"not null" json:"longitude" binding:"required"`
}
