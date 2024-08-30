package repository

import (
	"github.com/arjnep/gyanpass/internal/entity"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserRepository interface {
	Create(user *entity.User) error
	FindByEmail(email string) (*entity.User, error)
	FindByPhone(phone string) (*entity.User, error)
	FindByID(id uuid.UUID) (*entity.User, error)
	Update(*entity.User, map[string]interface{}) error
	Delete(*entity.User) error
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db}
}

func (r *userRepository) Create(user *entity.User) error {
	return r.db.Create(user).Error
}

func (r *userRepository) FindByID(id uuid.UUID) (*entity.User, error) {
	var user entity.User
	err := r.db.First(&user, id).Error
	return &user, err
}

func (r *userRepository) FindByEmail(email string) (*entity.User, error) {
	var user entity.User
	err := r.db.Where("email = ?", email).First(&user).Error
	return &user, err
}

func (r *userRepository) FindByPhone(phone string) (*entity.User, error) {
	var user entity.User
	err := r.db.Where("phone = ?", phone).First(&user).Error
	return &user, err
}

func (r *userRepository) Update(user *entity.User, updates map[string]interface{}) error {
	return r.db.Model(user).Updates(updates).Error
}

func (r *userRepository) Delete(user *entity.User) error {
	return r.db.Delete(user).Error
}
