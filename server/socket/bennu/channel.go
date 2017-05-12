package bennu

import (
	"strings"
	"fmt"
)

type (
	Handler func(Socket) error

	Channel struct {
		topic string
		handler *SocketHandler
		joinHandlers []*joinHandler
	}

	joinHandler struct {
		topic string
		handler Handler
	}

	// returned when a topic is not handled by a channel.
	errNotHandled struct {
		topic string
	}

	// returned when a channel.Join method returns c.Ok() or c.OkReply(reply)
	errOkReply struct {
		reply interface{}
		socket *socket
	}

	// returned when a channel.Join method returns c.Error() or c.ErrorReply(reply)
	errErrorReply struct {
		reply interface{}
		socket *socket
	}

	topic struct {
		topic string
		subtopic string
	}
)

func (err *errNotHandled) Error() string {
	return "topic is not handled by this channel, topic = " + err.topic
}

func (err *errErrorReply) Error() string {
	return fmt.Sprintf("error reply: %s", err.reply)
}

func (err *errOkReply) Error() string {
	return "ok"
}

func (ch *Channel) Join(topic string, handler Handler) {
	ch.joinHandlers = append(ch.joinHandlers, &joinHandler{
		topic: topic,
		handler: handler,
	})
}

func (ch *Channel) handleJoin(topic string, payload interface{}) error {
	handler := ch.searchJoinHandler(topic)
	if handler == nil {
		return &errNotHandled{topic}
	}

	s := &socket{
		topic: asTopic(topic),
		channel: ch,
		payload: payload,
	}

	err := handler.handler(s)

	if ok, isOk := err.(*errOkReply); isOk {
		return ch.join(ok.socket, ok.reply)
	}


	if err, isErr := err.(*errErrorReply); isErr {
		// TODO what kind of error should we return here?
		return err
	}

	return fmt.Errorf("join failed: it didn't return either Ok() or Error(), it just crashed with %s", err)
}

func (ch *Channel) join(s *socket, reply interface{}) error {
	// TODO subscribe to the broadcasting
	s.joined = true
	return &errOkReply{
		socket: s,
		reply: reply,
	}
}

func (ch *Channel) searchJoinHandler(topic string) *joinHandler {
	// XXX least found is returned, and note that it can still be nil if nothing is found
	// TODO diferentiate between "*" and specifics, so that order doesn't matter in this case: Join("room:private", ...); Join("room:*", ...);
	var found *joinHandler
	for _, h := range ch.joinHandlers {
		if isTopicMatch(h.topic, topic) {
			found = h
		}
	}
	return found
}

func asTopic(s string) *topic {
	p := strings.SplitN(s, ":", 2)
	t := &topic{}
	if len(p) >= 1 {
		t.topic = p[0]
	}
	if len(p) == 2 {
		t.subtopic = p[1]
	}
	return t
}

func (t *topic) isMatch(other *topic) bool {
	return t.topic == other.topic && (t.subtopic == "*" || t.subtopic == other.subtopic)
}

func isTopicMatch(a string, b string) bool {
	return asTopic(a).isMatch(asTopic(b))
}