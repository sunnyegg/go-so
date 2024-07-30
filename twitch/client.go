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
	oauthURL := "https://id.twitch.tv/oauth2"
	helixURL := "https://api.twitch.tv/helix"

	return &Client{
		clientID:     clientID,
		clientSecret: clientSecret,
		redirectURI:  redirectURI,
		oauthURL:     oauthURL,
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
	params.Set("redirect_uri", client.redirectURI+"/auth/login")

	req, err := http.NewRequest("POST", client.oauthURL+"/token", bytes.NewBufferString(params.Encode()))
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

func (client *Client) ValidateOAuthToken(accessToken string) error {
	var httpClient = &http.Client{}

	req, err := http.NewRequest("GET", client.oauthURL+"/validate", nil)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "OAuth "+accessToken)

	res, err := httpClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}

	if res.StatusCode == 401 {
		return errors.New(ErrExpiredToken)
	}

	if res.StatusCode != 200 {
		return errors.New("failed to validate oauth token: " + string(resBody))
	}

	return nil
}

func (client *Client) RefreshOAuthToken(refreshToken string) (*OAuthToken, error) {
	var httpClient = &http.Client{}

	params := url.Values{}
	params.Set("client_id", client.clientID)
	params.Set("client_secret", client.clientSecret)
	params.Set("grant_type", "refresh_token")
	params.Set("refresh_token", refreshToken)

	req, err := http.NewRequest("POST", client.oauthURL+"/token", bytes.NewBufferString(params.Encode()))
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
		return nil, errors.New("failed to refresh oauth token: " + string(resBody))
	}

	// convert bytes to struct
	var token = OAuthToken{}
	err = json.Unmarshal(resBody, &token)
	if err != nil {
		return nil, err
	}

	return &token, nil
}

func (client *Client) GetUserInfo(accessToken, userID, username string) (*UserInfoData, error) {
	var httpClient = &http.Client{}
	url := client.helixURL + "/users"
	if userID != "" {
		params := "?id=" + userID
		url += params
	}
	if username != "" {
		params := "?login=" + username
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

func (client *Client) ConnectTwitchChat(config ConnectConfig, username, accessToken string) {
	twitchClient := NewChatClient(username, accessToken)
	twitchClient.Connect(config)
	twitchClient.Join(username, username)
}

func (client *Client) GetStreamInfo(accessToken, userID string) (*StreamInfoData, error) {
	var httpClient = &http.Client{}
	url := client.helixURL + "/streams"
	if userID != "" {
		params := "?user_id=" + userID
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
		return nil, errors.New("failed to get stream info: " + string(resBody))
	}

	var streamInfo = StreamInfo{}
	err = json.Unmarshal(resBody, &streamInfo)
	if err != nil {
		return nil, err
	}
	if len(streamInfo.Data) == 0 {
		return nil, errors.New("stream not found")
	}

	return &streamInfo.Data[0], nil
}

func (client *Client) GetAppAccessToken(clientID, clientSecret string) (*AppAccessToken, error) {
	var httpClient = &http.Client{}

	params := url.Values{}
	params.Set("client_id", clientID)
	params.Set("client_secret", clientSecret)
	params.Set("grant_type", "client_credentials")

	req, err := http.NewRequest("POST", client.oauthURL+"/token", bytes.NewBufferString(params.Encode()))
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
		return nil, errors.New("failed to get app access token: " + string(resBody))
	}

	// convert bytes to struct
	var appAccessToken = AppAccessToken{}
	err = json.Unmarshal(resBody, &appAccessToken)
	if err != nil {
		return nil, err
	}

	return &appAccessToken, nil
}

func (client *Client) RegisterEventsub(accessToken string, subscription EventsubSubscription) error {
	var httpClient = &http.Client{}
	url := client.helixURL + "/eventsub/subscriptions"

	subscriptionBytes, err := json.Marshal(subscription)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(subscriptionBytes))
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Client-Id", client.clientID)
	req.Header.Set("Content-Type", "application/json")

	res, err := httpClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}

	// skip if already registered
	if res.StatusCode == 409 {
		return nil
	}

	// success
	if res.StatusCode == 202 {
		return nil
	}

	if res.StatusCode != 200 {
		return errors.New("failed to register eventsub: " + string(resBody))
	}

	return nil
}
