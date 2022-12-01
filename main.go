package main

import (
	"log"

	"github.com/tsawlen/emailService/common/sub"
	"github.com/tsawlen/emailService/controller"
)

func main() {

	sub.InitWebSocket()

	go controller.HandleIncomingEmails()

	log.Println("Email Service is online!")

	for {

	}

}
