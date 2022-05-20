package wishlistlib

import (
	"encoding/json"
)

type Token struct {
	Token     string `json:"token"`
	ExpiresAt string `json:"expires_at"`
}

type AuthReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Gets an access token from the server to be provided for all secured endpoints.
func (wc *WishClient) Authenticate(email, password string) error {

	// Serialise user data
	reqBody, err := json.Marshal(
		AuthReq{
			Email:    email,
			Password: password,
		},
	)
	if err != nil {
		return err
	}

	// Get Token
	resBody, err := wc.executeRequest("POST", "/token", nil, reqBody, false)
	if err != nil {
		return nil
	}
	var tok Token
	err = json.Unmarshal(resBody, &tok)
	if err != nil {
		return InvalidTokenError(resBody)
	}

	wc.Token = tok

	return nil
}
