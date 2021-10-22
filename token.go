package wishlistlib

import (
	"encoding/json"
	"io"
	"net/http"
	"time"
)

type Token struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   uint64 `json:"expires_in"`
	TokenType   string `json:"token_type"`
}

// Gets an access token from the server to be provided for all secured endpoints.
// `autoRenew` sets whether the token should automatically be renewed once it expires
func (ctx *Context) Authenticate(email, password string, autoRenew bool) error {
	// Send Request
	req, err := http.NewRequest("POST", ctx.BaseUrl+"/token", nil)
	if err != nil {
		return err
	}
	req.SetBasicAuth(email, password)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	// Read Response
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	// Parse Data
	err = json.Unmarshal(body, &ctx.Token)
	if err != nil {
		return err
	}

	if autoRenew {
		ctx.email = email
		ctx.password = password
		go ctx.renewToken()
	}

	return nil
}

func (ctx *Context) renewToken() {

	// Wait half of the token expiry time before renewing it
	time.Sleep(time.Duration(ctx.Token.ExpiresIn * uint64(time.Second)))
	err := ctx.Authenticate(ctx.email, ctx.password, false)
	if err == nil {
		ctx.renewToken()
	}
}
