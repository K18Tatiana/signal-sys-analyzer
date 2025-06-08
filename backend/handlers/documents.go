package handlers

import (
	"fmt"
	"log"
	"net/http"
	"sort"
	"strconv"
	"time"

	"backend/database"
	"backend/middleware"
	"backend/models"
	"backend/utils"
	"github.com/gin-gonic/gin"
)

func UploadDocumentHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Obtener ID del usuario del contexto (si está autenticado)
		var userID *uint
		if id, ok := middleware.GetUserIDFromGin(c); ok {
			userID = &id
		}

		// Limitar el tamaño máximo de la petición
		c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, utils.MaxFileSize)

		// Obtener el archivo del formulario
		file, header, err := c.Request.FormFile("file")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Error al obtener el archivo: " + err.Error()})
			return
		}
		defer file.Close()

		// Crear un manejador de almacenamiento Cloudinary
		cloudStorage := utils.NewCloudinaryStorage()

		// Subir el archivo a Cloudinary
		fileURL, err := cloudStorage.UploadFile(file, header)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al subir el archivo: " + err.Error()})
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
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al guardar el documento: " + err.Error()})
			return
		}

		// Preparar la respuesta
		response := document.ToDocumentResponse(0) // Nuevo documento, sin análisis aún

		// Enviar respuesta exitosa
		c.JSON(http.StatusCreated, response)
	}
}

func GetUserDocumentsHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Obtener ID del usuario del contexto de Gin
		userID, ok := middleware.GetUserIDFromGin(c)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuario no autenticado"})
			return
		}

		// Buscar documentos no eliminados del usuario
		var documents []models.Document
		if err := database.DB.Where("user_id = ? AND is_deleted = ?", userID, false).Find(&documents).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener documentos: " + err.Error()})
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
		c.JSON(http.StatusOK, response)
	}
}

func DeleteDocumentHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Obtener ID del usuario del contexto de Gin
		userID, ok := middleware.GetUserIDFromGin(c)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuario no autenticado"})
			return
		}

		// Obtener ID del documento de la URL
		documentIDStr := c.Param("id")
		documentID, err := strconv.ParseUint(documentIDStr, 10, 32)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "ID de documento inválido"})
			return
		}

		// Verificar que el documento pertenece al usuario y no está ya eliminado
		var document models.Document
		result := database.DB.Where("id = ? AND user_id = ? AND is_deleted = ?", documentID, userID, false).First(&document)
		if result.RowsAffected == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "Documento no encontrado o no pertenece al usuario"})
			return
		}

		// Marcar como eliminado (soft delete)
		document.IsDeleted = true
		if err := database.DB.Save(&document).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al marcar documento como eliminado: " + err.Error()})
			return
		}

		// Enviar respuesta exitosa
		c.Status(http.StatusNoContent)
	}
}

func PermanentDeleteDocumentHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Obtener ID del usuario del contexto de Gin
		userID, ok := middleware.GetUserIDFromGin(c)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuario no autenticado"})
			return
		}

		// Obtener ID del documento de la URL
		documentIDStr := c.Param("id")
		documentID, err := strconv.ParseUint(documentIDStr, 10, 32)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "ID de documento inválido"})
			return
		}

		// Buscar el documento
		var document models.Document
		result := database.DB.Where("id = ? AND user_id = ?", documentID, userID).First(&document)
		if result.Error != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Documento no encontrado o no pertenece al usuario"})
			return
		}

		// Guardar la URL para eliminar el archivo después
		fileURL := document.FilePath

		// Iniciar transacción
		tx := database.DB.Begin()

		// Eliminar todos los resultados asociados con análisis de este documento
		if err := tx.Where("analysis_request_id IN (SELECT id FROM analysis_requests WHERE document_id = ?)", document.ID).Delete(&models.Result{}).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al eliminar resultados asociados: " + err.Error()})
			return
		}

		// Eliminar todos los análisis asociados con este documento
		if err := tx.Where("document_id = ?", document.ID).Delete(&models.AnalysisRequest{}).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al eliminar análisis asociados: " + err.Error()})
			return
		}

		// Eliminar el documento
		if err := tx.Delete(&document).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al eliminar documento: " + err.Error()})
			return
		}

		// Confirmar transacción
		if err := tx.Commit().Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al confirmar eliminación: " + err.Error()})
			return
		}

		// Eliminar el archivo de Cloudinary
		cloudStorage := utils.NewCloudinaryStorage()
		if err := cloudStorage.DeleteFile(fileURL); err != nil {
			// Solo registrar el error, pero no fallar la respuesta
			log.Printf("Error al eliminar archivo de Cloudinary %s: %v", fileURL, err)
		}

		// Enviar respuesta exitosa
		c.Status(http.StatusNoContent)
	}
}

// GetDocumentWithAnalysisHandler obtiene un documento específico con todos sus análisis
func GetDocumentWithAnalysisHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Obtener ID del usuario del contexto
		userID, ok := middleware.GetUserIDFromGin(c)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuario no autenticado"})
			return
		}

		// Obtener ID del documento de la URL
		documentIDStr := c.Param("id")
		documentID, err := strconv.ParseUint(documentIDStr, 10, 32)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "ID de documento inválido"})
			return
		}

		// Verificar que el documento pertenece al usuario y no está eliminado
		var document models.Document
		result := database.DB.Where("id = ? AND user_id = ? AND is_deleted = ?", documentID, userID, false).First(&document)
		if result.RowsAffected == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "Documento no encontrado o no pertenece al usuario"})
			return
		}

		// Obtener todos los análisis de este documento (ordenados por fecha, más reciente primero)
		var analyses []models.AnalysisRequest
		if err := database.DB.Where("document_id = ?", documentID).Order("created_at DESC").Find(&analyses).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener análisis: " + err.Error()})
			return
		}

		// Para cada análisis, obtener su resultado (sin filtrar por is_latest)
		type AnalysisWithResult struct {
			Analysis models.AnalysisRequest `json:"analysis"`
			Result   *models.Result         `json:"result,omitempty"`
		}

		var analysesWithResults []AnalysisWithResult
		for _, analysis := range analyses {
			var result models.Result
			// CAMBIO CLAVE: Buscar cualquier resultado para este analysis_request_id
			// Sin filtrar por is_latest, para mostrar TODOS los análisis
			err := database.DB.Where("analysis_request_id = ?", analysis.ID).
				Order("created_at DESC").
				First(&result).Error

			analysisWithResult := AnalysisWithResult{
				Analysis: analysis,
			}

			// Si encontramos un resultado, agregarlo
			if err == nil {
				analysisWithResult.Result = &result
			}

			analysesWithResults = append(analysesWithResults, analysisWithResult)
		}

		// Preparar respuesta completa
		response := struct {
			Document models.DocumentResponse `json:"document"`
			Analyses []AnalysisWithResult    `json:"analyses"`
		}{
			Document: document.ToDocumentResponse(len(analyses)),
			Analyses: analysesWithResults,
		}

		c.JSON(http.StatusOK, response)
	}
}

// GetUserStatsHandler obtiene estadísticas del usuario
func GetUserStatsHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Obtener ID del usuario del contexto
		userID, ok := middleware.GetUserIDFromGin(c)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuario no autenticado"})
			return
		}

		// Contar documentos del usuario (no eliminados)
		var documentsCount int64
		if err := database.DB.Model(&models.Document{}).Where("user_id = ? AND is_deleted = ?", userID, false).Count(&documentsCount).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al contar documentos"})
			return
		}

		// Contar análisis del usuario
		var analysisCount int64
		if err := database.DB.Raw(`
			SELECT COUNT(*) 
			FROM analysis_requests ar 
			JOIN documents d ON ar.document_id = d.id 
			WHERE d.user_id = ? AND d.is_deleted = false
		`, userID).Scan(&analysisCount).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al contar análisis"})
			return
		}

		// Preparar respuesta
		stats := gin.H{
			"documents_count": documentsCount,
			"analysis_count":  analysisCount,
		}

		c.JSON(http.StatusOK, stats)
	}
}

// GetUserRecentActivityHandler obtiene la actividad reciente del usuario
func GetUserRecentActivityHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Obtener ID del usuario del contexto
		userID, ok := middleware.GetUserIDFromGin(c)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuario no autenticado"})
			return
		}

		type ActivityItem struct {
			Type        string    `json:"type"`
			Description string    `json:"description"`
			Date        time.Time `json:"date"`
		}

		var activities []ActivityItem

		// Obtener documentos recientes (últimos 10)
		var recentDocuments []models.Document
		if err := database.DB.Where("user_id = ? AND is_deleted = ?", userID, false).
			Order("upload_date DESC").
			Limit(5).
			Find(&recentDocuments).Error; err == nil {

			for _, doc := range recentDocuments {
				activities = append(activities, ActivityItem{
					Type:        "document",
					Description: fmt.Sprintf("Documento subido: %s", doc.OriginalFilename),
					Date:        doc.UploadDate,
				})
			}
		}

		// Obtener análisis recientes (últimos 10)
		var recentAnalyses []struct {
			ID           uint      `json:"id"`
			InputVoltage float64   `json:"input_voltage"`
			Comment      string    `json:"comment"`
			CreatedAt    time.Time `json:"created_at"`
			Filename     string    `json:"filename"`
		}

		if err := database.DB.Raw(`
			SELECT ar.id, ar.input_voltage, ar.comment, ar.created_at, d.original_filename as filename
			FROM analysis_requests ar
			JOIN documents d ON ar.document_id = d.id
			WHERE d.user_id = ? AND d.is_deleted = false
			ORDER BY ar.created_at DESC
			LIMIT 5
		`, userID).Scan(&recentAnalyses).Error; err == nil {

			for _, analysis := range recentAnalyses {
				description := fmt.Sprintf("Análisis realizado: %s (%.1fV)", analysis.Filename, analysis.InputVoltage)
				if analysis.Comment != "" {
					description += fmt.Sprintf(" - %s", analysis.Comment)
				}

				activities = append(activities, ActivityItem{
					Type:        "analysis",
					Description: description,
					Date:        analysis.CreatedAt,
				})
			}
		}

		// Ordenar todas las actividades por fecha (más reciente primero)
		sort.Slice(activities, func(i, j int) bool {
			return activities[i].Date.After(activities[j].Date)
		})

		// Limitar a las 8 más recientes
		if len(activities) > 8 {
			activities = activities[:8]
		}

		c.JSON(http.StatusOK, activities)
	}
}

// GetUserRecentDocumentsHandler obtiene los documentos recientes del usuario
func GetUserRecentDocumentsHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Obtener ID del usuario del contexto
		userID, ok := middleware.GetUserIDFromGin(c)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuario no autenticado"})
			return
		}

		// Obtener documentos recientes con conteo de análisis
		var documents []struct {
			ID               uint      `json:"id"`
			OriginalFilename string    `json:"name"`
			UploadDate       time.Time `json:"date"`
			AnalysisCount    int       `json:"analysis_count"`
		}

		err := database.DB.Raw(`
			SELECT 
				d.id,
				d.original_filename,
				d.upload_date,
				COALESCE(COUNT(ar.id), 0) as analysis_count
			FROM documents d
			LEFT JOIN analysis_requests ar ON d.id = ar.document_id
			WHERE d.user_id = ? AND d.is_deleted = false
			GROUP BY d.id, d.original_filename, d.upload_date
			ORDER BY d.upload_date DESC
			LIMIT 5
		`, userID).Scan(&documents).Error

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener documentos recientes"})
			return
		}

		c.JSON(http.StatusOK, documents)
	}
}
