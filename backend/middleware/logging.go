package middleware

import (
	"log"
	"net/http"
	"time"
)

// LoggingMiddleware registra información sobre cada solicitud HTTP
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Estructura para capturar detalles de la respuesta
		wrapped := wrapResponseWriter(w)

		// Llamar al siguiente handler
		next.ServeHTTP(wrapped, r)

		// Calcular tiempo transcurrido
		duration := time.Since(start)

		// Registrar la solicitud con más detalles
		log.Printf(
			"%s %s %d %s %s",
			r.Method,
			r.RequestURI,
			wrapped.status,
			duration,
			r.RemoteAddr,
		)
	})
}

// Wrapper para capturar el código de estado
type responseWriter struct {
	http.ResponseWriter
	status      int
	wroteHeader bool
}

func wrapResponseWriter(w http.ResponseWriter) *responseWriter {
	return &responseWriter{ResponseWriter: w}
}

func (rw *responseWriter) Status() int {
	return rw.status
}

func (rw *responseWriter) WriteHeader(code int) {
	if !rw.wroteHeader {
		rw.status = code
		rw.ResponseWriter.WriteHeader(code)
		rw.wroteHeader = true
	}
}

func (rw *responseWriter) Write(b []byte) (int, error) {
	if !rw.wroteHeader {
		rw.WriteHeader(http.StatusOK)
	}
	return rw.ResponseWriter.Write(b)
}
