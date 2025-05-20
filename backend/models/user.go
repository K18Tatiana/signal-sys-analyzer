package models

import (
	"time"
)

// User representa un usuario registrado en el sistema
type User struct {
	ID            uint           `gorm:"primaryKey;type:serial" json:"id"`
	Username      string         `gorm:"size:50;not null" json:"username"`
	Email         string         `gorm:"size:100;not null;uniqueIndex" json:"email"`
	PasswordHash  string         `gorm:"column:password_hash;size:100;not null" json:"-"`
	CreatedAt     time.Time      `gorm:"column:created_at;type:timestamp with time zone;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt     time.Time      `gorm:"column:updated_at;type:timestamp with time zone;default:CURRENT_TIMESTAMP" json:"updated_at"`
	Documents     []Document     `gorm:"foreignKey:UserID" json:"-"`
	ContactForms  []ContactForm  `gorm:"foreignKey:UserID" json:"-"`
	FeedbackForms []FeedbackForm `gorm:"foreignKey:UserID" json:"-"`
}

// UserLoginRequest para el login de usuarios
type UserLoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// UserRegisterRequest para el registro de usuarios
type UserRegisterRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// UserResponse es la respuesta enviada al cliente (sin datos sensibles)
type UserResponse struct {
	ID        uint      `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

// TokenResponse es la respuesta con el token JWT tras login/registro exitoso
type TokenResponse struct {
	Token string       `json:"token"`
	User  UserResponse `json:"user"`
}

// ToUserResponse convierte un User a UserResponse
func (u *User) ToUserResponse() UserResponse {
	return UserResponse{
		ID:        u.ID,
		Username:  u.Username,
		Email:     u.Email,
		CreatedAt: u.CreatedAt,
	}
}
