package notification

import (
	"log"
	"net/http"

	"github.com/arjnep/gyanpass/pkg/jwt"
	"github.com/arjnep/gyanpass/pkg/response"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (h *NotificationHandler) ReadNotification(c *gin.Context) {
	authUser, exists := c.Get("user")
	if !exists {
		log.Printf("Unable to extract user from request context for unknown reason: %v\n", c)
		err := response.NewInternalServerError()
		c.JSON(err.Status(), gin.H{
			"error": err,
		})
		return
	}
	loggedInUserID := authUser.(*jwt.TokenClaims).User.UID

	pathNotificationID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		if uuid.IsInvalidLengthError(err) {
			err := response.NewNotFoundError("notification", c.Param("id"))
			c.JSON(err.Status(), gin.H{
				"error": err,
			})
			return
		}
		log.Printf("Unable to Parse Notification ID From Param for unknown reason: %v\n", c)
		err := response.NewInternalServerError()
		c.JSON(err.Status(), gin.H{
			"error": err,
		})
		return
	}

	fetchedNotification, err := h.notificationService.GetNotificationByID(pathNotificationID)
	if err != nil {
		c.JSON(response.Status(err), gin.H{
			"error": err,
		})
		return
	}

	err = h.notificationService.MarkNotificationAsRead(fetchedNotification, loggedInUserID)
	if err != nil {
		log.Printf("Failed To Mark Notification Read %v\n", err)
		c.JSON(response.Status(err), gin.H{
			"error": err,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
	})
}
