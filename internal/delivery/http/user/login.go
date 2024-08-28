package user

import (
	"log"
	"net/http"

	"github.com/arjnep/gyanpass/internal/entity"
	"github.com/arjnep/gyanpass/pkg/response"
	"github.com/arjnep/gyanpass/pkg/utils"
	"github.com/gin-gonic/gin"
)

type loginReq struct {
	Email    string `gorm:"unique;not null" json:"email" binding:"required,email"`
	Password string `gorm:"not null" json:"password" binding:"required,min=8"`
}

func (h *UserHandler) LoginUser(c *gin.Context) {
	var req loginReq

	if ok := utils.BindData(c, &req); !ok {
		return
	}

	user := &entity.User{
		Email:    req.Email,
		Password: req.Password,
	}

	err := h.userUsecase.Login(user)
	if err != nil {
		log.Printf("Failed to log in user: %v\n", err.Error())
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

	c.JSON(http.StatusOK, gin.H{
		"token": token,
		"uid":   user.UID,
	})

}
