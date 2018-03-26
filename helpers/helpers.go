package helpers

import (
	"io/ioutil"
	"os"

	jwt "github.com/dgrijalva/jwt-go"
)

func ReadKey(path string) ([]byte, error) {
	dir, _ := os.Getwd()
	privateKey, err := ioutil.ReadFile(dir + path)
	if err != nil {
		return nil, err
	}
	return privateKey, nil
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
