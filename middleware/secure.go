package middleware

import (
	"github.com/gin-gonic/gin"
	"go-zentao-task/pkg/config"
)

type SecureConfig struct {
	XSSProtection           string
	ContentTypeOptions      string
	XFrameOptions           string
	ContentSecurityPolicy   string
	StrictTransportSecurity string
}

var (
	DefaultSecureConfig = SecureConfig{
		XSSProtection:      "1; mode=block",
		ContentTypeOptions: "nosniff",
		XFrameOptions:      "SAMEORIGIN",
	}
)

func Secure() gin.HandlerFunc {
	return SecureWithConfig(DefaultSecureConfig)
}

func SecureWithConfig(s SecureConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		if s.XSSProtection != "" {
			c.Header("X-XSS-Protection", s.XSSProtection)
		}
		if s.ContentTypeOptions != "" {
			c.Header("X-Content-Type-Options", s.ContentTypeOptions)
		}
		if s.XFrameOptions != "" {
			c.Header("X-Frame-Options", s.XFrameOptions)
		}
		if s.ContentSecurityPolicy != "" {
			c.Header("Content-Security-Policy", s.ContentSecurityPolicy)
		}
		if config.Get("app.scheme") == "https" && s.StrictTransportSecurity != "" {
			c.Header("Strict-Transport-Security", s.StrictTransportSecurity)
		}
		c.Next()
	}
}
