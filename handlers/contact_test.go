package handlers

import (
        "gopkg.in/gomail.v2"
        "testing"
        "fmt"
        "log"
        "crypto/tls"
        "net/smtp"
        "crypto/x509"
)

func TestPostEmailThroughLocalSmtp(t *testing.T) {
        m := gomail.NewMessage()
        m.SetHeader("From", "krishna.3lines@gmail.com")
        m.SetHeader("To", "kunapuli09@gmail.com")
        m.SetHeader("Subject", "New Contact is trying to reach you")
        m.SetBody("text/plain", fmt.Sprintf("%s \n %s \n %s \n %s", "hey", "hey", "hey", "hey@he.com"))
        d := gomail.Dialer{Host: "localhost", Port: 25}
        if err := d.DialAndSend(m); err != nil {
                log.Panic(err)
        }
}

func TestPostEmail(t *testing.T) {
        authPassword := "meyt2hn32hmy9sx6es"
        smtpHost := "smtp.gmail.com"
        messageBody := ""
        messageBody += fmt.Sprintf("From: %s\r\n", "krishna.3lines@gmail.com")
        messageBody += fmt.Sprintf("To: %s\r\n", "krishna.kunapuli@threelines.us")
        messageBody += fmt.Sprintf("Subject: %s\r\n", "New Contact is trying to reach you")
        messageBody += "\r\n" + fmt.Sprintf("%s \n %s \n %s \n %s", "test", "krishna", "7202802571", "test@test.com")
        //build an auth
        auth := smtp.PlainAuth("", "krishna.3lines@gmail.com", authPassword, smtpHost)

        // Gmail will reject connection if it's not secure
        // TLS config
        roots := x509.NewCertPool()
        ok := roots.AppendCertsFromPEM([]byte(rootPEM))
        if !ok {
                panic("failed to parse root certificate")
        }
        tlsconfig := &tls.Config{
                InsecureSkipVerify: true,
                ServerName:         "smtp.gmail.com:587",
        }

        conn, err := tls.Dial("tcp", "smtp.gmail.com:587", tlsconfig)
        if err != nil {
                log.Panic(err)
                return
        }

        client, err := smtp.NewClient(conn, smtpHost)
        if err != nil {
                log.Panic(err)
                return
        }

        // step 1: Use Auth
        if err = client.Auth(auth); err != nil {
                log.Panic(err)
                return
        }

        // step 2: add all from and to
        if err = client.Mail("krishna.3lines@gmail.com"); err != nil {
                log.Panic(err)
                return
        }
        if err = client.Rcpt("krishna.kunapuli@threelines.us"); err != nil {
                log.Panic(err)
                return
        }

        // Data
        c1, err := client.Data()
        if err != nil {
                log.Panic(err)
                return
        }

        _, err = c1.Write([]byte(messageBody))
        if err != nil {
                log.Panic(err)
                return
        }

        err = c1.Close()
        if err != nil {
                log.Panic(err)
                return
        }

        client.Quit()
}