package user

import (
	"log"
	"net/http"

	"github.com/arjnep/gyanpass/pkg/jwt"
	"github.com/arjnep/gyanpass/pkg/response"
	"github.com/arjnep/gyanpass/pkg/utils"
	"github.com/gin-gonic/gin"
)

func (h *UserHandler) GetUser(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists {
		log.Printf("Unable to extract user from request context for unknown reason: %v\n", c)
		err := response.NewInternalServerError()
		c.JSON(err.Status(), gin.H{
			"error": err,
		})

		return
	}

	uid := user.(*jwt.TokenClaims).User.UID
	log.Println("UID is :", uid)

	userFetched, err := h.userUsecase.GetUserByID(uid)
	if err != nil {
		log.Printf("Unable to find user: %v\n%v", uid, err)
		e := response.NewNotFoundError("user", uid.String())

		c.JSON(e.Status(), gin.H{
			"error": e,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user": userFetched,
	})
}

type updateReq struct {
	FirstName string `gorm:"not null" json:"first_name" binding:"omitempty"`
	LastName  string `gorm:"not null" json:"last_name" binding:"omitempty"`
	Email     string `gorm:"unique;not null" json:"email,omitempty" binding:"omitempty,email"`
	Phone     string `gorm:"unique;not null" json:"phone,omitempty" binding:"omitempty"`
}

func (h *UserHandler) UpdateUser(c *gin.Context) {
	authUser := c.MustGet("user").(*jwt.TokenClaims).User

	var req updateReq

	if ok := utils.BindData(c, &req); !ok {
		return
	}

	updates := make(map[string]interface{})

	if req.FirstName != "" && req.FirstName != authUser.FirstName {
		updates["first_name"] = req.FirstName
	}
	if req.LastName != "" && req.LastName != authUser.LastName {
		updates["last_name"] = req.LastName
	}
	if req.Email != "" && req.Email != authUser.Email {
		updates["email"] = req.Email
	}
	if req.Phone != "" && req.Phone != authUser.Phone {
		updates["phone"] = req.Phone
	}

	if len(updates) == 0 {
		err := response.NewBadRequestError("No fields to update")
		c.JSON(err.Status(), gin.H{
			"error": err,
		})
		return
	}

	// Update only the fields that have changed
	err := h.userUsecase.Update(authUser, updates)
	if err != nil {
		log.Printf("Failed to update user: %v\n", err.Error())
		c.JSON(response.Status(err), gin.H{
			"error": err,
		})
		return
	}

	updatedUser, err := h.userUsecase.GetUserByID(authUser.UID)
	if err != nil {
		log.Printf("Failed to retrieve updated user details: %v\n", err.Error())
		c.JSON(response.Status(err), gin.H{
			"error": err,
		})
		return
	}

	log.Println(updatedUser)

	c.JSON(http.StatusOK, gin.H{
		"user": updatedUser,
	})
}
