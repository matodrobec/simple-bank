package mail

import (
	"testing"

	"github.com/matodrobec/simplebank/util"
	"github.com/stretchr/testify/require"
)

func TestSendEmailWithGmail(t *testing.T) {
    if testing.Short() {
        t.Skip()
    }

	config, err := util.LoadConfig("..")
	require.NoError(t, err)

	var mailConfig  = MailSenderConfig(config)


	sender := NewGmailSender(
		mailConfig.GetFromName(),
		mailConfig.GetFromEmailAddress(),
		mailConfig.GetSmtpPassword(),
		// config.EmailSenderName,
		// config.EmailSenderAddress,
		// config.EmailSenderPassword,
	)

	subject := "A test email"
	content := `
    <h1>Test Email</h1>
    <p>This is a test messge</p>
    `
	to := []string{"mb@amazingh.sk"}
	attachFiles := []string{"../test_email_doc.txt"}

	err = sender.SendEmail(subject, content, to, nil, nil, attachFiles)
	require.NoError(t, err)
}
