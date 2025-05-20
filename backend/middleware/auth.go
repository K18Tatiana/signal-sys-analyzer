package middleware

import (
	"context"
	"net/http"
	"strings"

	"backend/utils"
)

// Clave para el contexto
type contextKey string

const UserIDKey contextKey = "userID"

// AuthMiddleware verifica que el usuario esté autenticado
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Obtener token del header Authorization
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Autorización requerida", http.StatusUnauthorized)
			return
		}

		// El token debe estar en formato "Bearer <token>"
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			http.Error(w, "Formato de autorización inválido", http.StatusUnauthorized)
			return
		}

		// Validar el token
		claims, err := utils.ValidateToken(parts[1])
		if err != nil {
			http.Error(w, "Token inválido", http.StatusUnauthorized)
			return
		}

		// Añadir userID al contexto
		ctx := context.WithValue(r.Context(), UserIDKey, claims.UserID)

		// Continuar con el siguiente handler
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// GetUserID extrae el ID de usuario del contexto
func GetUserID(r *http.Request) (uint, bool) {
	userID, ok := r.Context().Value(UserIDKey).(uint)
	return userID, ok
}

// OptionalAuthMiddleware intenta autenticar al usuario, pero no falla si no hay token
func OptionalAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Obtener token del header Authorization
		authHeader := r.Header.Get("Authorization")

		// Si no hay token, simplemente continuar
		if authHeader == "" {
			next.ServeHTTP(w, r)
			return
		}

		// El token debe estar en formato "Bearer <token>"
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			next.ServeHTTP(w, r)
			return
		}

		// Intentar validar el token
		claims, err := utils.ValidateToken(parts[1])
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}

		// Añadir userID al contexto
		ctx := context.WithValue(r.Context(), UserIDKey, claims.UserID)

		// Continuar con el siguiente handler
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
