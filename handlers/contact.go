package handlers

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"net/smtp"
	"os"
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
		log.Fatal(err)
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
	log.Println("Mail sent successfully")
	http.Redirect(w, r, "/", 302)
}
