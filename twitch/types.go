package twitch

const (
	ErrExpiredToken = "expired token"
)

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

type AppAccessToken struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
	TokenType   string `json:"token_type"`
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

type StreamInfo struct {
	Data []StreamInfoData `json:"data"`
}

type StreamInfoData struct {
	ID        string `json:"id"`
	GameName  string `json:"game_name"`
	Title     string `json:"title"`
	StartedAt string `json:"started_at"`
}

type EventsubSubscription struct {
	Type      string                        `json:"type"`
	Version   string                        `json:"version"`
	Condition EventsubSubscriptionCondition `json:"condition"`
	Transport EventsubSubscriptionTransport `json:"transport"`
}

type EventsubSubscriptionCondition struct {
	BroadcasterUserID string `json:"broadcaster_user_id"`
}

type EventsubSubscriptionTransport struct {
	Method   string `json:"method"`
	Callback string `json:"callback"`
	Secret   string `json:"secret"`
}

type ConnectConfig struct {
	StreamID string `json:"stream_id"`
	Delay    int    `json:"delay"`
	IsAutoSO bool   `json:"is_auto_so"`
}
