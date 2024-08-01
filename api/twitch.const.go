package api

const (
	EventsubMessageIDHeaderKey                = "Twitch-Eventsub-Message-Id"
	EventsubMessageTimestampHeaderKey         = "Twitch-Eventsub-Message-Timestamp"
	EventsubMessageSignatureHeaderKey         = "Twitch-Eventsub-Message-Signature"
	EventsubMessageTypeHeaderKey              = "Twitch-Eventsub-Message-Type"
	EventsubSubscriptionTypeChannelRedemption = "channel.channel_points_custom_reward_redemption.add"
	EventsubSubscriptionTypeStreamOnline      = "stream.online"
	EventsubSubscriptionTypeStreamOffline     = "stream.offline"
	EventsubSubscriptionTypeFollow            = "channel.follow"
)
