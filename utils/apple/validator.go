package apple

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	ValidationURL string = "https://appleid.apple.com/auth/token"
	ContentType   string = "application/x-www-form-urlencoded"
	UserAgent     string = "go-signin-with-apple"
	AcceptHeader  string = "application/json"
)

type ValidationClient interface {
	VerifyAppToken(ctx context.Context, reqBody AppValidationTokenRequest, result interface{}) error
}

type Client struct {
	validationURL string
	client        *http.Client
}

func New() *Client {
	client := &Client{
		validationURL: ValidationURL,
		client: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
	return client
}

func (c *Client) VerifyAppToken(ctx context.Context, reqBody AppValidationTokenRequest, result interface{}) error {
	data := url.Values{}
	data.Set("client_id", reqBody.ClientID)
	data.Set("client_secret", reqBody.ClientSecret)
	data.Set("code", reqBody.Code)
	data.Set("grant_type", "authorization_code")

	return doRequest(c.client, &result, c.validationURL, data)
}

func doRequest(client *http.Client, result interface{}, url string, data url.Values) error {
	req, err := http.NewRequest("POST", url, strings.NewReader(data.Encode()))
	if err != nil {
		return err
	}

	req.Header.Add("content-type", ContentType)
	req.Header.Add("accept", AcceptHeader)
	req.Header.Add("user-agent", UserAgent)

	res, err := client.Do(req)
	if err != nil {
		return err
	}

	defer res.Body.Close()

	return json.NewDecoder(res.Body).Decode(result)
}
