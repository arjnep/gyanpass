package user

import (
	"github.com/arjnep/gyanpass/config"
	"github.com/arjnep/gyanpass/internal/delivery/middleware"
	"github.com/arjnep/gyanpass/internal/usecase"
	"github.com/arjnep/gyanpass/pkg/jwt"
	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userUsecase usecase.UserUsecase
	jwtService  jwt.Service
	Cfg         *config.Configuration
}

type Config struct {
	R           *gin.Engine
	UserUsecase usecase.UserUsecase
	JwtService  jwt.Service
}

func NewUserHandler(c *Config) {
	h := &UserHandler{
		userUsecase: c.UserUsecase,
		jwtService:  c.JwtService,
	}

	authRoutes := c.R.Group("/api/auth")
	{
		authRoutes.POST("/register", h.RegisterUser)
		authRoutes.POST("/login", h.LoginUser)
		authRoutes.POST("/logout", middleware.AuthUser(h.jwtService), h.LogoutUser)
		// userAuth.POST("/forgot-password", middleware.AuthUser(h.jwtService), h.ForgotPassword)
	}

	userRoutes := c.R.Group("/api/users")
	{
		userRoutes.GET("/:id", middleware.AuthUser(h.jwtService), h.GetUser)
		userRoutes.PUT("/:id", middleware.AuthUser(h.jwtService), h.UpdateUser)
		// userRoutes.DELETE("/:id", middleware.AuthUser(h.jwtService), h.DeleteUser)
		// userRoutes.POST("/reset-password", middleware.AuthUser(h.jwtService), h.ResetPassword)
	}
}
