package middleware

import (
	"strconv"

	"github.com/arjnep/gyanpass/pkg/response"
	"github.com/gin-gonic/gin"
)

const (
	DEFAULT_PAGE_TEXT    = "page"
	DEFAULT_SIZE_TEXT    = "size"
	DEFAULT_PAGE         = "1"
	DEFAULT_PAGE_SIZE    = "10"
	DEFAULT_MIN_PAGESIZE = 10
	DEFAULT_MAX_PAGESIZE = 100
)

func Pagination() gin.HandlerFunc {
	return Paginate(
		DEFAULT_PAGE_TEXT,
		DEFAULT_SIZE_TEXT,
		DEFAULT_PAGE,
		DEFAULT_PAGE_SIZE,
		DEFAULT_MIN_PAGESIZE,
		DEFAULT_MAX_PAGESIZE,
	)
}

func Paginate(pageText, sizeText, defaultPage, defaultPageSize string, minPageSize, maxPageSize int) gin.HandlerFunc {
	return func(c *gin.Context) {
		pageStr := c.DefaultQuery(pageText, defaultPage)
		page, err := strconv.Atoi(pageStr)
		if err != nil {
			err := response.NewBadRequestError("page number must be an integer")
			c.JSON(err.Status(), gin.H{
				"error": err,
			})
			return
		}

		if page < 0 {
			err := response.NewBadRequestError("page number must be positive")
			c.JSON(err.Status(), gin.H{
				"error": err,
			})
			return
		}

		sizeStr := c.DefaultQuery(sizeText, defaultPageSize)
		size, err := strconv.Atoi(sizeStr)
		if err != nil {
			err := response.NewBadRequestError("page size must be an integer")
			c.JSON(err.Status(), gin.H{
				"error": err,
			})
			return
		}

		if size < minPageSize || size > maxPageSize {
			err := response.NewBadRequestError("page size must be between " + strconv.Itoa(minPageSize) + " and " + strconv.Itoa(maxPageSize))
			c.JSON(err.Status(), gin.H{
				"error": err,
			})
			return
		}

		c.Set(pageText, page)
		c.Set(sizeText, size)

		c.Next()
	}
}
