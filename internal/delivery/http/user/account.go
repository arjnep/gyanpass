package user

import (
	"log"
	"net/http"

	"github.com/arjnep/gyanpass/pkg/jwt"
	"github.com/arjnep/gyanpass/pkg/response"
	"github.com/arjnep/gyanpass/pkg/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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

	// Extract the user ID from path parameters
	pathUserID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		if uuid.IsInvalidLengthError(err) {
			err := response.NewNotFoundError("users", c.Param("id"))
			c.JSON(err.Status(), gin.H{
				"error": err,
			})
			return
		}
		log.Printf("Unable to Parse User ID From Param for unknown reason: %v\n", c)
		err := response.NewInternalServerError()
		c.JSON(err.Status(), gin.H{
			"error": err,
		})
		return
	}
	loggedInUserID := user.(*jwt.TokenClaims).User.UID

	if pathUserID != loggedInUserID {
		err := response.NewAuthorizationError("Unauthorized access to this user data")
		c.JSON(err.Status(), gin.H{
			"error": err,
		})
		return
	}

	userFetched, err := h.userUsecase.GetUserByID(loggedInUserID)
	if err != nil {
		log.Printf("Unable to find user: %v\n%v", loggedInUserID, err)
		c.JSON(response.Status(err), gin.H{
			"error": err,
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
	// Email     string `gorm:"unique;not null" json:"email,omitempty" binding:"omitempty,email"`
	Phone string `gorm:"unique;not null" json:"phone,omitempty" binding:"omitempty,len=10"`
}

func (h *UserHandler) UpdateUser(c *gin.Context) {
	authUser := c.MustGet("user").(*jwt.TokenClaims).User
	pathUserID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		if uuid.IsInvalidLengthError(err) {
			err := response.NewNotFoundError("users", c.Param("id"))
			c.JSON(err.Status(), gin.H{
				"error": err,
			})
			return
		}
		log.Printf("Unable to Parse User ID From Param for unknown reason: %v\n", c)
		err := response.NewInternalServerError()
		c.JSON(err.Status(), gin.H{
			"error": err,
		})
		return
	}

	if pathUserID != authUser.UID {
		err := response.NewAuthorizationError("Unauthorized access to update this user data")
		c.JSON(err.Status(), gin.H{
			"error": err,
		})
		return
	}

	existingUser, err := h.userUsecase.GetUserByID(authUser.UID)
	if err != nil {
		c.JSON(response.Status(err), gin.H{
			"error": err,
		})
		return
	}

	var req updateReq

	if ok := utils.BindData(c, &req); !ok {
		return
	}

	updates := make(map[string]interface{})

	if req.FirstName != "" && req.FirstName != existingUser.FirstName {
		updates["first_name"] = req.FirstName
	}
	if req.LastName != "" && req.LastName != existingUser.LastName {
		updates["last_name"] = req.LastName
	}
	// if req.Email != "" && req.Email != existingUser.Email {
	// 	updates["email"] = req.Email
	// }
	if req.Phone != "" && req.Phone != existingUser.Phone {
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
	err = h.userUsecase.Update(existingUser, updates)
	if err != nil {
		log.Printf("Failed to update user: %v\n", err.Error())
		c.JSON(response.Status(err), gin.H{
			"error": err,
		})
		return
	}

	updatedUser, err := h.userUsecase.GetUserByID(existingUser.UID)
	if err != nil {
		log.Printf("Failed to retrieve updated user details: %v\n", err.Error())
		c.JSON(response.Status(err), gin.H{
			"error": err,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user": updatedUser,
	})
}
