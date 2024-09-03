package sdk

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

func newRequest(ctx context.Context, method string, action string, body io.Reader) (*http.Request, error) {
	url := fmt.Sprintf("%s/%s", API_URL, action)
	req, err := http.NewRequest(method, url, body)

	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")
	token := ctx.Value(CONTEXT_KEY_TOKEN)
	if token != nil {
		req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token.(string)))
	}

	return req, nil
}

type AuthUser struct {
	Identifier string `json:"identifier"`
	Password   string `json:"password"`
}

type AuthSession struct {
	Active    bool   `json:"active"`
	AccessJWT string `json:"accessJwt"`
}

func Authenticate(authUser AuthUser) (context.Context, error) {
	client := &http.Client{}

	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(authUser)
	if err != nil {
		return nil, err
	}

	req, err := newRequest(context.Background(), "POST", "com.atproto.server.createSession", &buf)

	if err != nil {
		return nil, err
	}

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		return nil, errors.New("Invalid identifier or password")
	}

	var authSession AuthSession
	err = json.Unmarshal(body, &authSession)

	if err != nil {
		return nil, err
	}

	if !authSession.Active {
		return nil, errors.New("User not active")
	}

	ctx := context.WithValue(context.Background(), CONTEXT_KEY_TOKEN, authSession.AccessJWT)

	return ctx, nil
}

type Feed struct {
	Post struct {
		Author struct {
			ID          string `json:"did"`
			DisplayName string `json:"displayName"`
			Handle      string `json:"handle"`
		} `json:"author"`
		Record struct {
			Reply interface{} `json:"reply"`
			Langs []string    `json:"langs"`
			Text  string      `json:"text"`
		} `json:"record"`
		ReplyCount  int `json:"replyCount"`
		RepostCount int `json:"repostCount"`
		LikeCount   int `json:"likeCount"`
		QuoteCount  int `json:"quoteCount"`
	} `json:"post"`
}

type Timeline struct {
	Feed []Feed `json:"feed"`
}

func GetTimeline(ctx context.Context) (Timeline, error) {
	client := &http.Client{}

	req, err := newRequest(ctx, "GET", "app.bsky.feed.getTimeline", nil)
	if err != nil {
		return Timeline{}, err
	}

	res, err := client.Do(req)
	if err != nil {
		return Timeline{}, err
	}

	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return Timeline{}, err
	}

	var timeline Timeline

	err = json.Unmarshal(body, &timeline)
	if err != nil {
		return Timeline{}, err
	}

	return timeline, nil
}
