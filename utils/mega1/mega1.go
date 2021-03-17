package mega1

import (
	"auth/config"
	"auth/models"
	dt "auth/utils/datastore"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"time"

	"cloud.google.com/go/datastore"
	"github.com/labstack/gommon/log"
	"github.com/sendgrid/rest"
)

const (
	REGISTER = "register"
	LOGOUT   = "logout"
)

func GenerateMega1Code(client *datastore.Client, ctx context.Context) string {
	start := time.Now()
	rand.Seed(time.Now().UnixNano())
	mega1Code := fmt.Sprintf("MEGA1%010d", rand.Intn(9999999999))
	exist := true
	for exist == true {
		fmt.Println("Value of exist")
		filter := []dt.DatastoreFilter{
			{
				FilterString: "mega1_code =",
				Value:        mega1Code,
			}}
		docs := new([]models.AuthUser)
		err := dt.Get(client, ctx, "users", docs, filter, 1)
		if err != nil {
			log.Error(err)
		}
		if len(*docs) == 0 {
			exist = false
		} else {
			mega1Code = fmt.Sprintf("MEGA1%010d", rand.Intn(9999999999))
		}
	}
	fmt.Println("Value of mega1_code", mega1Code)
	fmt.Println("GenerateMega1Codetime: ", time.Since(start).Milliseconds())
	return mega1Code
}

func GetUserInfo(token string, white_list []string, enable bool) (*models.Mega1User, error) {

	Headers := make(map[string]string)
	Headers["Authorization"] = "Bearer " + token
	Headers["token"] = token
	var Body []byte
	appEnv := os.Getenv("APP_ENV")
	if appEnv == "localhost" || appEnv == "development" {
		Body = []byte(fmt.Sprintf(`{"token": "%s"}`, token))
	}

	request := rest.Request{
		Method:  rest.Post,
		BaseURL: config.C.Mega1Integrate.URL + "/api/yeah1/get-user",
		Headers: Headers,
		Body:    Body,
	}

	response, err := rest.Send(request)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	var res struct {
		Code    int               `json:"code"`
		Message string            `json:"message"`
		User    *models.Mega1User `json:"data"`
	}
	err = json.Unmarshal([]byte(response.Body), &res)

	if err != nil {
		log.Error(err)
		return nil, err
	}

	if res.Code != 0 {
		return nil, errors.New("Get user mega1 failed")
	}

	return res.User, nil

}
