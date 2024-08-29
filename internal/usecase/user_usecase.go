package usecase

import (
	"log"

	"github.com/arjnep/gyanpass/internal/entity"
	"github.com/arjnep/gyanpass/internal/repository"
	"github.com/arjnep/gyanpass/pkg/crypto"
	"github.com/arjnep/gyanpass/pkg/jwt"
	"github.com/arjnep/gyanpass/pkg/response"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserUsecase interface {
	Register(user *entity.User) error
	Login(user *entity.User) error
	GetUserByID(uid uuid.UUID) (*entity.User, error)
	Update(user *entity.User, updates map[string]interface{}) error
	Delete(user *entity.User) error
}

type userUsecase struct {
	userRepo   repository.UserRepository
	jwtService jwt.Service
}

func NewUserUsecase(userRepo repository.UserRepository, jwtService jwt.Service) UserUsecase {
	return &userUsecase{userRepo, jwtService}
}

func (u *userUsecase) GetUserByID(uid uuid.UUID) (*entity.User, error) {
	userFetched, err := u.userRepo.FindByID(uid)
	if err != nil && err == gorm.ErrRecordNotFound {
		return nil, response.NewNotFoundError("user", userFetched.Email)
	} else if err != nil && err != gorm.ErrRecordNotFound {
		return nil, response.NewInternalServerError()
	}
	return userFetched, nil
}

func (u *userUsecase) Register(user *entity.User) error {
	userFetched, err := u.userRepo.FindByEmail(user.Email)
	if err == nil && userFetched != nil {
		return response.NewConflictError("user", user.Email)
	} else if err != nil && err != gorm.ErrRecordNotFound {
		return response.NewInternalServerError()
	}

	userFetched, err = u.userRepo.FindByPhone(user.Phone)
	if err == nil && userFetched != nil {
		return response.NewConflictError("user", user.Phone)
	} else if err != nil && err != gorm.ErrRecordNotFound {
		return response.NewInternalServerError()
	}

	hashedPwd, err := crypto.HashPassword(user.Password)
	if err != nil {
		log.Printf("Unable to signup user for email: %v\n", user.Email)
		return response.NewInternalServerError()
	}
	user.Password = hashedPwd

	return u.userRepo.Create(user)

}

func (u *userUsecase) Login(user *entity.User) error {
	userFetched, err := u.userRepo.FindByEmail(user.Email)
	if err != nil && err == gorm.ErrRecordNotFound {
		return response.NewNotFoundError("user", user.Email)
	} else if err != nil && err != gorm.ErrRecordNotFound {
		return response.NewInternalServerError()
	}

	match, err := crypto.ComparePasswords(userFetched.Password, user.Password)
	if err != nil {
		return response.NewInternalServerError()
	}

	if !match {
		return response.NewAuthorizationError("invalid email or password")
	}

	*user = *userFetched
	return nil
}

func (u *userUsecase) Update(user *entity.User, updates map[string]interface{}) error {
	err := u.userRepo.Update(user, updates)
	if err != nil {
		return err
	}

	return nil
}

func (u *userUsecase) Delete(user *entity.User) error {
	err := u.userRepo.Delete(user)
	if err != nil {
		return err
	}
	return nil
}
