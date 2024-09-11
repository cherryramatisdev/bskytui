package sdk

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/keybase/go-keychain"
)

func newRequest(token, method, action string, body io.Reader) (*http.Request, error) {
	url := fmt.Sprintf("%s/%s", API_URL, action)
	req, err := http.NewRequest(method, url, body)

	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")

	if token != "" {
		req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))
	}

	return req, nil
}

type AuthUser struct {
	Identifier string `json:"identifier"`
	Password   string `json:"password"`
}

type AuthSession struct {
	AuthUser
	Active    bool   `json:"active"`
	AccessJWT string `json:"accessJwt"`
}

func SaveAuthInfo(authSession *AuthSession) error {
	data, err := json.Marshal(authSession)

	if err != nil {
		return err
	}

	authItem := keychain.NewItem()
	authItem.SetSecClass(keychain.SecClassGenericPassword)
	authItem.SetService("bsky_tui")
	authItem.SetAccount(authSession.Identifier)
	authItem.SetAuthenticationType(keychain.AuthenticationTypeKey)
	authItem.SetData(data)
	authItem.SetAccessible(keychain.AccessibleWhenUnlocked)

	keychain.DeleteItem(authItem)
	keychain.AddItem(authItem)

	return nil
}

func LoadAuthInfo() (AuthSession, error) {
	queryItem := keychain.NewItem()
	queryItem.SetSecClass(keychain.SecClassGenericPassword)
	queryItem.SetService("bsky_tui")
	queryItem.SetMatchLimit(keychain.MatchLimitOne)
	queryItem.SetReturnAttributes(true)
	queryItem.SetReturnData(true)

	results, err := keychain.QueryItem(queryItem)
	if err != nil || len(results) == 0 {
		return AuthSession{}, err
	}

	var auth AuthSession
	if err := json.Unmarshal(results[0].Data, &auth); err != nil {
		return AuthSession{}, err
	}

	return auth, nil
}

func Authenticate(username, password string) error {
	client := &http.Client{}

	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(&AuthUser{
		Identifier: username,
		Password:   password,
	})

	if err != nil {
		return err
	}

	req, err := newRequest("", "POST", "com.atproto.server.createSession", &buf)

	if err != nil {
		return err
	}

	res, err := client.Do(req)
	if err != nil {
		return err
	}

	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}

	if res.StatusCode != http.StatusOK {
		return errors.New("Invalid identifier or password")
	}

	var authSession AuthSession
	err = json.Unmarshal(body, &authSession)

	if err != nil {
		return err
	}

	if !authSession.Active {
		return errors.New("User not active")
	}

	err = SaveAuthInfo(&authSession)

	if err != nil {
		return err
	}

	return nil
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
	} `json:"post"`
}

type Timeline struct {
	Feed []Feed `json:"feed"`
}

func GetTimeline(session *AuthSession) (Timeline, error) {
	client := &http.Client{}

	req, err := newRequest(session.AccessJWT, "GET", "app.bsky.feed.getTimeline", nil)
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
