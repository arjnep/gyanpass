package book

import (
	"log"
	"net/http"
	"strconv"

	"github.com/arjnep/gyanpass/pkg/jwt"
	"github.com/arjnep/gyanpass/pkg/response"
	"github.com/arjnep/gyanpass/pkg/utils"
	"github.com/gin-gonic/gin"
)

type updateBookReq struct {
	Title       string  `gorm:"not null" json:"title" binding:"omitempty"`
	Author      string  `gorm:"not null" json:"author" binding:"omitempty"`
	Genre       string  `json:"genre" binding:"omitempty"`
	Description string  `json:"description" binding:"omitempty"`
	ImageUrl    string  `gorm:"not null" json:"image_url" binding:"omitempty"`
	Address     string  `json:"address" binding:"omitempty"`
	Latitude    float64 `gorm:"not null" json:"latitude" binding:"omitempty,latitude"`
	Longitude   float64 `gorm:"not null" json:"longitude" binding:"omitempty,longitude"`
}

func (h *BookHandler) UpdateBook(c *gin.Context) {
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
		c.JSON(response.Status(err), gin.H{
			"error": err,
		})
		return
	}

	if existingBook.UserID != authUser.UID {
		err := response.NewAuthorizationError("You are not authorized to update this book")
		c.JSON(err.Status(), gin.H{
			"error": err,
		})
		return
	}

	var req updateBookReq

	if ok := utils.BindData(c, &req); !ok {
		return
	}

	updates := make(map[string]interface{})

	if req.Title != "" && req.Title != existingBook.Title {
		updates["title"] = req.Title
	}
	if req.Author != "" && req.Author != existingBook.Author {
		updates["author"] = req.Author
	}
	if req.Genre != "" && req.Genre != existingBook.Genre {
		updates["genre"] = req.Genre
	}
	if req.Description != "" && req.Description != existingBook.Description {
		updates["description"] = req.Description
	}
	if req.ImageUrl != "" && req.ImageUrl != existingBook.ImageUrl {
		updates["image_url"] = req.ImageUrl
	}

	if req.Address != "" && req.Address != existingBook.PickupLocation.Address {
		updates["address"] = req.Address
	}

	if req.Latitude != 0 && req.Latitude != existingBook.PickupLocation.Latitude {
		updates["latitude"] = req.Latitude
	}

	if req.Longitude != 0 && req.Longitude != existingBook.PickupLocation.Longitude {
		updates["longitude"] = req.Longitude
	}

	if len(updates) == 0 {
		err := response.NewBadRequestError("No fields to update")
		c.JSON(err.Status(), gin.H{
			"error": err,
		})
		return
	}

	err = h.bookUsecase.UpdateBook(existingBook, updates)
	if err != nil {
		log.Printf("Failed to update book: %v\n", err.Error())
		c.JSON(response.Status(err), gin.H{
			"error": err,
		})
		return
	}

	updatedBook, err := h.bookUsecase.GetBookByID(existingBook.ID)
	if err != nil {
		log.Printf("Failed to retrieve updated book details: %v\n", err.Error())
		c.JSON(response.Status(err), gin.H{
			"error": err,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"book": updatedBook,
	})

}
