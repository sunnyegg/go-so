package twitch

import (
	"fmt"

	twitchClient "github.com/gempir/go-twitch-irc/v4"
	"github.com/sunnyegg/go-so/channel"
)

type ChatClient struct {
	username  string
	token     string
	ircClient *twitchClient.Client
}

var connectedClients = make(map[string]*twitchClient.Client)

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

	client.ircClient.OnPrivateMessage(func(message twitchClient.PrivateMessage) {
		fmt.Printf("[%s] %s: %s\n", message.Channel, message.User.DisplayName, message.Message)

		if newToken != "" && token != newToken {
			token = newToken
			client.ircClient.SetIRCToken(token)
			fmt.Printf("[%s] Token is updated\n", client.username)
		}

		go func() {
			chWs.Send(map[string]string{
				"username": message.User.Name,
				"channel":  message.Channel,
			})
		}()
	})

	go func() {
		err := client.ircClient.Connect()
		if err != nil {
			fmt.Printf("[%s] Error when connecting to twitch irc: %s\n", client.username, err)
		}
	}()
}

func Disconnect(username string) {
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
