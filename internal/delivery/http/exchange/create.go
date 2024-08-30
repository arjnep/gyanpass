package exchange

import (
	"log"
	"net/http"

	"github.com/arjnep/gyanpass/internal/entity"
	"github.com/arjnep/gyanpass/pkg/jwt"
	"github.com/arjnep/gyanpass/pkg/response"
	"github.com/arjnep/gyanpass/pkg/utils"
	"github.com/gin-gonic/gin"
)

type createReq struct {
	RequestedBookID uint `gorm:"not null" json:"requested_book_id" binding:"required"`
	OfferedBookID   uint `gorm:"not null" json:"offered_book_id" binding:"required"`
}

func (h *ExchangeHandler) CreateExchangeRequest(c *gin.Context) {
	var req createReq

	ok := utils.BindData(c, &req)
	if !ok {
		return
	}

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

	requestedBook, err := h.bookUsecase.GetBookByID(req.RequestedBookID)
	if err != nil {
		log.Printf("Unable to Get Book By id for unknown reason: %v\n", c)
		c.JSON(response.Status(err), gin.H{
			"error": err,
		})
		return
	}

	offeredBook, err := h.bookUsecase.GetBookByID(req.OfferedBookID)
	if err != nil {
		log.Printf("Unable to Get Book By id for unknown reason: %v\n", c)
		c.JSON(response.Status(err), gin.H{
			"error": err,
		})
		return
	}

	if offeredBook.UserID != loggedInUserID {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "offering book not found",
		})
		return
	}

	newExchangeRequest := entity.ExchangeRequest{
		RequestedByID:   loggedInUserID,
		RequestedToID:   requestedBook.Owner.UID,
		RequestedBookID: req.RequestedBookID,
		OfferedBookID:   req.OfferedBookID,
		RequestedBook:   *requestedBook,
		OfferedBook:     *offeredBook,
	}

	sanitized, err := h.exchangeUsecase.RequestExchange(&newExchangeRequest)
	if err != nil {
		log.Printf("Failed To Create New Exchange Request %v\n", err)
		c.JSON(response.Status(err), gin.H{
			"error": err,
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"request": sanitized,
	})
}
