package bennu

import (
	"net/http"
	"github.com/gorilla/websocket"
	"errors"
	"encoding/json"
	"fmt"
)

type session struct {
	conn *websocket.Conn
	request *http.Request
	handler *SocketHandler
}

func (sock *SocketHandler) newSession(r *http.Request, conn *websocket.Conn) *session {
	return &session{
		conn: conn,
		handler: sock,
		request: r,
	}

}

func (s *session) close() {
}

func (s *session) writer() {
}

type okReply struct {
	Status string `json:"status"`
	Response interface{} `json:"response"`
}

func (sesh *session) readUntilErrorOrClose() error {
	for {
		msgType, data, err := sesh.conn.ReadMessage()
		if err != nil {
			return err
		} else if msgType != websocket.TextMessage {
			return errors.New("Unsupported message type")
		}

		msg := &envelope{}
		if err := json.Unmarshal(data, msg); err != nil {
			return err
		}

		fmt.Println("INC", msg)

		if msg.Topic == "phoenix" && msg.Event == "heartbeat" {
			sesh.conn.WriteJSON(msg)
			continue
		}

		err = sesh.handler.handleMessage(msg)
		if err, isOk := err.(*errOkReply); isOk {
			sesh.conn.WriteJSON(&envelope{
				Ref: msg.Ref,
				Topic: msg.Topic,
				Event: fmt.Sprintf("chan_reply_%s", msg.Ref),
				Payload: &okReply{
					Status: "ok",
					Response: err.reply,
				},
			})
		} else if err != nil {
			return err
		} else {
			return fmt.Errorf("Unhandled message: %s", msg)
		}
	}
}
