package test

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/gistapp/api/user"
	"github.com/gofiber/fiber/v2"
)

type Method string

const (
	GET    Method = "GET"
	POST   Method = "POST"
	PUT    Method = "PUT"
	PATCH  Method = "PATCH"
	DELETE Method = "DELETE"
)

type TestRequest struct {
	T        *testing.T
	App      *fiber.App
	Url      string
	Method   Method
	Headers  map[string]string
	Payload  io.Reader
	Err      error
	Response *http.Response
}

func Client(t *testing.T, app *fiber.App) *TestRequest {
	return &TestRequest{
		T:   t,
		App: app,
	}
}

func (tr *TestRequest) Get(url string) *TestRequest {
	tr.Url = url
	tr.Method = GET
	tr.Payload = nil
	return tr
}

func (tr *TestRequest) Post(url string) *TestRequest {
	tr.Url = url
	tr.Method = POST
	return tr
}

func (tr *TestRequest) Put(url string) *TestRequest {
	tr.Url = url
	tr.Method = PUT
	return tr
}

func (tr *TestRequest) Patch(url string) *TestRequest {
	tr.Url = url
	tr.Method = PATCH
	return tr
}

func (tr *TestRequest) Delete(url string) *TestRequest {
	tr.Url = url
	tr.Method = DELETE
	return tr
}

func (tr *TestRequest) WithHeaders(headers map[string]string) *TestRequest {

	if tr.Headers == nil {
		tr.Headers = make(map[string]string)
	}
	for key, value := range headers {
		tr.Headers[key] = value
	}
	return tr
}

func (tr *TestRequest) WithPayload(payload interface{}) *TestRequest {
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		tr.Err = err
		return tr
	}
	tr.Payload = strings.NewReader(string(jsonPayload))
	return tr
}

func (tr *TestRequest) WithJson() *TestRequest {
	if tr.Headers == nil {
		tr.Headers = make(map[string]string)
	}
	tr.Headers["Content-Type"] = "application/json"
	return tr
}

func (tr *TestRequest) LoginAs(user user.AuthIdentityAndUser) *TestRequest {
	auth_token, err := user.GetAccessToken()
	tr.Err = err
	tr.Headers["Authorization"] = "Bearer " + auth_token
	return tr
}

func (tr *TestRequest) Test(test func(*http.Response, *testing.T)) (*http.Response, error) {
	resp, err := send(tr)
	if err != nil {
		return nil, err
	}
	test(resp, tr.T)
	return resp, nil
}

func (tr *TestRequest) Send() *TestRequest {
	tr.Response, tr.Err = send(tr)
	return tr
}

func send(client *TestRequest) (*http.Response, error) {
	if client.Err != nil {
		return nil, client.Err
	}

	req, err := http.NewRequest(string(client.Method), client.Url, client.Payload)

	if err != nil {
		return nil, err
	}

	for key, value := range client.Headers {
		req.Header.Add(key, value)
	}

	resp, err := client.App.Test(req)

	if err != nil {
		return nil, err
	}

	return resp, nil
}

func JSONHttpResponse(resp *http.Response) (map[string]string, error) {
	respBody := make(map[string]string)
	if err := json.NewDecoder(resp.Body).Decode(&respBody); err != nil {
		return nil, err
	}
	return respBody, nil
}

func (tr *TestRequest) ExpectStatus(status int) *TestRequest {
	if tr.Response.StatusCode != status {
		body, err := JSONHttpResponse(tr.Response)
		tr.Err = err
		tr.T.Fatalf("Expected status code %d, got %d : body is %s", status, tr.Response.StatusCode, body)
	}
	return tr
}
