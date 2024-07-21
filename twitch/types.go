package twitch

type Client struct {
	clientID     string
	clientSecret string
	redirectURI  string
	oauthURL     string
	helixURL     string
}

type OAuthToken struct {
	AccessToken  string   `json:"access_token"`
	RefreshToken string   `json:"refresh_token"`
	ExpiresIn    int      `json:"expires_in"`
	Scope        []string `json:"scope"`
	TokenType    string   `json:"token_type"`
}

type UserInfo struct {
	Data []UserInfoData `json:"data"`
}

type UserInfoData struct {
	ID              string `json:"id"`
	Login           string `json:"login"`
	DisplayName     string `json:"display_name"`
	Type            string `json:"type"`
	Broadcaster     bool   `json:"broadcaster"`
	Description     string `json:"description"`
	ProfileImageURL string `json:"profile_image_url"`
}

type ValidateOAuthToken struct {
	ClientID  string   `json:"client_id"`
	Login     string   `json:"login"`
	Scopes    []string `json:"scopes"`
	UserID    string   `json:"user_id"`
	ExpiresIn int      `json:"expires_in"`
}
