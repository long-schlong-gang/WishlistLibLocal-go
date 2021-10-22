package main

import (
	"encoding/json"
	"io"
	"net/http"
)

const DEFAULT_BASE_URL = "https://www.pearcenet.ch:2512"

type Context struct {
	BaseUrl string
	Token
	email    string
	password string
}

// Returns the default context
func DefaultContext() *Context {
	return &Context{
		BaseUrl:  DEFAULT_BASE_URL,
		Token:    Token{},
		email:    "",
		password: "",
	}
}

func (ctx *Context) parseObjectFromServer(path, method string, obj interface{}) error {
	body, err := ctx.getResponseBody(path, method)
	if err != nil {
		return err
	}

	err = json.Unmarshal([]byte(body), obj)
	if err != nil {
		return err
	}

	return nil
}

func (ctx *Context) getResponseBody(path, method string) ([]byte, error) {
	// Send Request
	req, err := http.NewRequest(method, ctx.BaseUrl+path, nil)
	if err != nil {
		return nil, err
	}
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