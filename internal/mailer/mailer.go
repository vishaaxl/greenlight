package mailer

import (
	"bytes"
	"embed"
	"html/template"
	"time"

	"github.com/go-mail/mail"
)

//go:embed "templates"
var templateFS embed.FS

type Mailer struct {
	dialer *mail.Dialer
	sender string
}

func New(host string, port int, username, password, sender string) Mailer {
	dialer := mail.NewDialer(host, port, username, password)
	dialer.Timeout = 5 * time.Second

	return Mailer{dialer: dialer, sender: sender}
}

func (m Mailer) Send(recipientEmail, subject, templateFile string, data interface{}) error {

	// Parse the HTML email template from the embedded file system
	tmpl, err := template.ParseFS(templateFS, "templates/"+templateFile)
	if err != nil {
		return err
	}

	// Create a buffer to hold the rendered template with dynamic data
	var body bytes.Buffer
	if err = tmpl.Execute(&body, data); err != nil {
		return err // Return an error if the template cannot be executed with the provided data
	}

	// Create a new email message
	msg := mail.NewMessage()

	msg.SetHeader("To", recipientEmail)     // Set the recipient's email address
	msg.SetHeader("From", m.sender)         // Set the sender's email address
	msg.SetHeader("Subject", subject)       // Set the email subject
	msg.SetBody("text/html", body.String()) // Set the email body as the rendered HTML template

	// Try sending the email up to three times before aborting and returning the final
	// error. We sleep for 500 milliseconds between each attempt.
	for i := 1; i <= 3; i++ {
		err = m.dialer.DialAndSend(msg)

		if nil == err {
			return nil
		}

		time.Sleep(500 * time.Millisecond)
	}

	// Send the email using the configured mail dialer

	return err // Return nil if the email is sent successfully
}
