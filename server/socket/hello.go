package socket

import (
	"github.com/ravishi/palantir/server/socket/bennu"
)

func ChatSocket() *bennu.SocketHandler {
	b := bennu.New()

	room := b.Channel("room:*")

	room.Join("room:lobby", func (s bennu.JoinSocket) error {
		return nil
	})

	room.HandleIn("new_msg", func (s bennu.InSocket) error {
		return nil
	})

	return b
}
