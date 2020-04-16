// This app is intented to be go-port of the defunckt's gist library in Ruby
// Currently, uploading single and multiple files are available.
// You can also create secret gists, and both anonymous and user gists.
//
// Author: Viyat Bhalodia
package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

//Defines the app version
const VERSION = "v0.4.0"

//#TODO: A list of clipboard commands with copy and paste support.
//This is intended for adding the gist URLs directly to the user clipboard,
//so that manual copying is not needed.
const (
	xclip   = "xclip -o"
	xsel    = "xsel -o"
	pbcopy  = "pbpaste"
	putclip = "getclip"
)

// Defines different constants used
// GIT_IO_URL is the Github's URL shortner
// API v3 is the current version of GitHub API
const (
	GITHUB_API_URL = "https://api.github.com/"
	BASE_PATH      = "/api/v3"
	GIT_IO_URL     = "https://git.io"
)

//User agent defines a custom agent (required by GitHub)
//`token` stores the GITHUB_TOKEN from the env variables
// GITHUB_TOKEN must be in format of `username:token`
var (
	USER_AGENT = "gist/" + VERSION
	username   = os.Getenv("GITHUB_USERNAME")
	token      = os.Getenv("GITHUB_TOKEN")
)

var (
	wantshort bool
	slug      string
)

// Variables used in `Gist` struct
var (
	public      bool
	description string
	anonymous   bool
	update      string
	responseObj map[string]interface{}
)

// The top-level struct for a gist file
type GistFile struct {
	Content string `json:"content"`
}

// The required structure for POST data for API purposes
type Gist struct {
	Description string              `json:"description",omitempty`
	Public      bool                `json:"public"`
	GistFile    map[string]GistFile `json:"files"`
}

//This function loads the GITHUB_TOKEN from a '$HOME/.gist' file
//from the user's home directory.
func loadTokenFromFile() (token string) {
	//get the tokenfile
	file := filepath.Join(os.Getenv("HOME"), ".gist")
	github_token, err := ioutil.ReadFile(file)
	if err != nil {
		log.Fatal(err)
	}
	return strings.TrimSpace(string(github_token))
}

// Defines basic usage when program is run with the help flag
func usage() {
	fmt.Fprintf(os.Stderr, "usage: gist [options] file...\n")
	flag.PrintDefaults()
	os.Exit(2)
}

// The main function parses the CLI args. It also checks the files, and
// loads them into an array.
// Then the files are separated into GistFile structs and collectively
// the files are saved in `files` field in the Gist struct.
// A request is then made to the GitHub api - it depends on whether it is
// anonymous gist or not.
// The response recieved is parsed and the Gist URL is printed to STDOUT.
func main() {
	flag.StringVar(&update, "update", "", "id of existing gist to update")
	flag.StringVar(&slug, "slug", "", "Set prefered short url")
	flag.BoolVar(&wantshort, "short", true, "Generate short url")
	flag.BoolVar(&public, "public", false, "Set to true for public gist.")
	flag.BoolVar(&anonymous, "anonymous", false, "Set to true for anonymous gist user")
	flag.StringVar(&description, "d", "", "Description for gist.")
	flag.Usage = usage
	flag.Parse()

	files_list := flag.Args()
	if len(files_list) == 0 {
		log.Fatal("Error: No files specified.")
	}

	files := map[string]GistFile{}

	for _, filename := range files_list {
		fmt.Println("Checking file:", filename)
		content, err := ioutil.ReadFile(filename)
		if err != nil {
			log.Fatal("File Error: ", err)
		}

		// gists api doesn't allow / on filenames
		name := filepath.Base(filename)
		files[name] = GistFile{string(content)}
	}

	if description == "" {
		description = strings.Join(files_list, ", ")
	}

	//create a gist from the files array
	gist := Gist{
		Description: description,
		Public:      public,
		GistFile:    files,
	}

	pfile, err := json.Marshal(gist)
	if err != nil {
		log.Fatal("Cannot marshal json: ", err)
	}

	b := bytes.NewBuffer(pfile)
	fmt.Println("Uploading...")

	// Send request to API
	post_to := GITHUB_API_URL + "gists"
	if update != "" {
		post_to += "/" + update
	}
	req, err := http.NewRequest("POST", post_to, b)
	if err != nil {
		log.Fatal("Cannot create request: ", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	if !anonymous {
		if token == "" {
			token = loadTokenFromFile()
		}
		words := strings.Split(token, ":")
		if len(words) != 2 {
			if username == "" {
				log.Fatalf("token must be in form `username:token`, was actually %s", token)
			}
			words = []string{username, token}
		}
		req.SetBasicAuth(words[0], words[1])
	}

	client := http.Client{}
	response, err := client.Do(req)
	if err != nil {
		log.Fatal("HTTP error: ", err)
	}
	defer response.Body.Close()
	err = json.NewDecoder(response.Body).Decode(&responseObj)
	if err != nil {
		log.Fatal("Response JSON error: ", err)
	}

	if _, ok := responseObj["html_url"]; !ok {
		// something went wrong
		fmt.Println(responseObj["message"])
		if a, ok := responseObj["errors"]; ok {
			for i, m := range a.([]interface{}) {
				for k, v := range m.(map[string]interface{}) {
					fmt.Printf("%d %s: %s\n", i, k, v)
				}
			}
		}
		os.Exit(1)
	}

	fmt.Println("===Gist URL===")
	fmt.Println(responseObj["html_url"])
	if wantshort {
		fmt.Println(shorten(responseObj["html_url"].(string)))
	}
}

func shorten(s string) string {

	form := url.Values{}
	form.Add("url", s)
	if slug != "" {
		form.Add("code", slug)
	}
	req, err := http.NewRequest("POST", GIT_IO_URL, strings.NewReader(form.Encode()))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	client := http.Client{}
	response, err := client.Do(req)
	if err != nil {
		log.Fatal("HTTP error: ", err)
	}
	defer response.Body.Close()

	switch response.StatusCode {
	case 200:
		// when we use /create we get 200 and the short url on the body
		b, _ := ioutil.ReadAll(response.Body)
		return GIT_IO_URL + "/" + string(b)
	case 201:
		// when we post to / we get 201 and the whole short url on the Location Header
		return string(response.Header["Location"][0])
	default:
		// epic fail!
		return s
	}
}
