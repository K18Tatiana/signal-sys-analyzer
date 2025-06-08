package handlers

import (
	"log"
	"net/http"
	"time"

	"backend/database"
	"backend/middleware"
	"backend/models"
	"backend/utils"
	"github.com/gin-gonic/gin"
)

func CreateContactFormHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Decodificar el cuerpo de la solicitud
		var req models.ContactFormRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Error al decodificar la solicitud: " + err.Error()})
			return
		}

		// Validar campos requeridos
		if req.Name == "" || req.Email == "" || req.Subject == "" || req.Message == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Todos los campos son requeridos"})
			return
		}

		// Obtener ID del usuario del contexto (si está autenticado)
		var userID *uint
		if id, ok := middleware.GetUserIDFromGin(c); ok {
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
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al guardar el formulario: " + err.Error()})
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
		c.JSON(http.StatusCreated, response)
	}
}

func CreateFeedbackFormHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Decodificar el cuerpo de la solicitud
		var req models.FeedbackFormRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Error al decodificar la solicitud: " + err.Error()})
			return
		}

		// Validar campos requeridos
		if req.Feedback == "" || req.Rating < 1 || req.Rating > 5 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "El feedback y un rating válido (1-5) son requeridos"})
			return
		}

		// Si no se proporciona email, usar "No proporcionado"
		if req.Email == "" {
			req.Email = "No proporcionado"
		}

		// Obtener ID del usuario del contexto (si está autenticado)
		var userID *uint
		if id, ok := middleware.GetUserIDFromGin(c); ok {
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
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al guardar el feedback: " + err.Error()})
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
		c.JSON(http.StatusCreated, response)
	}
}
