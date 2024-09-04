package exchange

import (
	"log"
	"net/http"
	"strconv"

	"github.com/arjnep/gyanpass/pkg/jwt"
	"github.com/arjnep/gyanpass/pkg/response"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (h *ExchangeHandler) GetUserExchangeRequests(c *gin.Context) {
	authUser, exists := c.Get("user")
	if !exists {
		log.Printf("Unable to extract user from request context: %v\n", c)
		err := response.NewInternalServerError()
		c.JSON(err.Status(), gin.H{
			"error": err,
		})
		return
	}
	loggedInUserID := authUser.(*jwt.TokenClaims).User.UID

	bookIDParam := c.Query("bookID")
	if bookIDParam != "" {
		parsedBookID, err := strconv.ParseUint(bookIDParam, 10, 32)
		if err != nil {
			log.Printf("Invalid book id: %v\n", err)
			err := response.NewBadRequestError("Invalid book id")
			c.JSON(err.Status(), gin.H{
				"error": err,
			})
			return
		}
		requests, err := h.exchangeUsecase.GetExchangeRequestsByBookIDAndUserID(uint(parsedBookID), loggedInUserID)
		if err != nil {
			log.Printf("Failed to get exchange requests by book ID %v: %v\n", loggedInUserID, err)
			c.JSON(response.Status(err), gin.H{
				"error": err,
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"requests": requests,
		})
		return
	}

	// Fetch exchange requests related to the logged-in user
	requests, err := h.exchangeUsecase.GetExchangeRequestsByUserID(loggedInUserID)
	if err != nil {
		log.Printf("Failed to get exchange requests by user ID %v: %v\n", loggedInUserID, err)
		c.JSON(response.Status(err), gin.H{
			"error": err,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"requests": requests,
	})
}

func (h *ExchangeHandler) GetExchangeRequestByID(c *gin.Context) {
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
		log.Printf("Failed To fetched exchange Request By id %v\n", err)
		c.JSON(response.Status(err), gin.H{
			"error": err,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"request": fetchedExchangeRequest,
	})

}

func (h *ExchangeHandler) GetExchangeRequestsMade(c *gin.Context) {
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

	requestsMade, err := h.exchangeUsecase.GetExchangeRequestsByRequestedByID(loggedInUserID)
	if err != nil {
		log.Printf("Failed To Get Exchange Requests By Requested By ID %v\n", err)
		c.JSON(response.Status(err), gin.H{
			"error": err,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"requests_made": requestsMade,
	})

}

func (h *ExchangeHandler) GetExchangeRequestsReceived(c *gin.Context) {
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

	requestsReceived, err := h.exchangeUsecase.GetExchangeRequestsByRequestedToID(loggedInUserID)
	if err != nil {
		log.Printf("Failed To Get Exchange Requests By Requested To ID %v\n", err)
		c.JSON(response.Status(err), gin.H{
			"error": err,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"requests_received": requestsReceived,
	})
}
