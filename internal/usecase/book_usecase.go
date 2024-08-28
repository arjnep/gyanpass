package usecase

import (
	"github.com/arjnep/gyanpass/internal/entity"
	"github.com/arjnep/gyanpass/internal/repository"
	"github.com/google/uuid"
)

type BookUsecase interface {
	AddBook(book *entity.Book) error
	GetBookByID(id uint) (*entity.Book, error)
	GetBooksByUserID(uid uuid.UUID) ([]entity.Book, error)
	SearchBooks(queryParams map[string]string, page, size int) ([]entity.Book, int, error)
	UpdateBook(book *entity.Book, updates map[string]interface{}) error
	DeleteBook(book *entity.Book) error
}

type bookUsecase struct {
	bookRepo repository.BookRepository
}

func NewBookUsecase(bookRepo repository.BookRepository) BookUsecase {
	return &bookUsecase{bookRepo}
}

func (u *bookUsecase) AddBook(book *entity.Book) error {
	return u.bookRepo.Create(book)
}

func (u *bookUsecase) GetBookByID(id uint) (*entity.Book, error) {
	return u.bookRepo.FindByID(id)
}

func (u *bookUsecase) GetBooksByUserID(uid uuid.UUID) ([]entity.Book, error) {
	return u.bookRepo.FindByUserID(uid)
}

func (u *bookUsecase) SearchBooks(queryParams map[string]string, page, size int) ([]entity.Book, int, error) {
	books, total, err := u.bookRepo.FindByQueryParams(queryParams, page, size)
	if err != nil {
		return nil, 0, err
	}
	return books, total, nil
}

func (u *bookUsecase) UpdateBook(book *entity.Book, updates map[string]interface{}) error {
	return u.bookRepo.Update(book, updates)
}

func (u *bookUsecase) DeleteBook(book *entity.Book) error {
	return u.bookRepo.Delete(book)
}
