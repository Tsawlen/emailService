package connector

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
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

	request.Header.Set("Authorization", "Bearer eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2NzAzNDk5NTgsInN1YiI6MSwidXNlciI6MX0.EFSj4-Aj95t9O8kWzgtDz2odJ5OiA2zrvHuJsiJkw0_U5w8IqIF6z3z_mLeR2uKVqfHl8XtELs0BGs3JuaANRvoSi1nviwf58oKuF7AwyY2DXT0cdtGVmUiMzi0CWg9BumjRsyL0M42oJV25sGpzwgWctk34yvNz0ScS0hBzvrhx2rSVHW3rJtRDevMp_UG9kZRDMPTKX9ax2jv_43FCFRtdcLdPO-CYJMQHhgMAZKO5nwAVqtOOtWXohSDrUPnSPgqOkbB8mOls6uckHoEgLFAUIJsxTckyh5Xt6_enyZ68W3ggsHDGPu0irvRfOrTYrB4fbaTMOQFZXqi4IJC_Xg")

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
