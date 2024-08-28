package book

import (
	"log"
	"net/http"
	"strconv"

	"github.com/arjnep/gyanpass/pkg/response"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func (h *BookHandler) GetBook(c *gin.Context) {
	pathBookID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		err := response.NewBadRequestError("id of book should be number")
		c.JSON(err.Status(), gin.H{
			"error": err,
		})
		return
	}

	book, err := h.bookUsecase.GetBookByID(uint(pathBookID))
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			err := response.NewNotFoundError("book", c.Param("id"))
			c.JSON(err.Status(), gin.H{
				"error": err,
			})
			return
		}
		log.Println("Failed Getting Book:", err)
		err := response.NewInternalServerError()
		c.JSON(err.Status(), gin.H{
			"error": err,
		})
		return
	}

	bookResponse := gin.H{
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
	}

	c.JSON(http.StatusOK, gin.H{
		"book": bookResponse,
	})

}
