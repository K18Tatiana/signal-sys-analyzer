package handlers

import (
	"net/http"
	"time"

	"backend/database"
	"backend/middleware"
	"backend/models"
	"backend/utils"
	"github.com/gin-gonic/gin"
)

func RegisterHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req models.UserRegisterRequest

		// Bind JSON del cuerpo de la solicitud
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Error al decodificar la solicitud"})
			return
		}

		// Validar campos requeridos
		if req.Username == "" || req.Email == "" || req.Password == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Todos los campos son requeridos"})
			return
		}

		// Verificar si el correo ya existe
		var existingUser models.User
		result := database.DB.Where("email = ?", req.Email).First(&existingUser)
		if result.RowsAffected > 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "El correo ya está registrado"})
			return
		}

		// Generar hash de la contraseña
		hashedPassword, err := utils.HashPassword(req.Password)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al procesar la contraseña"})
			return
		}

		// Crear nuevo usuario
		now := time.Now()
		user := models.User{
			Username:     req.Username,
			Email:        req.Email,
			PasswordHash: hashedPassword,
			CreatedAt:    now,
			UpdatedAt:    now,
		}

		result = database.DB.Create(&user)
		if result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al crear el usuario: " + result.Error.Error()})
			return
		}

		// Generar token JWT
		token, err := utils.GenerateToken(user.ID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al generar el token"})
			return
		}

		// Preparar respuesta
		response := models.TokenResponse{
			Token: token,
			User: models.UserResponse{
				ID:        user.ID,
				Username:  user.Username,
				Email:     user.Email,
				CreatedAt: user.CreatedAt,
			},
		}

		// Enviar respuesta
		c.JSON(http.StatusCreated, response)
	}
}

func LoginHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req models.UserLoginRequest

		// Bind JSON del cuerpo de la solicitud
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Error al decodificar la solicitud"})
			return
		}

		// Validar campos requeridos
		if req.Email == "" || req.Password == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Email y contraseña son requeridos"})
			return
		}

		// Buscar el usuario por email
		var user models.User
		result := database.DB.Where("email = ?", req.Email).First(&user)
		if result.Error != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Email o contraseña incorrectos"})
			return
		}

		// Verificar contraseña
		if !utils.CheckPasswordHash(req.Password, user.PasswordHash) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Email o contraseña incorrectos"})
			return
		}

		// Generar token JWT
		token, err := utils.GenerateToken(user.ID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al generar el token"})
			return
		}

		// Preparar respuesta
		response := models.TokenResponse{
			Token: token,
			User:  user.ToUserResponse(),
		}

		// Enviar respuesta
		c.JSON(http.StatusOK, response)
	}
}

func GetProfileHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Obtener ID del usuario del contexto de Gin
		userID, ok := middleware.GetUserIDFromGin(c)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuario no autenticado"})
			return
		}

		// Buscar el usuario por ID
		var user models.User
		result := database.DB.First(&user, userID)
		if result.Error != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Usuario no encontrado"})
			return
		}

		// Enviar respuesta
		c.JSON(http.StatusOK, user.ToUserResponse())
	}
}

func UpdateUserHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Obtener ID del usuario del contexto de Gin
		userID, ok := middleware.GetUserIDFromGin(c)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuario no autenticado"})
			return
		}

		// Decodificar la solicitud
		var req struct {
			Username    string `json:"username,omitempty"`
			Email       string `json:"email,omitempty"`
			NewPassword string `json:"new_password,omitempty"`
			OldPassword string `json:"old_password,omitempty"` // Requerido para cambiar contraseña
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Error al decodificar solicitud"})
			return
		}

		// Verificar si se proporciona suficiente información
		if req.Username == "" && req.Email == "" && req.NewPassword == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "No se proporcionaron datos para actualizar"})
			return
		}

		// Buscar el usuario actual
		var user models.User
		if err := database.DB.First(&user, userID).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al encontrar usuario"})
			return
		}

		// Iniciar transacción
		tx := database.DB.Begin()

		// SIEMPRE verificar contraseña actual para cualquier cambio
		if req.OldPassword == "" {
			tx.Rollback()
			c.JSON(http.StatusBadRequest, gin.H{"error": "Se requiere la contraseña actual para realizar cambios"})
			return
		}

		// Verificar contraseña actual
		if !utils.CheckPasswordHash(req.OldPassword, user.PasswordHash) {
			tx.Rollback()
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Contraseña actual incorrecta"})
			return
		}

		// Si se quiere cambiar la contraseña, generar nuevo hash
		if req.NewPassword != "" {
			newHash, err := utils.HashPassword(req.NewPassword)
			if err != nil {
				tx.Rollback()
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al procesar nueva contraseña"})
				return
			}

			user.PasswordHash = newHash
		}

		// Actualizar nombre de usuario si se proporciona
		if req.Username != "" {
			user.Username = req.Username
		}

		// Actualizar email si se proporciona
		if req.Email != "" {
			// Verificar que el email no exista ya
			var existingUser models.User
			result := database.DB.Where("email = ? AND id != ?", req.Email, userID).First(&existingUser)
			if result.RowsAffected > 0 {
				tx.Rollback()
				c.JSON(http.StatusBadRequest, gin.H{"error": "El email ya está en uso"})
				return
			}

			user.Email = req.Email
		}

		// Actualizar fecha de actualización
		user.UpdatedAt = time.Now()

		// Guardar cambios
		if err := tx.Save(&user).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al actualizar usuario: " + err.Error()})
			return
		}

		// Confirmar transacción
		if err := tx.Commit().Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al confirmar cambios"})
			return
		}

		// Enviar respuesta con datos actualizados
		c.JSON(http.StatusOK, user.ToUserResponse())
	}
}

func DeleteUserHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Obtener ID del usuario del contexto de Gin
		userID, ok := middleware.GetUserIDFromGin(c)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuario no autenticado"})
			return
		}

		// Decodificar la solicitud para confirmar con contraseña
		var req struct {
			Password string `json:"password"` // Requerido para confirmar eliminación
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Error al decodificar solicitud"})
			return
		}

		if req.Password == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Se requiere contraseña para confirmar eliminación"})
			return
		}

		// Buscar el usuario
		var user models.User
		if err := database.DB.First(&user, userID).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al encontrar usuario"})
			return
		}

		// Verificar contraseña
		if !utils.CheckPasswordHash(req.Password, user.PasswordHash) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Contraseña incorrecta"})
			return
		}

		// Iniciar transacción
		tx := database.DB.Begin()

		// Eliminar usuario (esto también eliminará todos los registros relacionados gracias a las restricciones de FK)
		if err := tx.Delete(&user).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al eliminar usuario: " + err.Error()})
			return
		}

		// Confirmar transacción
		if err := tx.Commit().Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al confirmar eliminación"})
			return
		}

		// Enviar respuesta
		c.Status(http.StatusNoContent)
	}
}
