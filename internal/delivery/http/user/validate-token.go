package user

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *UserHandler) ValidateToken(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"valid": true,
	})
}
