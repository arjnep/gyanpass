package exchange

import (
	"log"
	"net/http"

	"github.com/arjnep/gyanpass/pkg/jwt"
	"github.com/arjnep/gyanpass/pkg/response"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (h *ExchangeHandler) AcceptExchangeRequest(c *gin.Context) {

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

	pathExchangeRequestID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		if uuid.IsInvalidLengthError(err) {
			err := response.NewNotFoundError("exchange request", c.Param("id"))
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

	fetchedExchangeRequest, err := h.exchangeUsecase.GetExchangeRequestByID(pathExchangeRequestID, loggedInUserID)
	if err != nil {
		log.Printf("Failed To Get Exchange Request By ID to accept %v\n", fetchedExchangeRequest)
		c.JSON(response.Status(err), gin.H{
			"error": err,
		})
		return
	}

	err = h.exchangeUsecase.AcceptExchange(fetchedExchangeRequest, loggedInUserID)
	if err != nil {
		log.Printf("Failed To Accept Exchange Request %v\n", err)
		c.JSON(response.Status(err), gin.H{
			"error": err,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "request accepted",
	})
}
