package notification

import (
	"fmt"

	"github.com/arjnep/gyanpass/internal/entity"
	"github.com/arjnep/gyanpass/internal/repository"
	"github.com/arjnep/gyanpass/pkg/response"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Service interface {
	SendNotification(userID uuid.UUID, notificationType string, message string) error
	GetNotificationByID(id uuid.UUID) (*entity.Notification, error)
	GetUserNotifications(userID uuid.UUID) ([]entity.Notification, error)
	MarkNotificationAsRead(notification *entity.Notification, userID uuid.UUID) error
	RemoveNotification(notification *entity.Notification, userID uuid.UUID) error
}

type notificationService struct {
	// Todo: Add Email
	repo repository.NotificationRepository
}

func NewNotificationService(repo repository.NotificationRepository) Service {
	return &notificationService{repo}
}

func (s *notificationService) SendNotification(userID uuid.UUID, notificationType string, message string) error {
	notification := &entity.Notification{
		UserID:  userID,
		Type:    notificationType,
		Message: message,
	}

	return s.repo.Create(notification)
}

func (s *notificationService) GetNotificationByID(id uuid.UUID) (*entity.Notification, error) {
	notificationFetched, err := s.repo.GetByID(id)
	if err != nil && err == gorm.ErrRecordNotFound {
		return nil, response.NewNotFoundError("notification", fmt.Sprintf("%d", id))
	} else if err != nil && err != gorm.ErrRecordNotFound {
		return nil, response.NewInternalServerError()
	}
	return notificationFetched, nil
}

func (s *notificationService) GetUserNotifications(userID uuid.UUID) ([]entity.Notification, error) {
	return s.repo.GetByUserID(userID)
}

func (s *notificationService) MarkNotificationAsRead(notification *entity.Notification, userID uuid.UUID) error {
	if notification.UserID != userID {
		// return response.NewAuthorizationError("you do not have permission")
		return response.NewNotFoundError("notification", fmt.Sprintf("%v", notification.ID))
	}

	return s.repo.MarkAsRead(notification.ID)
}

func (s *notificationService) RemoveNotification(notification *entity.Notification, userID uuid.UUID) error {
	if notification.UserID != userID {
		// return response.NewAuthorizationError("you do not have permission")
		return response.NewNotFoundError("notification", fmt.Sprintf("%v", notification.ID))
	}
	return s.repo.Delete(notification)
}
