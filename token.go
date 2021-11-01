package wishlistlib

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type Token struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   uint64 `json:"expires_in"`
	TokenType   string `json:"token_type"`
}

// Sets the user to use for this context
func (ctx *Context) SetAuthenticatedUser(user User) {
	ctx.authUser = user
}

// Gets an access token from the server to be provided for all secured endpoints.
func (ctx *Context) authenticate() error {
	// Send Request
	data := url.Values{}
	data.Add("grant_type", "client_credentials")
	req, err := http.NewRequest("POST", ctx.BaseUrl+"/token", strings.NewReader(data.Encode()))
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	req.SetBasicAuth(ctx.authUser.Email, ctx.authUser.password)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	// Read Response
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return handleNonOkErrors(resp.StatusCode, resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	// Parse Data
	err = json.Unmarshal(body, &ctx.Token)
	if err != nil {
		return err
	}

	return nil
}
