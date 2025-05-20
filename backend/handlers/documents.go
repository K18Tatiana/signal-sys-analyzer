package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"backend/database"
	"backend/middleware"
	"backend/models"
	"backend/utils"
	"github.com/gorilla/mux"
)

// UploadDocumentHandler maneja la subida de archivos CSV
func UploadDocumentHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Obtener ID del usuario del contexto (si está autenticado)
		var userID *uint
		if id, ok := middleware.GetUserID(r); ok {
			userID = &id
		}

		// Limitar el tamaño máximo de la petición
		r.Body = http.MaxBytesReader(w, r.Body, utils.MaxFileSize)

		// Parsear el formulario multipart
		if err := r.ParseMultipartForm(utils.MaxFileSize); err != nil {
			http.Error(w, "Error al procesar el formulario: "+err.Error(), http.StatusBadRequest)
			return
		}

		// Obtener el archivo del formulario
		file, header, err := r.FormFile("file")
		if err != nil {
			http.Error(w, "Error al obtener el archivo: "+err.Error(), http.StatusBadRequest)
			return
		}
		defer file.Close()

		// Crear un manejador de almacenamiento Cloudinary
		cloudStorage := utils.NewCloudinaryStorage()

		// Subir el archivo a Cloudinary
		fileURL, err := cloudStorage.UploadFile(file, header)
		if err != nil {
			http.Error(w, "Error al subir el archivo: "+err.Error(), http.StatusInternalServerError)
			return
		}

		log.Printf("Archivo subido exitosamente a Cloudinary: %s", fileURL)

		// Crear documento en la base de datos
		document := models.Document{
			UserID:           userID,
			FilePath:         fileURL, // Guardamos la URL de Cloudinary
			OriginalFilename: header.Filename,
			UploadDate:       time.Now(),
			IsDeleted:        false,
		}

		// Guardar documento
		if err := database.DB.Create(&document).Error; err != nil {
			http.Error(w, "Error al guardar el documento: "+err.Error(), http.StatusInternalServerError)
			return
		}

		// Preparar la respuesta
		response := document.ToDocumentResponse(0) // Nuevo documento, sin análisis aún

		// Enviar respuesta exitosa
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(response)
	}
}

// GetUserDocumentsHandler obtiene los documentos no eliminados del usuario autenticado
func GetUserDocumentsHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Obtener ID del usuario del contexto
		userID, ok := middleware.GetUserID(r)
		if !ok {
			http.Error(w, "Usuario no autenticado", http.StatusUnauthorized)
			return
		}

		// Buscar documentos no eliminados del usuario
		var documents []models.Document
		if err := database.DB.Where("user_id = ? AND is_deleted = ?", userID, false).Find(&documents).Error; err != nil {
			http.Error(w, "Error al obtener documentos: "+err.Error(), http.StatusInternalServerError)
			return
		}

		// Preparar la respuesta
		response := []models.DocumentResponse{}
		for _, doc := range documents {
			// Contar análisis para este documento
			var analysisCount int64
			database.DB.Model(&models.AnalysisRequest{}).Where("document_id = ?", doc.ID).Count(&analysisCount)

			// Agregar a la respuesta
			response = append(response, doc.ToDocumentResponse(int(analysisCount)))
		}

		// Enviar respuesta
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}
}

// DeleteDocumentHandler marca un documento como eliminado (soft delete)
func DeleteDocumentHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Obtener ID del usuario del contexto
		userID, ok := middleware.GetUserID(r)
		if !ok {
			http.Error(w, "Usuario no autenticado", http.StatusUnauthorized)
			return
		}

		// Obtener ID del documento de la URL
		vars := mux.Vars(r)
		documentID, err := strconv.ParseUint(vars["id"], 10, 32)
		if err != nil {
			http.Error(w, "ID de documento inválido", http.StatusBadRequest)
			return
		}

		// Verificar que el documento pertenece al usuario y no está ya eliminado
		var document models.Document
		result := database.DB.Where("id = ? AND user_id = ? AND is_deleted = ?", documentID, userID, false).First(&document)
		if result.RowsAffected == 0 {
			http.Error(w, "Documento no encontrado o no pertenece al usuario", http.StatusNotFound)
			return
		}

		// Marcar como eliminado (soft delete)
		document.IsDeleted = true
		if err := database.DB.Save(&document).Error; err != nil {
			http.Error(w, "Error al marcar documento como eliminado: "+err.Error(), http.StatusInternalServerError)
			return
		}

		// Enviar respuesta exitosa
		w.WriteHeader(http.StatusNoContent)
	}
}

// PermanentDeleteDocumentHandler elimina físicamente un documento
func PermanentDeleteDocumentHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Obtener ID del usuario del contexto
		userID, ok := middleware.GetUserID(r)
		if !ok {
			http.Error(w, "Usuario no autenticado", http.StatusUnauthorized)
			return
		}

		// Obtener ID del documento de la URL
		vars := mux.Vars(r)
		documentID, err := strconv.ParseUint(vars["id"], 10, 32)
		if err != nil {
			http.Error(w, "ID de documento inválido", http.StatusBadRequest)
			return
		}

		// Buscar el documento
		var document models.Document
		result := database.DB.Where("id = ? AND user_id = ?", documentID, userID).First(&document)
		if result.Error != nil {
			http.Error(w, "Documento no encontrado o no pertenece al usuario", http.StatusNotFound)
			return
		}

		// Guardar la URL para eliminar el archivo después
		fileURL := document.FilePath

		// Iniciar transacción
		tx := database.DB.Begin()

		// Eliminar todos los resultados asociados con análisis de este documento
		if err := tx.Where("analysis_request_id IN (SELECT id FROM analysis_requests WHERE document_id = ?)", document.ID).Delete(&models.Result{}).Error; err != nil {
			tx.Rollback()
			http.Error(w, "Error al eliminar resultados asociados: "+err.Error(), http.StatusInternalServerError)
			return
		}

		// Eliminar todos los análisis asociados con este documento
		if err := tx.Where("document_id = ?", document.ID).Delete(&models.AnalysisRequest{}).Error; err != nil {
			tx.Rollback()
			http.Error(w, "Error al eliminar análisis asociados: "+err.Error(), http.StatusInternalServerError)
			return
		}

		// Eliminar el documento
		if err := tx.Delete(&document).Error; err != nil {
			tx.Rollback()
			http.Error(w, "Error al eliminar documento: "+err.Error(), http.StatusInternalServerError)
			return
		}

		// Confirmar transacción
		if err := tx.Commit().Error; err != nil {
			http.Error(w, "Error al confirmar eliminación: "+err.Error(), http.StatusInternalServerError)
			return
		}

		// Eliminar el archivo de Cloudinary
		cloudStorage := utils.NewCloudinaryStorage()
		if err := cloudStorage.DeleteFile(fileURL); err != nil {
			// Solo registrar el error, pero no fallar la respuesta
			log.Printf("Error al eliminar archivo de Cloudinary %s: %v", fileURL, err)
		}

		// Enviar respuesta exitosa
		w.WriteHeader(http.StatusNoContent)
	}
}
