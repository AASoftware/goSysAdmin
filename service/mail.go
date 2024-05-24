package service

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/smtp"
	"os"
)

type Mail struct {
	MailEnabled bool   `json:"mailenabled"`
	SmtpServer  string `json:"smtp"`
	Port        string `json:"port"`
	Subject     string `json:"subject"`
	SendTo      string `json:"to"`
	SendFrom    string `json:"from"`
}

func NewMail() (*Mail, error) {
	return &Mail{}, nil
}

func (m *Mail) LoadConfig(file string) {
	f, err := os.ReadFile(file)
	if err != nil {
		panic("Config can't be loaded")
	}

	err = json.Unmarshal(f, &m)
	if err != nil {
		panic("Failed to unmarshal config.json")
	}
}

func (m *Mail) SendMailWithAttachment(message []byte, attachment []byte, subject string) error {
	tlsconfig := &tls.Config{
		InsecureSkipVerify: true, // Skip server certificate verification
	}

	client, err := smtp.Dial(m.SmtpServer + ":" + m.Port)
	if err != nil {
		return fmt.Errorf("failed to connect with smtp server: %s", err)
	}
	defer client.Close()

	if err = client.StartTLS(tlsconfig); err != nil {
		return fmt.Errorf("failed to start tls: %s", err)
	}

	if err = client.Mail(m.SendFrom); err != nil {
		return fmt.Errorf("failed to set from: %s", err)
	}

	if err = client.Rcpt(m.SendTo); err != nil {
		return fmt.Errorf("failed to set sendto: %s", err)
	}

	w, err := client.Data()
	if err != nil {
		return fmt.Errorf("failed to get data writer: %s", err)
	}

	// Write email header
	header := fmt.Sprintf("Subject: %s\r\n\r\n", subject)
	_, err = w.Write([]byte(header))
	if err != nil {
		return fmt.Errorf("failed to write email header: %s", err)
	}

	// Write message
	_, err = w.Write(message)
	if err != nil {
		return fmt.Errorf("failed to write message: %s", err)
	}

	// Write attachment
	_, err = w.Write([]byte("\r\n"))
	if err != nil {
		return fmt.Errorf("failed to write newline: %s", err)
	}

	_, err = w.Write(attachment)
	if err != nil {
		return fmt.Errorf("failed to write attachment: %s", err)
	}

	err = w.Close()
	if err != nil {
		return fmt.Errorf("failed to close data writer: %s", err)
	}

	fmt.Println("Email sent")
	return nil
}
