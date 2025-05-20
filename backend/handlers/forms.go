package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"backend/database"
	"backend/middleware"
	"backend/models"
	"backend/utils"
)

// CreateContactFormHandler maneja la creación de un formulario de contacto
func CreateContactFormHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Decodificar el cuerpo de la solicitud
		var req models.ContactFormRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Error al decodificar la solicitud: "+err.Error(), http.StatusBadRequest)
			return
		}

		// Validar campos requeridos
		if req.Name == "" || req.Email == "" || req.Subject == "" || req.Message == "" {
			http.Error(w, "Todos los campos son requeridos", http.StatusBadRequest)
			return
		}

		// Obtener ID del usuario del contexto (si está autenticado)
		var userID *uint
		if id, ok := middleware.GetUserID(r); ok {
			userID = &id
		}

		// Crear el formulario de contacto
		contactForm := models.ContactForm{
			UserID:      userID,
			Name:        req.Name,
			Email:       req.Email,
			Subject:     req.Subject,
			Message:     req.Message,
			SubmittedAt: time.Now(),
			IsResponded: false,
		}

		// Guardar en la base de datos
		if err := database.DB.Create(&contactForm).Error; err != nil {
			http.Error(w, "Error al guardar el formulario: "+err.Error(), http.StatusInternalServerError)
			return
		}

		// Enviar notificación por correo (en segundo plano)
		go func() {
			if err := utils.SendContactFormNotification(req.Name, req.Email, req.Subject, req.Message); err != nil {
				log.Printf("Error al enviar notificación por correo: %v", err)
			} else {
				log.Printf("Notificación de contacto enviada correctamente para %s", req.Email)
			}
		}()

		// Preparar respuesta
		response := struct {
			ID      uint   `json:"id"`
			Message string `json:"message"`
		}{
			ID:      contactForm.ID,
			Message: "Formulario de contacto enviado exitosamente. Nos pondremos en contacto pronto.",
		}

		// Enviar respuesta
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(response)
	}
}

// CreateFeedbackFormHandler maneja la creación de un formulario de feedback
func CreateFeedbackFormHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Decodificar el cuerpo de la solicitud
		var req models.FeedbackFormRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Error al decodificar la solicitud: "+err.Error(), http.StatusBadRequest)
			return
		}

		// Validar campos requeridos
		if req.Feedback == "" || req.Rating < 1 || req.Rating > 5 {
			http.Error(w, "El feedback y un rating válido (1-5) son requeridos", http.StatusBadRequest)
			return
		}

		// Si no se proporciona email, usar "No proporcionado"
		if req.Email == "" {
			req.Email = "No proporcionado"
		}

		// Obtener ID del usuario del contexto (si está autenticado)
		var userID *uint
		if id, ok := middleware.GetUserID(r); ok {
			userID = &id
		}

		// Crear el formulario de feedback
		feedbackForm := models.FeedbackForm{
			UserID:      userID,
			Rating:      req.Rating,
			Feedback:    req.Feedback,
			Email:       req.Email,
			SubmittedAt: time.Now(),
		}

		// Guardar en la base de datos
		if err := database.DB.Create(&feedbackForm).Error; err != nil {
			http.Error(w, "Error al guardar el feedback: "+err.Error(), http.StatusInternalServerError)
			return
		}

		// Enviar notificación por correo (en segundo plano)
		go func() {
			if err := utils.SendFeedbackNotification(req.Rating, req.Feedback, req.Email); err != nil {
				log.Printf("Error al enviar notificación de feedback por correo: %v", err)
			} else {
				log.Printf("Notificación de feedback enviada correctamente")
			}
		}()

		// Preparar respuesta
		response := struct {
			ID      uint   `json:"id"`
			Message string `json:"message"`
		}{
			ID:      feedbackForm.ID,
			Message: "¡Gracias por tu feedback! Lo valoramos mucho.",
		}

		// Enviar respuesta
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(response)
	}
}
