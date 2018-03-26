package token

import (
	"testing"

	jwt "github.com/dgrijalva/jwt-go"
)

func TestToken_Valid(t *testing.T) {
	data := map[string]string{
		"client_id": "username",
		"username":  "user",
	}

	jj := Token("aaa", data)

	id := jj.Claims.(jwt.MapClaims)["client_id"]
	if id.(string) != "username" {
		t.Errorf("Expected %s as username, got %s", "username", id.(string))
	}
}

func TestToken_Invalid(t *testing.T) {
	data := map[string]string{}

	jj := Token("aaa", data)

	id := jj.Claims.(jwt.MapClaims)["client_id"]
	if id.(string) != "" {
		t.Errorf("Expected empty as username, got %s", id.(string))
	}
}

func TestToken_InvalidJTI(t *testing.T) {
	data := map[string]string{
		"username": "user",
	}

	jj := Token("", data)

	id := jj.Claims.(jwt.MapClaims)["jti"]
	if id.(string) != "" {
		t.Errorf("Expected empty as username, got %s", id.(string))
	}
}

func TestNewToken_EmptyBytes(t *testing.T) {
	_, err := NewToken(nil, []byte{})
	if err == nil {
		t.Errorf("Expecting an error for empty key")
	}
}
