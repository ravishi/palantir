package bennu

import (
	"net/http"
	"github.com/gorilla/websocket"
	"fmt"
)

type (
	SocketHandler struct {
		channels []*Channel
		upgrader websocket.Upgrader
	}
)

func New() *SocketHandler {
	return &SocketHandler{
		channels: make([]*Channel, 0),
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool { return true },
			ReadBufferSize: 1024,
			WriteBufferSize: 1024,
		},
	}
}

func (sock *SocketHandler) Channel(topic string) *Channel {
	ch := &Channel{
		topic: topic,
		handler: sock,

		joinHandlers: make([]*joinHandler, 0),
	}

	sock.channels = append(sock.channels, ch)

	return ch
}

func (sock *SocketHandler) Handle(w http.ResponseWriter, r *http.Request) error {
	conn, err := sock.upgrader.Upgrade(w, r, nil)
	if err != nil {
		return err
	}
	defer conn.Close()

	sesh := sock.newSession(r, conn)

	defer sesh.close()

	go sesh.writer()

	return sesh.readUntilErrorOrClose()
}

func (sock *SocketHandler) handleMessage(msg *envelope) error {
	if msg.Event == "phx_join" {
		for _, ch := range sock.channels {
			err := ch.handleJoin(msg.Topic, msg.Payload)
			if err == nil {
				continue
			} else {
				return err
			}
		}
	}

	return fmt.Errorf("TODO: handle other events, like: %s", msg)
}
