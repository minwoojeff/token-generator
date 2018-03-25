package main

import (
	"crypto/subtle"
	"encoding/json"
	"io/ioutil"
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

func Token(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
	return
}

func TokenHandler(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	AuthHandler(w, r, Token)
}

func AuthHandler(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	user, pass, ok := r.BasicAuth()
	// client_id and client_secert may be passed as part of payload
	if !ok {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}

		var data map[string]string
		err = json.Unmarshal(body, &data)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}

		user = data["client_id"]
		pass = data["client_secret"]
	}

	if subtle.ConstantTimeCompare([]byte(user), []byte(clientID)) != 1 || subtle.ConstantTimeCompare([]byte(pass), []byte(clientSecret)) != 1 {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Unauthorized"))
	} else {
		log.Println("Authorized request.")
		next(w, r)
	}
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

	// POST /oauth/token, Authenticated
	main.PathPrefix("/oauth/token").Handler(n).Methods("POST")

	log.Fatal(http.ListenAndServe(":3000", main))
}
