package channel

import "fmt"

type Channel struct {
	name string
}

var ch = make(map[string]chan map[string]string)

func NewChannel(name string) *Channel {
	return &Channel{
		name: name,
	}
}

func (channel *Channel) Create() {
	if _, ok := ch[channel.name]; ok {
		fmt.Printf("[%s] Channel already exists\n", channel.name)
		return
	}

	ch[channel.name] = make(chan map[string]string)
}

func (channel *Channel) Listen() <-chan map[string]string {
	return ch[channel.name]
}

func (channel *Channel) Send(message map[string]string) {
	ch[channel.name] <- message
}
