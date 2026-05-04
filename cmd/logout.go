package cmd

import (
	"insighta/internal/config"
	"os"
	"insighta/util"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"encoding/json"
	"bytes"
)

type PostData struct {
	UserID string `json:"user_id"`
}

func Logout(){
	backendURL := "http://localhost:3030/auth/logout"
	creds, err := config.GetCredential();
	if err != nil {
		fmt.Println("No session found")
		return
	}
	token, _, _ := new(jwt.Parser).ParseUnverified(creds.AccessToken, jwt.MapClaims{})
	claims, _ := token.Claims.(jwt.MapClaims)
	userID, _ := claims["sub"].(string)

	// Notify Backend (Optional but recommended for the "Gold Standard")
	// callBackendLogout(userID, creds.AccessToken)
	fmt.Println("Userid:",userID)
	postData,_ := json.Marshal(PostData{UserID: userID});

	res,err := http.Post(
		backendURL,
		"application/json",
		bytes.NewBuffer(postData),
	)
	if err != nil || res.StatusCode != http.StatusNoContent{
		fmt.Println("Failed Clear local session!!!")
		return;
	}
	defer res.Body.Close();


	configPath, _ := util.GetCredPath()
	err = os.Remove(configPath)
	if err != nil {
		fmt.Printf("Failed to clear local session: %v\n", err)
		return
	}
	fmt.Printf("Successfully logged out. See you next time, @%v\n",claims["username"])

}
