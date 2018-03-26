package main

import (
	"bytes"
	"encoding/json"
	"testing"

	"net/http"
	"net/http/httptest"
)

func TestAbout(t *testing.T) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "http://example.com", nil)

	About(w, r)

	if w.Code != http.StatusOK {
		t.Error("Expecting / to return 200")
	}

	if w.Body == nil {
		t.Error("Expecting / to return JSON body")
	}
}

func TestHealth(t *testing.T) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "http://example.com/health", nil)

	Health(w, r)

	if w.Code != http.StatusOK {
		t.Error("Expecting /health to return 200")
	}
}

func TestAuthenticate_Basic(t *testing.T) {
	data := map[string]string{}

	r := httptest.NewRequest("GET", "http://example.com", nil)
	r.SetBasicAuth("default", "password") // ENV not set. Default empty

	code, _ := Authenticate(r, data)
	if code != 0 {
		t.Errorf("Expecting 0 for returend code, got %d", code)
	}
}

func TestAuthenticate_Body(t *testing.T) {
	data := map[string]string{
		"client_id":     "default",
		"client_secret": "password",
	}

	r := httptest.NewRequest("GET", "http://example.com", nil)

	code, _ := Authenticate(r, data)
	if code != 0 {
		t.Errorf("Expecting 0 for returend code, got %d", code)
	}
}

func TestAuthenticate_Fail(t *testing.T) {
	data := map[string]string{
		"client_id":     "a",
		"client_secret": "b",
	}

	r := httptest.NewRequest("GET", "http://example.com", nil)

	code, _ := Authenticate(r, data)
	if code != 401 {
		t.Errorf("Expecting 401 for returend code, got %d", code)
	}
}

func TestTokenHandler_NoBody(t *testing.T) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "http://example.com", nil)
	r.SetBasicAuth("default", "password")

	TokenHandler(w, r, nil)

	if w.Code != 400 {
		t.Errorf("Expecting 400 for no body, got %d", w.Code)
	}
}

func TestTokenHandler_Valid(t *testing.T) {
	w := httptest.NewRecorder()

	data := map[string]string{
		"grant_type": "password",
		"username":   "somerandomemail",
	}
	b, _ := json.Marshal(data)
	r := httptest.NewRequest("GET", "http://example.com", bytes.NewReader(b))
	r.SetBasicAuth("default", "password")

	TokenHandler(w, r, nil)

	if w.Code != http.StatusOK {
		t.Errorf("Expecting 501 for no body, got %d", w.Code)
	}
}

func TestTokenHandler_AuthCode(t *testing.T) {
	w := httptest.NewRecorder()

	data := map[string]string{
		"grant_type": "authorization_code",
		"username":   "somerandomemail",
	}
	b, _ := json.Marshal(data)
	r := httptest.NewRequest("GET", "http://example.com", bytes.NewReader(b))
	r.SetBasicAuth("default", "password")

	TokenHandler(w, r, nil)

	if w.Code != http.StatusNotImplemented {
		t.Errorf("Expecting 501 for no body, got %d", w.Code)
	}
}

func TestTokenHandler_ClientCreds(t *testing.T) {
	w := httptest.NewRecorder()

	data := map[string]string{
		"grant_type": "client_credentials",
		"username":   "somerandomemail",
	}
	b, _ := json.Marshal(data)
	r := httptest.NewRequest("GET", "http://example.com", bytes.NewReader(b))
	r.SetBasicAuth("default", "password")

	TokenHandler(w, r, nil)

	if w.Code != http.StatusNotImplemented {
		t.Errorf("Expecting 501 for no body, got %d", w.Code)
	}
}

func TestTokenHandler_Unsupported(t *testing.T) {
	w := httptest.NewRecorder()

	data := map[string]string{
		"grant_type": "asdfasdf",
		"username":   "somerandomemail",
	}
	b, _ := json.Marshal(data)
	r := httptest.NewRequest("GET", "http://example.com", bytes.NewReader(b))
	r.SetBasicAuth("default", "password")

	TokenHandler(w, r, nil)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expecting 400 for no body, got %d", w.Code)
	}
}

func TestTokenHandler_AuthFail(t *testing.T) {
	w := httptest.NewRecorder()

	data := map[string]string{
		"grant_type": "asdfasdf",
		"username":   "somerandomemail",
	}
	b, _ := json.Marshal(data)
	r := httptest.NewRequest("GET", "http://example.com", bytes.NewReader(b))
	r.SetBasicAuth("a", "b")

	TokenHandler(w, r, nil)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expecting 401 for no body, got %d", w.Code)
	}
}
