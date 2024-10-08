package book

import (
	"log"
	"net/http"

	"github.com/arjnep/gyanpass/internal/entity"
	"github.com/arjnep/gyanpass/pkg/jwt"
	"github.com/arjnep/gyanpass/pkg/response"
	"github.com/arjnep/gyanpass/pkg/utils"
	"github.com/gin-gonic/gin"
)

type addBookReq struct {
	Title       string             `gorm:"not null" json:"title" binding:"required"`
	Author      string             `gorm:"not null" json:"author" binding:"required"`
	Genre       string             `json:"genre" binding:"omitempty"`
	ImageUrl    string             `gorm:"not null" json:"image_url" binding:"required"`
	Address     string             `json:"address" binding:"omitempty"`
	Description entity.Description `gorm:"embedded" json:"description" binding:"required"`
	Latitude    float64            `gorm:"not null" json:"latitude" binding:"required,latitude"`
	Longitude   float64            `gorm:"not null" json:"longitude" binding:"required,longitude"`
}

func (h *BookHandler) AddBook(c *gin.Context) {
	var req addBookReq
	if ok := utils.BindData(c, &req); !ok {
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

	newBook := entity.Book{
		Title:       req.Title,
		Author:      req.Author,
		Genre:       req.Genre,
		Description: req.Description,
		ImageUrl:    req.ImageUrl,
		Owner:       *authUser.(*jwt.TokenClaims).User,
		UserID:      authUser.(*jwt.TokenClaims).User.UID,
		PickupLocation: entity.Location{
			Address:   req.Address,
			Latitude:  req.Latitude,
			Longitude: req.Longitude,
		},
		IsActive: true,
	}

	err := h.bookUsecase.AddBook(&newBook)
	if err != nil {
		log.Printf("Failed to add new Book: %v", err)
		c.JSON(response.Status(err), gin.H{
			"error": err,
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"book": newBook,
	})

}
