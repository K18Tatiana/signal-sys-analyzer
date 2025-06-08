package utils

import (
	"fmt"
	"os"
	"strconv"

	"gopkg.in/gomail.v2"
)

// EmailConfig contiene la configuración para enviar correos
type EmailConfig struct {
	Host     string
	Port     int
	Username string
	Password string
	FromName string
	FromAddr string
}

// NewEmailConfig crea una nueva configuración para enviar correos
func NewEmailConfig() *EmailConfig {
	// Leer puerto como número
	portStr := os.Getenv("SMTP_PORT")
	port, err := strconv.Atoi(portStr)
	if err != nil || port == 0 {
		port = 587 // Puerto predeterminado para TLS
	}

	return &EmailConfig{
		Host:     os.Getenv("SMTP_HOST"),
		Port:     port,
		Username: os.Getenv("SMTP_USERNAME"),
		Password: os.Getenv("SMTP_PASSWORD"),
		FromName: os.Getenv("EMAIL_FROM_NAME"),
		FromAddr: os.Getenv("EMAIL_FROM_ADDRESS"),
	}
}

// SendEmail envía un correo electrónico
func (ec *EmailConfig) SendEmail(to, subject, body string) error {
	// Verificar configuración
	if ec.Host == "" || ec.Username == "" || ec.Password == "" {
		return fmt.Errorf("configuración de correo incompleta, revisa las variables de entorno SMTP_*")
	}

	// Configurar remitente
	from := ec.FromAddr
	if from == "" {
		from = ec.Username
	}

	// Crear mensaje
	m := gomail.NewMessage()

	if ec.FromName != "" {
		m.SetHeader("From", fmt.Sprintf("%s <%s>", ec.FromName, from))
	} else {
		m.SetHeader("From", from)
	}

	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)

	// Configurar dialer
	d := gomail.NewDialer(ec.Host, ec.Port, ec.Username, ec.Password)

	// Enviar
	return d.DialAndSend(m)
}

// SendContactFormNotification envía una notificación al administrador
// cuando se recibe un formulario de contacto
func SendContactFormNotification(name, email, subject, message string) error {
	// Crear configuración de email
	emailConfig := NewEmailConfig()

	// Email del administrador al que se enviará la notificación
	adminEmail := os.Getenv("ADMIN_EMAIL")
	if adminEmail == "" {
		adminEmail = "tvegac2003@gmail.com" // Cambiar a tu correo real
	}

	// Crear cuerpo del correo
	htmlBody := fmt.Sprintf(`
	<html>
	<body>
		<h2>Nuevo mensaje de contacto recibido</h2>
		<p><strong>Nombre:</strong> %s</p>
		<p><strong>Email:</strong> %s</p>
		<p><strong>Asunto:</strong> %s</p>
		<p><strong>Mensaje:</strong></p>
		<div style="background-color: #f5f5f5; padding: 15px; border-radius: 5px;">
			%s
		</div>
		<hr>
		<p>Este mensaje fue enviado desde el formulario de contacto de Signal System Analysis.</p>
	</body>
	</html>
	`, name, email, subject, message)

	// Enviar email
	return emailConfig.SendEmail(adminEmail, "Nuevo mensaje de contacto: "+subject, htmlBody)
}

// SendFeedbackNotification envía una notificación al administrador
// cuando se recibe un formulario de feedback
func SendFeedbackNotification(rating int, feedback, email string) error {
	// Crear configuración de email
	emailConfig := NewEmailConfig()

	// Email del administrador al que se enviará la notificación
	adminEmail := os.Getenv("ADMIN_EMAIL")
	if adminEmail == "" {
		adminEmail = "tvegac2003@gmail.com" // Cambiar a tu correo real
	}

	// Crear representación visual del rating
	starsHTML := ""
	for i := 0; i < rating; i++ {
		starsHTML += "★"
	}
	for i := rating; i < 5; i++ {
		starsHTML += "☆"
	}

	// Crear cuerpo del correo
	htmlBody := fmt.Sprintf(`
	<html>
	<body>
		<h2>Nuevo feedback recibido</h2>
		<p><strong>Calificación:</strong> %d/5 <span style="color: gold; font-size: 20px;">%s</span></p>
		<p><strong>Email:</strong> %s</p>
		<p><strong>Comentarios:</strong></p>
		<div style="background-color: #f5f5f5; padding: 15px; border-radius: 5px;">
			%s
		</div>
		<hr>
		<p>Este feedback fue enviado desde el formulario de Signal System Analysis.</p>
	</body>
	</html>
	`, rating, starsHTML, email, feedback)

	// Enviar email
	return emailConfig.SendEmail(adminEmail, fmt.Sprintf("Nuevo feedback (Calificación: %d/5)", rating), htmlBody)
}
