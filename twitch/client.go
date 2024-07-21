package twitch

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
)

func NewClient(clientID, clientSecret, redirectURI string) *Client {
	tokenURL := "https://id.twitch.tv/oauth2/token"
	helixURL := "https://api.twitch.tv/helix"

	return &Client{
		clientID:     clientID,
		clientSecret: clientSecret,
		redirectURI:  redirectURI,
		tokenURL:     tokenURL,
		helixURL:     helixURL,
	}
}

func (client *Client) GetOAuthToken(code string) (*OAuthToken, error) {
	var httpClient = &http.Client{}

	params := url.Values{}
	params.Set("client_id", client.clientID)
	params.Set("client_secret", client.clientSecret)
	params.Set("code", code)
	params.Set("grant_type", "authorization_code")
	params.Set("redirect_uri", "https://wild-grapes-flow.loca.lt/auth/login")

	req, err := http.NewRequest("POST", client.tokenURL, bytes.NewBufferString(params.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	res, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != 200 {
		return nil, errors.New("failed to get oauth token: " + string(resBody))
	}

	// convert bytes to struct
	var token = OAuthToken{}
	err = json.Unmarshal(resBody, &token)
	if err != nil {
		return nil, err
	}

	return &token, nil
}

func (client *Client) GetUserInfo(accessToken, userID string) (*UserInfoData, error) {
	var httpClient = &http.Client{}
	url := client.helixURL + "/users"
	if userID != "" {
		params := "?id=" + userID
		url += params
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Client-Id", client.clientID)

	res, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != 200 {
		return nil, errors.New("failed to get user info: " + string(resBody))
	}

	var userInfo = UserInfo{}
	err = json.Unmarshal(resBody, &userInfo)
	if err != nil {
		return nil, err
	}

	return &userInfo.Data[0], nil
}
