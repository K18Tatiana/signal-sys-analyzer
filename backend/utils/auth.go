package utils

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

const (
	// TokenExpiryTime define cuánto tiempo es válido el token JWT
	TokenExpiryTime = time.Hour * 24 // 24 horas

	// MinPasswordLength define la longitud mínima de contraseña
	MinPasswordLength = 6
)

// JWT signing key - en producción, esto debería venir de una variable de entorno
var jwtSecret = []byte("tu_clave_secreta_muy_larga_y_compleja_aqui")

// Claims estructura para el payload del JWT
type Claims struct {
	UserID uint `json:"user_id"` // Cambiado de int a uint para coincidir con GORM
	jwt.RegisteredClaims
}

// HashPassword cifra una contraseña usando bcrypt
func HashPassword(password string) (string, error) {
	if len(password) < MinPasswordLength {
		return "", errors.New("la contraseña debe tener al menos 6 caracteres")
	}

	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 12) // Costo 12
	return string(bytes), err
}

// CheckPasswordHash verifica si la contraseña coincide con el hash
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// GenerateToken genera un token JWT para el usuario
func GenerateToken(userID uint) (string, error) {
	// Configurar tiempo de expiración
	expirationTime := time.Now().Add(TokenExpiryTime)

	// Crear los claims
	claims := &Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	// Crear el token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Firmar el token con la clave secreta
	tokenString, err := token.SignedString(jwtSecret)

	return tokenString, err
}

// ValidateToken valida un token JWT
// Esta función es la que está buscando el middleware
func ValidateToken(tokenString string) (*Claims, error) {
	// Parsear el token
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("token inválido")
}
