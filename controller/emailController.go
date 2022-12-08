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
				SendSignUpMail(user)
			case "registerCode":
				user, err := connector.GetProfileById(messageObj.ToUser)
				if err != nil {
					log.Println(err)
				}
				SendRegisterCodeMail(user)
			case "payment":
				user, err := connector.GetProfileById(messageObj.ToUser)
				if err != nil {
					log.Println(err)
				}

				var invoice *dataStructures.InvoiceEmailMessage
				invoiceErr := json.Unmarshal([]byte(messageIn.Message), &messageObj)
				if invoiceErr != nil {
					log.Println(err)
				}

				SendRecievedPaymentMail(user, invoice)
			}
		}
	}
}

func SendRegisterCodeMail(user *dataStructures.User) {
	from := os.Getenv("EMAIL_ADDRESS")
	password := os.Getenv("EMAIL_PASSWORD")

	toEmail := user.Email
	to := []string{toEmail}

	host := os.Getenv("EMAIL_HOST")
	port := os.Getenv("EMAIL_PORT")
	address := host + ":" + port
	log.Println("Preparing to send email")
	body := registerRenderer(user)

	auth := smtp.PlainAuth("", from, password, host)

	err := smtp.SendMail(address, auth, from, to, body)
	if err != nil {
		panic(err)
	}
	log.Println("Email send!")
}

func SendSignUpMail(user *dataStructures.User) {
	from := os.Getenv("EMAIL_ADDRESS")
	password := os.Getenv("EMAIL_PASSWORD")

	toEmail := user.Email
	to := []string{toEmail}

	host := os.Getenv("EMAIL_HOST")
	port := os.Getenv("EMAIL_PORT")
	address := host + ":" + port
	log.Println("Preparing to send email")
	body := signupRenderer(user)

	auth := smtp.PlainAuth("", from, password, host)

	err := smtp.SendMail(address, auth, from, to, body)
	if err != nil {
		panic(err)
	}
	log.Println("Email send!")
}

func SendRecievedPaymentMail(user *dataStructures.User, invoice *dataStructures.InvoiceEmailMessage) {
	from := os.Getenv("EMAIL_ADDRESS")
	password := os.Getenv("EMAIL_PASSWORD")

	toEmail := user.Email
	to := []string{toEmail}

	host := os.Getenv("EMAIL_HOST")
	port := os.Getenv("EMAIL_PORT")
	address := host + ":" + port
	log.Println("Preparing to send email")
	body := paymentRecievedRenderer(user, invoice)

	auth := smtp.PlainAuth("", from, password, host)

	err := smtp.SendMail(address, auth, from, to, body)
	if err != nil {
		panic(err)
	}
	log.Println("Email send!")
}

// Helper

func signupRenderer(user *dataStructures.User) []byte {
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

func registerRenderer(user *dataStructures.User) []byte {
	t, _ := template.ParseFiles("./templates/Registercode.html")
	var body bytes.Buffer

	subject := "Subject: Finder Email authentification\n"
	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	message := []byte(subject + mime)

	t.Execute(&body, struct {
		Username string
		Authcode string
	}{
		Username: user.Username,
		Authcode: user.Name,
	})

	return append(message, body.Bytes()...)
}

func paymentRecievedRenderer(user *dataStructures.User, invoice *dataStructures.InvoiceEmailMessage) []byte {
	t, _ := template.ParseFiles("./templates/RecievedPayment.html")
	var body bytes.Buffer

	subject := "Subject: Finder Email Payment confirmation\n"
	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	message := []byte(subject + mime)

	t.Execute(&body, struct {
		FirstName         string
		Name              string
		Amount            float64
		RecieverFirstName string
		ReceiverName      string
		Hours             int
		Service           string
	}{
		FirstName:         user.First_name,
		Name:              user.Name,
		Amount:            invoice.Amount,
		RecieverFirstName: invoice.RecieverFirstName,
		ReceiverName:      invoice.ReceiverName,
		Hours:             invoice.Hours,
		Service:           invoice.Service,
	})

	return append(message, body.Bytes()...)
}
