package sdk

type contextKey string

func (c contextKey) String() string {
	return string(c)
}

const API_URL = "https://bsky.social/xrpc"
const CONTEXT_KEY_TOKEN = contextKey("token")
