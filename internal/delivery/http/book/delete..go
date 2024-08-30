package book

import (
	"log"
	"net/http"
	"strconv"

	"github.com/arjnep/gyanpass/pkg/jwt"
	"github.com/arjnep/gyanpass/pkg/response"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func (h *BookHandler) DeleteBook(c *gin.Context) {
	authUser := c.MustGet("user").(*jwt.TokenClaims).User

	pathBookID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		err := response.NewBadRequestError("id of book should be number")
		c.JSON(err.Status(), gin.H{
			"error": err,
		})
		return
	}

	existingBook, err := h.bookUsecase.GetBookByID(uint(pathBookID))
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			err := response.NewNotFoundError("book", c.Param("id"))
			c.JSON(err.Status(), gin.H{
				"error": err,
			})
			return
		}
		log.Println("Failed Getting Book To Delete:", err)
		err := response.NewInternalServerError()
		c.JSON(err.Status(), gin.H{
			"error": err,
		})
		return
	}

	if existingBook.UserID != authUser.UID {
		err := response.NewAuthorizationError("You are not authorized to delete this book")
		c.JSON(err.Status(), gin.H{
			"error": err,
		})
		return
	}

	err = h.bookUsecase.DeleteBook(existingBook)
	if err != nil {
		log.Println("Failed Deleting Book: ", err)
		err := response.NewInternalServerError()
		c.JSON(err.Status(), gin.H{
			"error": err,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Book Deleted",
	})

}
