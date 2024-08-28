package utils

import (
	"fmt"
	"log"

	"github.com/arjnep/gyanpass/pkg/response"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type invalidArgument struct {
	Field string      `json:"field"`
	Value interface{} `json:"value"`
	Tag   string      `json:"tag"`
	Param string      `json:"param"`
}

func BindData(c *gin.Context, req interface{}) bool {
	log.Println("Step 1")
	if c.ContentType() != "application/json" {
		msg := fmt.Sprintf("%s only accepts Content-Type application/json", c.FullPath())

		err := response.NewUnsupportedMediaTypeError(msg)

		c.JSON(err.Status(), gin.H{
			"error": err,
		})
		return false
	}

	if err := c.ShouldBind(req); err != nil {
		log.Printf("Error binding data: %+v\n", err)
		if errs, ok := err.(validator.ValidationErrors); ok {
			var invalidArgs []invalidArgument

			for _, err := range errs {
				invalidArgs = append(invalidArgs, invalidArgument{
					err.Field(),
					err.Value(),
					err.Tag(),
					err.Param(),
				})
			}

			err := response.NewBadRequestError("Invalid request parameters. See invalidArgs")

			c.JSON(err.Status(), gin.H{
				"error":       err,
				"invalidArgs": invalidArgs,
			})
			return false
		}

		fallBack := response.NewBadRequestError("Malformed Payload")

		c.JSON(fallBack.Status(), gin.H{
			"error": fallBack,
		})
		return false
	}
	return true
}
