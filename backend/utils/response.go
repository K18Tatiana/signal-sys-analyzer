package utils

import (
	"encoding/json"
	"net/http"
)

// Response es una estructura para respuestas HTTP estándar
type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

// RespondWithJSON envía una respuesta HTTP JSON estándar
func RespondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	// Preparar respuesta
	response, err := json.Marshal(payload)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"success":false,"error":"Error al procesar la respuesta"}`))
		return
	}

	// Enviar respuesta
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

// RespondWithError envía una respuesta de error HTTP JSON
func RespondWithError(w http.ResponseWriter, code int, message string) {
	RespondWithJSON(w, code, Response{
		Success: false,
		Error:   message,
	})
}

// RespondWithSuccess envía una respuesta de éxito HTTP JSON
func RespondWithSuccess(w http.ResponseWriter, code int, data interface{}, message string) {
	RespondWithJSON(w, code, Response{
		Success: true,
		Data:    data,
		Message: message,
	})
}
