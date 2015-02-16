package gist

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
)

const VERSION = "__0.10__"

//clipboard commands
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
	GITHUB_BASE_PATH = ""
	URL_ENV_NAME     = "GITHUB URL"
	USER_AGENT       = "gist/#" + VERSION //Github requires this, else rejects API request
)
