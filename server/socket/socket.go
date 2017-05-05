package socket

import "github.com/ravishi/palantir/server/socket/bennu"

func HelloSocket() *bennu.SocketHandler {
	hello := bennu.NewChannelBuilder()

	hello.Join("room:lobby", func (s bennu.JoiningSocket) error {
		return s.Ok()
	})

	hello.HandleIn("hello", func (e bennu.Event) error {
		return e.Reply("hello, there!")
	})

	return bennu.NewSocketHandler(hello.Build())
}
