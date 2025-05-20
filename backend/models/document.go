package models

import (
	"time"
)

// Document representa un archivo CSV subido al sistema
type Document struct {
	ID               uint              `gorm:"primaryKey;type:serial" json:"id"`
	UserID           *uint             `gorm:"column:user_id;type:integer" json:"user_id"` // Puede ser nulo para usuarios no registrados
	User             *User             `gorm:"foreignKey:UserID" json:"-"`
	FilePath         string            `gorm:"column:file_path;size:255;not null" json:"-"` // No enviar al cliente
	OriginalFilename string            `gorm:"column:original_filename;size:255;not null" json:"original_filename"`
	UploadDate       time.Time         `gorm:"column:upload_date;type:timestamp with time zone;default:CURRENT_TIMESTAMP" json:"upload_date"`
	IsDeleted        bool              `gorm:"column:is_deleted;default:false" json:"is_deleted"` // Indica si el documento está "eliminado" para la vista
	AnalysisRequests []AnalysisRequest `gorm:"foreignKey:DocumentID" json:"-"`
}

// DocumentUploadRequest para subir un nuevo documento
type DocumentUploadRequest struct {
	// El archivo se procesa desde FormFile
	UserID *uint `json:"user_id,omitempty"` // Opcional
}

// DocumentResponse es la respuesta enviada al cliente
type DocumentResponse struct {
	ID               uint      `json:"id"`
	OriginalFilename string    `json:"original_filename"`
	UploadDate       time.Time `json:"upload_date"`
	AnalysisCount    int       `json:"analysis_count"` // Número de análisis realizados
	IsDeleted        bool      `json:"is_deleted"`
}

// ToDocumentResponse convierte un Document a DocumentResponse
func (d *Document) ToDocumentResponse(analysisCount int) DocumentResponse {
	return DocumentResponse{
		ID:               d.ID,
		OriginalFilename: d.OriginalFilename,
		UploadDate:       d.UploadDate,
		AnalysisCount:    analysisCount,
		IsDeleted:        d.IsDeleted,
	}
}
