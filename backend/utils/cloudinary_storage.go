package utils

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

// CloudinaryStorage proporciona métodos para manejar archivos en Cloudinary
type CloudinaryStorage struct {
	CloudName string
	APIKey    string
	APISecret string
	Folder    string
}

// NewCloudinaryStorage crea una nueva instancia de CloudinaryStorage
func NewCloudinaryStorage() *CloudinaryStorage {
	return &CloudinaryStorage{
		CloudName: os.Getenv("CLOUDINARY_CLOUD_NAME"),
		APIKey:    os.Getenv("CLOUDINARY_API_KEY"),
		APISecret: os.Getenv("CLOUDINARY_API_SECRET"),
		Folder:    "Signal System Analysis", // Nombre exacto de la carpeta que creaste
	}
}

// UploadFile sube un archivo a Cloudinary
func (cs *CloudinaryStorage) UploadFile(file multipart.File, header *multipart.FileHeader) (string, error) {
	// Verificar credenciales
	if cs.CloudName == "" || cs.APIKey == "" || cs.APISecret == "" {
		return "", fmt.Errorf("faltan credenciales de Cloudinary, asegúrate de configurar las variables de entorno CLOUDINARY_CLOUD_NAME, CLOUDINARY_API_KEY y CLOUDINARY_API_SECRET")
	}

	// Verificar la extensión del archivo
	ext := filepath.Ext(header.Filename)
	if strings.ToLower(ext) != ".csv" {
		return "", fmt.Errorf("solo se permiten archivos CSV")
	}

	// Inicializar Cloudinary
	cld, err := cloudinary.NewFromParams(cs.CloudName, cs.APIKey, cs.APISecret)
	if err != nil {
		return "", fmt.Errorf("error al inicializar Cloudinary: %v", err)
	}

	// Generar nombre único basado en timestamp y nombre original
	timestamp := time.Now().Unix()
	filenameWithoutExt := strings.TrimSuffix(header.Filename, ext)

	// Usar solo el nombre de archivo en el publicID (sin carpeta)
	baseFilename := fmt.Sprintf("%d_%s", timestamp, filenameWithoutExt)

	// Subir archivo
	uploadResult, err := cld.Upload.Upload(
		context.Background(),
		file,
		uploader.UploadParams{
			PublicID:     baseFilename, // Solo el nombre base
			ResourceType: "raw",
			Folder:       cs.Folder, // La carpeta se especifica por separado
		},
	)

	if err != nil {
		return "", fmt.Errorf("error al subir archivo a Cloudinary: %v", err)
	}

	return uploadResult.SecureURL, nil
}

// GetFile obtiene un archivo de Cloudinary
func (cs *CloudinaryStorage) GetFile(url string) (io.ReadCloser, error) {
	// Obtener archivo directamente a través de la URL
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("error al obtener archivo de Cloudinary: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return nil, fmt.Errorf("error al obtener archivo, código de estado: %d", resp.StatusCode)
	}

	return resp.Body, nil
}

// DeleteFile elimina un archivo de Cloudinary
func (cs *CloudinaryStorage) DeleteFile(url string) error {
	// Verificar credenciales
	if cs.CloudName == "" || cs.APIKey == "" || cs.APISecret == "" {
		return fmt.Errorf("faltan credenciales de Cloudinary")
	}

	// Inicializar Cloudinary
	cld, err := cloudinary.NewFromParams(cs.CloudName, cs.APIKey, cs.APISecret)
	if err != nil {
		return fmt.Errorf("error al inicializar Cloudinary: %v", err)
	}

	// Extraer el nombre de la carpeta y archivo desde la URL
	// Formato de URL: https://res.cloudinary.com/cloud_name/raw/upload/v1234567890/folder/filename.csv

	// Obtenemos el último segmento de la URL (el nombre del archivo con extensión)
	filename := filepath.Base(url)

	// Recuperamos la carpeta desde la estructura de la URL
	// Dividimos la URL por "/upload/" y tomamos la segunda parte
	urlParts := strings.Split(url, "/upload/")
	if len(urlParts) < 2 {
		return fmt.Errorf("URL de Cloudinary inválida: %s", url)
	}

	// La parte después de /upload/ puede tener un componente de versión (v1234567890)
	// Dividimos por "/" y buscamos la parte que corresponde a la carpeta
	pathParts := strings.Split(urlParts[1], "/")

	// Necesitamos obtener el public ID completo incluyendo la carpeta
	var publicID string

	// Si la URL tiene estructura de carpetas
	if len(pathParts) >= 3 { // [version, folder, filename]
		// El public ID es "folder/filename_sin_extension"
		publicID = fmt.Sprintf("%s/%s", cs.Folder, strings.TrimSuffix(filename, filepath.Ext(filename)))
	} else {
		// Si no hay estructura de carpetas (poco probable dado nuestro código de subida)
		publicID = strings.TrimSuffix(filename, filepath.Ext(filename))
	}

	// Eliminar archivo
	_, err = cld.Upload.Destroy(context.Background(), uploader.DestroyParams{
		PublicID:     publicID,
		ResourceType: "raw",
	})

	if err != nil {
		return fmt.Errorf("error al eliminar archivo de Cloudinary: %v", err)
	}

	return nil
}
