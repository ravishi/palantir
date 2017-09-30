package bennu

import (
	"net/http"

	"github.com/gorilla/websocket"
)

type (
	SocketHandler struct {
		bc *broadcaster
		channels []*Channel
		upgrader websocket.Upgrader
	}
)

func New() *SocketHandler {
	return &SocketHandler{
		bc: newBroadcaster(),
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

func (sock *SocketHandler) findChannel(topic string) *Channel {
	// XXX least found is returned, and note that it can still be nil if nothing is found
	// TODO make "*" a special case, so that order doesn't matter when '*' is used after
	// everything else, like: Join("room:private", ...); Join("room:*", ...);
	// NOTE that it would also make the interface more complex to explain and implement
	var found *Channel
	for _, ch := range sock.channels {
			if isTopicMatch(ch.topic, topic) {
				found = ch
			}
		}
	return found
}
