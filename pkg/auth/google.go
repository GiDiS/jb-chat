package auth

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"io/ioutil"
	"jb_chat/pkg/models"
	"net/http"
	"net/url"
	"os"
)

var googleOauthConfig *oauth2.Config
var randState = "JBChat"

const userInfoUrl = "https://www.googleapis.com/oauth2/v3/userinfo"

func init() {

	googleOauthConfig = &oauth2.Config{
		RedirectURL:  "http://localhost:8888/api/auth/google",
		ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
		ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email"},
		Endpoint:     google.Endpoint,
	}
}

func RandToken(l int) (string, error) {
	b := make([]byte, l)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(b), nil
}

func GetProfileByAccessToken(ctx context.Context, accessToken string) (*models.GoogleProfile, error) {
	v := url.Values{}
	v.Set("access_token", accessToken)
	req, err := http.NewRequest("GET", userInfoUrl+"?"+v.Encode(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed make user info request: %s", err.Error())
	}
	req = req.WithContext(ctx)
	userinfo, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed getting user info: %s", err.Error())
	}
	defer func() { _ = userinfo.Body.Close() }()
	data, _ := ioutil.ReadAll(userinfo.Body)
	u := models.GoogleProfile{}
	if err = json.Unmarshal(data, &u); err != nil {
		return nil, err
	}
	return &u, nil
}

func GetProfile(ctx context.Context, state string, code string) (*models.GoogleProfile, error) {
	if randState != state {
		return nil, errors.New("invalid session state")
	}
	tok, err := googleOauthConfig.Exchange(ctx, code)
	if err != nil {
		return nil, err
	}
	client := googleOauthConfig.Client(ctx, tok)
	userinfo, err := client.Get(userInfoUrl)
	if err != nil {
		return nil, err
	}

	defer func() { _ = userinfo.Body.Close() }()
	data, _ := ioutil.ReadAll(userinfo.Body)
	u := models.GoogleProfile{}
	if err = json.Unmarshal(data, &u); err != nil {
		return nil, err
	}
	return &u, nil
}
