package twitch

import (
	"fmt"
	"strings"
	"time"

	twitchClient "github.com/gempir/go-twitch-irc/v4"
	"github.com/sunnyegg/go-so/channel"
)

type ChatClient struct {
	username  string
	token     string
	ircClient *twitchClient.Client
}

var ConnectedClients = make(map[string]*twitchClient.Client)
var AlreadyPresent = make(map[string]bool)

func NewChatClient(username, token string) *ChatClient {
	if _, ok := ConnectedClients[username]; !ok {
		ConnectedClients[username] = twitchClient.NewClient(username, "oauth:"+token)
	}

	return &ChatClient{
		username:  username,
		token:     token,
		ircClient: ConnectedClients[username],
	}
}

func (client *ChatClient) Connect(config ConnectConfig) {
	client.ircClient.OnPrivateMessage(func(message twitchClient.PrivateMessage) {
		fmt.Printf("[%s] %s: %s\n", message.Channel, message.User.DisplayName, message.Message)

		// skip blacklist
		if len(config.Blacklist) > 0 {
			for _, blacklist := range config.Blacklist {
				if strings.EqualFold(blacklist, message.User.Name) {
					return
				}
			}
		}

		user := config.StreamID + message.User.Name

		if _, ok := AlreadyPresent[user]; ok {
			return
		}

		AlreadyPresent[user] = true
		ch := channel.NewChannel(channel.ChannelWebsocket)
		ch.Send(map[string]string{
			"stream_id": config.StreamID,
			"username":  message.User.Name,
		})

		if config.IsAutoSO {
			go func() {
				time.Sleep(time.Second * time.Duration(config.Delay))

				// !so message to twitch chat
				client.ircClient.Say(message.Channel, "!so "+message.User.Name)
			}()
		}
	})

	go func() {
		err := client.ircClient.Connect()
		if err != nil {
			fmt.Printf("[%s] Error when connecting to twitch irc: %s\n", client.username, err)
		}
	}()
}

func (client *ChatClient) Disconnect(username string) {
	err := client.ircClient.Disconnect()
	if err != nil {
		fmt.Printf("[%s] Error when disconnecting from twitch irc: %s\n", username, err)
	}
}

func (client *ChatClient) Join(username, channel string) {
	client.ircClient.Join(channel)
	fmt.Printf("[%s] Joined channel: %s\n", client.username, channel)
}

func (client *ChatClient) Depart(username, channel string) {
	client.ircClient.Depart(channel)
	fmt.Printf("[%s] Departed channel: %s\n", client.username, channel)
}
