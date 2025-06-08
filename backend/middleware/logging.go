package middleware

import (
	"github.com/gin-gonic/gin"
	"log"
)

// LoggingMiddleware es un middleware personalizado de logging para Gin
func LoggingMiddleware() gin.HandlerFunc {
	return gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		// Personalizar el formato de log
		log.Printf(
			"%s %s %d %v %s",
			param.Method,
			param.Path,
			param.StatusCode,
			param.Latency,
			param.ClientIP,
		)
		return ""
	})
}
