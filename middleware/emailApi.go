package middleware

import (
	"mini_project_p2/models"
	"os"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

func SendConfirmationEmail(user *models.User) error {
	from := mail.NewEmail("Handy Library", "mhandyalfurqon@gmail.com")
	subject := "Konfirmasi Pendaftaran"
	to := mail.NewEmail(user.Username, user.Email)
	plainTextContent := "Terima kasih telah mendaftar!"
	htmlContent := "<strong>Terima kasih telah mendaftar!</strong>"
	message := mail.NewSingleEmail(from, subject, to, plainTextContent, htmlContent)

	client := sendgrid.NewSendClient(os.Getenv("emailApi")) // Ganti dengan API key SendGrid kamu
	_, err := client.Send(message)

	return err
}
