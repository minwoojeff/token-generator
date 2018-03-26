package token

import (
	"encoding/json"
	"net/http"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"github.com/token-generator/helpers"
	"github.com/token-generator/models"
)

const PRIVATE_KEY = "/tmp/private_key.rsa"

func Token(jti string, data map[string]string) *jwt.Token {
	token := jwt.New(jwt.SigningMethodRS256)
	claims := make(jwt.MapClaims)

	now := time.Now().Unix()
	expiry := time.Now().Add(time.Hour * 1).Unix()

	claims["grant_type"] = "password"
	claims["username"] = data["username"]
	claims["iat"] = now    // issusedAt
	claims["exp"] = expiry // Expiry, One Hour
	claims["scope"] = []string{}
	claims["client_id"] = data["client_id"]
	claims["aud"] = []string{"password"}
	claims["jti"] = jti

	token.Claims = claims
	return token
}

func NewToken(token *jwt.Token, privateKey []byte) (string, error) {
	signKey, err := jwt.ParseRSAPrivateKeyFromPEM(privateKey)
	if err != nil {
		return "", err
	}

	// JWT Creation
	jwt, err := token.SignedString(signKey)
	if err != nil {
		return "", err
	}

	return jwt, nil
}

func PasswordGrant(w http.ResponseWriter, data map[string]string) {
	// UUID4
	identifier, err := uuid.NewRandom()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	// Set token claims
	token := Token(identifier.String(), data)

	r, err := helpers.Open(PRIVATE_KEY)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	privateKey, err := helpers.Read(r)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	jj, err := NewToken(token, privateKey)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	// Repsonse Token
	exp := token.Claims.(jwt.MapClaims)["exp"]
	iat := token.Claims.(jwt.MapClaims)["iat"]

	refreshToken, _ := uuid.NewRandom()
	resp := &models.Token{
		AccessToken:  jj,
		TokenType:    "password",
		RefreshToken: refreshToken.String(),
		ExpiresIn:    exp.(int64) - iat.(int64),
		Scope:        []string{},
		JTI:          identifier.String(),
	}

	j, err := json.Marshal(resp)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(j))
	return
}
