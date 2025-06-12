package mail

import mail "github.com/xhit/go-simple-mail/v2"

type SmtpSender struct {
	MailSenderConfig
}

func NewSmtpSender(config MailSenderConfig) EmaiSender {
	return &SmtpSender{
		config,
	}
}

func (sender *SmtpSender) SendEmail(
	subject string,
	content string,
	to []string,
	cc []string,
	bcc []string,
	attacheFiles []string,
) error {
	server := mail.NewSMTPClient()
	server.Host = sender.GetSmtpHost()
	server.Port = sender.GetSmtpPort()
	server.Username = sender.GetSmtpUser()
	server.Password = sender.GetSmtpPassword()
	server.Encryption = getEncryption(sender.GetSmtpEncryption())
	server.KeepAlive = false

	smtpClient, err := server.Connect()
	if err != nil {
		return err
	}

	email := mail.NewMSG()
	email.SetFrom(sender.GetFromEmailAddress()).
		AddTo(to...).
		SetSubject(subject).
		AddBcc(bcc...).
		AddCc(cc...)

	email.SetBody(mail.TextPlain, content)

	if len(attacheFiles) > 0 {
		for _, x := range attacheFiles {
			email.AddAttachment(x)
		}
	}

	return email.Send(smtpClient)
}

func getEncryption(s string) mail.Encryption {
	switch s {
	case "tls":
		return mail.EncryptionSTARTTLS
	case "ssl":
		return mail.EncryptionSSLTLS
	case "none", "":
		return mail.EncryptionNone
	default:
		return mail.EncryptionSTARTTLS
	}
}
