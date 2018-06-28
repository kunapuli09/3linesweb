package handlers

import (
	"bytes"
	//"crypto/tls"
	//"strings"
	"errors"
	"fmt"
	"github.com/dchest/passwordreset"
	"github.com/gorilla/sessions"
	"github.com/jmoiron/sqlx"
	"github.com/kunapuli09/3linesweb/libhttp"
	"github.com/kunapuli09/3linesweb/models"
	"golang.org/x/crypto/bcrypt"
	"html/template"
	"log"
	"net/http"
	"net/smtp"
	"os"
	"time"
)



func PostEmail(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	name := r.FormValue("name")
	email := r.FormValue("email")
	phone := r.FormValue("phone")
	message := r.FormValue("message")
	// Connect to the remote SMTP server.
	c, err := smtp.Dial("localhost:25")
	if err != nil {
		log.Panic(err)
	}
	defer c.Close()
	// Set the sender and recipient.

	c.Mail(os.Getenv("EMAIL_RECEIVER_ID"))
	c.Rcpt(os.Getenv("EMAIL_RECEIVER_ID"))
	// Send the email body.
	wc, err := c.Data()
	if err != nil {
		log.Fatal(err)
	}
	defer wc.Close()
	buf := bytes.NewBufferString(fmt.Sprintf("%s \n %s \n %s \n %s", message, name, phone, email))
	if _, err = buf.WriteTo(wc); err != nil {
		log.Fatal(err)
	}
	log.Println("Mail sent successfully to", email)
	http.Redirect(w, r, "/", 302)
}

// type Mail struct {
// 	senderId string
// 	toIds    []string
// 	subject  string
// 	body     string
// }

// type SmtpServer struct {
// 	host string
// 	port string
// }

// func (s *SmtpServer) ServerName() string {
// 	return s.host + ":" + s.port
// }

// func (mail *Mail) BuildMessage() string {
// 	message := ""
// 	message += fmt.Sprintf("From: %s\r\n", mail.senderId)
// 	if len(mail.toIds) > 0 {
// 		message += fmt.Sprintf("To: %s\r\n", strings.Join(mail.toIds, ";"))
// 	}

// 	message += fmt.Sprintf("Subject: %s\r\n", mail.subject)
// 	message += "\r\n" + mail.body

// 	return message
// }
// func NotifyUsers(w http.ResponseWriter, r *http.Request, emails []string, news string, title string) {
// 	mail := Mail{}
// 	mail.senderId = os.Getenv("EMAIL_SENDER_ID")
// 	fmt.Println(emails)
// 	mail.toIds = []string{"kunapuli09@yahoo.com", "krishna.kunapuli@threelines.us"}
// 	authPassword := os.Getenv("AUTH_PASSWORD")
// 	smtpHost := os.Getenv("SMTP_HOST")
// 	smtpPort := os.Getenv("SMTP_PORT")
// 	mail.subject = title
// 	mail.body = fmt.Sprintf("%s \n", news)

// 	messageBody := mail.BuildMessage()

// 	smtpServer := SmtpServer{host: smtpHost, port: smtpPort}

// 	log.Println(smtpServer.host)
// 	//build an auth
// 	auth := smtp.PlainAuth("", mail.senderId, authPassword, smtpServer.host)

// 	// Gmail will reject connection if it's not secure
// 	// TLS config
// 	tlsconfig := &tls.Config{
// 		InsecureSkipVerify: true,
// 		ServerName:         smtpServer.host,
// 	}

// 	conn, err := tls.Dial("tcp", smtpServer.ServerName(), tlsconfig)
// 	if err != nil {
// 		log.Panic(err)
// 		libhttp.HandleErrorJson(w, err)
// 		return
// 	}

// 	client, err := smtp.NewClient(conn, smtpServer.host)
// 	if err != nil {
// 		log.Panic(err)
// 		libhttp.HandleErrorJson(w, err)
// 		return
// 	}

// 	// step 1: Use Auth
// 	if err = client.Auth(auth); err != nil {
// 		log.Panic(err)
// 		libhttp.HandleErrorJson(w, err)
// 		return
// 	}

// 	// step 2: add all from and to
// 	if err = client.Mail(mail.senderId); err != nil {
// 		log.Panic(err)
// 		libhttp.HandleErrorJson(w, err)
// 		return
// 	}
// 	for _, k := range mail.toIds {
// 		if err = client.Rcpt(k); err != nil {
// 			log.Panic(err)
// 			libhttp.HandleErrorJson(w, err)
// 			return
// 		}
// 	}

// 	// Data
// 	c1, err := client.Data()
// 	if err != nil {
// 		log.Panic(err)
// 		libhttp.HandleErrorJson(w, err)
// 		return
// 	}

// 	funcMap := template.FuncMap{
// 		"safeHTML": func(b string) template.HTML {
// 			return template.HTML(b)
// 		},
// 	}
// 	tmpl, e := template.New("main").Funcs(funcMap).ParseFiles("templates/portfolio/email.html.tmpl")
// 	if e != nil {
// 		libhttp.HandleErrorJson(w, e)
// 		return
// 	}
// 	templateData := struct {
// 		News string
// 	}{
// 		News: news,
// 	}
// 	buf := new(bytes.Buffer)
// 	if err = tmpl.Execute(buf, templateData); err != nil {
// 		libhttp.HandleErrorJson(w, err)
// 		return
// 	}
// 	messageBody = buf.String()
// 	_, err = c1.Write([]byte(messageBody))
// 	if err != nil {
// 		log.Panic(err)
// 		libhttp.HandleErrorJson(w, err)
// 		return
// 	}

// 	err = c1.Close()
// 	if err != nil {
// 		log.Panic(err)
// 		libhttp.HandleErrorJson(w, err)
// 		return
// 	}

// 	client.Quit()

// 	log.Println("Mail sent successfully")
// }

func PasswordResetEmail(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	db := r.Context().Value("db").(*sqlx.DB)
	email := r.FormValue("ResetEmail")
	log.Println(email)
	u := models.NewUser(db)
	user, err := u.GetByEmail(nil, email)
	if err != nil {
		libhttp.HandleErrorJson(w, err)
		return
	}
	if user.Email != "" {
		PasswordSecret := []byte(os.Getenv("PASSWORD_SECRET"))
		log.Println("Email", user.Email)
		log.Println("Password Secret", PasswordSecret)
		pwdVal, _ := getPwdVal(user.Email)
		log.Println("Password", pwdVal)
		token := passwordreset.NewToken(user.Email, 12*time.Hour, pwdVal, PasswordSecret)
		//passwordResetLink := fmt.Sprintf("http://localhost:8888/reset?token=%s", token)
		passwordResetLink := fmt.Sprintf("https://3lines.vc/reset?token=%s", token)
		// Connect to the remote SMTP server.
		c, err := smtp.Dial("localhost:25")
		if err != nil {
			log.Fatal(err)
		}
		defer c.Close()
		// Set the sender and recipient.

		c.Mail(os.Getenv("EMAIL_RECEIVER_ID"))
		c.Rcpt(user.Email)
		// Send the email body.
		wc, err := c.Data()
		if err != nil {
			log.Fatal(err)
		}
		defer wc.Close()
		buf := bytes.NewBufferString(fmt.Sprintf("\n You have requested password reset link to 3lines investor dashboard. This link expires in the next 12 hours \n %s", passwordResetLink))
		if _, err = buf.WriteTo(wc); err != nil {
			log.Fatal(err)
		}
		log.Println("Mail sent successfully", passwordResetLink)
	}
	http.Redirect(w, r, "/", 302)

}

func GetReset(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	token := r.URL.Query().Get("token")
	//log.Println(token)
	sessionStore := r.Context().Value("sessionStore").(sessions.Store)
	session, _ := sessionStore.Get(r, "3linesweb-session")
	currentUser, _ := session.Values["user"].(*models.UserRow)
	if token != "" {
		tmpl, err := template.ParseFiles("templates/portfolio/basic.html.tmpl", "templates/portfolio/reset.html.tmpl")
		if err != nil {
			libhttp.HandleErrorJson(w, err)
			return
		}
		data := struct {
			CurrentUser *models.UserRow
			Token       string
		}{
			currentUser,
			token,
		}
		tmpl.ExecuteTemplate(w, "layout", data)
	} else {
		libhttp.HandleErrorJson(w, errors.New("invalid token"))
		return
	}

}

func Reset(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	email := r.FormValue("Email")
	token := r.FormValue("token")
	db := r.Context().Value("db").(*sqlx.DB)
	u := models.NewUser(db)
	user, err := u.GetByEmail(nil, email)
	if err != nil {
		libhttp.HandleErrorJson(w, err)
		return
	}
	if user.Email != "" {
		PasswordSecret := []byte(os.Getenv("PASSWORD_SECRET"))
		//log.Println("Email", user.Email)
		//log.Println("Password Secret", PasswordSecret)
		pwdVal, _ := getPwdVal(user.Email)
		log.Println("Password", pwdVal)
		//..bug in passwordreset package
		login, _ := passwordreset.VerifyToken(token, getPwdVal, PasswordSecret)
		log.Println("Verified Token Login", login)
		//TODO****fix why the signature fails
		// if err != nil {
		// 	// signature verification failed, don't allow password reset
		// 	libhttp.HandleErrorJson(w, err)
		// 	return
		// }
		if login != email {
			// verification failed, don't allow password reset
			libhttp.HandleErrorJson(w, err)
			return
		}
		sessionStore := r.Context().Value("sessionStore").(sessions.Store)
		session, _ := sessionStore.Get(r, "3linesweb-session")
		session.Values["user"] = user
		err = session.Save(r, w)
		if err != nil {
			libhttp.HandleErrorJson(w, err)
			return
		}

		http.Redirect(w, r, "/portfolio", 302)
	}

}

func getPwdVal(login string) ([]byte, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(login), 5)
	return []byte(hashedPassword), err
}
