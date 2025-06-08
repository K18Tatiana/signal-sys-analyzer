package middleware

import (
	"net/http"
	"strings"

	"backend/utils"
	"github.com/gin-gonic/gin"
)

// AuthMiddleware verifica que el usuario esté autenticado (obligatorio)
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Obtener token del header Authorization
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Autorización requerida"})
			c.Abort()
			return
		}

		// El token debe estar en formato "Bearer <token>"
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Formato de autorización inválido"})
			c.Abort()
			return
		}

		// Validar el token
		claims, err := utils.ValidateToken(parts[1])
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token inválido"})
			c.Abort()
			return
		}

		// Guardar userID en el contexto de Gin
		c.Set("userID", claims.UserID)

		// Continuar con el siguiente handler
		c.Next()
	}
}

// OptionalAuthMiddleware intenta autenticar al usuario, pero no falla si no hay token
func OptionalAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Obtener token del header Authorization
		authHeader := c.GetHeader("Authorization")

		// Si no hay token, simplemente continuar
		if authHeader == "" {
			c.Next()
			return
		}

		// El token debe estar en formato "Bearer <token>"
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.Next()
			return
		}

		// Intentar validar el token
		claims, err := utils.ValidateToken(parts[1])
		if err != nil {
			c.Next()
			return
		}

		// Guardar userID en el contexto de Gin
		c.Set("userID", claims.UserID)

		// Continuar con el siguiente handler
		c.Next()
	}
}

// GetUserIDFromGin extrae el ID de usuario del contexto de Gin
func GetUserIDFromGin(c *gin.Context) (uint, bool) {
	if userID, exists := c.Get("userID"); exists {
		if id, ok := userID.(uint); ok {
			return id, true
		}
	}
	return 0, false
}
