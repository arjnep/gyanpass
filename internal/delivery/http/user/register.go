package user

import (
	"log"
	"net/http"

	"github.com/arjnep/gyanpass/internal/entity"
	"github.com/arjnep/gyanpass/pkg/response"
	"github.com/arjnep/gyanpass/pkg/utils"
	"github.com/gin-gonic/gin"
)

type registerReq struct {
	FirstName string `gorm:"not null" json:"first_name" binding:"required"`
	LastName  string `gorm:"not null" json:"last_name" binding:"required"`
	Email     string `gorm:"unique;not null" json:"email" binding:"required,email"`
	Phone     string `gorm:"unique;not null" json:"phone" binding:"required"`
	Password  string `gorm:"not null" json:"password" binding:"required,min=8"`
}

func (h *UserHandler) RegisterUser(c *gin.Context) {
	var req registerReq
	if ok := utils.BindData(c, &req); !ok {
		return
	}

	user := &entity.User{
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Email:     req.Email,
		Phone:     req.Phone,
		Password:  req.Password,
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

	c.JSON(http.StatusCreated, gin.H{
		"tokens": token,
		"uid":    user.UID,
	})

}
