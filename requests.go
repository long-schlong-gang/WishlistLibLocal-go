package wishlistlib

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
)

const DEFAULT_PORT = 2512

type WishClient struct {
	BaseURL string
	Port    int
	Token
}

// Returns the default context
func DefaultWishClient(baseURL string) *WishClient {
	return &WishClient{
		BaseURL: baseURL,
		Port:    DEFAULT_PORT,
		Token:   Token{},
	}
}

func (c *WishClient) executeRequest(method, path string, params map[string]string, reqBodyData []byte, addAuth bool) ([]byte, error) {
	// Put body data into a reader
	reqBody := bytes.NewReader(reqBodyData)
	req, err := http.NewRequest(method, fmt.Sprint(c.BaseURL, ":", c.Port, path), reqBody)
	if err != nil {
		return nil, err
	}

	// Add access-token in header
	if addAuth {
		req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", c.Token.Token))
	}

	// Add query parameters if present
	if params != nil {
		q := req.URL.Query()
		for key, value := range params {
			q.Add(key, value)
		}
		req.URL.RawQuery = q.Encode()
	}

	// Set request data type if present
	if reqBodyData != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	// Create HTTP client
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Read response body
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// Check for error responses
	err = handleNonOkErrors(resp.StatusCode, string(respBody))
	if err != nil {
		return nil, err
	}

	return respBody, nil
}

func handleNonOkErrors(code int, status string) error {
	switch code {
	case 200:
		fallthrough
	case 204:
		return nil
	case 400:
		return BadRequestError(status)
	case 401:
		return InvalidCredentialsError(status)
	case 403:
		return ForbiddenError(status)
	case 404:
		return NotFoundError(status)
	case 409:
		return EmailExistsError(status)
	case 500:
		return InternalServerError(status)
	}
	return UnknownHttpError(status)
}
