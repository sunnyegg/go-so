package api

type getTwitchUserRequest struct {
	UserLogin string `form:"user_login" binding:"required"`
}

type connectChatRequest struct {
	StreamID  string `json:"stream_id" binding:"required"`
	UserLogin string `json:"user_login" binding:"required"`
	Channel   string `json:"channel" binding:"required"`
}

type eventsubRequest struct {
	Challenge    string `json:"challenge"`
	Subscription struct {
		Type      string `json:"type"`
		Condition struct {
			BroadcasterUserID string `json:"broadcaster_user_id"`
		}
	} `json:"subscription"`
	Event struct {
		BroadcasterUserID    string `json:"broadcaster_user_id"`    // streamer
		BroadcasterUserLogin string `json:"broadcaster_user_login"` // streamer
		UserLogin            string `json:"user_login"`             // chatter
		Reward               struct {
			Title string `json:"title"`
		} `json:"reward"`
	} `json:"event"`
}

type getChannelInfoRequest struct {
	UserLogin string `form:"user_login" binding:"required"`
}

type sendChatMessageRequest struct {
	Channel string `json:"channel" binding:"required"`
	Message string `json:"message" binding:"required"`
}

type sendShoutoutRequest struct {
	FromID      string `json:"from_id" binding:"required"`
	ToID        string `json:"to_id" binding:"required"`
	ModeratorID string `json:"moderator_id" binding:"required"`
}
