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

	userAuth := c.R.Group("/api/auth")
	{
		userAuth.POST("/register", h.RegisterUser)
		userAuth.POST("/login", h.LoginUser)
		userAuth.POST("/logout", middleware.AuthUser(h.jwtService), h.LogoutUser)
		// userAuth.POST("/forgot-password", middleware.AuthUser(h.jwtService), h.ForgotPassword)
		// userAuth.POST("/reset-password", middleware.AuthUser(h.jwtService), h.ResetPassword)
	}

	userProfile := c.R.Group("/api")
	userProfile.Use(middleware.AuthUser(h.jwtService))
	{
		userProfile.GET("/user", h.GetUser)
		userProfile.PUT("/user", h.UpdateUser)
		// userProfile.DELETE("/user", h.DeleteUser)

	}
}
