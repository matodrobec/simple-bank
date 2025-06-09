package mail

import (
	"fmt"
	"net/smtp"

	"github.com/jordan-wright/email"
)

type GenericSender struct {
	smtpUser       string
	smtpPassord    string
	smtpHost       string
	smtpPort       int
	smtpEncryption string

	fromName         string
	fromEmailAddress string
}

// MAIL_DOMAIN=localhost
// MAIL_HOST=mailhog
// MAIL_PORT=1025
// MAIL_ENCRYPTION=none
// MAIL_USERNAME=""
// MAIL_PASSWORD=""
// FROM_NAME="Martin Balaz"
// FROM_ADDRESS="martin@test.loc"

func NewGenericSender(config MailSenderConfig) EmaiSender {
	return &GenericSender{
		smtpUser:       config.GetSmtpUser(),
		smtpPassord:    config.GetSmtpPassword(),
		smtpHost:       config.GetSmtpHost(),
		smtpPort:       config.GetSmtpPort(),
		smtpEncryption: config.GetSmtpEncryption(),

		fromName:         config.GetFromName(),
		fromEmailAddress: config.GetFromEmailAddress(),
	}
}

func (sender *GenericSender) SendEmail(
	subject string,
	content string,
	to []string,
	cc []string,
	bcc []string,
	attachFiles []string,
) error {
	email := email.NewEmail()
	email.From = fmt.Sprintf("%s <%s>", sender.fromName, sender.fromEmailAddress)
	email.Subject = subject
	email.HTML = []byte(content)
	email.To = to
	email.Cc = cc
	email.Bcc = bcc

	for _, f := range attachFiles {
		if _, err := email.AttachFile(f); err != nil {
			return fmt.Errorf("failed to attache file %s: %w", f, err)
		}

	}

	smtpAuth := smtp.PlainAuth("", sender.smtpUser, sender.smtpPassord, sender.smtpHost)
	return email.Send(fmt.Sprintf("%s:%d", sender.smtpHost, sender.smtpPort), smtpAuth)

}
