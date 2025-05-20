package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"

	"backend/database"
	"backend/handlers"
	"backend/middleware"
)

func main() {
	// Cargar variables de entorno desde el archivo .env si existe
	godotenv.Load()

	// Inicializar la conexión a la base de datos
	_, err := database.InitDB()
	if err != nil {
		log.Fatalf("Error al inicializar la base de datos: %v", err)
	}
	// No es necesario cerrar la conexión explícitamente con GORM

	// Crear router
	router := mux.NewRouter()

	// Middleware global para todos los endpoints
	router.Use(middleware.LoggingMiddleware)

	// API endpoints
	api := router.PathPrefix("/api").Subrouter()

	// Rutas públicas (no requieren autenticación)
	api.HandleFunc("/register", handlers.RegisterHandler()).Methods("POST")
	api.HandleFunc("/login", handlers.LoginHandler()).Methods("POST")

	// Rutas para formularios (autenticación opcional)
	formsRouter := api.PathPrefix("/forms").Subrouter()
	formsRouter.Use(middleware.OptionalAuthMiddleware)
	formsRouter.HandleFunc("/contact", handlers.CreateContactFormHandler()).Methods("POST")
	formsRouter.HandleFunc("/feedback", handlers.CreateFeedbackFormHandler()).Methods("POST")

	// Rutas para documentos (subida disponible para todos, autenticación opcional)
	documentsRouter := api.PathPrefix("/documents").Subrouter()
	documentsRouter.Use(middleware.OptionalAuthMiddleware)
	documentsRouter.HandleFunc("", handlers.UploadDocumentHandler()).Methods("POST")

	// Rutas para análisis (disponible para todos, autenticación opcional)
	analysisRouter := api.PathPrefix("/analysis").Subrouter()
	analysisRouter.Use(middleware.OptionalAuthMiddleware)
	analysisRouter.HandleFunc("", handlers.CreateAnalysisRequestHandler()).Methods("POST")
	analysisRouter.HandleFunc("/{id:[0-9]+}", handlers.GetAnalysisResultHandler()).Methods("GET")

	// Rutas protegidas (requieren autenticación)
	protectedRouter := api.PathPrefix("").Subrouter()
	protectedRouter.Use(middleware.AuthMiddleware)

	// Perfil de usuario
	protectedRouter.HandleFunc("/profile", handlers.GetProfileHandler()).Methods("GET")
	protectedRouter.HandleFunc("/user/update", handlers.UpdateUserHandler()).Methods("PUT")
	protectedRouter.HandleFunc("/user/delete", handlers.DeleteUserHandler()).Methods("DELETE")

	// Documentos del usuario
	protectedRouter.HandleFunc("/user/documents", handlers.GetUserDocumentsHandler()).Methods("GET")
	protectedRouter.HandleFunc("/user/documents/{id:[0-9]+}", handlers.DeleteDocumentHandler()).Methods("DELETE")
	// Ruta para eliminación permanente
	protectedRouter.HandleFunc("/user/documents/{id:[0-9]+}/permanent", handlers.PermanentDeleteDocumentHandler()).Methods("DELETE")

	// Análisis del usuario
	protectedRouter.HandleFunc("/user/analysis", handlers.GetUserAnalysisRequestsHandler()).Methods("GET")

	// Configurar y iniciar el servidor
	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}

	srv := &http.Server{
		Handler:      router,
		Addr:         "0.0.0.0:" + port,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	log.Printf("Servidor iniciado en el puerto %s\n", port)
	log.Fatal(srv.ListenAndServe())
}
