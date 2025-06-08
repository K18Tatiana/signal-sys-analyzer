package database

import (
	"fmt"
	"log"

	"backend/models"
	"gorm.io/gorm"
)

// MigrateDatabase configura GORM para trabajar con la base de datos
func MigrateDatabase(db *gorm.DB) error {
	log.Println("Configurando GORM para trabajar con el esquema de base de datos...")

	// Verificar si la base de datos está accesible
	var tableCount int64
	if err := db.Raw("SELECT COUNT(*) FROM information_schema.tables WHERE table_schema = 'public'").Scan(&tableCount).Error; err != nil {
		return err
	}

	log.Printf("Encontradas %d tablas en la base de datos.", tableCount)

	if tableCount == 0 {
		log.Println("No se encontraron tablas. Creando esquema completo desde los modelos...")

		if err := db.AutoMigrate(
			&models.User{},
			&models.Document{},
			&models.AnalysisRequest{},
			&models.Result{},
			&models.ContactForm{},
			&models.FeedbackForm{},
		); err != nil {
			log.Printf("Error al crear las tablas: %v", err)
			return err
		}

		log.Println("Esquema de base de datos creado exitosamente")
		return nil
	}

	// Si ya hay tablas, aplicar migraciones para nuevas columnas
	if err := applyMigrations(db); err != nil {
		log.Printf("Error aplicando migraciones: %v", err)
		return err
	}

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

// applyMigrations aplica migraciones específicas para nuevas columnas
func applyMigrations(db *gorm.DB) error {
	log.Println("Aplicando migraciones de esquema...")

	// Verificar y agregar columna graph_data en la tabla results
	if err := addNewColumnIfNotExists(db, "results", "graph_data", "JSONB"); err != nil {
		return err
	}

	// Verificar y agregar columna comment en la tabla analysis_requests
	if err := addNewColumnIfNotExists(db, "analysis_requests", "comment", "VARCHAR(500)"); err != nil {
		return err
	}

	// Agregar columnas ML en la tabla results
	if err := addNewColumnIfNotExists(db, "results", "ml_predicted_type", "INTEGER"); err != nil {
		return err
	}
	if err := addNewColumnIfNotExists(db, "results", "ml_polo1_real", "REAL"); err != nil {
		return err
	}
	if err := addNewColumnIfNotExists(db, "results", "ml_polo1_imag", "REAL"); err != nil {
		return err
	}
	if err := addNewColumnIfNotExists(db, "results", "ml_polo2_real", "REAL"); err != nil {
		return err
	}
	if err := addNewColumnIfNotExists(db, "results", "ml_polo2_imag", "REAL"); err != nil {
		return err
	}
	if err := addNewColumnIfNotExists(db, "results", "ml_confidence", "REAL"); err != nil {
		return err
	}
	if err := addNewColumnIfNotExists(db, "results", "description", "VARCHAR(1000)"); err != nil {
		return err
	}
	if err := addNewColumnIfNotExists(db, "results", "technical_summary", "JSONB"); err != nil {
		return err
	}

	log.Println("Todas las migraciones aplicadas correctamente")
	return nil
}

// Función auxiliar para agregar columnas de forma segura
func addNewColumnIfNotExists(db *gorm.DB, tableName, columnName, columnType string) error {
	var columnExists bool
	err := db.Raw(`
        SELECT EXISTS (
            SELECT 1 
            FROM information_schema.columns 
            WHERE table_name = ? 
            AND column_name = ?
            AND table_schema = 'public'
        )
    `, tableName, columnName).Scan(&columnExists).Error

	if err != nil {
		return err
	}

	if !columnExists {
		log.Printf("Agregando columna %s.%s...", tableName, columnName)

		sql := fmt.Sprintf("ALTER TABLE %s ADD COLUMN %s %s", tableName, columnName, columnType)
		err = db.Exec(sql).Error
		if err != nil {
			return err
		}

		log.Printf("Columna %s.%s agregada exitosamente", tableName, columnName)
	} else {
		log.Printf("Columna %s.%s ya existe", tableName, columnName)
	}

	return nil
}
