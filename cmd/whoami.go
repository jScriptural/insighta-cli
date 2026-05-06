package cmd

import (
	"fmt"
	"insighta/internal/config"
	"insighta/util"
)

func Whoami() {
	cred, err := config.GetCredential()
	if err != nil {
		fmt.Println("No user found")
		return
	}

	username := util.GetUsername(cred.AccessToken)

	if username == "" {
		fmt.Println("No user found")
		return
	}
	role := util.GetRole(cred.AccessToken)

	fmt.Println("\n\n─────── INSIGHTA SESSION ───────")
	fmt.Printf("👤 Username: @%s\n", username)
	fmt.Printf("🛡️  Role:     %s\n", role)
	fmt.Println("────────────────────────────────")

	//fmt.Printf("📧 Email:    %s\n", email)
	//fmt.Printf("\tUsername: @%v\n\tRole: %v\n",username,role);

}
