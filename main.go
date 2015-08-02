// currently only supports anonymous gists
// or use env variables for Github auth tokens
package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

const VERSION = "v0.1.0"

//A list of clipboard commands with copy and paste support.
const (
	xclip   = "xclip -o"
	xsel    = "xsel -o"
	pbcopy  = "pbpaste"
	putclip = "getclip"
)

// github API urls
const (
	GITHUB_API_URL = "https://api.github.com/"
	GIT_IO_URL     = "http://git.io"
	GHE_BASE_PATH  = "/api/v3"
)

var (
	USER_AGENT = "gist/#" + VERSION //Github requires this, else rejects API request
	token      = os.Getenv("GITHUB_TOKEN")
)

var (
	publicFlag  bool
	description string
	anonymous   bool
	responseObj map[string]interface{}
)

type GistFile struct {
	Content string `json:"content"`
}

type Gist struct {
	Description string              `json:"description"`
	publicFile  bool                `json:"public"`
	GistFile    map[string]GistFile `json:"files"`
}

func loadTokenFromFile() (token string) {
	//get the tokenfile
	file := os.Getenv("$HOME") + "/.gist"
	github_token, err := ioutil.ReadFile(file)
	if err != nil {
		log.Fatal(err)
	}
	return string(github_token)
}

func usage() {
	fmt.Fprintf(os.Stderr, "usage: gist [-p] [-d] [-u] example.go\n")
	flag.PrintDefaults()
	os.Exit(2)
}

func main() {
	flag.BoolVar(&publicFlag, "p", true, "Set to false for private gist.")
	flag.BoolVar(&anonymous, "u", true, "Set false if gist should be not anonymous")
	flag.StringVar(&description, "d", "This is a gist", "Description for gist.")
	flag.Usage = usage
	flag.Parse()

	files_list := flag.Args()
	if len(files_list) == 0 {
		log.Fatal("Error: No files specified.")
	}

	//fmt.Println(files_list)
	//fmt.Println(token)

	files := map[string]GistFile{}

	for _, filename := range files_list {
		fmt.Println("Checking file:", filename)
		content, err := ioutil.ReadFile(filename)
		if err != nil {
			log.Fatal("File Error: ", err)
		}
		files[filename] = GistFile{string(content)}
	}

	if description == "" {
		description = strings.Join(files_list, ", ")
	}

	gist := Gist{
		description,
		publicFlag,
		files,
	}

	pfile, err := json.Marshal(gist)
	if err != nil {
		log.Fatal("Cannot marshal json: ", err)
	}

	fmt.Println("OK")

	b := bytes.NewBuffer(pfile)
	fmt.Println("Uploading...")

	if anonymous == true {
		response, err := http.Post("https://api.github.com/gists", "application/json", b)
		if err != nil {
			log.Fatal("HTTP error: ", err)
		}

		err = json.NewDecoder(response.Body).Decode(&responseObj)
		if err != nil {
			log.Fatal("Response JSON error: ", err)
		}

		fmt.Println("--- Gist URL ---")
		fmt.Println(responseObj["html_url"])
	} else {
		req, err := http.NewRequest("POST", "https://api.github.com/gists", b)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Accept", "application/json")
		req.SetBasicAuth(token, "x-oauth-basic")

		client := http.Client{}
		response, err := client.Do(req)
		if err != nil {
			log.Fatal("HTTP error: ", err)
		}
		err = json.NewDecoder(response.Body).Decode(&responseObj)
		if err != nil {
			log.Fatal("Response JSON error: ", err)
		}

		fmt.Println("--- Gist URL ---")
		fmt.Println(responseObj["html_url"])
	}
}
