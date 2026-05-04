package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"insighta/util"
)

type Credential struct {
	Status string `json:"status"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type RefreshToken struct {
	Token string `json:"refresh_token"`
}

func SaveCredential(cred Credential) error {
	p, err := util.GetCredPath()
	if err != nil {
		return nil;
	}


	err  = os.MkdirAll(filepath.Dir(p),0o700);
	if err != nil {
		return err;
	}

	f, err := os.OpenFile(p,os.O_RDWR|os.O_CREATE|os.O_TRUNC,0o600)
	if err != nil {
		return err;
	}
	defer f.Close()

	encoder := json.NewEncoder(f)
	encoder.SetIndent("", " ")
	if err = encoder.Encode(cred); err != nil {
		return err;
	}

	return nil;
}


func GetCredential() (*Credential,error) {
	path,err:= util.GetCredPath();
	if err != nil {
		return nil,err;
	}

	f,err := os.Open(path);
	if err != nil {
		return nil,err;
	}
	defer f.Close()

	cred := Credential{}
	if err := json.NewDecoder(f).Decode(&cred); err != nil {
		return nil,err;
	}

	return &cred,nil;
}
