package mail

type EmaiSender interface {
	SendEmail(
		subject string,
		content string,
		to []string,
		cc []string,
		bcc []string,
		sttacheFiles []string,
	) error
}

type MailSenderConfig interface {
	GetSmtpHost() string
	GetSmtpPort() int
	GetSmtpUser() string
	GetSmtpPassword() string
	GetSmtpEncryption() string
	GetFromName() string
	GetFromEmailAddress() string
}
