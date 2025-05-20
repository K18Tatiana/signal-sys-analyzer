package database

import (
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

var DB *gorm.DB

// InitDB inicializa la conexión a la base de datos
func InitDB() (*gorm.DB, error) {
	// Configuración de la conexión a PostgreSQL
	host := getEnv("DB_HOST", "localhost")
	port := getEnv("DB_PORT", "5432")
	user := getEnv("DB_USER", "postgres")
	password := getEnv("DB_PASSWORD", "sasoftco0212")
	dbname := getEnv("DB_NAME", "signal_sys_analysis")

	// Cadena de conexión
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	// Configurar logger de GORM
	gormLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold:             time.Second, // Umbral para considerar una consulta como lenta
			LogLevel:                  logger.Info, // Nivel de log por defecto
			IgnoreRecordNotFoundError: true,        // Ignorar errores de "registro no encontrado"
			Colorful:                  true,        // Habilitar colores
		},
	)

	if getEnv("ENV", "development") == "production" {
		gormLogger = logger.New(
			log.New(os.Stdout, "\r\n", log.LstdFlags),
			logger.Config{
				SlowThreshold:             time.Second,
				LogLevel:                  logger.Error,
				IgnoreRecordNotFoundError: true,
				Colorful:                  false,
			},
		)
	}

	// Abrir la conexión con configuración importante para base de datos existente
	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: gormLogger,
		// No crear restricciones de clave foránea cuando se hace migración
		DisableForeignKeyConstraintWhenMigrating: true,
		// Permitir que los nombres de tablas sean plurales (como ya existen)
		NamingStrategy: schema.NamingStrategy{
			SingularTable: false,
		},
	})

	if err != nil {
		return nil, err
	}

	// Verificar si es primera ejecución o migración
	var tableCount int64
	DB.Raw("SELECT COUNT(*) FROM information_schema.tables WHERE table_schema = 'public'").Scan(&tableCount)

	if tableCount == 0 {
		log.Println("ADVERTENCIA: No se encontraron tablas en la base de datos. Esto puede indicar un problema de conexión.")
		return DB, nil
	}

	// Configurar GORM para trabajar con el esquema existente
	if err := MigrateDatabase(DB); err != nil {
		log.Printf("Advertencia: Problemas al configurar GORM para el esquema existente: %v", err)
		log.Println("Continuando a pesar de advertencias. Algunas funcionalidades podrían no operar correctamente.")
	}

	// Configurar pool de conexiones
	sqlDB, err := DB.DB()
	if err != nil {
		return nil, err
	}

	// Establecer el número máximo de conexiones abiertas
	sqlDB.SetMaxOpenConns(100)
	// Establecer el número máximo de conexiones inactivas
	sqlDB.SetMaxIdleConns(10)
	// Establecer el tiempo máximo de vida de una conexión
	sqlDB.SetConnMaxLifetime(time.Hour)

	log.Println("Conexión exitosa a la base de datos PostgreSQL")
	return DB, nil
}

// getEnv obtiene una variable de entorno, o devuelve un valor predeterminado
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

// CloseDB cierra la conexión a la base de datos
func CloseDB() {
	if DB != nil {
		sqlDB, err := DB.DB()
		if err == nil {
			sqlDB.Close()
			log.Println("Conexión a la base de datos cerrada")
		}
	}
}
