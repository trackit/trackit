//   Copyright 2017 MSolution.IO
//
//   Licensed under the Apache License, Version 2.0 (the "License");
//   you may not use this file except in compliance with the License.
//   You may obtain a copy of the License at
//
//       http://www.apache.org/licenses/LICENSE-2.0
//
//   Unless required by applicable law or agreed to in writing, software
//   distributed under the License is distributed on an "AS IS" BASIS,
//   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//   See the License for the specific language governing permissions and
//   limitations under the License.

package mail

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"
	"net/smtp"

	"github.com/trackit/jsonlog"

	"github.com/trackit/trackit/config"
)

// Mail contains the data necessary to send a mail.
type Mail struct {
	SmtpAddress  string
	SmtpPort     string
	SmtpUser     string
	SmtpPassword string
	Sender       string
	Recipient    string
	Subject      string
	Body         string
}

// SendMail is the easiest way to send a mail.
// It gets the SMTP information from the config file.
func SendMail(recipient string, subject, body string, ctx context.Context) error {
	mail := Mail{
		config.SmtpAddress,
		config.SmtpPort,
		config.SmtpUser,
		config.SmtpPassword,
		config.SmtpSender,
		recipient,
		subject,
		body,
	}
	return mail.Send(ctx)
}

func (m Mail) buildMessage() []byte {
	message := ""
	message += fmt.Sprintf("From: %s\r\n", m.Sender)
	message += fmt.Sprintf("To: %s\r\n", m.Recipient)
	message += fmt.Sprintf("Subject: %s\r\n", m.Subject)
	message += "\r\n" + m.Body
	return []byte(message)
}

func (m Mail) setTlsConfig(client *smtp.Client) error {
	tlsconfig := &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         m.SmtpAddress,
	}
	if err := client.StartTLS(tlsconfig); err != nil {
		return err
	}
	return nil
}

func (m Mail) setAuth(client *smtp.Client) error {
	auth := smtp.PlainAuth(
		"",
		m.SmtpUser,
		m.SmtpPassword,
		m.SmtpAddress,
	)
	if err := client.Auth(auth); err != nil {
		return err
	}
	return nil
}

// getSmtpClient creates a smtp.Client with information
// from the Mail structure.
func (m Mail) getSmtpClient(ctx context.Context) (*smtp.Client, error) {

	conn, err := net.Dial("tcp", m.SmtpAddress+":"+m.SmtpPort)
	if err != nil {
		return nil, err
	}
	client, err := smtp.NewClient(conn, m.SmtpAddress)
	if err != nil {
		return client, err
	}
	if err := m.setTlsConfig(client); err != nil {
		return client, err
	}
	if err := m.setAuth(client); err != nil {
		return client, err
	}
	return client, nil
}

func (m Mail) setAddresses(client *smtp.Client) error {
	if err := client.Mail(m.Sender); err != nil {
		return err
	}
	if err := client.Rcpt(m.Recipient); err != nil {
		return err
	}
	return nil
}

func (m Mail) setMessage(client *smtp.Client) error {
	w, err := client.Data()
	if err != nil {
		return err
	}
	if _, err := w.Write(m.buildMessage()); err != nil {
		return err
	}
	if err := w.Close(); err != nil {
		return err
	}
	return nil
}

// Send provides a way to send a mail with SMTP information
// from the Mail structure.
func (m Mail) Send(ctx context.Context) error {
	dataLogged := map[string]interface{}{"subject": m.Subject, "recipient": m.Recipient}
	logger := jsonlog.LoggerFromContextOrDefault(ctx)
	logger.Info("Sending mail.", dataLogged)

	client, err := m.getSmtpClient(ctx)
	if err != nil {
		return err
	}
	if err := m.setAddresses(client); err != nil {
		return err
	}
	if err := m.setMessage(client); err != nil {
		return err
	}
	if err := client.Quit(); err != nil {
		return err
	}
	logger.Info("Mail successfully sent.", dataLogged)
	return nil
}
