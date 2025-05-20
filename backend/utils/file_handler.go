package utils

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
)

const (
	// Configuración para el manejo de archivos
	MaxFileSize = 10 * 1024 * 1024 // 10MB
	UploadDir   = "./uploads"
)

// FileHandler proporciona métodos para manejar archivos de forma segura
type FileHandler struct {
	UploadDir string
	MaxSize   int64
}

// NewFileHandler crea una nueva instancia de FileHandler
func NewFileHandler() *FileHandler {
	// Obtener directorio de uploads desde variable de entorno o usar valor por defecto
	uploadDir := os.Getenv("UPLOAD_DIR")
	if uploadDir == "" {
		uploadDir = UploadDir
	}

	// Crear el directorio de uploads si no existe
	os.MkdirAll(uploadDir, 0755)

	// Obtener tamaño máximo desde variable de entorno o usar valor por defecto
	var maxSize int64 = MaxFileSize
	if sizeStr := os.Getenv("MAX_FILE_SIZE"); sizeStr != "" {
		if size, err := parseInt64(sizeStr); err == nil {
			maxSize = size
		}
	}

	return &FileHandler{
		UploadDir: uploadDir,
		MaxSize:   maxSize,
	}
}

// SaveFile guarda un archivo subido de forma segura
func (h *FileHandler) SaveFile(file multipart.File, header *multipart.FileHeader) (string, error) {
	// Verificar el tamaño del archivo
	if header.Size > h.MaxSize {
		return "", errors.New("el archivo es demasiado grande (máximo 10MB)")
	}

	// Verificar la extensión del archivo
	ext := strings.ToLower(filepath.Ext(header.Filename))
	if ext != ".csv" {
		return "", errors.New("solo se permiten archivos CSV")
	}

	// Generar un nombre de archivo único
	filename, err := generateSecureFilename(ext)
	if err != nil {
		return "", errors.New("error al generar nombre de archivo")
	}

	// Ruta completa del archivo
	filePath := filepath.Join(h.UploadDir, filename)

	// Crear el archivo en el sistema de archivos
	dst, err := os.Create(filePath)
	if err != nil {
		return "", errors.New("error al crear el archivo")
	}
	defer dst.Close()

	// Copiar el contenido del archivo subido
	if _, err = io.Copy(dst, file); err != nil {
		os.Remove(filePath) // Limpiar en caso de error
		return "", errors.New("error al guardar el archivo")
	}

	return filePath, nil
}

// ValidateCSV realiza validaciones básicas del archivo CSV
func (h *FileHandler) ValidateCSV(filePath string) error {
	// Abrir el archivo
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	// Leer las primeras líneas para validación básica
	// Esta es una validación mínima, puedes ampliarla según tus necesidades
	buffer := make([]byte, 512)
	_, err = file.Read(buffer)
	if err != nil && err != io.EOF {
		return err
	}

	// Verificar si parece un CSV (contiene comas o punto y coma)
	content := string(buffer)
	if !strings.Contains(content, ",") && !strings.Contains(content, ";") {
		return errors.New("el archivo no parece ser un CSV válido")
	}

	return nil
}

// generateSecureFilename genera un nombre de archivo aleatorio y seguro
func generateSecureFilename(extension string) (string, error) {
	// Generar 16 bytes aleatorios (32 caracteres hex)
	randomBytes := make([]byte, 16)
	if _, err := rand.Read(randomBytes); err != nil {
		return "", err
	}

	// Convertir a string hexadecimal y añadir la extensión
	return hex.EncodeToString(randomBytes) + extension, nil
}

// parseInt64 convierte un string a int64
func parseInt64(s string) (int64, error) {
	var result int64
	_, err := fmt.Sscanf(s, "%d", &result)
	return result, err
}
