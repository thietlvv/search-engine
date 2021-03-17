package smtp

import (
	"auth/config"
	"auth/constants"
	"auth/models"
	"auth/utils/jwt"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/sirupsen/logrus"
)

func SendEmail(user *models.AuthUser, stype int, app_type string) {
	var url string
	var subject string

	if app_type == constants.KOLS_ECOMMERCE {
		if stype == constants.VERIFY_EMAIL {
			url = config.C.SMTPService.KolsEcommerce.Email.Url
			subject = config.C.SMTPService.KolsEcommerce.Email.Subject
		} else {
			url = config.C.SMTPService.KolsEcommerce.Password.Url
			subject = config.C.SMTPService.KolsEcommerce.Password.Subject
		}

	}
	if stype == constants.VERIFY_EMAIL {
		token, _ := jwt.GenerateVerifyEmailToken(user)
		url = fmt.Sprintf(`%s=%s`, url, token)
	} else {
		token, _ := jwt.GenerateResetPassByEmailToken(user)
		url = fmt.Sprintf(`%s=%s`, url, token)
	}

	uri := config.C.SMTPService.Uri
	method := "POST"

	payload := map[string]interface{}{
		"email":   user.Email,
		"subject": subject,
		"content": map[string]string{
			"username": user.Name,
			"url":      url,
		},
		"type":   stype,
		"method": 1,
	}

	bytesRepresentation, err := json.Marshal(payload)
	if err != nil {
		logrus.Error(err)
	}

	client := &http.Client{}
	req, err := http.NewRequest(method, uri, bytes.NewBuffer(bytesRepresentation))

	if err != nil {
		logrus.Error(err)
	}
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		logrus.Error(err)
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		logrus.Error(err)
	}
	fmt.Println(string(body))
}
