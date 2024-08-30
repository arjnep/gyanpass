package entity

import "github.com/google/uuid"

type ExchangeRequest struct {
	ID                   uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	RequestedByID        uuid.UUID `gorm:"not null" json:"requested_by_id" binding:"required"`
	RequestedToID        uuid.UUID `gorm:"not null" json:"requested_to_id" binding:"required"`
	RequestedBy          User      `gorm:"foreignKey:RequestedByID"`
	RequestedTo          User      `gorm:"foreignKey:RequestedToID"`
	RequestedBookID      uint      `gorm:"not null" json:"requested_book_id" binding:"required"`
	RequestedBook        Book      `gorm:"foreignKey:RequestedBookID"`
	OfferedBookID        uint      `gorm:"not null" json:"offered_book_id" binding:"required"`
	OfferedBook          Book      `gorm:"foreignKey:OfferedBookID"`
	Status               string    `gorm:"not null" json:"status"` // "pending", "accepted", "declined", "exchanged"
	RequestedByConfirmed bool      `json:"requested_by_confirmed"`
	RequestedToConfirmed bool      `json:"requested_to_confirmed"`
}
