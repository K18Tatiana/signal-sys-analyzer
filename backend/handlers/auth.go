package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"backend/database"
	"backend/middleware"
	"backend/models"
	"backend/utils"
)

// RegisterHandler maneja el registro de nuevos usuarios
func RegisterHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Decodificar el cuerpo de la solicitud
		var req models.UserRegisterRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Error al decodificar la solicitud", http.StatusBadRequest)
			return
		}

		// Validar campos requeridos
		if req.Username == "" || req.Email == "" || req.Password == "" {
			http.Error(w, "Todos los campos son requeridos", http.StatusBadRequest)
			return
		}

		// Verificar si el correo ya existe
		var existingUser models.User
		result := database.DB.Where("email = ?", req.Email).First(&existingUser)
		if result.RowsAffected > 0 {
			http.Error(w, "El correo ya está registrado", http.StatusBadRequest)
			return
		}

		// Generar hash de la contraseña
		hashedPassword, err := utils.HashPassword(req.Password)
		if err != nil {
			http.Error(w, "Error al procesar la contraseña", http.StatusInternalServerError)
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
			http.Error(w, "Error al crear el usuario: "+result.Error.Error(), http.StatusInternalServerError)
			return
		}

		// Generar token JWT
		token, err := utils.GenerateToken(user.ID)
		if err != nil {
			http.Error(w, "Error al generar el token", http.StatusInternalServerError)
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
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(response)
	}
}

// LoginHandler maneja el inicio de sesión de usuarios
func LoginHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Decodificar el cuerpo de la solicitud
		var req models.UserLoginRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Error al decodificar la solicitud", http.StatusBadRequest)
			return
		}

		// Validar campos requeridos
		if req.Email == "" || req.Password == "" {
			http.Error(w, "Email y contraseña son requeridos", http.StatusBadRequest)
			return
		}

		// Buscar el usuario por email
		var user models.User
		result := database.DB.Where("email = ?", req.Email).First(&user)
		if result.Error != nil {
			http.Error(w, "Email o contraseña incorrectos", http.StatusUnauthorized)
			return
		}

		// Verificar contraseña
		if !utils.CheckPasswordHash(req.Password, user.PasswordHash) {
			http.Error(w, "Email o contraseña incorrectos", http.StatusUnauthorized)
			return
		}

		// Generar token JWT
		token, err := utils.GenerateToken(user.ID)
		if err != nil {
			http.Error(w, "Error al generar el token", http.StatusInternalServerError)
			return
		}

		// Preparar respuesta
		response := models.TokenResponse{
			Token: token,
			User:  user.ToUserResponse(),
		}

		// Enviar respuesta
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}
}

// GetProfileHandler obtiene el perfil del usuario actual
func GetProfileHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Obtener ID del usuario del contexto (establecido por AuthMiddleware)
		userID, ok := middleware.GetUserID(r)
		if !ok {
			http.Error(w, "Usuario no autenticado", http.StatusUnauthorized)
			return
		}

		// Buscar el usuario por ID
		var user models.User
		result := database.DB.First(&user, userID)
		if result.Error != nil {
			http.Error(w, "Usuario no encontrado", http.StatusNotFound)
			return
		}

		// Enviar respuesta
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(user.ToUserResponse())
	}
}

// UpdateUserHandler maneja la actualización de datos de usuario
func UpdateUserHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Obtener ID del usuario del contexto (debe estar autenticado)
		userID, ok := middleware.GetUserID(r)
		if !ok {
			http.Error(w, "Usuario no autenticado", http.StatusUnauthorized)
			return
		}

		// Decodificar la solicitud
		var req struct {
			Username    string `json:"username,omitempty"`
			Email       string `json:"email,omitempty"`
			NewPassword string `json:"new_password,omitempty"`
			OldPassword string `json:"old_password,omitempty"` // Requerido para cambiar contraseña
		}

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Error al decodificar solicitud", http.StatusBadRequest)
			return
		}

		// Verificar si se proporciona suficiente información
		if req.Username == "" && req.Email == "" && req.NewPassword == "" {
			http.Error(w, "No se proporcionaron datos para actualizar", http.StatusBadRequest)
			return
		}

		// Buscar el usuario actual
		var user models.User
		if err := database.DB.First(&user, userID).Error; err != nil {
			http.Error(w, "Error al encontrar usuario", http.StatusInternalServerError)
			return
		}

		// Iniciar transacción
		tx := database.DB.Begin()

		// Si se quiere cambiar la contraseña, verificar la antigua
		if req.NewPassword != "" {
			if req.OldPassword == "" {
				tx.Rollback()
				http.Error(w, "Se requiere la contraseña actual para cambiarla", http.StatusBadRequest)
				return
			}

			// Verificar contraseña actual
			if !utils.CheckPasswordHash(req.OldPassword, user.PasswordHash) {
				tx.Rollback()
				http.Error(w, "Contraseña actual incorrecta", http.StatusUnauthorized)
				return
			}

			// Generar nuevo hash
			newHash, err := utils.HashPassword(req.NewPassword)
			if err != nil {
				tx.Rollback()
				http.Error(w, "Error al procesar nueva contraseña", http.StatusInternalServerError)
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
				http.Error(w, "El email ya está en uso", http.StatusBadRequest)
				return
			}

			user.Email = req.Email
		}

		// Actualizar fecha de actualización
		user.UpdatedAt = time.Now()

		// Guardar cambios
		if err := tx.Save(&user).Error; err != nil {
			tx.Rollback()
			http.Error(w, "Error al actualizar usuario: "+err.Error(), http.StatusInternalServerError)
			return
		}

		// Confirmar transacción
		if err := tx.Commit().Error; err != nil {
			http.Error(w, "Error al confirmar cambios", http.StatusInternalServerError)
			return
		}

		// Enviar respuesta con datos actualizados
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(user.ToUserResponse())
	}
}

// DeleteUserHandler maneja la eliminación de la cuenta de usuario
func DeleteUserHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Obtener ID del usuario del contexto
		userID, ok := middleware.GetUserID(r)
		if !ok {
			http.Error(w, "Usuario no autenticado", http.StatusUnauthorized)
			return
		}

		// Decodificar la solicitud para confirmar con contraseña
		var req struct {
			Password string `json:"password"` // Requerido para confirmar eliminación
		}

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Error al decodificar solicitud", http.StatusBadRequest)
			return
		}

		if req.Password == "" {
			http.Error(w, "Se requiere contraseña para confirmar eliminación", http.StatusBadRequest)
			return
		}

		// Buscar el usuario
		var user models.User
		if err := database.DB.First(&user, userID).Error; err != nil {
			http.Error(w, "Error al encontrar usuario", http.StatusInternalServerError)
			return
		}

		// Verificar contraseña
		if !utils.CheckPasswordHash(req.Password, user.PasswordHash) {
			http.Error(w, "Contraseña incorrecta", http.StatusUnauthorized)
			return
		}

		// Iniciar transacción
		tx := database.DB.Begin()

		// Eliminar usuario (esto también eliminará todos los registros relacionados gracias a las restricciones de FK)
		if err := tx.Delete(&user).Error; err != nil {
			tx.Rollback()
			http.Error(w, "Error al eliminar usuario: "+err.Error(), http.StatusInternalServerError)
			return
		}

		// Confirmar transacción
		if err := tx.Commit().Error; err != nil {
			http.Error(w, "Error al confirmar eliminación", http.StatusInternalServerError)
			return
		}

		// Enviar respuesta
		w.WriteHeader(http.StatusNoContent)
	}
}
