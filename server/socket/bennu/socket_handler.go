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

func New(options ...func(*SocketHandler) error) (*SocketHandler, error) {
	s := &SocketHandler{
		bc: newBroadcaster(),
		channels: make([]*Channel, 0),
		upgrader: websocket.Upgrader{
			ReadBufferSize: 1024,
			WriteBufferSize: 1024,
		},
	}

	if err := s.setOption(options...); err != nil {
		return nil, err
	}

	return s, nil
}

func (h *SocketHandler) setOption(options ...func(*SocketHandler) error) error {
	for _, opt := range options {
		if err := opt(h); err != nil {
			return err
		}
	}
	return nil
}

func (h *SocketHandler) Channel(topic string) *Channel {
	ch := &Channel{
		topic: topic,
		handler: h,

		joinHandlers: make([]*joinHandler, 0),
	}

	h.channels = append(h.channels, ch)

	return ch
}

func (h *SocketHandler) Handle(w http.ResponseWriter, r *http.Request) error {
	conn, err := h.upgrader.Upgrade(w, r, nil)
	if err != nil {
		return err
	}
	defer conn.Close()

	sesh := h.newSession(r, conn)
	defer sesh.close()

	go sesh.writer()

	return sesh.readUntilErrorOrClose()
}

func (h *SocketHandler) findChannel(topic string) *Channel {
	// XXX least found is returned, and note that it can still be nil if nothing is found
	// TODO make "*" a special case, so that order doesn't matter when '*' is used after
	// everything else, like: Join("room:private", ...); Join("room:*", ...);
	// NOTE that it would also make the interface more complex to explain and implement
	var found *Channel
	for _, ch := range h.channels {
			if isTopicMatch(ch.topic, topic) {
				found = ch
			}
		}
	return found
}
