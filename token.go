package wishlistlib

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt"
)

const (
	TOKEN_EXPIRY    = 30 * 24 * 60 * 60 // 30 Days
	ENV_SIGNING_KEY = "WISHLIST_TOK_SIGNING_KEY"
)

type Token struct {
	Token     string `json:"token"`
	ExpiresAt string `json:"expires_at"`
}

type TokenClaims struct {
	Email string `json:"email"`
	jwt.StandardClaims
}

type AuthReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Gets an access token from the server to be provided for all secured endpoints.
func (wc *WishClient) Authenticate(email, password string) error {
	password, err := hashPassword(password)
	if err != nil {
		return err
	}

	// Authenticate
	u, err := wc.GetUserByEmail(email)
	if err != nil {
		return err
	}
	if wc.Users[u.ID].PasswordHash != password {
		return NotAuthenticatedError("Invalid email/password provided")
	}

	// Create JWT and store it
	expiry := time.Now().Unix() + TOKEN_EXPIRY
	tok := generateToken(email, expiry)
	signd, err := tok.SignedString([]byte(os.Getenv(ENV_SIGNING_KEY)))
	if err != nil {
		return err
	}
	wc.Token = Token{
		Token:     signd,
		ExpiresAt: time.Unix(expiry, 0).Local().String(),
	}

	return nil
}

// Generates a JWT with the user's email and expiry time
func generateToken(email string, expiresAt int64) *jwt.Token {
	cl := TokenClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expiresAt,
		},
		Email: email,
	}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, cl)
}
