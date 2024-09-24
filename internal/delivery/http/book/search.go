package book

import (
	"log"
	"net/http"

	"github.com/arjnep/gyanpass/pkg/response"
	"github.com/gin-gonic/gin"
)

func (h *BookHandler) SearchBooks(c *gin.Context) {
	queryParams := map[string]string{
		"title":   c.Query("title"),
		"address": c.Query("address"),
	}

	page, _ := c.Get("page")
	size, _ := c.Get("size")

	pageInt, ok := page.(int)
	if !ok {
		pageInt = 1
	}
	sizeInt, ok := size.(int)
	if !ok {
		sizeInt = 10
	}

	books, total, err := h.bookUsecase.SearchBooks(queryParams, pageInt, sizeInt)
	if err != nil {
		log.Printf("Failed to Search Book: %v", err)
		c.JSON(response.Status(err), gin.H{
			"error": err,
		})
	}

	var booksResponse []gin.H
	for _, book := range books {
		booksResponse = append(booksResponse, gin.H{
			"id":        book.ID,
			"title":     book.Title,
			"author":    book.Author,
			"genre":     book.Genre,
			"image_url": book.ImageUrl,
		})
	}

	totalPages := (total + sizeInt - 1) / sizeInt

	c.JSON(http.StatusOK, gin.H{
		"books":       booksResponse,
		"page":        pageInt,
		"size":        sizeInt,
		"total":       total,
		"total_pages": totalPages,
	})
}
