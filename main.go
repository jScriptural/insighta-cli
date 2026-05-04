package main

import (
	//"fmt"
	"insighta/cmd"
	"log"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Insufficient argument")
	}

	args := os.Args[1:]
	c := args[0]
	switch {
	case c == "login":
		clientID := os.Getenv("GITHUB_CLIENT_ID")
		backendURL := "http://localhost:3030/auth/github"
		cmd.Login(clientID, backendURL)
	case c == "whoami":
		cmd.Whoami()
	case c == "logout":
		cmd.Logout()
	case c == "profiles" && len(args) >= 2:
		cmd.Profiles(args[1:])
	case c == "export" && len(args) >= 3:
		cmd.Export(args[1:])
	default:
		log.Fatalf("Unknown command: %v", c)
	}
}
