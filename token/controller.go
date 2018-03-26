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

func PasswordGrant(w http.ResponseWriter, data map[string]string) {
	// Validate required fields for Password Grant Token
	username, userOk := data["username"]
	if !userOk {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Missing required params"))
		return
	}

	// UUID4
	identifier, err := uuid.NewRandom()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	// Set token claims
	token := jwt.New(jwt.SigningMethodRS256)
	claims := token.Claims.(jwt.MapClaims)

	now := time.Now().Unix()
	expiry := time.Now().Add(time.Hour * 1).Unix()

	claims["grant_type"] = "password"
	claims["username"] = username
	claims["iat"] = now    // issusedAt
	claims["exp"] = expiry // Expiry, One Hour
	claims["scope"] = []string{}
	claims["client_id"] = data["client_id"]
	claims["aud"] = []string{"password"}
	claims["jti"] = identifier.String()

	privateKey, err := helpers.ReadKey(PRIVATE_KEY)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	jwt, err := helpers.NewToken(token, privateKey)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	// Repsonse Token
	refreshToken, _ := uuid.NewRandom()
	resp := &models.Token{
		AccessToken:  jwt,
		TokenType:    "password",
		RefreshToken: refreshToken.String(),
		ExpiresIn:    expiry - now,
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
