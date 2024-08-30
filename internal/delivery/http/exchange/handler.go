package exchange

import (
	"github.com/arjnep/gyanpass/config"
	"github.com/arjnep/gyanpass/internal/delivery/middleware"
	"github.com/arjnep/gyanpass/internal/usecase"
	"github.com/arjnep/gyanpass/pkg/jwt"
	"github.com/gin-gonic/gin"
)

type ExchangeHandler struct {
	bookUsecase     usecase.BookUsecase
	exchangeUsecase usecase.ExchangeUsecase
	jwtService      jwt.Service
	Cfg             *config.Configuration
}

type Config struct {
	R               *gin.Engine
	BookUsecase     usecase.BookUsecase
	ExchangeUsecase usecase.ExchangeUsecase
	JwtService      jwt.Service
}

func NewExchangeHandler(c *Config) {
	h := &ExchangeHandler{
		bookUsecase:     c.BookUsecase,
		exchangeUsecase: c.ExchangeUsecase,
		jwtService:      c.JwtService,
	}

	exchangeRoutes := c.R.Group("/api/exchange/requests")
	{
		exchangeRoutes.POST("/", middleware.AuthUser(h.jwtService), h.CreateExchangeRequest)
		exchangeRoutes.GET("/:id", middleware.AuthUser(h.jwtService), h.GetExchangeRequestByID)
		exchangeRoutes.GET("/made", middleware.AuthUser(h.jwtService), h.GetExchangeRequestsMade)
		exchangeRoutes.GET("/received", middleware.AuthUser(h.jwtService), h.GetExchangeRequestsReceived)
		exchangeRoutes.POST("/:id/accept", middleware.AuthUser(h.jwtService), h.AcceptExchangeRequest)
		exchangeRoutes.POST("/:id/decline", middleware.AuthUser(h.jwtService), h.DeclineExchangeRequest)
		exchangeRoutes.POST("/:id/confirm", middleware.AuthUser(h.jwtService), h.ConfirmExchangeRequest)
		exchangeRoutes.DELETE("/:id/delete", middleware.AuthUser(h.jwtService), h.DeleteExchangeRequest)

	}
}
