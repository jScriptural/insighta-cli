package util

import (
	"github.com/golang-jwt/jwt/v5"
	"os"
	"path/filepath"
	"fmt"
	"log"
)

func GetUsername(authCred string) string {
	if authCred == "" {
		return ""
	}
	token, _, err := new(jwt.Parser).ParseUnverified(authCred, jwt.MapClaims{})
	if err != nil {
		return "Guest"
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		if username, exists := claims["username"]; exists {
			return fmt.Sprintf("%v", username)
		}
	}

	return ""
}


func GetRole(authCred string) string {
	if authCred == "" {
		return ""
	}
	token, _, err := new(jwt.Parser).ParseUnverified(authCred, jwt.MapClaims{})
	if err != nil {
		return ""
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		if role, exists := claims["role"]; exists {
			return fmt.Sprintf("%v", role)
		}
	}

	return ""
	 
}

func GetHomeDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		log.Printf("No home dir: %v", err)
		return "", err
	}

	return home,nil
}

func GetCredPath() (string,error) {
	home,err := GetHomeDir();
	if err != nil {
		return "",err;
	}

	return filepath.Join(home,".insighta","credentials.json"),nil;
}
