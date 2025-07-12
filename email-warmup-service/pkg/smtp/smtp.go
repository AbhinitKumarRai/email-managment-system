package smtp

import (
	"crypto/tls"
	"log"
	"net"
	"net/smtp"
	"os"
)

func CreateSmtpClient() (net.Conn, error) {
	server := "mailpit:1025"
	conn, err := net.Dial("tcp", server)
	if err != nil {
		log.Printf("SMTP connect error: %v", err)
		return nil, err
	}

	log.Println("SMTP connection established")
	return conn, nil
}

func CreateGmailSmtpClient() (*smtp.Client, error) {
	host := "smtp.gmail.com"
	port := "587"
	username := os.Getenv("GMAIL_SMTP_USERNAME")
	password := os.Getenv("GMAIL_SMTP_APP_PASSWORD")

	auth := smtp.PlainAuth("", username, password, host)
	conn, err := smtp.Dial(host + ":" + port)
	if err != nil {
		return nil, err
	}

	tlsconfig := &tls.Config{ServerName: host}
	if err = conn.StartTLS(tlsconfig); err != nil {
		return nil, err
	}

	if err = conn.Auth(auth); err != nil {
		return nil, err
	}

	return conn, nil
}
