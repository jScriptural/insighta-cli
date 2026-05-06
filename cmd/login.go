package cmd

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/pkg/browser"
	"insighta/internal/auth"
	"insighta/internal/config"
	"insighta/util"
	"log"
	"net/http"
	"time"
)

func Login(clientID string, backendURL string) {
	verifier, challenge := auth.GeneratePKCE()

	url := fmt.Sprintf("https://github.com/login/oauth/authorize?client_id=%s&code_challenge=%s&code_challenge_method=S256&scope=user", clientID, challenge)

	fmt.Println("Opening browser for Github login...")
	err := browser.OpenURL(url)
	if err != nil {
		log.Printf("browser: %v", err)
		fmt.Printf("Please open this URL: %s\n", url)
	}

	codeChan := make(chan string)
	srv := &http.Server{Addr: ":8080"}

	http.HandleFunc("/callback", func(w http.ResponseWriter, r *http.Request) {
		code := r.URL.Query().Get("code")
		if code != "" {
			fmt.Fprintf(w, "Authentication successful: Return to your terminal.")
		} else {
			fmt.Fprintf(w, "Authentication Failed: Return to your terminal.")
		}
		codeChan <- code
	})

	go srv.ListenAndServe()
	var code string
	select {
	case <-time.After(1 * time.Hour):
	case code = <-codeChan:
	}
	srv.Close()

	if code == "" {
		log.Fatal("Authentication Failed or Request timeout")
	}

	fmt.Println("Exchanging code with backend...")

	cred, err := exchangeCodeWithBackend(code, verifier, backendURL)
	if err != nil {
		log.Fatal(err)
	}

	if err := config.SaveCredential(cred); err != nil {
		log.Fatal(err)
	}

	username := util.GetUsername(cred.AccessToken)
	if username == "" {
		log.Println("Login Failed")
		return
	}
	fmt.Printf("\nLogin successful: @%v\n", username)
}

func exchangeCodeWithBackend(code, verifier, backendURL string) (config.Credential, error) {
	client := &http.Client{Timeout: 20 * time.Second}

	ctx, cancelFunc := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancelFunc()

	data := struct {
		Code     string `json:"code"`
		Verifier string `json:"verifier"`
	}{Code: code, Verifier: verifier}

	buffer := &bytes.Buffer{}
	err := json.NewEncoder(buffer).Encode(data)
	if err != nil {
		log.Fatal(err)
	}
	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		backendURL,
		buffer,
	)
	if err != nil {
		log.Fatal(err)
	}

	res, err := client.Do(req)
	if err != nil {
		return config.Credential{}, err
	}

	if res.StatusCode != http.StatusOK {
		log.Printf("Error: %v - %v", res.StatusCode, http.StatusText(res.StatusCode))
		return config.Credential{}, errors.New("Non-200 response from server")
	}
	defer res.Body.Close()

	cred := config.Credential{}
	err = json.NewDecoder(res.Body).Decode(&cred)
	if err != nil {
		log.Fatal(err)
	}

	return cred, nil
}
