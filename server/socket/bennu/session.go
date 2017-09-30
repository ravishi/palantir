package bennu

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/gorilla/websocket"
)

type session struct {
	inc chan *envelope
	conn *websocket.Conn
	request *http.Request
	handler *SocketHandler
}

func (sock *SocketHandler) newSession(r *http.Request, conn *websocket.Conn) *session {
	s := &session{
		inc: make(chan *envelope),
		conn: conn,
		handler: sock,
		request: r,
	}
	sock.bc.Subscribe(s.inc)
	return s
}

func (sesh *session) close() {
	// TODO
}

type okReply struct {
	Status string `json:"status"`
	Response interface{} `json:"response"`
}

type errorReply struct {
	Status string `json:"status"`
	Reason interface{} `json:"reason"`
}

func (sesh *session) readUntilErrorOrClose() error {
	for {
		msgType, data, err := sesh.conn.ReadMessage()
		if err != nil {
			return err
		} else if msgType != websocket.TextMessage {
			return errors.New("unsupported message type")
		} else {
			fmt.Println("RCV", string(data))
		}

		msg := &envelope{}
		if err := json.Unmarshal(data, msg); err != nil {
			return err
		}

		if msg.Topic == "phoenix" && msg.Event == "heartbeat" {
			if err := sesh.writeMessage(msg); err != nil {
				// Ignore heartbeat error ?
			}
			continue
		}

		err = sesh.handleMessage(msg)
		if _, isNoReply := err.(*errNoReply); !isNoReply && err != nil {
			return err
		}
	}
}

// Returns nil if everything is OK, errors otherwise
func (sesh *session) handleMessage(msg *envelope) error {
	if msg.Event == "phx_join" {
		err := sesh.handleJoin(msg)
		if ok, isOk := err.(*errOkReply); isOk {
			return sesh.writeMessage(&envelope{
				Ref: msg.Ref,
				Topic: msg.Topic,
				Event: fmt.Sprintf("chan_reply_%s", msg.Ref),
				JoinRef: msg.JoinRef,
				Payload: &okReply{
					Status: "ok",
					Response: ok.reply,
				},
			})
		} else if errReply, isErrorReply := err.(*errErrorReply); isErrorReply {
			return sesh.writeMessage(&envelope{
				Ref: msg.Ref,
				Topic: msg.Topic,
				Event: "phx_error",
				Payload: &errorReply{
					Status: "error",
					Reason: errReply.reason,
				},
			})
		} else {
			return err
		}
	} else if !strings.HasPrefix(msg.Event, "phx_") {
		err := sesh.handleIn(msg)
		if _, isNoReply := err.(*errNoReply); isNoReply {
			return nil
		} else if ok, isOk := err.(*errOkReply); isOk {
			return sesh.writeMessage(&envelope{
				Ref: msg.Ref,
				Topic: msg.Topic,
				Event: fmt.Sprintf("chan_reply_%s", msg.Ref),
				JoinRef: msg.JoinRef,
				Payload: &okReply{
					Status: "ok",
					Response: ok.reply,
				},
			})
		} else if errReply, isErrorReply := err.(*errErrorReply); isErrorReply {
			return sesh.writeMessage(&envelope{
				Ref: msg.Ref,
				Topic: msg.Topic,
				Event: "phx_error",
				Payload: &errorReply{
					Status: "error",
					Reason: errReply.reason,
				},
			})
		} else {
			return err
		}
	}
	return fmt.Errorf("TODO: handle other control events, like: %s", msg)
}

func (sesh *session) writer() {
	for {
		select {
		case msg, ok := <-sesh.inc:
			if !ok {
				return
			}
			err := sesh.writeMessage(msg)
			if err != nil {
				// TODO Notify someone, do something
				fmt.Errorf("error while writing message: %s\n", err)
			}

			// TODO Case we are closed
		}
	}
}

func (sesh *session) writeMessage(e *envelope) error {
	data, err := json.Marshal(e)
	if err != nil {
		return err
	}

	w, err := sesh.conn.NextWriter(websocket.TextMessage)
	if err != nil {
		return err
	}

	_, err1 := w.Write(data)
	err2 := w.Close()
	if err1 != nil {
		return err1
	}

	if err2 == nil {
		fmt.Println("SND", string(data))
	}
	return err2
}

func (sesh *session) handleJoin(msg *envelope) error {
	s := &socket{
		topic: asTopic(msg.Topic),
		payload: msg.Payload,
		session: sesh,
	}
	for _, ch := range sesh.handler.channels {
		s.channel = ch
		var err error
		s.joined, err = ch.handleJoin(s)
		if err != nil || s.joined {
			return err
		}
	}
	return nil
}

func (sesh *session) handleIn(msg *envelope) error {
	s := &socket{
		event: msg.Event,
		topic: asTopic(msg.Topic),
		payload: msg.Payload,
		session: sesh,
	}
	for _, ch := range sesh.handler.channels {
		s.channel = ch
		if err := ch.handleIn(s); err != nil {
			return err
		}
	}
	return nil
}
