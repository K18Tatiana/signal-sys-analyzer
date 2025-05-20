package models

import (
	"time"
)

// ContactForm representa un mensaje del formulario de contacto
type ContactForm struct {
	ID          uint      `gorm:"primaryKey;type:serial" json:"id"`
	UserID      *uint     `gorm:"column:user_id;type:integer" json:"user_id,omitempty"` // Puede ser nulo para usuarios no registrados
	User        *User     `gorm:"foreignKey:UserID" json:"-"`
	Name        string    `gorm:"column:name;size:100;not null" json:"name"`
	Email       string    `gorm:"column:email;size:100;not null" json:"email"`
	Subject     string    `gorm:"column:subject;size:200;not null" json:"subject"`
	Message     string    `gorm:"column:message;type:text;not null" json:"message"`
	SubmittedAt time.Time `gorm:"column:submitted_at;type:timestamp with time zone;default:CURRENT_TIMESTAMP" json:"submitted_at"`
	IsResponded bool      `gorm:"column:is_responded;default:false" json:"is_responded"`
}

// ContactFormRequest para enviar un formulario de contacto
type ContactFormRequest struct {
	Name    string `json:"name"`
	Email   string `json:"email"`
	Subject string `json:"subject"`
	Message string `json:"message"`
}

// FeedbackForm representa un mensaje de feedback
type FeedbackForm struct {
	ID          uint      `gorm:"primaryKey;type:serial" json:"id"`
	UserID      *uint     `gorm:"column:user_id;type:integer" json:"user_id,omitempty"` // Puede ser nulo para usuarios no registrados
	User        *User     `gorm:"foreignKey:UserID" json:"-"`
	Rating      int       `gorm:"column:rating;not null;check:rating >= 1 AND rating <= 5" json:"rating"`
	Feedback    string    `gorm:"column:feedback;type:text;not null" json:"feedback"`
	Email       string    `gorm:"column:email;size:100" json:"email,omitempty"` // Puede ser 'No proporcionado'
	SubmittedAt time.Time `gorm:"column:submitted_at;type:timestamp with time zone;default:CURRENT_TIMESTAMP" json:"submitted_at"`
}

// FeedbackFormRequest para enviar un formulario de feedback
type FeedbackFormRequest struct {
	Rating   int    `json:"rating"`
	Feedback string `json:"feedback"`
	Email    string `json:"email,omitempty"` // Opcional
}
