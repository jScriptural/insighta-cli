package cmd


import (
	"insighta/internal/config"
	"insighta/util"
	"fmt"
)


func Whoami(){
	cred,err := config.GetCredential();
	if err != nil {
		fmt.Println("No user found")
		return;
	}

	username := util.GetUsername(cred.AccessToken);

	if username == "" {
		fmt.Println("No user found")
		return;
	}
	role := util.GetRole(cred.AccessToken)


	fmt.Println("─────── INSIGHTA SESSION ───────")
	fmt.Printf("👤 Username: @%s\n", username)
	fmt.Printf("🛡️  Role:     %s\n", role)
	fmt.Println("────────────────────────────────")



	//fmt.Printf("📧 Email:    %s\n", email)
	//fmt.Printf("\tUsername: @%v\n\tRole: %v\n",username,role);

}
