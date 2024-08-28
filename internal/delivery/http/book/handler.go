package book

import (
	"github.com/arjnep/gyanpass/config"
	"github.com/arjnep/gyanpass/internal/delivery/middleware"
	"github.com/arjnep/gyanpass/internal/usecase"
	"github.com/arjnep/gyanpass/pkg/jwt"
	"github.com/gin-gonic/gin"
)

type BookHandler struct {
	bookUsecase usecase.BookUsecase
	jwtService  jwt.Service
	Cfg         *config.Configuration
}

type Config struct {
	R           *gin.Engine
	BookUsecase usecase.BookUsecase
	JwtService  jwt.Service
}

func NewBookHandler(c *Config) {
	h := &BookHandler{
		bookUsecase: c.BookUsecase,
		jwtService:  c.JwtService,
	}

	bookRoutes := c.R.Group("/api/books")
	{
		bookRoutes.POST("/", middleware.AuthUser(h.jwtService), h.AddBook)
		bookRoutes.GET("/search", middleware.Pagination(), h.SearchBooks)
		bookRoutes.GET("/:id", h.GetBook)
		bookRoutes.PUT("/:id", middleware.AuthUser(h.jwtService), h.UpdateBook)
		bookRoutes.DELETE("/:id", middleware.AuthUser(h.jwtService), h.DeleteBook)
	}
}
