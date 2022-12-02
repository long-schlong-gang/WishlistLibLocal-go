package wishlistlib

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
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
	// Authenticate
	u, err := wc.GetUserByEmail(email)
	if err != nil {
		return NotAuthenticatedError("Invalid email/password provided")
	}
	if comparePassword(wc.Users[u.ID].PasswordHash, password) != nil {
		return NotAuthenticatedError("Invalid email/password provided")
	}

	// Create JWT and store it
	expiry := time.Now().Unix() + TOKEN_EXPIRY
	tok := generateToken(email, expiry)
	signd, err := tok.SignedString([]byte(os.Getenv(ENV_SIGNING_KEY)))
	if err != nil {
		return err
	}
	wc.token = Token{
		Token:     signd,
		ExpiresAt: time.Unix(expiry, 0).Local().String(),
	}

	return nil
}

func (wc *WishClient) GetToken() Token {
	return wc.token
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

func hashPassword(password string) (string, error) {
	pw := []byte(password)
	result, err := bcrypt.GenerateFromPassword(pw, bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(result), nil
}

func comparePassword(hashPassword string, password string) error {
	pw := []byte(password)
	hw := []byte(hashPassword)
	err := bcrypt.CompareHashAndPassword(hw, pw)
	return err
}
