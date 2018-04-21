package handlers

import (
	"crypto/tls"
	"fmt"
	"github.com/kunapuli09/3linesweb/libhttp"
	"log"
	"net/http"
	"net/smtp"
	"os"
	"strings"
)

type Mail struct {
	senderId string
	toIds    []string
	subject  string
	body     string
}

type SmtpServer struct {
	host string
	port string
}

func (s *SmtpServer) ServerName() string {
	return s.host + ":" + s.port
}

func (mail *Mail) BuildMessage() string {
	message := ""
	message += fmt.Sprintf("From: %s\r\n", mail.senderId)
	if len(mail.toIds) > 0 {
		message += fmt.Sprintf("To: %s\r\n", strings.Join(mail.toIds, ";"))
	}

	message += fmt.Sprintf("Subject: %s\r\n", mail.subject)
	message += "\r\n" + mail.body

	return message
}

func PostEmail(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	name := r.FormValue("name")
	email := r.FormValue("email")
	phone := r.FormValue("phone")
	message := r.FormValue("message")
	mail := Mail{}
	mail.senderId = os.Getenv("EMAIL_SENDER_ID")
	mail.toIds = strings.Split(os.Getenv("EMAIL_RECEIVER_ID"), ",")
	authPassword := os.Getenv("AUTH_PASSWORD")
	smtpHost := os.Getenv("SMTP_HOST")
	smtpPort := os.Getenv("SMTP_PORT")
	mail.subject = "New Contact is trying to reach you"
	mail.body = fmt.Sprintf("%s \n %s \n %s \n %s", message, name, phone, email)

	messageBody := mail.BuildMessage()

	smtpServer := SmtpServer{host: smtpHost, port: smtpPort}

	log.Println(smtpServer.host)
	//build an auth
	auth := smtp.PlainAuth("", mail.senderId, authPassword, smtpServer.host)

	// Gmail will reject connection if it's not secure
	// TLS config
	tlsconfig := &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         smtpServer.host,
	}

	conn, err := tls.Dial("tcp", smtpServer.ServerName(), tlsconfig)
	if err != nil {
		log.Panic(err)
		libhttp.HandleErrorJson(w, err)
		return
	}

	client, err := smtp.NewClient(conn, smtpServer.host)
	if err != nil {
		log.Panic(err)
		libhttp.HandleErrorJson(w, err)
		return
	}

	// step 1: Use Auth
	if err = client.Auth(auth); err != nil {
		log.Panic(err)
		libhttp.HandleErrorJson(w, err)
		return
	}

	// step 2: add all from and to
	if err = client.Mail(mail.senderId); err != nil {
		log.Panic(err)
		libhttp.HandleErrorJson(w, err)
		return
	}
	for _, k := range mail.toIds {
		if err = client.Rcpt(k); err != nil {
			log.Panic(err)
			libhttp.HandleErrorJson(w, err)
			return
		}
	}

	// Data
	c1, err := client.Data()
	if err != nil {
		log.Panic(err)
		libhttp.HandleErrorJson(w, err)
		return
	}

	_, err = c1.Write([]byte(messageBody))
	if err != nil {
		log.Panic(err)
		libhttp.HandleErrorJson(w, err)
		return
	}

	err = c1.Close()
	if err != nil {
		log.Panic(err)
		libhttp.HandleErrorJson(w, err)
		return
	}

	client.Quit()

	log.Println("Mail sent successfully")
	http.Redirect(w, r, "/", 302)
}
