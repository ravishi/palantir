package bennu

import (
	"net/http"
	"github.com/gorilla/websocket"
)

type (
	Reason string

	Socket interface {
	}

	JoiningSocket interface {
		Socket

		// Like {:ok, socket}
		Ok() error

		// Like {:ok, reply :: map, socket}
		OkReply(reply interface{}) error

		// Like {:error, reply :: map}
		ErrorReply(reply interface{}) error
	}

	Event interface {
		Socket

		// Like {:noreply, socket}
		NoReply() error

		// {:reply, reply, socket}
		Reply(reply interface{}) error

		// {:stop, :normal, socket}
		Stop() error

		// {:stop, {:shutdown, term}, socket}
		Shutdown(error) error
	}

	SocketHandler struct {
		ch *Channel

		upgrader *websocket.Upgrader

		broadcaster *broadcaster
	}
)

const (
	Reason
)

func NewSocketHandler(ch *Channel) *SocketHandler {
	return &SocketHandler{
		ch: ch,
		upgrader: &websocket.Upgrader{
			// TODO CheckOrigin, make sure origin is checked by default and is only disabled when wanted (by flag, config, different default for dev server, whatever)
			CheckOrigin: func (r *http.Request) { return true },
			ReadBufferSize: 1024,
			WriteBufferSize: 1024,
		},
		broadcaster: createBroadcaster(),
	}
}

func (s *SocketHandler) Handle(w http.Response, r http.Request) error {
	ws, err := s.upgrader.Upgrade(w, r, nil)
	if err != nil {
		return err
	}
	defer ws.Close()

	return nil
}

func (s *SocketHandler) Close() {
	s.broadcaster.Close()
}
