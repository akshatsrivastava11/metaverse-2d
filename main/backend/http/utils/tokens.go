package utils

import (
	"fmt"
	"log"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var secretKey = []byte("secret-key")

func GetToken(username string) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": username,
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
	})
	fmt.Println(token)
	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(tokenString)
	return tokenString

}
func VerifyToken(tokenString string) (string, error) {
	data, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	if !data.Valid {
		return "", fmt.Errorf("Not a valid token")
	}
	// claims, err := jwt.Claims.(jwt.MapClaims)
	claims, ok := data.Claims.(jwt.MapClaims)
	if !ok {
		return "", fmt.Errorf("Could not parse claims")
	}
	username, _ := claims["username"].(string)
	return username, nil
}
