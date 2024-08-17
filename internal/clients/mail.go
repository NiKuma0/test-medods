package clients

import (
	"fmt"
	"net/smtp"

	smtpmock "github.com/mocktools/go-smtp-mock/v2"
)

type MailClient struct {
	client *smtp.Client
	sender string
}

func NewMailClient(sender, smtpAddress string, auth smtp.Auth) (s MailClient, err error) {
	c, err := smtp.Dial(smtpAddress)
	s = MailClient{client: c, sender: sender}
	if err != nil {
		return
	}
	err = c.Auth(auth)
	return
}

func NewMailClientNoAuth(sender, smtpAddress string) (s MailClient, err error) {
	c, err := smtp.Dial(smtpAddress)
	s = MailClient{client: c, sender: sender}
	return
}

func StartNewMockServer() (server *smtpmock.Server) {
	server = smtpmock.New(smtpmock.ConfigurationAttr{
		LogToStdout:       true,
		LogServerActivity: true,
	})
	if err := server.Start(); err != nil {
		panic(err)
	}
	return
}

func (s *MailClient) Shutdown() {
	s.client.Quit()
	s.client.Close()
}

func (s *MailClient) SendNotification(msg, recipient string) (err error) {
	if err = s.client.Mail(s.sender); err != nil {
		return
	}
	if err = s.client.Rcpt(recipient); err != nil {
		return
	}
	wc, err := s.client.Data()
	if err != nil {
		return
	}
	defer wc.Close()
	_, err = fmt.Fprint(wc, msg)
	return
}
