package main

import (
	"log"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"backend/database"
	"backend/handlers"
	"backend/middleware"
)

func main() {
	godotenv.Load()

	// Inicializar la conexión a la base de datos
	_, err := database.InitDB()
	if err != nil {
		log.Fatalf("Error al inicializar la base de datos: %v", err)
	}

	// Configurar Gin para producción si es necesario
	if os.Getenv("ENV") == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.Default()

	// Configurar CORS middleware
	frontendURL := os.Getenv("FRONTEND_URL")
	allowedOrigins := []string{"http://localhost:3000"} // Desarrollo por defecto

	if frontendURL != "" {
		allowedOrigins = append(allowedOrigins, frontendURL)
	}

	router.Use(cors.New(cors.Config{
		AllowOrigins:     allowedOrigins,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	}))

	router.Use(middleware.LoggingMiddleware())

	// API endpoints
	api := router.Group("/api")

	// Rutas públicas (no requieren autenticación)
	api.POST("/register", handlers.RegisterHandler())
	api.POST("/login", handlers.LoginHandler())

	// Rutas para formularios (autenticación opcional)
	forms := api.Group("/forms")
	forms.Use(middleware.OptionalAuthMiddleware())
	{
		forms.POST("/contact", handlers.CreateContactFormHandler())
		forms.POST("/feedback", handlers.CreateFeedbackFormHandler())
	}

	// Rutas para documentos (subida disponible para todos, autenticación opcional)
	documents := api.Group("/documents")
	documents.Use(middleware.OptionalAuthMiddleware())
	{
		documents.POST("", handlers.UploadDocumentHandler())
	}

	// Rutas para análisis (disponible para todos, autenticación opcional)
	analysis := api.Group("/analysis")
	analysis.Use(middleware.OptionalAuthMiddleware())
	{
		analysis.POST("", handlers.CreateAnalysisRequestHandler())
		analysis.GET("/:id", handlers.GetAnalysisResultHandler())
	}

	// Rutas protegidas (requieren autenticación)
	protected := api.Group("")
	protected.Use(middleware.AuthMiddleware())
	{
		// Perfil de usuario
		protected.GET("/profile", handlers.GetProfileHandler())
		protected.PUT("/user/update", handlers.UpdateUserHandler())
		protected.DELETE("/user/delete", handlers.DeleteUserHandler())

		// Dashboard y estadísticas
		protected.GET("/user/stats", handlers.GetUserStatsHandler())
		protected.GET("/user/recent-activity", handlers.GetUserRecentActivityHandler())
		protected.GET("/user/recent-documents", handlers.GetUserRecentDocumentsHandler())

		// Documentos del usuario
		protected.GET("/user/documents", handlers.GetUserDocumentsHandler())
		protected.GET("/user/documents/:id", handlers.GetDocumentWithAnalysisHandler())
		protected.DELETE("/user/documents/:id", handlers.DeleteDocumentHandler())
		protected.DELETE("/user/documents/:id/permanent", handlers.PermanentDeleteDocumentHandler())

		// Análisis del usuario
		protected.GET("/user/analysis", handlers.GetUserAnalysisRequestsHandler())
	}

	// Configurar puerto
	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}

	log.Printf("Servidor Gin iniciado en el puerto %s\n", port)
	log.Fatal(router.Run(":" + port))
}
