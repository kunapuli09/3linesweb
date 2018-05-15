package handlers

import (
        "testing"
        "fmt"
        "log"
        "net/smtp"
        "bytes"
)

func TestPostEmailLocalSmtp(t *testing.T) {
        // Connect to the remote SMTP server.
        c, err := smtp.Dial("localhost:25")
        if err != nil {
                log.Fatal(err)
        }
        defer c.Close()
        // Set the sender and recipient.
        c.Mail("krishna.kunapuli@threelines.us")
        c.Rcpt("krishna.kunapuli@threelines.us")
        // Send the email body.
        wc, err := c.Data()
        if err != nil {
                log.Fatal(err)
        }
        defer wc.Close()
        buf := bytes.NewBufferString("This is the email body.")
        if _, err = buf.WriteTo(wc); err != nil {
                log.Fatal(err)
        }
}