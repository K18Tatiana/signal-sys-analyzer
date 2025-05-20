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
	IsProcessed  bool      `gorm:"column:is_processed;default:false" json:"is_processed"`
	CreatedAt    time.Time `gorm:"column:created_at;type:timestamp with time zone;default:CURRENT_TIMESTAMP" json:"created_at"`
	Results      []Result  `gorm:"foreignKey:AnalysisRequestID" json:"-"`
}

// AnalysisRequestCreate para solicitar un nuevo análisis
type AnalysisRequestCreate struct {
	DocumentID   uint    `json:"document_id"`
	InputVoltage float64 `json:"input_voltage"`
}

// Result representa el resultado del análisis ML de un documento
type Result struct {
	ID                uint            `gorm:"primaryKey;type:serial" json:"id"`
	AnalysisRequestID uint            `gorm:"column:analysis_request_id;not null;index" json:"analysis_request_id"`
	AnalysisRequest   AnalysisRequest `gorm:"foreignKey:AnalysisRequestID" json:"-"`
	SystemType        string          `gorm:"column:system_type;size:50;not null" json:"system_type"`
	Poles             datatypes.JSON  `gorm:"column:poles;type:jsonb;not null" json:"poles"` // JSONB en PostgreSQL
	RawData           datatypes.JSON  `gorm:"column:raw_data;type:jsonb" json:"raw_data,omitempty"`
	IsLatest          bool            `gorm:"column:is_latest;default:true" json:"is_latest"`
	CreatedAt         time.Time       `gorm:"column:created_at;type:timestamp with time zone;default:CURRENT_TIMESTAMP" json:"created_at"`
}

// ResultResponse es la respuesta completa de un análisis
type ResultResponse struct {
	Result          Result           `json:"result"`
	AnalysisRequest AnalysisRequest  `json:"analysis_request"`
	Document        DocumentResponse `json:"document"`
}
