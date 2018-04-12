package handlers

import (
	"github.com/kunapuli09/3linesweb/libhttp"
	"crypto/tls"
    "fmt"
    "log"
    "net/smtp"
    "strings"
    "net/http"
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
    mail.senderId = "krishna.3lines@gmail.com"
    mail.toIds = []string{"krishna.kunapuli@threelines.us"}
    mail.subject = "New Contact is trying to reach you"
    mail.body = fmt.Sprintf("%s \n %s \n %s \n %s", message, name, phone, email)

    messageBody := mail.BuildMessage()

    smtpServer := SmtpServer{host: "smtp.gmail.com", port: "465"}

    log.Println(smtpServer.host)
    //build an auth
    auth := smtp.PlainAuth("", mail.senderId, "meyt2hn32hmy9sx6es", smtpServer.host)

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


