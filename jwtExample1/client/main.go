package main

import (
	"fmt"
	jwt "github.com/dgrijalva/jwt-go"
	"time"
)

func main() {
	fmt.Println("Simple Client")

	tokenString, err := GenerateJWT()
	if err != nil {
		fmt.Print(tokenString)
		fmt.Printf("error: %v\n", err)
	}
	fmt.Printf("token string: %s\n", tokenString)
}

func GenerateJWT() (string, error) {
	token := jwt.New(jwt.SigningMethodES256)

	claims := token.Claims.(jwt.MapClaims)
	claims["authorized"] = true
	claims["user"] = "username here"
	claims["exp"] = time.Now().Add(time.Minute * 30).Unix()

	var mySigningKey = []byte("supersecretkeyphrase")

	fmt.Printf("token raw: %v\n", token.Claims)
	tokenString, err := token.SignedString(mySigningKey)
	fmt.Printf("token string: %v\n", tokenString)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}
