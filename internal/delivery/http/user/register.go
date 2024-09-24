package user

import (
	"log"
	"net/http"
	"unicode"

	"github.com/arjnep/gyanpass/internal/entity"
	"github.com/arjnep/gyanpass/pkg/response"
	"github.com/arjnep/gyanpass/pkg/utils"
	"github.com/gin-gonic/gin"
)

type registerReq struct {
	FirstName string `gorm:"not null" json:"first_name" binding:"required"`
	LastName  string `gorm:"not null" json:"last_name" binding:"required"`
	Email     string `gorm:"unique;not null" json:"email" binding:"required,email"`
	Phone     string `gorm:"unique;not null" json:"phone" binding:"required,len=10"`
	Password  string `gorm:"not null" json:"password" binding:"required,min=8"`
}

func (h *UserHandler) RegisterUser(c *gin.Context) {
	var req registerReq
	if ok := utils.BindData(c, &req); !ok {
		return
	}

	if !isPasswordValid(req.Password) {
		err := response.NewBadRequestError("password must contain at least 1 uppercase, 1 lowercase, 1 alphanumeric, 1 number and should be above 8 character long")
		c.JSON(err.Status(), gin.H{
			"error": err,
		})
		return
	}

	user := &entity.User{
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Email:     req.Email,
		Phone:     req.Phone,
		Password:  req.Password,
		Role:      "user",
	}

	err := h.userUsecase.Register(user)
	if err != nil {
		log.Printf("Failed to register user: %v\n", err.Error())
		c.JSON(response.Status(err), gin.H{
			"error": err,
		})
		return
	}

	token, err := h.jwtService.GenerateToken(user)
	if err != nil {
		log.Printf("Failed to create tokens for user: %v\n", err.Error())
		c.JSON(response.Status(err), gin.H{
			"error": err,
		})
		return
	}

	c.SetCookie("token", token, 900, "/", "", true, true)
	c.SetCookie("uid", user.UID.String(), 900, "/", "", true, true)

	user.Password = ""
	c.JSON(http.StatusCreated, gin.H{
		"tokens": token,
		"user":   user,
	})

}

func isPasswordValid(s string) bool {
	var (
		hasMinLen  = false
		hasUpper   = false
		hasLower   = false
		hasNumber  = false
		hasSpecial = false
	)
	if len(s) >= 8 {
		hasMinLen = true
	}
	for _, char := range s {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsNumber(char):
			hasNumber = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}
	return hasMinLen && hasUpper && hasLower && hasNumber && hasSpecial
}
