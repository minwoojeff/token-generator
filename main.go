package main

import (
	"crypto/subtle"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/urfave/negroni"
)

// Global vars
var clientID = os.Getenv("CLIENT_ID")
var clientSecret = os.Getenv("CLIENT_SECRET")

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

func PasswordGrantToken(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("password grant here~"))
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
		PasswordGrantToken(w, r)
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
	r.Close = false
	// client_id and client_secret may be passed as part of payload
	if !ok {
		user = data["client_id"]
		pass = data["client_secret"]
	}

	if subtle.ConstantTimeCompare([]byte(user), []byte(clientID)) != 1 || subtle.ConstantTimeCompare([]byte(pass), []byte(clientSecret)) != 1 {
		return http.StatusUnauthorized, errors.New("Unauthorized")
	}

	log.Println("Successfully authenticated")
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
