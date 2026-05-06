package cmd

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"insighta/internal/models"
	"io"
	"log"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"slices"
	"strconv"
	"strings"
	"time"
)

var (
	supportedFormat = []string{
		"csv",
		"json",
	}
)

func Export(args []string) {
	backendURL := "http://localhost:3030/api/profiles/export"
	query := ParseArgs(args)
	format := query.Get("format")

	format = strings.ToLower(format)
	if format == "" || !slices.Contains(supportedFormat, format) {
		fmt.Fprintf(os.Stderr, "\n\033[1;31mError\033[0m:Invalid format: %v\n", format)
		return
	}

	res, err := SendRequest(
		http.MethodGet,
		backendURL+"?"+query.Encode(),
		nil,
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "\n\033[1;31mError\033[0m:%v\n", err)
		return
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		errRes := models.ErrResponse{}
		err := json.NewDecoder(res.Body).Decode(&errRes)
		if err != nil {
			log.Printf("export: %v", err)
			return
		}
		fmt.Fprintf(os.Stderr, "\n\033[1;31mError\033[0m:%v\n", errRes.Message)
		return
	}

	cd := res.Header.Get("Content-Disposition")
	filename := getFileNameFromHeader(cd)
	if filename == "" {
		filename = "data_" + strconv.FormatInt(time.Now().Unix(), 10) + "." + format
	}

	cwd, err := os.Getwd()
	if err != nil {
		fmt.Fprintf(os.Stderr, "\n\033[1;31mError\033[0m:%v\n", err)
		return
	}
	file := filepath.Join(cwd, filename)

	f, err := os.OpenFile(file, os.O_RDWR|os.O_CREATE, 0o644)
	if err != nil {
		fmt.Fprintf(os.Stderr, "\n\033[1;31mError\033[0m:%v\n", err)
		return
	}
	defer f.Close()

	format = strings.ToLower(format)
	switch format {
	case "csv":
		err := importCSV(res.Body, f)
		if err != nil {
			os.Remove(file)
			fmt.Fprintf(os.Stderr, "\n\033[1;31mError\033[0m:%v\n", err)
			return
		}
		fmt.Println("\nFilename: ", filename)
	case "json":
		err := importJSON(res.Body, f)
		if err != nil {
			os.Remove(file)
			fmt.Fprintf(os.Stderr, "\n\033[1;31mError\033[0m:%v\n", err)
			return
		}
		fmt.Println("\nFilename: ", filename)
	}
}

func importCSV(r io.Reader, w io.Writer) error {
	records, err := csv.NewReader(r).ReadAll()

	if err != nil {
		return err
	}

	cw := csv.NewWriter(w)
	err = cw.WriteAll(records)
	if err != nil {
		return err
	}
	cw.Flush()
	return nil
}

func importJSON(r io.Reader, w io.Writer) error {
	data := []models.Profile{}
	err := json.NewDecoder(r).Decode(&data)
	if err != nil {
		return err
	}

	enc := json.NewEncoder(w)
	enc.SetIndent("", " ")
	err = enc.Encode(data)
	if err != nil {
		return err
	}

	return nil
}

func getFileNameFromHeader(s string) string {
	if s == "" {
		return ""
	}

	_, params, err := mime.ParseMediaType(s)
	if err != nil {
		log.Println(err)
		return ""
	}

	filename, exist := params["filename"]
	if !exist {
		return ""
	}

	if filename == "" {
		return ""
	}

	return filename
}
