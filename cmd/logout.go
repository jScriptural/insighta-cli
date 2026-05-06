package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"insighta/internal/config"
	"insighta/util"
	"net/http"
	"os"
)

type PostData struct {
	UserID string `json:"user_id"`
}

func Logout() {
	backendURL := "http://localhost:3030/auth/logout"
	creds, err := config.GetCredential()
	if err != nil {
		fmt.Println("No session found")
		return
	}
	token, _, _ := new(jwt.Parser).ParseUnverified(creds.AccessToken, jwt.MapClaims{})
	claims, _ := token.Claims.(jwt.MapClaims)
	userID, _ := claims["sub"].(string)

	fmt.Println("Userid:", userID)
	postData, _ := json.Marshal(PostData{UserID: userID})

	res, err := http.Post(
		backendURL,
		"application/json",
		bytes.NewBuffer(postData),
	)
	if err != nil || res.StatusCode != http.StatusNoContent {
		fmt.Println("Failed Clear local session!!!")
		return
	}
	defer res.Body.Close()

	configPath, _ := util.GetCredPath()
	err = os.Remove(configPath)
	if err != nil {
		fmt.Printf("\nFailed to clear local session: %v\n", err)
		return
	}
	fmt.Printf("\nSuccessfully logged out. See you next time, @%v\n", claims["username"])

}
