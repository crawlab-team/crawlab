package notification

import (
	"errors"
	"github.com/apex/log"
	"github.com/matcornic/hermes"
	"gopkg.in/gomail.v2"
	"net/mail"
	"os"
	"runtime/debug"
	"strconv"
)

func SendMail(toEmail string, subject string, content string) error {
	// hermes instance
	h := hermes.Hermes{
		Product: hermes.Product{
			Name: "Hermes",
			Link: "https://example-hermes.com/",
			Logo: "http://www.duchess-france.org/wp-content/uploads/2016/01/gopher.png",
		},
	}

	// config
	port, _ := strconv.Atoi(os.Getenv("CRAWLAB_NOTIFICATION_MAIL_PORT"))
	password := os.Getenv("CRAWLAB_NOTIFICATION_MAIL_SMTP_PASSWORD")
	SMTPUser := os.Getenv("CRAWLAB_NOTIFICATION_MAIL_SMTP_USER")
	smtpConfig := smtpAuthentication{
		Server:         os.Getenv("CRAWLAB_NOTIFICATION_MAIL_SERVER"),
		Port:           port,
		SenderEmail:    os.Getenv("CRAWLAB_NOTIFICATION_MAIL_SENDEREMAIL"),
		SenderIdentity: os.Getenv("CRAWLAB_NOTIFICATION_MAIL_SENDERIDENTITY"),
		SMTPPassword:   password,
		SMTPUser:       SMTPUser,
	}
	options := sendOptions{
		To:      toEmail,
		Subject: subject,
	}

	// email instance
	email := hermes.Email{
		Body: hermes.Body{
			FreeMarkdown: hermes.Markdown(content),
		},
	}

	// generate html
	html, err := h.GenerateHTML(email)
	if err != nil {
		log.Errorf(err.Error())
		debug.PrintStack()
		return err
	}

	// generate text
	text, err := h.GeneratePlainText(email)
	if err != nil {
		log.Errorf(err.Error())
		debug.PrintStack()
		return err
	}

	// send the email
	if err := send(smtpConfig, options, html, text); err != nil {
		log.Errorf(err.Error())
		debug.PrintStack()
		return err
	}

	return nil
}

type smtpAuthentication struct {
	Server         string
	Port           int
	SenderEmail    string
	SenderIdentity string
	SMTPUser       string
	SMTPPassword   string
}

// sendOptions are options for sending an email
type sendOptions struct {
	To      string
	Subject string
}

// send sends the email
func send(smtpConfig smtpAuthentication, options sendOptions, htmlBody string, txtBody string) error {

	if smtpConfig.Server == "" {
		return errors.New("SMTP server config is empty")
	}
	if smtpConfig.Port == 0 {
		return errors.New("SMTP port config is empty")
	}

	if smtpConfig.SMTPUser == "" {
		return errors.New("SMTP user is empty")
	}

	if smtpConfig.SenderIdentity == "" {
		return errors.New("SMTP sender identity is empty")
	}

	if smtpConfig.SenderEmail == "" {
		return errors.New("SMTP sender email is empty")
	}

	if options.To == "" {
		return errors.New("no receiver emails configured")
	}

	from := mail.Address{
		Name:    smtpConfig.SenderIdentity,
		Address: smtpConfig.SenderEmail,
	}

	m := gomail.NewMessage()
	m.SetHeader("From", from.String())
	m.SetHeader("To", options.To)
	m.SetHeader("Subject", options.Subject)

	m.SetBody("text/plain", txtBody)
	m.AddAlternative("text/html", htmlBody)

	d := gomail.NewPlainDialer(smtpConfig.Server, smtpConfig.Port, smtpConfig.SMTPUser, smtpConfig.SMTPPassword)

	return d.DialAndSend(m)
}
