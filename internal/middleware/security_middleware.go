package middleware

import (
	"net/http"

	"github.com/HasanNugroho/starter-golang/internal/configs"
	"github.com/HasanNugroho/starter-golang/internal/pkg/response"
	"github.com/gin-gonic/gin"
)

func SecurityMiddleware(cfg *configs.Configuration) gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Host != cfg.Security.ExpectedHost {
			response.SendError(c, http.StatusBadRequest, "Invalid host header", nil)
			c.Abort()
			return
		}
		c.Header("X-Frame-Options", cfg.Security.XFrameOptions)
		c.Header("Content-Security-Policy", cfg.Security.ContentSecurity)
		c.Header("X-XSS-Protection", cfg.Security.XXSSProtection)
		c.Header("Strict-Transport-Security", cfg.Security.StrictTransport)
		c.Header("Referrer-Policy", cfg.Security.ReferrerPolicy)
		c.Header("X-Content-Type-Options", cfg.Security.XContentTypeOpts)
		c.Header("Permissions-Policy", cfg.Security.PermissionsPolicy)
		c.Next()
	}
}
