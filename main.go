package gist

const VERSION = "v0.1.0"

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
	GITHUB_BASE_PATH string
	URL_ENV_NAME     string
	USER_AGENT       = "gist/#" + VERSION //Github requires this, else rejects API request
)
