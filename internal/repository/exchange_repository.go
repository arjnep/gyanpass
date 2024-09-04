package repository

import (
	"github.com/arjnep/gyanpass/internal/entity"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ExchangeRepository interface {
	Create(exchangeRequest *entity.ExchangeRequest) error
	FindByID(id uuid.UUID) (*entity.ExchangeRequest, error)
	FindByRequestedByID(requestedByID uuid.UUID) ([]entity.ExchangeRequest, error)
	FindByRequestedToID(requestedToID uuid.UUID) ([]entity.ExchangeRequest, error)
	Update(exchangeRequest *entity.ExchangeRequest) error
	Delete(exchangeRequest *entity.ExchangeRequest) error
	CanRequest(userID, requestedToID uuid.UUID) (bool, error)
	IsSelfRequest(requestedByID, requestedToID uuid.UUID) bool
	FindPendingRequests(requestedByID, requestedToID uuid.UUID) ([]entity.ExchangeRequest, error)
	FindPendingRequestsByBookID(bookID uint) ([]entity.ExchangeRequest, error)
	FindRequestsByBookIDAndUserID(bookID uint, userID uuid.UUID) ([]entity.ExchangeRequest, error)
	FindRequestsByUserID(userID uuid.UUID) ([]entity.ExchangeRequest, error)
}

type exchangeRepository struct {
	db *gorm.DB
}

func NewExchangeRepository(db *gorm.DB) ExchangeRepository {
	return &exchangeRepository{db}
}

func (r *exchangeRepository) Create(exchangeRequest *entity.ExchangeRequest) error {
	return r.db.Create(exchangeRequest).Error
}

func (r *exchangeRepository) FindByID(id uuid.UUID) (*entity.ExchangeRequest, error) {
	var exchangeRequest entity.ExchangeRequest
	err := r.db.Preload("RequestedBy").Preload("RequestedTo").Preload("RequestedBook").Preload("RequestedBook.Owner").Preload("OfferedBook").Preload("OfferedBook.Owner").First(&exchangeRequest, id).Error
	return &exchangeRequest, err
}

func (r *exchangeRepository) FindByRequestedByID(userID uuid.UUID) ([]entity.ExchangeRequest, error) {
	var exchangeRequests []entity.ExchangeRequest
	err := r.db.Preload("RequestedBy").Preload("RequestedTo").Preload("RequestedBook").Preload("RequestedBook.Owner").Preload("OfferedBook").Preload("OfferedBook.Owner").Where("requested_by_id = ?", userID).Find(&exchangeRequests).Error
	return exchangeRequests, err
}

func (r *exchangeRepository) FindByRequestedToID(userID uuid.UUID) ([]entity.ExchangeRequest, error) {
	var exchangeRequests []entity.ExchangeRequest
	err := r.db.Preload("RequestedBy").Preload("RequestedTo").Preload("RequestedBook").Preload("RequestedBook.Owner").Preload("OfferedBook").Preload("OfferedBook.Owner").Where("requested_to_id = ?", userID).Find(&exchangeRequests).Error
	return exchangeRequests, err
}

func (r *exchangeRepository) CanRequest(requestedByID, requestedToID uuid.UUID) (bool, error) {
	var count int64
	err := r.db.Model(&entity.ExchangeRequest{}).Where("requested_by_id = ? AND requested_to_id = ? AND status IN (?, ?)",
		requestedByID, requestedToID, "pending", "accepted").
		Count(&count).Error
	if err != nil {
		return false, err
	}

	return count == 0, nil
}

func (r *exchangeRepository) IsSelfRequest(requestedByID, requestedToID uuid.UUID) bool {
	return requestedByID == requestedToID
}

func (r *exchangeRepository) FindPendingRequests(requestedByID, requestedToID uuid.UUID) ([]entity.ExchangeRequest, error) {
	var requests []entity.ExchangeRequest
	err := r.db.Where("requested_by_id = ? AND requested_to_id = ? AND status = ?", requestedByID, requestedToID, "pending").Find(&requests).Error
	return requests, err
}

func (r *exchangeRepository) FindPendingRequestsByBookID(bookID uint) ([]entity.ExchangeRequest, error) {
	var requests []entity.ExchangeRequest
	err := r.db.Where("(requested_book_id = ? OR offered_book_id = ?) AND status = ?", bookID, bookID, "pending").Find(&requests).Error
	return requests, err
}

func (r *exchangeRepository) Update(exchangeRequest *entity.ExchangeRequest) error {
	return r.db.Save(exchangeRequest).Error
}

func (r *exchangeRepository) Delete(exchangeRequest *entity.ExchangeRequest) error {
	return r.db.Delete(exchangeRequest).Error
}

func (r *exchangeRepository) FindRequestsByBookIDAndUserID(bookID uint, userID uuid.UUID) ([]entity.ExchangeRequest, error) {
	var requests []entity.ExchangeRequest
	err := r.db.Where("(requested_book_id = ? OR offered_book_id = ?) AND (requested_by_id = ? OR requested_to_id = ?)", bookID, bookID, userID, userID).Find(&requests).Error
	return requests, err
}

func (r *exchangeRepository) FindRequestsByUserID(userID uuid.UUID) ([]entity.ExchangeRequest, error) {
	var exchangeRequests []entity.ExchangeRequest
	err := r.db.Preload("RequestedBy").Preload("RequestedTo").Preload("RequestedBook").Preload("RequestedBook.Owner").Preload("OfferedBook").Preload("OfferedBook.Owner").
		Where("requested_by_id = ? OR requested_to_id = ?", userID, userID).Find(&exchangeRequests).Error
	return exchangeRequests, err
}
