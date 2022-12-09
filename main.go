package main

import (
	"log"

	"github.com/tsawlen/emailService/common/initializer"
	"github.com/tsawlen/emailService/common/sub"
	"github.com/tsawlen/emailService/controller"
)

func main() {
	done := make(chan bool)
	go initializer.LoadEnvVariables(done)
	<-done
	go controller.HandleIncomingEmails()
	sub.InitWebSocket()

	log.Println("Email Service is online!")

	for {

	}

}
