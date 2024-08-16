package twitch

import (
	"fmt"
	"time"

	twitchClient "github.com/gempir/go-twitch-irc/v4"
	"github.com/sunnyegg/go-so/channel"
)

type ChatClient struct {
	username  string
	token     string
	ircClient *twitchClient.Client
}

var connectedClients = make(map[string]*twitchClient.Client)
var alreadyPresent = make(map[string]map[string]bool)

func NewChatClient(username, token string) *ChatClient {
	if _, ok := connectedClients[username]; !ok {
		connectedClients[username] = twitchClient.NewClient(username, "oauth:"+token)
	}

	return &ChatClient{
		username:  username,
		token:     token,
		ircClient: connectedClients[username],
	}
}

func (client *ChatClient) Connect(config ConnectConfig) {
	var msgGeneral map[string]string
	var newToken string
	chGeneral := channel.NewChannel(channel.ChannelGeneral)
	chWs := channel.NewChannel(channel.ChannelWebsocket)
	token := client.token

	go func() {
		for {
			msgGeneral = <-chGeneral.Listen()

			if msgGeneral["channel"] == client.username {
				newToken = msgGeneral["token"]
			}
		}
	}()

	if _, ok := alreadyPresent[config.StreamID]; !ok {
		alreadyPresent[config.StreamID] = make(map[string]bool)
	}

	client.ircClient.OnPrivateMessage(func(message twitchClient.PrivateMessage) {
		fmt.Printf("[%s] %s: %s\n", message.Channel, message.User.DisplayName, message.Message)

		if newToken != "" && token != newToken {
			token = newToken
			client.ircClient.SetIRCToken(token)
			fmt.Printf("[%s] Token is updated\n", client.username)
		}

		if _, ok := alreadyPresent[config.StreamID][message.User.Name]; ok {
			return
		}

		alreadyPresent[config.StreamID][message.User.Name] = true

		go func() {
			chWs.Send(map[string]string{
				"stream_id": config.StreamID,
				"username":  message.User.Name,
				"channel":   message.Channel,
				"token":     token,
			})
		}()

		if config.IsAutoSO {
			go func() {
				time.Sleep(time.Second * time.Duration(config.Delay))

				// !so message to twitch chat
				client.ircClient.Say(message.Channel, "!so @"+message.User.Name)
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

func (client *ChatClient) Disconnect(streamId, username string) {
	err := client.ircClient.Disconnect()
	if err != nil {
		fmt.Printf("[%s] Error when disconnecting from twitch irc: %s\n", username, err)
	}

	delete(alreadyPresent, streamId)
	delete(connectedClients, username)
}

func (client *ChatClient) Join(username, channel string) {
	client.ircClient.Join(channel)
	fmt.Printf("[%s] Joined channel: %s\n", client.username, channel)
}

func (client *ChatClient) Depart(username, channel string) {
	client.ircClient.Depart(channel)
	fmt.Printf("[%s] Departed channel: %s\n", client.username, channel)
}
