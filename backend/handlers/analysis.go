package handlers

import (
	"encoding/csv"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"

	"backend/database"
	"backend/middleware"
	"backend/models"
	"backend/utils"
	"github.com/gorilla/mux"
	"gorm.io/datatypes"
)

// CreateAnalysisRequestHandler maneja la creación de una solicitud de análisis
func CreateAnalysisRequestHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Decodificar el cuerpo de la solicitud
		var req models.AnalysisRequestCreate
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Error al decodificar la solicitud: "+err.Error(), http.StatusBadRequest)
			return
		}

		// Verificar que el documento existe y no está eliminado
		var document models.Document
		result := database.DB.Where("id = ? AND is_deleted = ?", req.DocumentID, false).First(&document)
		if result.Error != nil {
			http.Error(w, "Documento no encontrado o ha sido eliminado", http.StatusNotFound)
			return
		}

		// Si el documento pertenece a un usuario registrado, verificar si el usuario actual
		// es el propietario (solo aplicable si hay un usuario autenticado)
		if document.UserID != nil {
			// Si hay un usuario autenticado
			if currentUserID, ok := middleware.GetUserID(r); ok {
				// Si el documento pertenece a otro usuario, denegar acceso
				if *document.UserID != currentUserID {
					http.Error(w, "No tienes permiso para analizar este documento", http.StatusForbidden)
					return
				}
			} else {
				// Si el documento pertenece a un usuario, pero no hay usuario autenticado
				http.Error(w, "Este documento solo puede ser analizado por su propietario", http.StatusForbidden)
				return
			}
		}

		// Para usuarios no registrados, verificar si ya hay un análisis existente
		// (solo pueden hacer un análisis por documento)
		if _, ok := middleware.GetUserID(r); !ok {
			var count int64
			database.DB.Model(&models.AnalysisRequest{}).Where("document_id = ?", req.DocumentID).Count(&count)
			if count > 0 {
				http.Error(w, "Los usuarios no registrados solo pueden realizar un análisis por documento", http.StatusForbidden)
				return
			}
		}

		// Crear la solicitud de análisis
		analysis := models.AnalysisRequest{
			DocumentID:   req.DocumentID,
			InputVoltage: req.InputVoltage,
			IsProcessed:  false,
			CreatedAt:    time.Now(),
		}

		if err := database.DB.Create(&analysis).Error; err != nil {
			http.Error(w, "Error al crear solicitud de análisis: "+err.Error(), http.StatusInternalServerError)
			return
		}

		// Aquí normalmente enviarías la solicitud a un worker/cola para su procesamiento
		// Para este ejemplo, simularemos el procesamiento inmediatamente
		go processAnalysisRequest(analysis.ID, req.DocumentID, req.InputVoltage)

		// Preparar respuesta
		response := struct {
			ID           uint    `json:"id"`
			DocumentID   uint    `json:"document_id"`
			InputVoltage float64 `json:"input_voltage"`
			Message      string  `json:"message"`
		}{
			ID:           analysis.ID,
			DocumentID:   analysis.DocumentID,
			InputVoltage: analysis.InputVoltage,
			Message:      "Solicitud de análisis creada. El procesamiento comenzará en breve.",
		}

		// Enviar respuesta exitosa
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(response)
	}
}

// GetAnalysisResultHandler obtiene el resultado de un análisis específico
func GetAnalysisResultHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Obtener ID del análisis de la URL
		vars := mux.Vars(r)
		analysisID, err := strconv.ParseUint(vars["id"], 10, 32)
		if err != nil {
			http.Error(w, "ID de análisis inválido", http.StatusBadRequest)
			return
		}

		// Buscar el análisis
		var analysis models.AnalysisRequest
		if err := database.DB.First(&analysis, analysisID).Error; err != nil {
			http.Error(w, "Análisis no encontrado", http.StatusNotFound)
			return
		}

		// Buscar el documento asociado
		var document models.Document
		if err := database.DB.First(&document, analysis.DocumentID).Error; err != nil {
			http.Error(w, "Error al obtener información del documento", http.StatusInternalServerError)
			return
		}

		// Verificar si el documento está eliminado
		if document.IsDeleted {
			http.Error(w, "El documento asociado ha sido eliminado", http.StatusNotFound)
			return
		}

		// Verificar si el usuario tiene permisos para ver este análisis
		if document.UserID != nil {
			// Si hay un usuario autenticado, verificar que sea el propietario
			if currentUserID, ok := middleware.GetUserID(r); ok {
				if *document.UserID != currentUserID {
					http.Error(w, "No tienes permiso para ver este análisis", http.StatusForbidden)
					return
				}
			} else {
				// Si el documento pertenece a un usuario pero no hay usuario autenticado
				http.Error(w, "Este análisis solo puede ser visto por su propietario", http.StatusForbidden)
				return
			}
		}

		// Contar el número de análisis para este documento
		var analysisCount int64
		database.DB.Model(&models.AnalysisRequest{}).Where("document_id = ?", document.ID).Count(&analysisCount)

		// Crear respuesta con documento
		documentResponse := document.ToDocumentResponse(int(analysisCount))

		// Buscar el resultado más reciente para este análisis
		var result models.Result
		err = database.DB.Where("analysis_request_id = ?", analysisID).Where("is_latest = ?", true).First(&result).Error

		// Verificar si el resultado está disponible
		if err != nil {
			// El análisis está en proceso, enviar estado pendiente
			response := struct {
				AnalysisRequest models.AnalysisRequest  `json:"analysis_request"`
				Document        models.DocumentResponse `json:"document"`
				Status          string                  `json:"status"`
			}{
				AnalysisRequest: analysis,
				Document:        documentResponse,
				Status:          "pending",
			}

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(response)
			return
		}

		// Resultado disponible, enviar respuesta completa
		resultResponse := models.ResultResponse{
			Result:          result,
			AnalysisRequest: analysis,
			Document:        documentResponse,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resultResponse)
	}
}

// GetUserAnalysisRequestsHandler obtiene todos los análisis del usuario autenticado
func GetUserAnalysisRequestsHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Obtener ID del usuario del contexto
		userID, ok := middleware.GetUserID(r)
		if !ok {
			http.Error(w, "Usuario no autenticado", http.StatusUnauthorized)
			return
		}

		// Buscar documentos del usuario
		var documents []models.Document
		if err := database.DB.Where("user_id = ? AND is_deleted = ?", userID, false).Find(&documents).Error; err != nil {
			http.Error(w, "Error al obtener documentos del usuario", http.StatusInternalServerError)
			return
		}

		// Si no hay documentos, devolver lista vacía
		if len(documents) == 0 {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode([]struct{}{})
			return
		}

		// Extraer IDs de documentos
		var documentIDs []uint
		for _, doc := range documents {
			documentIDs = append(documentIDs, doc.ID)
		}

		// Buscar análisis para estos documentos
		type AnalysisWithFilename struct {
			ID           uint      `json:"id"`
			DocumentID   uint      `json:"document_id"`
			InputVoltage float64   `json:"input_voltage"`
			IsProcessed  bool      `json:"is_processed"`
			CreatedAt    time.Time `json:"created_at"`
			Filename     string    `json:"filename"`
		}

		// Ejecutar consulta raw para obtener análisis con nombres de archivos
		var analyses []AnalysisWithFilename

		rows, err := database.DB.Raw(`
			SELECT ar.id, ar.document_id, ar.input_voltage, ar.is_processed, ar.created_at,
				   d.original_filename as filename
			FROM analysis_requests ar
			JOIN documents d ON ar.document_id = d.id
			WHERE d.user_id = ? AND d.is_deleted = FALSE
			ORDER BY ar.created_at DESC
		`, userID).Rows()

		if err != nil {
			http.Error(w, "Error al obtener análisis: "+err.Error(), http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		// Recorrer resultados
		for rows.Next() {
			var analysis AnalysisWithFilename
			err := rows.Scan(
				&analysis.ID, &analysis.DocumentID, &analysis.InputVoltage,
				&analysis.IsProcessed, &analysis.CreatedAt, &analysis.Filename,
			)
			if err != nil {
				http.Error(w, "Error al procesar análisis: "+err.Error(), http.StatusInternalServerError)
				return
			}
			analyses = append(analyses, analysis)
		}

		// Verificar errores del iterador
		if err := rows.Err(); err != nil {
			http.Error(w, "Error al iterar análisis: "+err.Error(), http.StatusInternalServerError)
			return
		}

		// Enviar respuesta
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(analyses)
	}
}

// processAnalysisRequest procesa una solicitud de análisis
func processAnalysisRequest(analysisID, documentID uint, inputVoltage float64) {
	// Esperar un poco para simular procesamiento
	time.Sleep(2 * time.Second)

	// Obtener la URL del archivo
	var document models.Document
	if err := database.DB.First(&document, documentID).Error; err != nil {
		log.Printf("Error al obtener documento: %v", err)
		return
	}

	// Crear un manejador de almacenamiento en Cloudinary
	cloudStorage := utils.NewCloudinaryStorage()

	// Obtener el archivo de Cloudinary
	fileReader, err := cloudStorage.GetFile(document.FilePath)
	if err != nil {
		log.Printf("Error al obtener archivo de Cloudinary: %v", err)
		return
	}
	defer fileReader.Close()

	// Aquí implementarías tu lógica real de ML para procesar el CSV
	// Por ahora, leeremos el archivo para simular el procesamiento
	reader := csv.NewReader(fileReader)

	// Leer encabezados
	_, err = reader.Read()
	if err != nil {
		log.Printf("Error al leer encabezados CSV: %v", err)
		return
	}

	// Leer datos (simplemente contamos las filas para este ejemplo)
	var dataPoints [][]string
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Printf("Error al leer datos CSV: %v", err)
			continue
		}
		dataPoints = append(dataPoints, record)
	}

	log.Printf("Leídos %d puntos de datos del archivo", len(dataPoints))

	// Simulación simple de identificación de sistema y polos
	systemType := "Sistema de segundo orden subamortiguado"

	// Crear un JSON con los polos simulados
	polesData := map[string]interface{}{
		"polos": []map[string]float64{
			{"real": -0.5, "imag": 0.866},
			{"real": -0.5, "imag": -0.866},
		},
	}

	// Datos adicionales
	rawData := map[string]interface{}{
		"ganancia":               2.0,
		"frecuencia_natural":     1.0,
		"factor_amortiguamiento": 0.5,
		"voltaje_entrada":        inputVoltage,
		"puntos_analizados":      len(dataPoints),
	}

	// Convertir a JSON
	polesJSON, _ := json.Marshal(polesData)
	rawDataJSON, _ := json.Marshal(rawData)

	// Comenzar transacción
	tx := database.DB.Begin()

	// Marcar todos los resultados previos como no-latest
	if err := tx.Model(&models.Result{}).Where(
		"analysis_request_id IN (SELECT id FROM analysis_requests WHERE document_id = ?)", documentID,
	).Update("is_latest", false).Error; err != nil {
		tx.Rollback()
		log.Printf("Error al actualizar resultados previos: %v", err)
		return
	}

	// Insertar el resultado
	result := models.Result{
		AnalysisRequestID: analysisID,
		SystemType:        systemType,
		Poles:             datatypes.JSON(polesJSON),
		RawData:           datatypes.JSON(rawDataJSON),
		IsLatest:          true,
		CreatedAt:         time.Now(),
	}

	if err := tx.Create(&result).Error; err != nil {
		tx.Rollback()
		log.Printf("Error al guardar resultado: %v", err)
		return
	}

	// Marcar la solicitud como procesada
	if err := tx.Model(&models.AnalysisRequest{}).Where("id = ?", analysisID).Update("is_processed", true).Error; err != nil {
		tx.Rollback()
		log.Printf("Error al actualizar solicitud: %v", err)
		return
	}

	// Confirmar transacción
	if err := tx.Commit().Error; err != nil {
		log.Printf("Error al confirmar transacción: %v", err)
		return
	}

	log.Printf("Análisis %d completado exitosamente", analysisID)
}
