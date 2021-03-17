package sms

import (
	"auth/config"
	"bytes"
	"io/ioutil"
	"net/http"

	"github.com/labstack/gommon/log"
)

func SendOTP(data string) []byte {
	var jsonBody = []byte(data)
	req, err := http.NewRequest("POST", config.C.SMSServices.GAPIT.Url, bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth(config.C.SMSServices.GAPIT.Username, config.C.SMSServices.GAPIT.Password)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Error(err)
		return nil
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	return body
}
