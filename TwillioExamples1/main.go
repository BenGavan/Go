package main

import (
	"github.com/twilio/twilio-go"
	openapi "github.com/twilio/twilio-go/rest/api/v2010"
)

func main() {
	err := sendText()
	if err != nil {
		panic(err)
	}
}

func sendText() error {
	accountSid := ""
	authToken := ""
	phoneNumber := ""

	client := twilio.NewRestClient(accountSid, authToken)


}
