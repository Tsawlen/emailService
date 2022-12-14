package connector

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/tsawlen/emailService/common/dataStructures"
)

func GetPaymentByUserId(id int) (*[]dataStructures.Invoice, error) {
	var invoices []dataStructures.Invoice
	query := "http://0.0.0.0:8080/invoice/" + strconv.Itoa(id)
	restClient := http.Client{
		Timeout: time.Second * 40,
	}

	request, errReq := http.NewRequest(http.MethodGet, query, nil)
	if errReq != nil {
		log.Println("Could not query invoices!")
		return nil, errReq
	}

	request.Header.Set("Authorization", "Bearer "+os.Getenv("JWT"))

	result, errRes := restClient.Do(request)
	if errRes != nil {
		log.Println("Could not query invoices!")
		return nil, errRes
	}
	if result.Body != nil {
		defer result.Body.Close()
	}
	body, errRead := ioutil.ReadAll(result.Body)
	if errRead != nil {
		log.Println("Could not read body")
		return nil, errRead
	}
	if errJson := json.Unmarshal(body, &invoices); errJson != nil {
		log.Println(errJson)
		return nil, errJson
	}
	return &invoices, nil
}
