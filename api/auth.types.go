package api

import "github.com/sunnyegg/go-so/twitch"

type loginUserRequest struct {
	Code             string `form:"code"`
	Scope            string `form:"scope"`
	State            string `form:"state" binding:"required"`
	Error            string `form:"error"`
	ErrorDescription string `form:"error_description"`
}

type userResponse struct {
	UserLogin       string `json:"user_login"`
	UserName        string `json:"user_name"`
	ProfileImageUrl string `json:"profile_image_url"`
}

type loginUserResponse struct {
	AccessToken  string       `json:"access_token"`
	RefreshToken string       `json:"refresh_token"`
	User         userResponse `json:"user"`
}

type createOrUpdateUserArg struct {
	UserID          string
	UserLogin       string
	UserName        string
	ProfileImageUrl string
	Token           *twitch.OAuthToken
}

type refreshUserRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type createStateResponse struct {
	URL string `json:"url"`
}

type logoutUserRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}
