package notification

import (
	"github.com/arjnep/gyanpass/internal/delivery/middleware"
	"github.com/arjnep/gyanpass/pkg/jwt"
	"github.com/arjnep/gyanpass/pkg/notification"
	"github.com/gin-gonic/gin"
)

type NotificationHandler struct {
	jwtService          jwt.Service
	notificationService notification.Service
}

type Config struct {
	R                   *gin.Engine
	NotificationService notification.Service
	JWTService          jwt.Service
}

func NewNotificationHandler(c *Config) {
	h := &NotificationHandler{
		notificationService: c.NotificationService,
		jwtService:          c.JWTService,
	}

	notificationRoutes := c.R.Group("/api/notifications")
	{
		notificationRoutes.GET("/", middleware.AuthUser(h.jwtService), h.GetUserNotifications)
		notificationRoutes.POST("/:id/read", middleware.AuthUser(h.jwtService), h.ReadNotification)
		notificationRoutes.DELETE("/:id/remove", middleware.AuthUser(h.jwtService), h.RemoveNotification)
	}
}
