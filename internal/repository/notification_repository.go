package repository

import (
	"log"

	"github.com/arjnep/gyanpass/internal/entity"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type NotificationRepository interface {
	Create(notification *entity.Notification) error
	GetByID(id uuid.UUID) (*entity.Notification, error)
	GetByUserID(userID uuid.UUID) ([]entity.Notification, error)
	MarkAsRead(id uuid.UUID) error
	Delete(notification *entity.Notification) error
}

type notificationRepository struct {
	db *gorm.DB
}

func NewNotificationRepository(db *gorm.DB) NotificationRepository {
	return &notificationRepository{db}
}

func (r *notificationRepository) Create(notification *entity.Notification) error {
	log.Println("Checkpoint 3:", notification)
	return r.db.Create(notification).Error
}

func (r *notificationRepository) GetByID(id uuid.UUID) (*entity.Notification, error) {
	var notification entity.Notification
	err := r.db.First(&notification, id).Error
	if err != nil {
		return nil, err
	}
	return &notification, err
}

func (r *notificationRepository) GetByUserID(userID uuid.UUID) ([]entity.Notification, error) {
	var notifications []entity.Notification
	err := r.db.Where("user_id = ?", userID).Order("created_at desc").Find(&notifications).Error
	if err != nil {
		return nil, err
	}
	return notifications, err
}

func (r *notificationRepository) MarkAsRead(id uuid.UUID) error {
	return r.db.Model(&entity.Notification{}).Where("id = ?", id).Update("is_read", true).Error
}

func (r *notificationRepository) Delete(notification *entity.Notification) error {
	return r.db.Delete(notification).Error
}
