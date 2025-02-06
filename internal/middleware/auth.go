package middleware

import (
	"github.com/digitranslab/kozmo-sandbox/internal/static"
	"github.com/gin-gonic/gin"
)

func Auth() gin.HandlerFunc {
	config := static.GetKozmoSandboxGlobalConfigurations()
	return func(c *gin.Context) {
		if config.App.Key != c.GetHeader("X-Api-Key") {
			c.AbortWithStatus(401)
			return
		}
	}
}
