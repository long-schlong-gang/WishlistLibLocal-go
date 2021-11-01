package wishlistlib

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"
)

const DEFAULT_BASE_URL = "https://www.pearcenet.ch:2512"

type Context struct {
	BaseUrl string
	Token
	authUser User
}

// Returns the default context
func DefaultContext() *Context {
	return &Context{
		BaseUrl:  DEFAULT_BASE_URL,
		Token:    Token{},
		authUser: User{},
	}
}

func (ctx *Context) parseObjectFromServer(path, method string, obj interface{}, params map[string]string, isAuth bool) error {
	body, err := ctx.getResponseBody(path, method, params, isAuth)
	if err != nil {
		return err
	}

	err = json.Unmarshal(body, obj)
	if err != nil {
		return err
	}

	return nil
}

func (ctx *Context) sendObjectToServer(path, method string, obj interface{}, isAuth bool) error {

	// Marshal Object to JSON
	jsonBytes, err := json.Marshal(obj)
	if err != nil {
		return err
	}

	bodyReader := strings.NewReader(string(jsonBytes))
	req, err := http.NewRequest(method, ctx.BaseUrl+path, bodyReader)
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "application/json")

	if isAuth {
		q := req.URL.Query()
		q.Add("token", ctx.Token.AccessToken)
		req.URL.RawQuery = q.Encode()
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return handleNonOkErrors(resp.StatusCode, resp.Status)
	}

	defer resp.Body.Close()

	return nil
}

func (ctx *Context) getResponseBody(path, method string, params map[string]string, isAuth bool) ([]byte, error) {
	req, err := http.NewRequest(method, ctx.BaseUrl+path, nil)
	if err != nil {
		return nil, err
	}

	// Add parameters
	if params != nil {
		q := req.URL.Query()

		if isAuth {
			q.Add("token", ctx.Token.AccessToken)
		}

		for k, v := range params {
			q.Add(k, v)
		}
		req.URL.RawQuery = q.Encode()
	}

	// Execute request
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, handleNonOkErrors(resp.StatusCode, resp.Status)
	}

	// Read Response
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func (ctx *Context) simpleRequest(path, method string, isAuth bool) error {
	req, err := http.NewRequest(method, ctx.BaseUrl+path, nil)
	if err != nil {
		return err
	}

	if isAuth {
		q := req.URL.Query()
		q.Add("token", ctx.Token.AccessToken)
		req.URL.RawQuery = q.Encode()
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return handleNonOkErrors(resp.StatusCode, resp.Status)
	}

	return err
}

func handleNonOkErrors(code int, status string) error {
	switch code {
	case 400:
		return BadRequestError(status)
	case 409:
		return EmailExistsError(status)
	case 401:
		return InvalidCredentialsError(status)
	}
	return nil
}
