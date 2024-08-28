package user

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *UserHandler) LogoutUser(c *gin.Context) {
	// user := c.MustGet("user")
	c.SetCookie("token", "", -1, "/", "", true, true)
	c.SetCookie("uid", "", -1, "/", "", true, true)
	c.JSON(http.StatusOK, gin.H{
		"message": "user loggged out",
	})
}
