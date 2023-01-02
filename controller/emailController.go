package controller

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"mime/multipart"

	"net/http"
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
				return
			}
			switch messageObj.Type {
			case "register":
				user, err := connector.GetProfileById(messageObj.ToUser)
				if err != nil {
					log.Println(err)
					return
				}
				SendSignUpMail(user)
			case "registerCode":
				user, err := connector.GetProfileById(messageObj.ToUser)
				if err != nil {
					log.Println(err)
					return
				}
				type codeStruct struct {
					Code string
				}
				var activationCode codeStruct
				json.Unmarshal([]byte(messageObj.Message), &activationCode)
				log.Println(activationCode.Code)
				SendRegisterCodeMail(user, activationCode.Code)
			case "payment":
				user, err := connector.GetProfileById(messageObj.ToUser)
				if err != nil {
					log.Println(err)
					return
				}

				var invoice *dataStructures.InvoiceEmailMessage
				invoiceErr := json.Unmarshal([]byte(messageIn.Message), &messageObj)
				if invoiceErr != nil {
					log.Println(err)
				}
				SendRecievedPaymentMail(user, invoice)
			case "invoice":
				user, err := connector.GetProfileById(messageObj.ToUser)
				if err != nil {
					log.Println(err)
					return
				}
				var invoice *dataStructures.Invoice
				invoiceErr := json.Unmarshal([]byte(messageObj.Message), &invoice)
				if invoiceErr != nil {
					log.Println(invoiceErr)
					return
				}
				SendInvoiceMail(user, invoice)
			}

		}
	}
}

func SendRegisterCodeMail(user *dataStructures.User, activationCode string) {
	from := os.Getenv("EMAIL_ADDRESS")
	password := os.Getenv("EMAIL_PASSWORD")

	toEmail := user.Email
	to := []string{toEmail}

	host := os.Getenv("EMAIL_HOST")
	port := os.Getenv("EMAIL_PORT")
	address := host + ":" + port
	log.Println("Preparing to send email")
	body := registerRenderer(user, activationCode)

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

func SendInvoiceMail(user *dataStructures.User, invoice *dataStructures.Invoice) {
	from := os.Getenv("EMAIL_ADDRESS")
	password := os.Getenv("EMAIL_PASSWORD")

	toEmail := user.Email
	to := []string{toEmail}

	host := os.Getenv("EMAIL_HOST")
	port := os.Getenv("EMAIL_PORT")
	address := host + ":" + port
	log.Println("Preparing to send email")
	body := invoiceRenderer(user, invoice)
	auth := smtp.PlainAuth("", from, password, host)
	err := smtp.SendMail(address, auth, from, to, body)
	if err != nil {
		panic(err)
	}
	log.Println("Email send!")
	return
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

func registerRenderer(user *dataStructures.User, activationCode string) []byte {
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
		Authcode: activationCode,
	})

	return append(message, body.Bytes()...)
}

func invoiceRenderer(user *dataStructures.User, invoice *dataStructures.Invoice) []byte {
	buf := bytes.NewBuffer(nil)
	writer := multipart.NewWriter(buf)
	boundary := writer.Boundary()
	buf.WriteString(fmt.Sprintf("Content-Type: multipart/mixed; boundary=%s \n", boundary))
	buf.WriteString("MIME-Version: 1.0\n")
	buf.WriteString(fmt.Sprintf("To: %s\n", user.Email))
	buf.WriteString(fmt.Sprintf("Subject: Finder Invoice \n\n"))
	buf.WriteString(fmt.Sprintf("--%s\n", boundary))

	t, _ := template.ParseFiles("./templates/Invoice.html")
	var body bytes.Buffer
	t.Execute(&body, struct {
		FirstName string
		Name      string
		Amount    float64
		Hours     int
		Service   string
	}{
		FirstName: user.First_name,
		Name:      user.Name,
		Amount:    invoice.Amount,
		Hours:     invoice.Hours,
		Service:   invoice.Service,
	})
	buf.WriteString("Content-Type: text/html; charset=\"UTF-8\"\n")
	buf.WriteString("MIME-Version: 1.0\n")
	buf.WriteString("Content-Transfer-Encoding: 7bit\n\n")
	buf.Write(body.Bytes())
	buf.WriteString(fmt.Sprintf("\n\n--%s\n", boundary))

	buf.WriteString(fmt.Sprintf("Content-Type: %s; name=%s \n", http.DetectContentType(invoice.InvoicePDF), "Invoice.pdf"))
	buf.WriteString("MIME-Version: 1.0\n")
	buf.WriteString("Content-Transfer-Encoding: base64\n")
	buf.WriteString(fmt.Sprintf("Content-Disposition: attachment; filename=Invoice.pdf\n\n"))
	b := make([]byte, base64.StdEncoding.EncodedLen(len(invoice.InvoicePDF)))
	base64.StdEncoding.Encode(b, invoice.InvoicePDF)
	buf.Write(b)
	buf.WriteString(fmt.Sprintf("\n--%s", boundary+"--"))

	return buf.Bytes()

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

// Generate pdf
