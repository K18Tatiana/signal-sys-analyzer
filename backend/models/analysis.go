package models

import (
	"time"

	"gorm.io/datatypes"
)

// AnalysisRequest representa una solicitud de análisis de un documento
type AnalysisRequest struct {
	ID           uint      `gorm:"primaryKey;type:serial" json:"id"`
	DocumentID   uint      `gorm:"column:document_id;not null;index" json:"document_id"`
	Document     Document  `gorm:"foreignKey:DocumentID" json:"-"`
	InputVoltage float64   `gorm:"column:input_voltage;not null" json:"input_voltage"`
	Comment      string    `gorm:"column:comment;size:500" json:"comment,omitempty"`
	IsProcessed  bool      `gorm:"column:is_processed;default:false" json:"is_processed"`
	CreatedAt    time.Time `gorm:"column:created_at;type:timestamp with time zone;default:CURRENT_TIMESTAMP" json:"created_at"`
	Results      []Result  `gorm:"foreignKey:AnalysisRequestID" json:"-"`
}

// AnalysisRequestCreate para solicitar un nuevo análisis
type AnalysisRequestCreate struct {
	DocumentID   uint    `json:"document_id"`
	InputVoltage float64 `json:"input_voltage"`
	Comment      string  `json:"comment,omitempty"`
}

// Result representa el resultado del análisis ML de un documento
type Result struct {
	ID                uint            `gorm:"primaryKey;type:serial" json:"id"`
	AnalysisRequestID uint            `gorm:"column:analysis_request_id;not null;index" json:"analysis_request_id"`
	AnalysisRequest   AnalysisRequest `gorm:"foreignKey:AnalysisRequestID" json:"-"`
	SystemType        string          `gorm:"column:system_type;size:50;not null" json:"system_type"`
	Description       string          `gorm:"column:description;size:1000" json:"description"`
	Poles             datatypes.JSON  `gorm:"column:poles;type:jsonb;not null" json:"poles"`
	RawData           datatypes.JSON  `gorm:"column:raw_data;type:jsonb" json:"raw_data,omitempty"`
	GraphData         datatypes.JSON  `gorm:"column:graph_data;type:jsonb" json:"graph_data,omitempty"`
	TechnicalSummary  datatypes.JSON  `gorm:"column:technical_summary;type:jsonb" json:"technical_summary,omitempty"`
	IsLatest          bool            `gorm:"column:is_latest;default:true" json:"is_latest"`
	CreatedAt         time.Time       `gorm:"column:created_at;type:timestamp with time zone;default:CURRENT_TIMESTAMP" json:"created_at"`

	MLPredictedType *int     `gorm:"column:ml_predicted_type" json:"ml_predicted_type,omitempty"`
	MLPolo1Real     *float64 `gorm:"column:ml_polo1_real" json:"ml_polo1_real,omitempty"`
	MLPolo1Imag     *float64 `gorm:"column:ml_polo1_imag" json:"ml_polo1_imag,omitempty"`
	MLPolo2Real     *float64 `gorm:"column:ml_polo2_real" json:"ml_polo2_real,omitempty"`
	MLPolo2Imag     *float64 `gorm:"column:ml_polo2_imag" json:"ml_polo2_imag,omitempty"`
	MLConfidence    *float64 `gorm:"column:ml_confidence" json:"ml_confidence,omitempty"`
}

// GraphData estructura para almacenar datos de tiempo y salida para gráficas
type GraphData struct {
	Time   []float64 `json:"time"`
	Output []float64 `json:"output"`
}

// ResultResponse es la respuesta completa de un análisis
type ResultResponse struct {
	Result          Result           `json:"result"`
	AnalysisRequest AnalysisRequest  `json:"analysis_request"`
	Document        DocumentResponse `json:"document"`
}
