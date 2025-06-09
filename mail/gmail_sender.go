package mail

import (
	"fmt"
	"net/smtp"

	"github.com/jordan-wright/email"
)

const (
	gmailSmtpAuthAddress   = "smtp.gmail.com"
	gmailSmtpServerAddress = "smtp.gmail.com:587"
)

type GmailSender struct {
	name              string
	fromEmailAddress  string
	fromEmailPassword string
}

func NewGmailSender(name, fromEmailAdress, fromEmailPassword string) EmaiSender {
	return &GmailSender{
		name:              name,
		fromEmailAddress:  fromEmailAdress,
		fromEmailPassword: fromEmailPassword,
	}
}

func (sender *GmailSender) SendEmail(
	subject string,
	content string,
	to []string,
	cc []string,
	bcc []string,
	attachFiles []string,
) error {
	email := email.NewEmail()
	email.From = fmt.Sprintf("%s <%s>", sender.name, sender.fromEmailAddress)
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

	smtpAuth := smtp.PlainAuth("", sender.fromEmailAddress, sender.fromEmailPassword, gmailSmtpAuthAddress)
	return email.Send(gmailSmtpServerAddress, smtpAuth)

}
