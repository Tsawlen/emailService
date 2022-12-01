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

func GetProfileById(id int) (*dataStructures.User, error) {
	var user dataStructures.User
	query := "http://0.0.0.0:8080/profile/" + strconv.Itoa(id)
	restClient := http.Client{
		Timeout: time.Second * 40,
	}

	request, errReq := http.NewRequest(http.MethodGet, query, nil)
	if errReq != nil {
		log.Println("Could not query user!")
		return nil, errReq
	}

	request.Header.Set("Authorization", "Bearer eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2Njk5MjAxNjAsInN1YiI6MSwidXNlciI6MX0.CDGheYGiW8sUakVZOcv2X2XCALvVyHk3SU2J7eHua06hzmJRRB7bHcPFHaQ8-fWb7W6h5UXPd2woSevNVkreswzXU24FPyN33Nakwve_y4jTp948vzcj2rfTAQj9TiRfOtR1Wk9___cLCpLu1Q92Jq4e460s2KX4BqdgcSG3ePGN3Xu7duz6zhnGBi7_nnFQXQuJ2CxMXrM9yFlVo6OJnadtIl5XKIyBHO04YSUWt9pFLMgr_riuYYW4qkwFMc6uSH3eBkgpMrJy7VgnblJueKI6PQx-QY5zYFY8Kc8VIp9wAyZzrPHKfhMZCdGH-8m4UxOhzxV2epoTM97o0ngMZg")

	result, errRes := restClient.Do(request)
	if errRes != nil {
		log.Println("Could not query user!")
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
	if errJson := json.Unmarshal(body, &user); errJson != nil {
		log.Println(errJson)
		return nil, errJson
	}
	return &user, nil
}
