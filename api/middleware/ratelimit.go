package middleware

import (
	"net/http"
	"web-metric/config"
	"web-metric/service/ratelimiter"

	"github.com/gin-gonic/gin"
)

func RateLimitMiddleware(s *ratelimiter.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.Query("userID")
		if !s.RateLimit(userID, config.UserRateLimit()) {
			c.AbortWithStatusJSON(
				http.StatusTooManyRequests,
				gin.H{
					"success": false,
					"error": gin.H{
						"code":    http.StatusTooManyRequests,
						"message": "Rate limit exceeded",
					},
				},
			)
			return
		}
		c.Next()
	}
}
