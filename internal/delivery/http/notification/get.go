package notification

import (
	"log"
	"net/http"

	"github.com/arjnep/gyanpass/pkg/jwt"
	"github.com/arjnep/gyanpass/pkg/response"
	"github.com/gin-gonic/gin"
)

func (h *NotificationHandler) GetUserNotifications(c *gin.Context) {
	authUser, exists := c.Get("user")
	if !exists {
		log.Printf("Unable to extract user from request context for unknown reason: %v\n", c)
		err := response.NewInternalServerError()
		c.JSON(err.Status(), gin.H{
			"error": err,
		})

		return
	}

	notifications, err := h.notificationService.GetUserNotifications(authUser.(*jwt.TokenClaims).User.UID)
	if err != nil {
		log.Printf("Failed to Get Notifications: %v", err)
		c.JSON(response.Status(err), gin.H{
			"error": err,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"notifications": notifications,
	})

}
