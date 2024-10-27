package middleware

import (
	"api/internal/entity"
	"errors"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Error() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		errs := c.Errors
		if len(errs) > 0 {
			err := errs[0]

			slog.Error(err.Error())

			var bErr *entity.BusinessError
			if errors.As(err, &bErr) {
				c.JSON(http.StatusConflict, gin.H{
					"error":   "Conflict",
					"message": err.Error(),
				})
				return
			}

			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "Internal Server Error",
				"message": err.Error(),
			})
		}
	}
}
