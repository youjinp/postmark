package postmark

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

var (
	postmarkURL = `https://api.postmarkapp.com`
)

type HttpClientAPI interface {
	Do(req *http.Request) (*http.Response, error)
}

type ClientAPI interface {
	SendEmail(email *Email) (*EmailResponse, error)
	SendEmailBatch(emails *[]Email) (*[]EmailResponse, error)
	SendEmailWithTemplate(email *EmailWithTemplate) (*EmailResponse, error)
	SendBatchEmailWithTemplate(emails *[]EmailWithTemplate) (*[]EmailResponse, error)
}

// Client provides a connection to the Postmark API
type Client struct {
	// HTTPClient
	HTTPClient HttpClientAPI
	// Server Token: Used for requests that require server level privileges. This token can be found on the Credentials tab under your Postmark server.
	ServerToken string
	// AccountToken: Used for requests that require account level privileges. This token is only accessible by the account owner, and can be found on the Account tab of your Postmark account.
	AccountToken string
	// BaseURL is the root API endpoint
	BaseURL string
}

const (
	serverToken  = "server"
	accountToken = "account"
)

// an object to hold variable parameters to perform request.
type parameters struct {
	// Method is HTTP method type.
	Method string
	// Path is postfix for URI.
	Path string
	// Payload for the request.
	Payload interface{}
	// TokenType defines which token to use
	TokenType string
}

// NewClient builds a new Client pointer using the provided tokens, a default HTTPClient, and a default API base URL
// Accepts `Server Token`, and `Account Token` as arguments
// http://developer.postmarkapp.com/developer-api-overview.html#authentication
func NewClient(httpClient HttpClientAPI, serverToken string, accountToken string) *Client {
	return &Client{
		HTTPClient:   httpClient,
		ServerToken:  serverToken,
		AccountToken: accountToken,
		BaseURL:      postmarkURL,
	}
}

func (client *Client) doRequest(opts parameters, dst interface{}) error {
	url := fmt.Sprintf("%s/%s", client.BaseURL, opts.Path)

	req, err := http.NewRequest(opts.Method, url, nil)
	if err != nil {
		return fmt.Errorf("Failed to create request: %v", err)
	}

	if opts.Payload != nil {
		payloadData, err := json.Marshal(opts.Payload)
		if err != nil {
			return fmt.Errorf("Failed to marshal payload: %v", err)
		}
		req.Body = ioutil.NopCloser(bytes.NewBuffer(payloadData))
	}

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")

	switch opts.TokenType {
	case accountToken:
		req.Header.Add("X-Postmark-Account-Token", client.AccountToken)

	default:
		req.Header.Add("X-Postmark-Server-Token", client.ServerToken)
	}

	res, err := client.HTTPClient.Do(req)
	if err != nil {
		return fmt.Errorf("Failed to perform request: %v", err)
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("Failed to read body: %v", err)
	}

	err = json.Unmarshal(body, dst)
	if err != nil {
		return fmt.Errorf("Failed to unmarshal body: %v", err)
	}

	return nil
}
