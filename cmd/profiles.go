package cmd

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"insighta/internal/config"
	"insighta/internal/models"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"
)


func Profiles(args []string) {
	//log.Printf("args: %#v\n", args)

	subcmd := args[0]
	//log.Println("subcmd: ", subcmd)
	switch subcmd {
	case "list":
		if len(args) < 2 {
			List(ParseArgs([]string{}).Encode())
			return
		}
		List(ParseArgs(args[1:]).Encode())
	case "search":
		if len(args) < 2 {
			Search([]string{})
			return
		}
		Search(args[1:])
	case "get":
		if len(args) < 2 {
			os.Stderr.WriteString("Insufficient args to command 'get'\nUSAGE: insighta profiles get <id>")
			return
		}
		Get(args[1]);
	case "create":
		if len(args) < 2 {
			os.Stderr.WriteString("Insufficient flag to command 'create'\nUSAGE: insighta profiles create --name <name>")
			return
		}
		Create(args[1:])
	default:
		log.Println("Unknown subcommand")
	}

}

func List(query string) {
	backendURL := "http://localhost:3030/api/profiles"
	log.Printf("query: %v\n", query)

	res, err := SendRequest(
		http.MethodGet,
		backendURL+"?"+query,
		nil,
	)

	WriteOutput(res,err,&models.Response{})
}

func Search(args []string) {
	backendURL := "http://localhost:3030/api/profiles/search"

	query := url.Values{};
	switch len(args) {
	case 0:
		query = ParseArgs(args);
	case 1:
		query.Set("q",args[0]);
	default:
		query = ParseArgs(args[1:])
		query.Set("q",args[0]);
	}

	res, err := SendRequest(
		http.MethodGet,
		backendURL+"?"+query.Encode(),
		nil,
	)

	WriteOutput(res,err,&models.Response{})
}


func Get(id string){
	backendURL := "http://localhost:3030/api/profiles/"
	res,err := SendRequest(
		http.MethodGet,
		backendURL+id,
		nil,
	)

	WriteOutput(res,err,&models.Response{})
}


func Create(args []string) {
	backendURL := "http://localhost:3030/api/profiles"
	fs := flag.NewFlagSet("create",flag.ExitOnError);
	name := fs.String("name","","Create profile for given name")
	_ = fs.Bool("help",false,"insighta profiles create --name <name>")

	fs.Parse(args);
	if *name == "" {
		os.Stderr.WriteString("Bad call to create\nUSAGE: insighta profiles create --name <name>")
		return
	}
	data := models.Name{Name: *name};
	payload := &bytes.Buffer{};
	err := json.NewEncoder(payload).Encode(data)
	if err != nil {
		os.Stderr.WriteString("Fail to encode payload")
		return
	}

	res,err := SendRequest(
		http.MethodPost,
		backendURL,
		payload,
	)

	//d := []models.Profile{}
	WriteOutput(res,err,&models.Response{});
}

func RefreshSession(t string) (*config.Credential, error) {
	backendURL := "http://localhost:3030/auth/refresh"
	tk := config.RefreshToken{Token: t}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	payload := &bytes.Buffer{}
	if err := json.NewEncoder(payload).Encode(tk); err != nil {
		log.Println(err)
		return nil, fmt.Errorf("RefreshSession: %w", err)
	}

	log.Printf("payload: %v", payload.String())
	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		backendURL,
		payload,
	)
	if err != nil {
		return nil, err
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	cred := &config.Credential{}
	if res.StatusCode != http.StatusOK {
		errRes := models.ErrResponse{}
		if err := json.NewDecoder(res.Body).Decode(&errRes); err != nil {
			log.Println(err)
			return nil, nil
		}
		return nil, fmt.Errorf("RefreshSession: %w", errors.New(errRes.Message))
	}

	if err := json.NewDecoder(res.Body).Decode(cred); err != nil {
		log.Println(err)
		return nil, err
	}

	return cred, nil
}


func SendRequest(method string, url string, payload io.Reader) (*http.Response, error) {
	client := &http.Client{}
	cred, err := config.GetCredential()
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(
		ctx,
		method,
		url,
		payload,
	)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("X-API-Version", "1")
	req.Header.Set("Authorization", "Bearer "+cred.AccessToken)

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	if res.StatusCode == http.StatusUnauthorized {
		log.Println("Starting New Session")
		cred, err = RefreshSession(cred.RefreshToken)
		if err != nil {
			return nil, err
		}
		log.Println("New Session")
		go config.SaveCredential(*cred)
		req.Header.Set("Authorization", "Bearer "+cred.AccessToken)
		return client.Do(req)
	}

	return res, err
}



func WriteOutput(res *http.Response, err error, data any){
	if err != nil {
		os.Stderr.WriteString(err.Error())
		return
	}
	defer res.Body.Close()

	//data := &models.Response{}
	errRes := models.ErrResponse{}


	if res.StatusCode > 308{
		if err := json.NewDecoder(res.Body).Decode(&errRes); err != nil {
			log.Printf("Error decoding response: %v", err)
			return
		}
		os.Stderr.WriteString(errRes.Message)
		return
	}

	if err := json.NewDecoder(res.Body).Decode(data); err != nil {
		os.Stderr.WriteString(err.Error())
		return
	}

	if err := json.NewEncoder(os.Stdout).Encode(data); err != nil {
		os.Stderr.WriteString(err.Error())
		return
	}

}


func ParseArgs(args []string) url.Values {
	fs := flag.NewFlagSet("profiles", flag.ExitOnError)

	minAge := fs.Int("min-age", -1, "minimum age to return")
	maxAge := fs.Int("max-age", -1, "maximum age to return")
	gender := fs.String("gender", "", "Gender to return")
	minGenderProb := fs.Float64("min-gender-probability", -1, "minimum gender probability to return")
	countryID := fs.String("country-id", "", "filter by country id")
	countryName := fs.String("country-name", "", "filter by country name")
	ageGroup := fs.String("age-group", "", "filter by age group")
	minCountryProb := fs.Float64("min-country-probability", -1, "minimum country probability to return")
	sortBy := fs.String("sort-by", "age", "criteria to sort the return profiles")
	order := fs.String("order", "DESC", "criteria to order the return profiles")
	page := fs.Int("page", 1, "pagination")
	limit := fs.Int("limit", 10, "limit per page")

	fs.Parse(args)

	query := url.Values{}
	if *maxAge != -1 {
		query.Set("max_age", strconv.Itoa(*maxAge))
	}

	if *minAge != -1 {
		query.Set("min_age", strconv.Itoa(*minAge))
	}

	if *minCountryProb != -1 {
		query.Set("min_country_probability", fmt.Sprintf("%.2f", *minCountryProb))
	}

	query.Set("sort_by", *sortBy)
	query.Set("order", *order)

	if *minGenderProb != -1 {
		query.Set("min_gender_probability", fmt.Sprintf("%.2f", *minGenderProb))
	}

	if *ageGroup != "" {
		query.Set("age_group", *ageGroup)
	}

	if *gender != "" {
		query.Set("gender", *gender)
	}

	if *countryID != "" {
		query.Set("country_id", *countryID)
	}

	if *countryName != "" {
		query.Set("country_name", *countryName)
	}

	query.Set("page", strconv.Itoa(*page))
	query.Set("limit", strconv.Itoa(*limit))

	return query
}

