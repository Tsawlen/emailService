package sub

import (
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/gorilla/websocket"
	"github.com/tsawlen/emailService/common/dataStructures"
)

var EmailChannel = make(chan *dataStructures.IncomingSubMessage)

func readClientMessages(webSocket *websocket.Conn, incomingMessages chan dataStructures.IncomingSubMessage) {
	for {
		var messageRec dataStructures.IncomingSubMessage
		err := webSocket.ReadJSON(&messageRec)
		if err != nil {
			log.Println(err)
		}
		log.Println("Received request")
		EmailChannel <- &messageRec
	}

}

func InitWebSocket() {
	url := url.URL{Scheme: "ws", Host: os.Getenv("CLOUD_RELAY_HOST"), Path: "/subscribe"}
	header := http.Header{}
	header.Add("topic", "email")
	client, _, errCon := websocket.DefaultDialer.Dial(url.String(), header)
	if errCon != nil {
		log.Println(errCon)
	}
	log.Println("Subscribed to email topic!")
	channel := make(chan dataStructures.IncomingSubMessage)
	go readClientMessages(client, channel)
}
