package controller

import (
	"bytes"
	"encoding/json"
	"html/template"
	"log"
	"net/smtp"
	"os"

	"github.com/tsawlen/emailService/common/dataStructures"
	"github.com/tsawlen/emailService/common/sub"
	"github.com/tsawlen/emailService/connector"
)

func HandleIncomingEmails() {
	for {
		select {
		case messageIn := <-sub.EmailChannel:
			var messageObj dataStructures.EmailMessage
			err := json.Unmarshal([]byte(messageIn.Message), &messageObj)
			if err != nil {
				log.Println(err)
			}
			switch messageObj.Type {
			case "register":
				user, err := connector.GetProfileById(messageObj.ToUser)
				if err != nil {
					log.Println(err)
				}
				SendRegistrationMail(user)
			}
		}
	}
}

func SendRegistrationMail(user *dataStructures.User) {
	from := os.Getenv("EMAIL_ADDRESS")
	password := os.Getenv("EMAIL_PASSWORD")

	toEmail := user.Email
	to := []string{toEmail}

	host := os.Getenv("EMAIL_HOST")
	port := os.Getenv("EMAIL_PORT")
	address := host + ":" + port
	log.Println("Preparing to send email")
	body := renderer(user)

	auth := smtp.PlainAuth("", from, password, host)

	err := smtp.SendMail(address, auth, from, to, body)
	if err != nil {
		panic(err)
	}
	log.Println("Email send!")
}

// Helper

func renderer(user *dataStructures.User) []byte {
	t, _ := template.ParseFiles("./templates/Signup.html")
	var body bytes.Buffer

	subject := "Subject: Your registration at Finder\n"
	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	message := []byte(subject + mime)

	t.Execute(&body, struct {
		FirstName string
		Name      string
	}{
		FirstName: user.First_name,
		Name:      user.Name,
	})

	return append(message, body.Bytes()...)
}
