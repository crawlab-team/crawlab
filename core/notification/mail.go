package notification

import (
	"errors"
	"github.com/apex/log"
	"github.com/crawlab-team/crawlab/core/models/models"
	"github.com/matcornic/hermes/v2"
	"gopkg.in/gomail.v2"
	"net/mail"
	"runtime/debug"
	"strconv"
	"strings"
)

func SendMail(m *models.NotificationSettingMail, to, cc, title, content string) error {
	// theme
	theme := new(MailThemeFlat)

	// hermes instance
	h := hermes.Hermes{
		Theme: theme,
		Product: hermes.Product{
			Logo:      "data:image/svg+xml;base64,PHN2ZyB3aWR0aD0iMzAwIiBoZWlnaHQ9IjMwMCIgeG1sbnM9Imh0dHA6Ly93d3cudzMub3JnLzIwMDAvc3ZnIj4KICAgIDxnIGZpbGw9Im5vbmUiPgogICAgICAgIDxjaXJjbGUgY3g9IjE1MCIgY3k9IjE1MCIgcj0iMTMwIiBmaWxsPSJub25lIiBzdHJva2Utd2lkdGg9IjQwIiBzdHJva2U9IiM0MDllZmYiPgogICAgICAgIDwvY2lyY2xlPgogICAgICAgIDxjaXJjbGUgY3g9IjE1MCIgY3k9IjE1MCIgcj0iMTEwIiBmaWxsPSJ3aGl0ZSI+CiAgICAgICAgPC9jaXJjbGU+CiAgICAgICAgPGNpcmNsZSBjeD0iMTUwIiBjeT0iMTUwIiByPSI3MCIgZmlsbD0iIzQwOWVmZiI+CiAgICAgICAgPC9jaXJjbGU+CiAgICAgICAgPHBhdGggZD0iCiAgICAgICAgICAgIE0gMTUwLDE1MAogICAgICAgICAgICBMIDI4MCwyMjUKICAgICAgICAgICAgQSAxNTAsMTUwIDkwIDAgMCAyODAsNzUKICAgICAgICAgICAgIiBmaWxsPSIjNDA5ZWZmIj4KICAgICAgICA8L3BhdGg+CiAgICA8L2c+Cjwvc3ZnPgo=",
			Name:      "Crawlab",
			Copyright: "© 2024 Crawlab-Team",
		},
	}

	// config
	port, _ := strconv.Atoi(m.Port)
	password := m.Password // test password: ALWVDPRHBEXOENXD
	SMTPUser := m.User
	smtpConfig := smtpAuthentication{
		Server:         m.Server,
		Port:           port,
		SenderEmail:    m.SenderEmail,
		SenderIdentity: m.SenderIdentity,
		SMTPPassword:   password,
		SMTPUser:       SMTPUser,
	}
	options := sendOptions{
		To:      to,
		Cc:      cc,
		Subject: title,
	}

	// add style
	content += theme.GetStyle()

	// markdown
	markdown := hermes.Markdown(content + GetFooter())

	// email instance
	email := hermes.Email{
		Body: hermes.Body{
			Signature:    "Happy Crawling ☺",
			FreeMarkdown: markdown,
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
	Cc      string
}

// send email
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

	var toList []string
	if strings.Contains(options.To, ";") {
		toList = strings.Split(options.To, ";")
		// trim space
		for i, to := range toList {
			toList[i] = strings.TrimSpace(to)
		}
	} else {
		toList = []string{options.To}
	}

	m := gomail.NewMessage()
	m.SetHeader("From", from.String())
	m.SetHeader("To", getRecipientList(options.To)...)
	m.SetHeader("Subject", options.Subject)
	if options.Cc != "" {
		m.SetHeader("Cc", getRecipientList(options.Cc)...)
	}

	m.SetBody("text/plain", txtBody)
	m.AddAlternative("text/html", htmlBody)

	d := gomail.NewDialer(smtpConfig.Server, smtpConfig.Port, smtpConfig.SMTPUser, smtpConfig.SMTPPassword)

	return d.DialAndSend(m)
}

func getRecipientList(value string) (values []string) {
	if strings.Contains(value, ";") {
		values = strings.Split(value, ";")
		// trim space
		for i, v := range values {
			values[i] = strings.TrimSpace(v)
		}
	} else {
		values = []string{value}
	}
	return values
}

func GetFooter() string {
	return `
[Github](https://github.com/crawlab-team/crawlab) | [Documentation](http://docs.crawlab.cn) | [Docker](https://hub.docker.com/r/tikazyq/crawlab)
`
}
