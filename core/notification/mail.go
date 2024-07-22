package notification

import (
	"errors"
	"github.com/PuerkitoBio/goquery"
	"github.com/apex/log"
	"github.com/crawlab-team/crawlab/core/models/models/v2"
	"github.com/crawlab-team/crawlab/trace"
	"gopkg.in/gomail.v2"
	"net/mail"
	"regexp"
	"runtime/debug"
	"strings"
)

func SendMail(s *models.NotificationSettingV2, ch *models.NotificationChannelV2, to, cc, bcc, title, content string) error {
	// config
	smtpConfig := smtpAuthentication{
		Server:         ch.SMTPServer,
		Port:           ch.SMTPPort,
		SenderIdentity: s.SenderName,
		SenderEmail:    s.SenderEmail,
		SMTPUser:       ch.SMTPUsername,
		SMTPPassword:   ch.SMTPPassword,
	}

	options := sendOptions{
		Subject: title,
		To:      to,
		Cc:      cc,
		Bcc:     bcc,
	}

	// convert html to text
	text := content
	if isHtml(text) {
		text = convertHtmlToText(text)
	}

	// apply theme
	if isHtml(content) {
		content = GetTheme() + content
	}

	// send the email
	if err := send(smtpConfig, options, content, text); err != nil {
		log.Errorf(err.Error())
		debug.PrintStack()
		return err
	}

	return nil
}

func isHtml(content string) bool {
	regex := regexp.MustCompile(`(?i)<\s*(html|head|body|div|span|p|a|img|table|tr|td|th|tbody|thead|tfoot|ul|ol|li|dl|dt|dd|form|input|textarea|button|select|option|optgroup|fieldset|legend|label|iframe|embed|object|param|video|audio|source|canvas|svg|math|style|link|script|meta|base|title|br|hr|b|strong|i|em|u|s|strike|del|ins|mark|small|sub|sup|big|pre|code|q|blockquote|abbr|address|bdo|cite|dfn|kbd|var|samp|ruby|rt|rp|time|progress|meter|output|area|map)`)
	return regex.MatchString(content)
}

func convertHtmlToText(content string) string {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(content))
	if err != nil {
		log.Errorf("failed to convert html to text: %v", err)
		trace.PrintError(err)
		return ""
	}
	return doc.Text()
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
	Subject string
	To      string
	Cc      string
	Bcc     string
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
	if options.Bcc != "" {
		m.SetHeader("Bcc", getRecipientList(options.Bcc)...)
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
