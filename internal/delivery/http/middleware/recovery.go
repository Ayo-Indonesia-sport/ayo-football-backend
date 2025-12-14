package middleware

import (
	"log"
	"net/http"
	"runtime/debug"

	"github.com/gin-gonic/gin"
	"github.com/zenkriztao/ayo-football-backend/pkg/response"
)

// RecoveryMiddleware handles panic recovery
func RecoveryMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// Log the stack trace
				log.Printf("Panic recovered: %v\n%s", err, debug.Stack())

				response.Error(c, http.StatusInternalServerError, "Internal server error", nil)
				c.Abort()
			}
		}()

		c.Next()
	}
}
