package main

import (
	"crypto/subtle"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/token-generator/models"
	"github.com/urfave/negroni"
)

// Global vars
var clientID = os.Getenv("CLIENT_ID")
var clientSecret = os.Getenv("CLIENT_SECRET")

const PRIVATE_KEY = "/tmp/private_key.rsa"

type Info struct {
	Version     string `json:"version"`
	Description string `json:"description"`
}

func About(w http.ResponseWriter, r *http.Request) {
	about := &Info{
		Version:     "0.0.1",
		Description: "Token generating microservice",
	}

	j, err := json.Marshal(about)
	if err != nil {
		log.Fatalf("ERROR: Internal error, %s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(j)
	return
}

func Health(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	return
}

func PasswordGrantToken(w http.ResponseWriter, data map[string]string) {
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

	expiry := time.Now().Add(time.Hour * 1).Unix()

	claims["grant_type"] = "password"
	claims["username"] = username
	claims["iat"] = time.Now().Unix() // issusedAt
	claims["exp"] = expiry            // Expiry, One Hour
	claims["scope"] = []string{}
	claims["client_id"] = clientID
	claims["aud"] = []string{"password"}
	claims["jti"] = identifier.String()

	dir, _ := os.Getwd()
	privateKey, err := ioutil.ReadFile(dir + PRIVATE_KEY)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	signKey, err := jwt.ParseRSAPrivateKeyFromPEM(privateKey)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	// JWT Creation
	jwt, err := token.SignedString(signKey)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	refreshToken, _ := uuid.NewRandom()
	// Repsonse Token
	resp := &models.Token{
		AccessToken:  jwt,
		TokenType:    "password",
		RefreshToken: refreshToken.String(),
		ExpiresIn:    expiry,
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

func TokenHandler(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	// Parse Body
	var data map[string]string
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil && err == io.EOF {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Missing required params"))
		return
	} else if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	// Auth
	status, err := Authenticate(r, data)
	if err != nil {
		w.WriteHeader(status)
		w.Write([]byte(err.Error()))
		return
	}

	switch grantType := data["grant_type"]; grantType {
	case "password":
		PasswordGrantToken(w, data)
	case "authorization_code":
		w.WriteHeader(http.StatusNotImplemented)
		return
	case "client_credentials":
		w.WriteHeader(http.StatusNotImplemented)
		return
	default:
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Unsupported grant type"))
		return
	}
}

func Authenticate(r *http.Request, data map[string]string) (int, error) {
	user, pass, ok := r.BasicAuth()

	// client_id and client_secret may be passed as part of payload
	if !ok {
		user = data["client_id"]
		pass = data["client_secret"]
	}

	if subtle.ConstantTimeCompare([]byte(user), []byte(clientID)) != 1 || subtle.ConstantTimeCompare([]byte(pass), []byte(clientSecret)) != 1 {
		return http.StatusUnauthorized, errors.New("Unauthorized")
	}

	return 0, nil
}

func main() {
	// Main Router
	main := mux.NewRouter().StrictSlash(true)
	// Authentication Router
	authRouter := mux.NewRouter().PathPrefix("").Subrouter().StrictSlash(true)
	n := negroni.New(negroni.HandlerFunc(TokenHandler), negroni.Wrap(authRouter))

	log.Printf("CREDS: %s, %s", clientID, clientSecret)

	// GET /, No Auth
	main.HandleFunc("/", About).Methods("GET")
	// GET /health, No Auth
	main.HandleFunc("/health", Health).Methods("GET")

	// POST /oauth/token, Authenticated
	main.PathPrefix("/oauth/token").Handler(n).Methods("POST")

	log.Fatal(http.ListenAndServe(":3000", main))
}
