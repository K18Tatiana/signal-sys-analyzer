package database

import (
	"log"

	"backend/models"
	"gorm.io/gorm"
)

// MigrateDatabase configura GORM para trabajar con la base de datos existente
func MigrateDatabase(db *gorm.DB) error {
	log.Println("Configurando GORM para trabajar con el esquema de base de datos existente...")

	// Verificar si la base de datos está accesible
	var tableCount int64
	if err := db.Raw("SELECT COUNT(*) FROM information_schema.tables WHERE table_schema = 'public'").Scan(&tableCount).Error; err != nil {
		return err
	}

	log.Printf("Encontradas %d tablas en la base de datos.", tableCount)

	// Desactivar las funcionalidades de auto-migración que podrían modificar el esquema
	db = db.Session(&gorm.Session{
		DryRun: true, // No ejecutar consultas reales, solo preparar
	})

	// Verificar si podemos acceder a los usuarios
	var userCount int64
	if err := db.Model(&models.User{}).Count(&userCount).Error; err != nil {
		log.Printf("Error al acceder a la tabla users: %v", err)
		log.Println("Intentando consulta SQL directa...")

		// Intentar consulta directa en caso de fallo
		if err := db.Raw("SELECT COUNT(*) FROM users").Scan(&userCount).Error; err != nil {
			log.Printf("También falló la consulta SQL directa: %v", err)
			return err
		}
	}

	log.Printf("Hay %d usuarios en el sistema.", userCount)

	// Verificar conexiones a otras tablas
	var docCount int64
	if err := db.Raw("SELECT COUNT(*) FROM documents").Scan(&docCount).Error; err != nil {
		log.Printf("Advertencia: Error al acceder a la tabla documents: %v", err)
	} else {
		log.Printf("Hay %d documentos en el sistema.", docCount)
	}

	log.Println("Configuración de base de datos completada con éxito.")
	return nil
}
