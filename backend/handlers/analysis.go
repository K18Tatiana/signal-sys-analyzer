package handlers

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"backend/database"
	"backend/middleware"
	"backend/models"
	"backend/utils"
	"github.com/gin-gonic/gin"
	"gorm.io/datatypes"
)

func CreateAnalysisRequestHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Decodificar el cuerpo de la solicitud
		var req models.AnalysisRequestCreate
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Error al decodificar la solicitud: " + err.Error()})
			return
		}

		// Verificar que el documento existe y no está eliminado
		var document models.Document
		result := database.DB.Where("id = ? AND is_deleted = ?", req.DocumentID, false).First(&document)
		if result.Error != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Documento no encontrado o ha sido eliminado"})
			return
		}

		// Variable para almacenar el comentario (solo si está autenticado)
		var analysisComment string

		// Si el documento pertenece a un usuario registrado, verificar si el usuario actual
		// es el propietario (solo aplicable si hay un usuario autenticado)
		if document.UserID != nil {
			// Si hay un usuario autenticado
			if currentUserID, ok := middleware.GetUserIDFromGin(c); ok {
				// Si el documento pertenece a otro usuario, denegar acceso
				if *document.UserID != currentUserID {
					c.JSON(http.StatusForbidden, gin.H{"error": "No tienes permiso para analizar este documento"})
					return
				}
				// Usuario autenticado y propietario: puede agregar comentario
				analysisComment = req.Comment
			} else {
				// Si el documento pertenece a un usuario, pero no hay usuario autenticado
				c.JSON(http.StatusForbidden, gin.H{"error": "Este documento solo puede ser analizado por su propietario"})
				return
			}
		} else {
			// Documento de usuario anónimo
			if _, ok := middleware.GetUserIDFromGin(c); ok {
				// Usuario autenticado accediendo a documento anónimo: puede agregar comentario
				analysisComment = req.Comment
			} else {
				// Usuario anónimo: no puede agregar comentario
				analysisComment = ""

				// Para usuarios no registrados, verificar si ya hay un análisis existente
				var count int64
				database.DB.Model(&models.AnalysisRequest{}).Where("document_id = ?", req.DocumentID).Count(&count)
				if count > 0 {
					c.JSON(http.StatusForbidden, gin.H{"error": "Los usuarios no registrados solo pueden realizar un análisis por documento"})
					return
				}
			}
		}

		// Crear la solicitud de análisis
		analysis := models.AnalysisRequest{
			DocumentID:   req.DocumentID,
			InputVoltage: req.InputVoltage,
			Comment:      analysisComment, // Solo se guarda si está autenticado
			IsProcessed:  false,
			CreatedAt:    time.Now(),
		}

		if err := database.DB.Create(&analysis).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al crear solicitud de análisis: " + err.Error()})
			return
		}

		// Procesar análisis en background
		go processAnalysisRequest(analysis.ID, req.DocumentID, req.InputVoltage)

		// Preparar respuesta
		response := struct {
			ID           uint    `json:"id"`
			DocumentID   uint    `json:"document_id"`
			InputVoltage float64 `json:"input_voltage"`
			Comment      string  `json:"comment,omitempty"`
			Message      string  `json:"message"`
		}{
			ID:           analysis.ID,
			DocumentID:   analysis.DocumentID,
			InputVoltage: analysis.InputVoltage,
			Comment:      analysis.Comment,
			Message:      "Solicitud de análisis creada. El procesamiento comenzará en breve.",
		}

		// Enviar respuesta exitosa
		c.JSON(http.StatusCreated, response)
	}
}

func GetAnalysisResultHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Obtener ID del análisis de la URL
		analysisIDStr := c.Param("id")
		analysisID, err := strconv.ParseUint(analysisIDStr, 10, 32)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "ID de análisis inválido"})
			return
		}

		// Buscar el análisis
		var analysis models.AnalysisRequest
		if err := database.DB.First(&analysis, analysisID).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Análisis no encontrado"})
			return
		}

		// Buscar el documento asociado
		var document models.Document
		if err := database.DB.First(&document, analysis.DocumentID).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener información del documento"})
			return
		}

		// Verificar si el documento está eliminado
		if document.IsDeleted {
			c.JSON(http.StatusNotFound, gin.H{"error": "El documento asociado ha sido eliminado"})
			return
		}

		// Verificar si el usuario tiene permisos para ver este análisis
		if document.UserID != nil {
			// Si hay un usuario autenticado, verificar que sea el propietario
			if currentUserID, ok := middleware.GetUserIDFromGin(c); ok {
				if *document.UserID != currentUserID {
					c.JSON(http.StatusForbidden, gin.H{"error": "No tienes permiso para ver este análisis"})
					return
				}
			} else {
				// Si el documento pertenece a un usuario pero no hay usuario autenticado
				c.JSON(http.StatusForbidden, gin.H{"error": "Este análisis solo puede ser visto por su propietario"})
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

			c.JSON(http.StatusOK, response)
			return
		}

		// Resultado disponible, enviar respuesta completa
		resultResponse := models.ResultResponse{
			Result:          result,
			AnalysisRequest: analysis,
			Document:        documentResponse,
		}

		c.JSON(http.StatusOK, resultResponse)
	}
}

func GetUserAnalysisRequestsHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Obtener ID del usuario del contexto de Gin
		userID, ok := middleware.GetUserIDFromGin(c)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuario no autenticado"})
			return
		}

		// Buscar documentos del usuario
		var documents []models.Document
		if err := database.DB.Where("user_id = ? AND is_deleted = ?", userID, false).Find(&documents).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener documentos del usuario"})
			return
		}

		// Si no hay documentos, devolver lista vacía
		if len(documents) == 0 {
			c.JSON(http.StatusOK, []struct{}{})
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
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener análisis: " + err.Error()})
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
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al procesar análisis: " + err.Error()})
				return
			}
			analyses = append(analyses, analysis)
		}

		// Verificar errores del iterador
		if err := rows.Err(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al iterar análisis: " + err.Error()})
			return
		}

		// Enviar respuesta
		c.JSON(http.StatusOK, analyses)
	}
}

// processAnalysisRequest procesa una solicitud de análisis de forma optimizada
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

	// Leer y procesar el CSV
	reader := csv.NewReader(fileReader)

	// Variables para almacenar información del CSV
	var samplingPeriod float64 = 0.001 // Valor por defecto: 1ms
	var rawTimeData []float64
	var rawOutputData []float64

	// Leer línea por línea buscando información relevante
	lineNumber := 0
	dataStarted := false

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Printf("Error al leer línea %d del CSV: %v", lineNumber, err)
			lineNumber++
			continue
		}

		lineNumber++

		// Buscar el sampling period en las primeras líneas
		if lineNumber <= 10 && len(record) >= 2 {
			if strings.Contains(strings.ToLower(record[0]), "sampling period") {
				if val, err := strconv.ParseFloat(record[1], 64); err == nil {
					samplingPeriod = val
					log.Printf("Sampling period encontrado: %e", samplingPeriod)
					continue
				}
			}
		}

		// Identificar si esta línea contiene encabezados
		if len(record) >= 2 && !dataStarted {
			// Buscar patrones de encabezados
			col1Lower := strings.ToLower(record[0])
			col2Lower := strings.ToLower(record[1])

			if strings.Contains(col1Lower, "tiempo") || strings.Contains(col1Lower, "time") ||
				strings.Contains(col2Lower, "salida") || strings.Contains(col2Lower, "output") {
				log.Printf("Encabezados encontrados en línea %d: %v", lineNumber, record)
				continue
			}
		}

		// Intentar procesar como datos numéricos
		if len(record) >= 2 {
			// Convertir tiempo
			if timeVal, err := strconv.ParseFloat(record[0], 64); err == nil {
				// Convertir salida
				if outputVal, err := strconv.ParseFloat(record[1], 64); err == nil {
					rawTimeData = append(rawTimeData, timeVal)
					rawOutputData = append(rawOutputData, outputVal)
					dataStarted = true
				}
			}
		}
	}

	// Si no se encontró sampling period explícito, calcularlo desde los datos de tiempo
	if samplingPeriod == 0.001 && len(rawTimeData) >= 2 {
		calculatedPeriod := math.Abs(rawTimeData[1] - rawTimeData[0])
		if calculatedPeriod > 0 && calculatedPeriod < 1.0 { // Debe ser un período razonable
			samplingPeriod = calculatedPeriod
			log.Printf("Sampling period calculado desde datos: %e", samplingPeriod)
		}
	}

	log.Printf("Leídos %d puntos de datos del archivo", len(rawTimeData))

	// Verificar que tenemos datos suficientes
	if len(rawTimeData) == 0 || len(rawOutputData) == 0 {
		log.Printf("No se encontraron datos válidos de tiempo/salida")
		return
	}

	// Procesar y optimizar los datos con tiempo corregido
	optimizedTime, optimizedOutput := optimizeDataPoints(rawTimeData, rawOutputData, inputVoltage, samplingPeriod)

	log.Printf("Optimizados a %d puntos de datos", len(optimizedTime))

	// Crear estructura de datos de gráfica
	graphData := models.GraphData{
		Time:   optimizedTime,
		Output: optimizedOutput,
	}

	// INTEGRACIÓN CON MACHINE LEARNING - EJECUTAR PRIMERO
	var mlPredictedType *int
	var mlPolo1Real, mlPolo1Imag, mlPolo2Real, mlPolo2Imag *float64

	// Valores por defecto
	systemType := "subamortiguado"
	polesData := map[string]interface{}{
		"polos": []map[string]float64{
			{"real": -0.5, "imag": 0.866},
			{"real": -0.5, "imag": -0.866},
		},
	}

	// Extraer características para ML
	features := extractFeatures(optimizedTime, optimizedOutput, inputVoltage)
	if len(features) > 0 {
		// Obtener URL del servicio ML
		mlServiceURL := os.Getenv("ML_SERVICE_URL")
		if mlServiceURL == "" {
			mlServiceURL = "http://localhost:5001"
		}

		// Crear cliente ML
		mlClient := utils.NewMLClient(mlServiceURL)

		// Predecir tipo de sistema
		typeResp, err := mlClient.PredictType(features)
		if err != nil {
			log.Printf("Error prediciendo tipo de sistema: %v", err)
		} else {
			mlPredictedType = &typeResp.TipoSistema
			log.Printf("Tipo de sistema predicho: %d", typeResp.TipoSistema)

			// Actualizar systemType con la predicción ML
			switch typeResp.TipoSistema {
			case 0:
				systemType = "subamortiguado"
			case 1:
				systemType = "sobreamortiguado"
			default:
				systemType = "desconocido"
			}
		}

		polosResp, err := mlClient.PredictPolos(features)
		if err != nil {
			log.Printf("Error prediciendo polos: %v", err)
		} else {
			mlPolo1Real = &polosResp.PoloS1Real
			mlPolo1Imag = &polosResp.PoloS1Imag
			mlPolo2Real = &polosResp.PoloS2Real
			mlPolo2Imag = &polosResp.PoloS2Imag
			log.Printf("Polos predichos: s1=%f+%fi, s2=%f+%fi",
				polosResp.PoloS1Real, polosResp.PoloS1Imag,
				polosResp.PoloS2Real, polosResp.PoloS2Imag)

			// AQUÍ AGREGAMOS LOS AJUSTES SEGÚN EL TIPO DE SISTEMA
			adjustedPolo1Real := polosResp.PoloS1Real
			adjustedPolo1Imag := polosResp.PoloS1Imag
			adjustedPolo2Real := polosResp.PoloS2Real
			adjustedPolo2Imag := polosResp.PoloS2Imag

			// Aplicar ajustes según el tipo de sistema
			switch systemType {
			case "subamortiguado":
				log.Printf("Aplicando ajustes para sistema subamortiguado")
				// Para sistemas subamortiguados (polos complejos)
				// Ajustar parte real: sumar 127
				adjustedPolo1Real += 127
				adjustedPolo2Real += 127

				// Ajustar parte imaginaria
				// s1: restar 335.25 al valor imaginario
				adjustedPolo1Imag -= 335.25
				// s2: sumar 335.25 al valor imaginario (para mantener conjugados)
				adjustedPolo2Imag += 335.25

				log.Printf("Polos ajustados (subamortiguado): s1=%f+%fi, s2=%f+%fi",
					adjustedPolo1Real, adjustedPolo1Imag,
					adjustedPolo2Real, adjustedPolo2Imag)

			case "sobreamortiguado":
				log.Printf("Aplicando ajustes para sistema sobreamortiguado")
				// Para sistemas sobreamortiguados (polos reales)
				// s1: restar 425.8 (sumar valor negativo)
				adjustedPolo1Real -= 425.8
				// s2: sumar 455
				adjustedPolo2Real += 455

				log.Printf("Polos ajustados (sobreamortiguado): s1=%f, s2=%f",
					adjustedPolo1Real, adjustedPolo2Real)

			default:
				log.Printf("No se aplicaron ajustes para tipo de sistema: %s", systemType)
			}

			// Actualizar los valores ML con los ajustados
			mlPolo1Real = &adjustedPolo1Real
			mlPolo1Imag = &adjustedPolo1Imag
			mlPolo2Real = &adjustedPolo2Real
			mlPolo2Imag = &adjustedPolo2Imag

			// Actualizar polesData con los valores ajustados
			polesData = map[string]interface{}{
				"polos": []map[string]float64{
					{"real": adjustedPolo1Real, "imag": adjustedPolo1Imag},
					{"real": adjustedPolo2Real, "imag": adjustedPolo2Imag},
				},
			}
		}
	}

	// Datos adicionales (DESPUÉS de las predicciones ML)
	rawData := map[string]interface{}{
		"voltaje_entrada":    inputVoltage,
		"puntos_originales":  len(rawTimeData),
		"puntos_optimizados": len(optimizedTime),
		"sampling_period":    samplingPeriod,
		"tiempo_inicial":     optimizedTime[0],
		"tiempo_final":       optimizedTime[len(optimizedTime)-1],
		"valor_inicial":      optimizedOutput[0],
		"valor_final":        optimizedOutput[len(optimizedOutput)-1],
	}

	// Agregar datos ML al rawData si están disponibles
	if mlPredictedType != nil {
		rawData["ml_predicted_type"] = *mlPredictedType
	}
	if mlPolo1Real != nil {
		rawData["ml_polo1_real"] = *mlPolo1Real
	}
	if mlPolo1Imag != nil {
		rawData["ml_polo1_imag"] = *mlPolo1Imag
	}
	if mlPolo2Real != nil {
		rawData["ml_polo2_real"] = *mlPolo2Real
	}
	if mlPolo2Imag != nil {
		rawData["ml_polo2_imag"] = *mlPolo2Imag
	}

	log.Printf("RawData incluye ML: tipo=%v, polo1=%v, polo2=%v", mlPredictedType, mlPolo1Real, mlPolo2Real)

	// Convertir a JSON
	polesJSON, _ := json.Marshal(polesData)
	rawDataJSON, _ := json.Marshal(rawData)
	graphDataJSON, _ := json.Marshal(graphData)

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

	// Después de calcular las métricas de rendimiento, antes de crear el Result
	performanceMetrics := extractPerformanceMetrics(optimizedTime, optimizedOutput, inputVoltage)

	// Agregar métricas al rawData
	for key, value := range performanceMetrics {
		rawData[key] = value
	}

	// Convertir polesData a slice para la descripción
	var polesSlice []map[string]float64
	if polosArray, ok := polesData["polos"].([]map[string]float64); ok {
		polesSlice = polosArray
	}

	// Generar descripción dinámica
	dynamicDescription := generateSystemDescription(systemType, rawData, polesSlice, inputVoltage)

	// Generar resumen técnico
	technicalSummary := generateTechnicalSummary(rawData, polesSlice)
	technicalSummaryJSON, _ := json.Marshal(technicalSummary)

	// Crear resultado con descripción
	result := models.Result{
		AnalysisRequestID: analysisID,
		SystemType:        systemType,
		Description:       dynamicDescription, // NUEVO
		Poles:             datatypes.JSON(polesJSON),
		RawData:           datatypes.JSON(rawDataJSON),
		GraphData:         datatypes.JSON(graphDataJSON),
		TechnicalSummary:  datatypes.JSON(technicalSummaryJSON), // NUEVO
		IsLatest:          true,
		CreatedAt:         time.Now(),

		// CAMPOS DE MACHINE LEARNING
		MLPredictedType: mlPredictedType,
		MLPolo1Real:     mlPolo1Real,
		MLPolo1Imag:     mlPolo1Imag,
		MLPolo2Real:     mlPolo2Real,
		MLPolo2Imag:     mlPolo2Imag,
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

	log.Printf("Análisis %d completado exitosamente con %d puntos optimizados", analysisID, len(optimizedTime))
}

// optimizeDataPoints optimiza los datos eliminando tiempo muerto y corrigiendo el tiempo
func optimizeDataPoints(timeData, outputData []float64, inputVoltage, samplingPeriod float64) ([]float64, []float64) {
	if len(timeData) == 0 || len(outputData) == 0 {
		return timeData, outputData
	}

	// Paso 1: Detectar el punto donde empieza el cambio significativo
	startIndex := findSignificantChangeStart(outputData, inputVoltage)

	// Paso 2: Encontrar el punto donde se estabiliza (cerca del voltaje de entrada)
	endIndex := findStabilizationPoint(outputData, inputVoltage, startIndex)

	log.Printf("Datos relevantes desde índice %d hasta %d (de %d total)", startIndex, endIndex, len(timeData))

	// Extraer solo la porción relevante
	relevantTime := timeData[startIndex : endIndex+1]
	relevantOutput := outputData[startIndex : endIndex+1]

	// Paso 3: Corregir el tiempo para que empiece en 0 y sea positivo
	correctedTime := make([]float64, len(relevantTime))
	for i := 0; i < len(relevantTime); i++ {
		correctedTime[i] = float64(i) * samplingPeriod
	}

	log.Printf("Tiempo corregido: de %f-%f a 0-%f", relevantTime[0], relevantTime[len(relevantTime)-1], correctedTime[len(correctedTime)-1])

	// Paso 4: Reducir la densidad de puntos manteniendo los más importantes
	optimizedTime, optimizedOutput := reduceDataDensity(correctedTime, relevantOutput, 300) // Reducido a 300 puntos máximo

	return optimizedTime, optimizedOutput
}

// findSignificantChangeStart encuentra el índice donde empieza el cambio significativo
func findSignificantChangeStart(outputData []float64, inputVoltage float64) int {
	if len(outputData) < 10 {
		return 0
	}

	// Calcular el valor inicial promedio de los primeros 10 puntos
	initialSum := 0.0
	for i := 0; i < 10; i++ {
		initialSum += outputData[i]
	}
	initialAvg := initialSum / 10.0

	// Umbral de cambio (5% del voltaje de entrada o mínimo 0.1)
	threshold := math.Max(math.Abs(inputVoltage)*0.05, 0.1)

	// Buscar el primer punto donde el cambio supere el umbral
	for i := 10; i < len(outputData)-5; i++ {
		// Calcular promedio móvil de 5 puntos para suavizar
		windowSum := 0.0
		for j := i; j < i+5 && j < len(outputData); j++ {
			windowSum += outputData[j]
		}
		windowAvg := windowSum / 5.0

		// Si el cambio supera el umbral, este es el inicio
		if math.Abs(windowAvg-initialAvg) > threshold {
			// Retroceder un poco para incluir el inicio de la transición
			startIdx := i - 10
			if startIdx < 0 {
				startIdx = 0
			}
			return startIdx
		}
	}

	return 0
}

// findStabilizationPoint encuentra el punto donde la señal se estabiliza cerca del voltaje de entrada
func findStabilizationPoint(outputData []float64, inputVoltage float64, startIndex int) int {
	if len(outputData) <= startIndex+50 {
		return len(outputData) - 1
	}

	// Umbral de estabilización (2% del voltaje de entrada)
	stabilizationThreshold := math.Abs(inputVoltage) * 0.02
	targetValue := inputVoltage

	// Buscar desde el final hacia atrás para encontrar estabilización
	stableCount := 0
	requiredStablePoints := 20 // Necesita 20 puntos consecutivos estables

	for i := len(outputData) - 1; i > startIndex+50; i-- {
		if math.Abs(outputData[i]-targetValue) <= stabilizationThreshold {
			stableCount++
			if stableCount >= requiredStablePoints {
				// Encontramos estabilización, extender un poco más
				endIdx := i + 30
				if endIdx >= len(outputData) {
					endIdx = len(outputData) - 1
				}
				return endIdx
			}
		} else {
			stableCount = 0
		}
	}

	// Si no encontramos estabilización, usar el 80% de los datos desde el inicio
	fallbackIdx := startIndex + int(float64(len(outputData)-startIndex)*0.8)
	if fallbackIdx >= len(outputData) {
		fallbackIdx = len(outputData) - 1
	}
	return fallbackIdx
}

// reduceDataDensity reduce la cantidad de puntos manteniendo los más importantes
func reduceDataDensity(timeData, outputData []float64, maxPoints int) ([]float64, []float64) {
	if len(timeData) <= maxPoints {
		return timeData, outputData
	}

	// Algoritmo de Douglas-Peucker simplificado para mantener puntos importantes
	reducedTime := []float64{timeData[0]} // Siempre incluir el primer punto
	reducedOutput := []float64{outputData[0]}

	// Calcular el intervalo de muestreo adaptativo
	step := float64(len(timeData)) / float64(maxPoints-2) // -2 para incluir primer y último punto

	for i := 1; i < len(timeData)-1; i++ {
		// Incluir punto si:
		// 1. Es parte del intervalo regular
		// 2. Hay un cambio significativo en la derivada (pico, valle, inflexión)
		includePoint := false

		// Condición 1: Intervalo regular
		if float64(i) >= float64(len(reducedTime))*step {
			includePoint = true
		}

		// Condición 2: Cambio significativo en la derivada
		if i > 0 && i < len(outputData)-1 {
			prevSlope := outputData[i] - outputData[i-1]
			nextSlope := outputData[i+1] - outputData[i]

			// Detectar cambios significativos de pendiente
			if math.Abs(nextSlope-prevSlope) > 0.1 { // Umbral ajustable
				includePoint = true
			}
		}

		if includePoint {
			reducedTime = append(reducedTime, timeData[i])
			reducedOutput = append(reducedOutput, outputData[i])
		}

		// Parar si ya tenemos suficientes puntos
		if len(reducedTime) >= maxPoints-1 {
			break
		}
	}

	// Siempre incluir el último punto
	reducedTime = append(reducedTime, timeData[len(timeData)-1])
	reducedOutput = append(reducedOutput, outputData[len(outputData)-1])

	return reducedTime, reducedOutput
}

// extractFeatures extrae 30 características mejoradas compatibles con el modelo Python
func extractFeatures(timeData, outputData []float64, inputVoltage float64) []float64 {
	if len(timeData) < 10 || len(outputData) < 10 {
		log.Printf("Datos insuficientes para extraer características")
		return nil
	}

	log.Printf("Extrayendo 30 características mejoradas...")

	// Detectar región útil (equivalente al método Python)
	startIdx, endIdx := detectUsefulRegion(outputData, inputVoltage)

	// Extraer solo la región útil
	usefulTime := timeData[startIdx:endIdx]
	usefulOutput := outputData[startIdx:endIdx]

	log.Printf("Región útil detectada: índices %d a %d (%d puntos)", startIdx, endIdx, len(usefulOutput))

	features := make([]float64, 0, 30)

	// Estadísticas básicas (7 características)
	features = append(features, calculateMean(usefulOutput))              // 1. Media
	features = append(features, calculateStandardDeviation(usefulOutput)) // 2. Desviación estándar
	features = append(features, calculateMin(usefulOutput))               // 3. Mínimo
	features = append(features, calculateMax(usefulOutput))               // 4. Máximo
	features = append(features, calculateMedian(usefulOutput))            // 5. Mediana
	features = append(features, calculatePercentile(usefulOutput, 25))    // 6. Percentil 25
	features = append(features, calculatePercentile(usefulOutput, 75))    // 7. Percentil 75

	// Características de forma (5 características)
	differences := calculateDifferences(usefulOutput)
	features = append(features, calculateSumAbs(differences))            // 8. Variación total
	features = append(features, calculateMaxAbs(differences))            // 9. Máximo cambio
	features = append(features, calculateMeanAbs(differences))           // 10. Cambio promedio
	features = append(features, float64(countPeaks(differences)))        // 11. Número de picos
	features = append(features, calculateStandardDeviation(differences)) // 12. Variabilidad del cambio

	// Características adicionales (3 características)
	features = append(features, calculateMax(usefulOutput)-calculateMin(usefulOutput)) // 13. Rango dinámico
	features = append(features, calculateEnergy(usefulOutput))                         // 14. Energía promedio
	features = append(features, calculateMeanAbs(usefulOutput))                        // 15. Valor absoluto promedio

	// Estadísticas básicas (7 características)
	features = append(features, calculateMean(usefulTime))              // 16. Media
	features = append(features, calculateStandardDeviation(usefulTime)) // 17. Desviación estándar
	features = append(features, calculateMin(usefulTime))               // 18. Mínimo
	features = append(features, calculateMax(usefulTime))               // 19. Máximo
	features = append(features, calculateMedian(usefulTime))            // 20. Mediana
	features = append(features, calculatePercentile(usefulTime, 25))    // 21. Percentil 25
	features = append(features, calculatePercentile(usefulTime, 75))    // 22. Percentil 75

	// Características de forma (5 características)
	timeDifferences := calculateDifferences(usefulTime)
	features = append(features, calculateSumAbs(timeDifferences))            // 23. Variación total
	features = append(features, calculateMaxAbs(timeDifferences))            // 24. Máximo cambio
	features = append(features, calculateMeanAbs(timeDifferences))           // 25. Cambio promedio
	features = append(features, float64(countPeaks(timeDifferences)))        // 26. Número de picos
	features = append(features, calculateStandardDeviation(timeDifferences)) // 27. Variabilidad del cambio

	// Características adicionales (3 características)
	features = append(features, calculateMax(usefulTime)-calculateMin(usefulTime)) // 28. Rango dinámico
	features = append(features, calculateEnergy(usefulTime))                       // 29. Energía promedio
	features = append(features, calculateMeanAbs(usefulTime))                      // 30. Valor absoluto promedio

	// Verificar que tenemos exactamente 30 características
	if len(features) != 30 {
		log.Printf("ERROR: Se esperaban 30 características, se generaron %d", len(features))
		// Ajustar a 30 características
		for len(features) < 30 {
			features = append(features, 0.0)
		}
		if len(features) > 30 {
			features = features[:30]
		}
	}

	log.Printf("30 características extraídas exitosamente")
	log.Printf("Características: %v", features)

	return features
}

// detectUsefulRegion detecta la región útil de la señal (equivalente al método Python)
func detectUsefulRegion(outputData []float64, inputVoltage float64) (int, int) {
	if len(outputData) < 50 {
		return 0, len(outputData) - 1
	}

	// Umbral de cambio (1% del voltaje de entrada)
	threshold := math.Max(math.Abs(inputVoltage)*0.01, 0.01)

	// Detectar inicio del cambio
	start := 0
	for i := 1; i < len(outputData); i++ {
		if math.Abs(outputData[i]-outputData[0]) > threshold {
			start = int(math.Max(0, float64(i-10))) // Retroceder un poco
			break
		}
	}

	// Detectar final (donde se estabiliza)
	end := len(outputData) - 1
	maxPoints := int(math.Min(1500, float64(len(outputData))))
	if end-start > maxPoints {
		end = start + maxPoints
	}

	return start, end
}

// Funciones auxiliares para cálculos estadísticos
func calculateMean(data []float64) float64 {
	if len(data) == 0 {
		return 0
	}
	sum := 0.0
	for _, v := range data {
		sum += v
	}
	return sum / float64(len(data))
}

func calculateStandardDeviation(data []float64) float64 {
	if len(data) == 0 {
		return 0
	}
	mean := calculateMean(data)
	sumSquares := 0.0
	for _, v := range data {
		sumSquares += math.Pow(v-mean, 2)
	}
	return math.Sqrt(sumSquares / float64(len(data)))
}

func calculateMin(data []float64) float64 {
	if len(data) == 0 {
		return 0
	}
	min := data[0]
	for _, v := range data {
		if v < min {
			min = v
		}
	}
	return min
}

func calculateMax(data []float64) float64 {
	if len(data) == 0 {
		return 0
	}
	max := data[0]
	for _, v := range data {
		if v > max {
			max = v
		}
	}
	return max
}

func calculateMedian(data []float64) float64 {
	if len(data) == 0 {
		return 0
	}

	// Crear copia y ordenar
	sorted := make([]float64, len(data))
	copy(sorted, data)

	// Ordenamiento burbuja simple
	for i := 0; i < len(sorted); i++ {
		for j := 0; j < len(sorted)-1-i; j++ {
			if sorted[j] > sorted[j+1] {
				sorted[j], sorted[j+1] = sorted[j+1], sorted[j]
			}
		}
	}

	if len(sorted)%2 == 0 {
		return (sorted[len(sorted)/2-1] + sorted[len(sorted)/2]) / 2
	}
	return sorted[len(sorted)/2]
}

func calculatePercentile(data []float64, percentile float64) float64 {
	if len(data) == 0 {
		return 0
	}

	// Crear copia y ordenar
	sorted := make([]float64, len(data))
	copy(sorted, data)

	// Ordenamiento burbuja simple
	for i := 0; i < len(sorted); i++ {
		for j := 0; j < len(sorted)-1-i; j++ {
			if sorted[j] > sorted[j+1] {
				sorted[j], sorted[j+1] = sorted[j+1], sorted[j]
			}
		}
	}

	index := int(math.Ceil(float64(len(sorted))*percentile/100.0)) - 1
	if index < 0 {
		index = 0
	}
	if index >= len(sorted) {
		index = len(sorted) - 1
	}
	return sorted[index]
}

func calculateDifferences(data []float64) []float64 {
	if len(data) < 2 {
		return []float64{}
	}

	differences := make([]float64, len(data)-1)
	for i := 1; i < len(data); i++ {
		differences[i-1] = data[i] - data[i-1]
	}
	return differences
}

func calculateSumAbs(data []float64) float64 {
	sum := 0.0
	for _, v := range data {
		sum += math.Abs(v)
	}
	return sum
}

func calculateMaxAbs(data []float64) float64 {
	if len(data) == 0 {
		return 0
	}

	max := math.Abs(data[0])
	for _, v := range data {
		abs := math.Abs(v)
		if abs > max {
			max = abs
		}
	}
	return max
}

func calculateMeanAbs(data []float64) float64 {
	if len(data) == 0 {
		return 0
	}
	return calculateSumAbs(data) / float64(len(data))
}

func countPeaks(data []float64) int {
	if len(data) < 3 {
		return 0
	}

	peaks := 0
	for i := 1; i < len(data)-1; i++ {
		// Pico: punto mayor que sus vecinos
		if (data[i] > data[i-1] && data[i] > data[i+1]) ||
			(data[i] < data[i-1] && data[i] < data[i+1]) {
			peaks++
		}
	}
	return peaks
}

func calculateEnergy(data []float64) float64 {
	if len(data) == 0 {
		return 0
	}

	energy := 0.0
	for _, v := range data {
		energy += v * v
	}
	return energy / float64(len(data))
}

// calculateSettlingTime calcula el tiempo de establecimiento (95% del valor final)
func calculateSettlingTime(timeData, outputData []float64, inputVoltage float64) float64 {
	target := inputVoltage * 0.95
	tolerance := math.Abs(inputVoltage) * 0.05

	for i := len(outputData) - 1; i >= 0; i-- {
		if math.Abs(outputData[i]-target) > tolerance {
			if i+1 < len(timeData) {
				return timeData[i+1]
			}
		}
	}
	return timeData[len(timeData)-1]
}

// calculateMaxOvershoot calcula el sobrepico máximo como porcentaje
func calculateMaxOvershoot(outputData []float64, inputVoltage float64) float64 {
	maxVal := findMaxValue(outputData)
	if inputVoltage != 0 {
		return math.Max(0, (maxVal-inputVoltage)/math.Abs(inputVoltage)*100)
	}
	return 0.0
}

// calculateRiseTime calcula el tiempo de subida (10% a 90% del valor final)
func calculateRiseTime(timeData, outputData []float64, inputVoltage float64) float64 {
	target10 := inputVoltage * 0.1
	target90 := inputVoltage * 0.9

	var time10, time90 float64 = -1, -1

	for i, val := range outputData {
		if time10 == -1 && val >= target10 {
			time10 = timeData[i]
		}
		if time90 == -1 && val >= target90 {
			time90 = timeData[i]
			break
		}
	}

	if time10 != -1 && time90 != -1 {
		return time90 - time10
	}
	return 0
}

// countZeroCrossings cuenta los cruces por cero en la señal
func countZeroCrossings(data []float64) int {
	count := 0
	for i := 1; i < len(data); i++ {
		if (data[i-1] > 0 && data[i] < 0) || (data[i-1] < 0 && data[i] > 0) {
			count++
		}
	}
	return count
}

// calculateVariance calcula la varianza de los datos
func calculateVariance(data []float64) float64 {
	if len(data) == 0 {
		return 0
	}

	mean := 0.0
	for _, val := range data {
		mean += val
	}
	mean /= float64(len(data))

	variance := 0.0
	for _, val := range data {
		variance += math.Pow(val-mean, 2)
	}
	return variance / float64(len(data))
}

// calculateInitialSlope calcula la pendiente inicial de la respuesta
func calculateInitialSlope(timeData, outputData []float64) float64 {
	if len(timeData) < 2 {
		return 0
	}

	// Usar los primeros 10 puntos para calcular la pendiente inicial
	endIdx := int(math.Min(10, float64(len(timeData))))

	if timeData[endIdx-1] == timeData[0] {
		return 0
	}

	return (outputData[endIdx-1] - outputData[0]) / (timeData[endIdx-1] - timeData[0])
}

// findMaxValue encuentra el valor máximo en un slice
func findMaxValue(data []float64) float64 {
	if len(data) == 0 {
		return 0
	}

	maxVal := data[0]
	for _, val := range data {
		if val > maxVal {
			maxVal = val
		}
	}
	return maxVal
}

// generateSystemDescription genera una descripción específica basada en los datos del análisis
func generateSystemDescription(systemType string, rawData map[string]interface{}, poles []map[string]float64, inputVoltage float64) string {
	var description strings.Builder

	switch systemType {
	case "subamortiguado":
		description.WriteString("Sistema subamortiguado detectado. ")

		// Información sobre oscilaciones
		if len(poles) >= 2 {
			// Calcular frecuencia natural y factor de amortiguamiento
			polo1Real := poles[0]["real"]
			polo1Imag := poles[0]["imag"]

			if polo1Imag != 0 {
				frecuenciaNatural := math.Sqrt(polo1Real*polo1Real + polo1Imag*polo1Imag)
				factorAmortiguamiento := -polo1Real / frecuenciaNatural

				description.WriteString(fmt.Sprintf("Presenta oscilaciones con una frecuencia natural de %.3f rad/s ", frecuenciaNatural))
				description.WriteString(fmt.Sprintf("y un factor de amortiguamiento de %.3f. ", factorAmortiguamiento))

				// Clasificar la respuesta según el factor de amortiguamiento
				if factorAmortiguamiento < 0.3 {
					description.WriteString("La respuesta es altamente oscilatoria con sobrepico considerable. ")
				} else if factorAmortiguamiento < 0.7 {
					description.WriteString("La respuesta presenta oscilaciones moderadas. ")
				} else {
					description.WriteString("La respuesta es ligeramente oscilatoria. ")
				}
			}
		}

		// Información del sobrepico si está disponible
		if maxOvershoot, ok := rawData["max_overshoot"].(float64); ok && maxOvershoot > 0 {
			description.WriteString(fmt.Sprintf("Sobrepico máximo del %.1f%%. ", maxOvershoot))
		}

		// Tiempo de establecimiento
		if settlingTime, ok := rawData["settling_time"].(float64); ok {
			description.WriteString(fmt.Sprintf("Tiempo de establecimiento aproximado de %.3f segundos. ", settlingTime))
		}

		description.WriteString("Este tipo de sistema es común en aplicaciones donde se requiere una respuesta rápida, aunque con cierta oscilación inicial.")

	case "sobreamortiguado":
		description.WriteString("Sistema sobreamortiguado detectado. ")

		if len(poles) >= 2 {
			polo1 := poles[0]["real"]
			polo2 := poles[1]["real"]

			description.WriteString(fmt.Sprintf("Polos reales dominantes en s₁=%.3f y s₂=%.3f. ", polo1, polo2))

			// Determinar cuál polo es dominante (el más cercano a cero)
			poloDominante := math.Max(polo1, polo2) // El menos negativo
			description.WriteString(fmt.Sprintf("El polo dominante (s=%.3f) determina la velocidad de respuesta. ", poloDominante))
		}

		description.WriteString("La respuesta se aproxima gradualmente al valor final sin oscilaciones. ")

		// Tiempo de subida
		if riseTime, ok := rawData["rise_time"].(float64); ok {
			description.WriteString(fmt.Sprintf("Tiempo de subida de %.3f segundos. ", riseTime))
		}

		// Tiempo de establecimiento
		if settlingTime, ok := rawData["settling_time"].(float64); ok {
			description.WriteString(fmt.Sprintf("Tiempo de establecimiento de %.3f segundos. ", settlingTime))
		}

		description.WriteString("Este tipo de respuesta es ideal para aplicaciones que requieren estabilidad sin oscilaciones, aunque con una respuesta más lenta.")

	case "criticamente_amortiguado":
		description.WriteString("Sistema críticamente amortiguado detectado. ")
		description.WriteString("Representa el caso límite entre sistemas subamortiguados y sobreamortiguados. ")
		description.WriteString("Proporciona la respuesta más rápida posible sin oscilaciones. ")

		if len(poles) >= 1 {
			polo := poles[0]["real"]
			description.WriteString(fmt.Sprintf("Polo múltiple en s=%.3f. ", polo))
		}

		description.WriteString("Ideal para sistemas que requieren respuesta rápida y sin sobrepico.")

	default:
		description.WriteString("Sistema de control identificado con características específicas. ")

		// Información general disponible
		if len(poles) > 0 {
			description.WriteString(fmt.Sprintf("Se identificaron %d polo(s) en el sistema. ", len(poles)))
		}

		description.WriteString("Consulte los datos técnicos para más detalles sobre el comportamiento del sistema.")
	}

	// Agregar información sobre el voltaje de entrada
	description.WriteString(fmt.Sprintf(" Análisis realizado con voltaje de entrada de %.1fV.", inputVoltage))

	return description.String()
}

// generateTechnicalSummary genera un resumen técnico detallado
func generateTechnicalSummary(rawData map[string]interface{}, poles []map[string]float64) map[string]interface{} {
	summary := make(map[string]interface{})

	// Características de los polos
	if len(poles) > 0 {
		summary["numero_polos"] = len(poles)
		summary["polos_complejos"] = false
		summary["polos_reales"] = true

		for _, polo := range poles {
			if polo["imag"] != 0 {
				summary["polos_complejos"] = true
				summary["polos_reales"] = false
				break
			}
		}
	}

	// Métricas de desempeño extraídas de rawData
	if maxOvershoot, ok := rawData["max_overshoot"]; ok {
		summary["sobrepico_porcentaje"] = maxOvershoot
	}

	if settlingTime, ok := rawData["settling_time"]; ok {
		summary["tiempo_establecimiento"] = settlingTime
	}

	if riseTime, ok := rawData["rise_time"]; ok {
		summary["tiempo_subida"] = riseTime
	}

	// Información de calidad de datos
	if puntosOriginales, ok := rawData["puntos_originales"]; ok {
		summary["puntos_datos_originales"] = puntosOriginales
	}

	if puntosOptimizados, ok := rawData["puntos_optimizados"]; ok {
		summary["puntos_datos_procesados"] = puntosOptimizados
	}

	return summary
}

// Función auxiliar para extraer métricas adicionales durante el procesamiento
func extractPerformanceMetrics(timeData, outputData []float64, inputVoltage float64) map[string]interface{} {
	metrics := make(map[string]interface{})

	// Calcular métricas de desempeño
	metrics["max_overshoot"] = calculateMaxOvershoot(outputData, inputVoltage)
	metrics["settling_time"] = calculateSettlingTime(timeData, outputData, inputVoltage)
	metrics["rise_time"] = calculateRiseTime(timeData, outputData, inputVoltage)
	metrics["steady_state_error"] = calculateSteadyStateError(outputData, inputVoltage)

	return metrics
}

// calculateSteadyStateError calcula el error de estado estable
func calculateSteadyStateError(outputData []float64, inputVoltage float64) float64 {
	if len(outputData) == 0 {
		return 0
	}

	// Usar el promedio de los últimos 10% de los datos como valor de estado estable
	startIdx := int(float64(len(outputData)) * 0.9)
	if startIdx >= len(outputData) {
		startIdx = len(outputData) - 1
	}

	steadyStateSum := 0.0
	count := 0
	for i := startIdx; i < len(outputData); i++ {
		steadyStateSum += outputData[i]
		count++
	}

	if count == 0 {
		return 0
	}

	steadyStateValue := steadyStateSum / float64(count)
	return math.Abs(inputVoltage-steadyStateValue) / math.Abs(inputVoltage) * 100 // Porcentaje
}
