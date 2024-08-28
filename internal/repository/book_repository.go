package repository

import (
	"github.com/arjnep/gyanpass/internal/entity"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BookRepository interface {
	Create(book *entity.Book) error
	FindByID(id uint) (*entity.Book, error)
	FindByUserID(uid uuid.UUID) ([]entity.Book, error)
	FindByQueryParams(queryParams map[string]string, page, size int) ([]entity.Book, int, error)
	Update(book *entity.Book, updates map[string]interface{}) error
	Delete(book *entity.Book) error
}

type bookRepository struct {
	db *gorm.DB
}

func NewBookRepository(db *gorm.DB) BookRepository {
	return &bookRepository{db}
}

func (r *bookRepository) Create(book *entity.Book) error {
	return r.db.Create(book).Error
}

func (r *bookRepository) FindByID(id uint) (*entity.Book, error) {
	var book entity.Book
	err := r.db.Preload("Owner").First(&book, id).Error
	if err != nil {
		return nil, err
	}
	return &book, nil
}

func (r *bookRepository) FindByUserID(uid uuid.UUID) ([]entity.Book, error) {
	var books []entity.Book
	err := r.db.Where("uid = ?", uid).Find(&books).Error
	if err != nil {
		return nil, err
	}
	return books, err
}

func (r *bookRepository) FindByQueryParams(queryParams map[string]string, page, size int) ([]entity.Book, int, error) {
	var books []entity.Book
	var total int64

	query := r.db.Model(&entity.Book{})
	for key, value := range queryParams {
		if value != "" {
			switch key {
			case "title":
				query = query.Where("title ILIKE ?", "%"+value+"%")
			case "address":
				query = query.Where("address ILIKE ?", "%"+value+"%")
			}
		}
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * size
	query = query.Limit(size).Offset(offset).Preload("Owner")

	if err := query.Find(&books).Error; err != nil {
		return nil, 0, err
	}

	return books, int(total), nil
}

func (r *bookRepository) Update(book *entity.Book, updates map[string]interface{}) error {
	return r.db.Model(book).Updates(updates).Error
}

func (r *bookRepository) Delete(book *entity.Book) error {
	return r.db.Delete(book).Error
}
