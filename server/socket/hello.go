package socket

import (
	"github.com/ravishi/palantir/server/socket/bennu"
)

func ChatSocket() (*bennu.SocketHandler, error) {
	b, err := bennu.New()
	if err != nil {
		return nil, err
	}

	room := b.Channel("room:*")

	room.Join("room:lobby", func (s bennu.JoinSocket) error {
		return nil
	})

	room.HandleIn("new_msg", func (s bennu.InSocket) error {
		s.Broadcast("new_msg", s.Payload())
		return nil
	})

	return b, nil
}
