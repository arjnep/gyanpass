package middleware

import (
	"strings"

	"github.com/arjnep/gyanpass/pkg/jwt"
	"github.com/arjnep/gyanpass/pkg/response"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type authHeader struct {
	Token string `header:"Authorization"`
}

type invalidArgument struct {
	Field string `json:"field"`
	Value string `json:"value"`
	Tag   string `json:"tag"`
	Param string `json:"param"`
}

func AuthUser(s jwt.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		h := authHeader{}

		if err := c.ShouldBindHeader(&h); err != nil {
			if errs, ok := err.(validator.ValidationErrors); ok {
				var invalidArgs []invalidArgument

				for _, err := range errs {
					invalidArgs = append(invalidArgs, invalidArgument{
						err.Field(),
						err.Value().(string),
						err.Tag(),
						err.Param(),
					})
				}

				err := response.NewBadRequestError("Invalid request parameters. See invalidArgs")

				c.JSON(err.Status(), gin.H{
					"error":       err,
					"invalidArgs": invalidArgs,
				})
				c.Abort()
				return
			}

			err := response.NewInternalServerError()
			c.JSON(err.Status(), gin.H{
				"error": err,
			})
			c.Abort()
			return
		}

		token := h.Token
		if strings.HasPrefix(token, "Bearer ") {
			token = strings.TrimPrefix(token, "Bearer ")
		} else {
			token = ""
		}

		if token == "" {
			var err error
			token, err = c.Cookie("token")
			if err != nil {
				err := response.NewAuthorizationError("Authorization Token Required")
				c.JSON(err.Status(), gin.H{
					"error": err,
				})
				c.Abort()
				return
			}
		}

		user, err := s.ValidateToken(token)

		if err != nil {
			err := response.NewAuthorizationError("Token Invalid")
			c.JSON(err.Status(), gin.H{
				"error": err,
			})
			c.Abort()
			return
		}

		c.Set("user", user)
		c.Next()
	}
}
