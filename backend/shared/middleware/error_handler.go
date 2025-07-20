package middleware

import (
	"certitrack/backend/shared/errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) > 0 {
			err := c.Errors.Last().Err

			if e, ok := err.(*errors.ErrorResponse); ok {
				c.JSON(e.Status, e)
				return
			}

			c.JSON(http.StatusInternalServerError, errors.ErrInternalServer)
		}
	}
}
