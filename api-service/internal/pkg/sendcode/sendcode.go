package sendcode

import (
	"bytes"
	"fmt"
	"html/template"
	"net/smtp"
)

func SendCode(email string, code string) {
	// sender data
	from := "touristanbookingsystem@gmail.com"
	password := "wjct gbxs flag pfol"
  
	// Receiver email address
	to := []string{
	  email,
	}
  
	// smtp server configuration.
	smtpHost := "smtp.gmail.com"
	smtpPort := "587"
  
	// Authentication.
	auth := smtp.PlainAuth("", from, password, smtpHost)
  
	t, _ := template.ParseFiles("template.html")
  
	var body bytes.Buffer
  
	mimeHeaders := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	body.Write([]byte(fmt.Sprintf("Subject: Your verification code \n%s\n\n", mimeHeaders)))
  
	t.Execute(&body, struct {
	  Passwd string
	}{
  
	  Passwd: code,
	})
  
	// Sending email.
	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, body.Bytes())
	if err != nil {
	  fmt.Println(err)
	  return
	}
	fmt.Println("Email sended to:", email)

	return
  }
  