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
	"html/template"
	"log"
	"net/http"
	"net/smtp"
	"os"
	"strings"
	"time"
	"github.com/haisum/recaptcha"
)

func RSVP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	name := r.FormValue("FullName")
	email := r.FormValue("Email")
	phone := r.FormValue("Phone")
	companyname := r.FormValue("CompanyName")
	//Connect to the remote SMTP server.
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
	msg := fmt.Sprintf("%s \n %s \n %s \n %s \n %s", "Future Of Work Webinar Registration", companyname, name, phone, email)
	buf := bytes.NewBufferString(msg)
	if _, err = buf.WriteTo(wc); err != nil {
		log.Fatal(err)
	}
	log.Println("Mail sent successfully to", msg)
	http.Redirect(w, r, "/", 302)
}

// func RSVP(w http.ResponseWriter, r *http.Request) {
// 	w.Header().Set("Content-Type", "text/html")
// 	name := r.FormValue("FullName")
// 	email := r.FormValue("Email")
// 	phone := r.FormValue("Phone")
// 	companyname := r.FormValue("CompanyName")
// 	log.Println("Mail sent successfully to", fmt.Sprintf("%s \n %s \n %s \n %s \n %s", "Future Of Work Webinar Registration", companyname, name, phone, email))
// 	http.Redirect(w, r, "/", 302)
// }

func PostEmail(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	re := recaptcha.R{
    	Secret: os.Getenv("CAPTCHA_SITE_SECRET"),
	}
	token := r.FormValue("rcres")
	log.Println("Verifying Captcha token", token)
	isValid := re.VerifyResponse(token)
	if !isValid {
    	log.Printf("Invalid Captcha! These errors ocurred: %v", re.LastError())
        libhttp.HandleErrorJson(w, errors.New("Invalid Captcha!"))
		return
    }
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

func PasswordResetEmail(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	db := r.Context().Value("db").(*sqlx.DB)
	email := r.FormValue("ResetEmail")
	log.Println(email)
	if email != "" {
		u := models.NewUser(db)
		user, err := u.GetByEmail(nil, email)
		if err != nil {
			libhttp.HandleErrorJson(w, err)
			return
		}
		PasswordSecret := []byte(os.Getenv("PASSWORD_SECRET"))
		if len(PasswordSecret) <= 0 {
			log.Println("PasswordSecret Environment Variable is missing")
		}
		pwdVal, err := getPwdVal(user.Password)
		token := passwordreset.NewToken(user.Email, 12*time.Hour, pwdVal, PasswordSecret)
		//passwordResetLink := fmt.Sprintf("http://localhost:8888/reset?token=%s", token)
		passwordResetLink := fmt.Sprintf("https://3lines.vc/reset?token=%s", token)
		log.Printf("PasswordResetLink %s", passwordResetLink)
		// Connect to the remote SMTP server.
		c, err := smtp.Dial("localhost:25")
		if err != nil {
			log.Fatal(err)
			err := errors.New("Sorry, mail server is not up.")
			libhttp.HandleErrorJson(w, err)
			return
		}
		defer c.Close()
		// Set the sender and recipient.
		log.Printf("A Password Reset Email was Requested By %s", user.Email)
		c.Mail(os.Getenv("EMAIL_RECEIVER_ID"))
		c.Rcpt(user.Email)
		// Send the email body.
		wc, err := c.Data()
		if err != nil {
			log.Fatal(err)
			err := errors.New("Sorry, issue with sending an email to your address")
			libhttp.HandleErrorJson(w, err)
			return
		}
		defer wc.Close()
		buf := bytes.NewBufferString(fmt.Sprintf("Subject: 3Lines Dashboard Password Reset!\r\n"+
			"\r\n"+"\n You have requested password reset link to 3lines investor dashboard. This link expires in the next 12 hours \n %s", passwordResetLink))
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

	if email != "" {
		u := models.NewUser(db)
		user, err := u.GetByEmail(nil, email)
		if err != nil {
			fmt.Println("User Doesn't exist.")
			e := errors.New("Invalid Reset Link. Reset Link is Not Generated for the Email Entered")
			libhttp.HandleErrorJson(w, e)
			return
		}
		PasswordSecret := []byte(os.Getenv("PASSWORD_SECRET"))
		log.Println("Verifying Token", token)
		login, err := passwordreset.VerifyToken(token, getPwdVal, PasswordSecret)
		log.Println("Verifying Token. Hashed Login", login)
		log.Println("Verifying Token. Email Reset By User", user.Email)
		if !(strings.EqualFold(strings.Trim(login, " "), strings.Trim(email, " "))) {
			// verification failed, don't allow password reset
			err := errors.New("Invalid Reset Link. Reset Link is Not Generated for the Email Entered")
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

		http.Redirect(w, r, "/entryaccess", 302)
	}

}

func getPwdVal(dbPassword string) ([]byte, error) {
	return []byte(dbPassword), nil
}
