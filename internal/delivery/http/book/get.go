package book

import (
	"log"
	"net/http"
	"strconv"

	"github.com/arjnep/gyanpass/pkg/jwt"
	"github.com/arjnep/gyanpass/pkg/response"
	"github.com/gin-gonic/gin"
)

func (h *BookHandler) GetUserBooks(c *gin.Context) {
	authUser, exists := c.Get("user")
	if !exists {
		log.Printf("Unable to extract user from request context for unknown reason: %v\n", c)
		err := response.NewInternalServerError()
		c.JSON(err.Status(), gin.H{
			"error": err,
		})

		return
	}

	books, err := h.bookUsecase.GetBooksByUserID(authUser.(*jwt.TokenClaims).User.UID)
	if err != nil {
		log.Printf("Failed to Get Books: %v", err)
		c.JSON(response.Status(err), gin.H{
			"error": err,
		})
		return
	}

	if len(books) == 0 {
		c.JSON(http.StatusOK, gin.H{
			"message": "No Books Available",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"books": books,
	})

}

func (h *BookHandler) GetBook(c *gin.Context) {
	pathBookID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		err := response.NewBadRequestError("id of book should be number")
		c.JSON(err.Status(), gin.H{
			"error": err,
		})
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

	book, err := h.bookUsecase.GetBookByID(uint(pathBookID))
	if err != nil {
		log.Println("Failed Getting Book:", err)
		c.JSON(response.Status(err), gin.H{
			"error": err,
		})
		return
	}

	var bookResponse interface{}
	if authUser.(*jwt.TokenClaims).User.UID == book.Owner.UID {
		bookResponse = book
	} else {
		bookResponse = gin.H{
			"id":          book.ID,
			"title":       book.Title,
			"author":      book.Author,
			"genre":       book.Genre,
			"description": book.Description,
			"owner": gin.H{
				"user_id":    book.Owner.UID,
				"first_name": book.Owner.FirstName,
				"last_name":  book.Owner.LastName,
			},
			"is_active": book.IsActive,
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"book": bookResponse,
	})

}
