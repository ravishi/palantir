package bennu

import (
	"strings"
	"fmt"
)

type (
	Handler func(JoinSocket) error

	Channel struct {
		topic string
		handler *SocketHandler
		joinHandlers []*joinHandler
	}
)

func (ch *Channel) Join(topic string, handler Handler) {
	ch.joinHandlers = append(ch.joinHandlers, &joinHandler{
		topic: topic,
		handler: handler,
	})
}

func (ch *Channel) handleJoin(topic string, payload interface{}) error {
	handler := ch.searchJoinHandler(topic)
	if handler == nil {
		return nil
	}

	s := &socket{
		topic: asTopic(topic),
		channel: ch,
		payload: payload,
	}

	err := handler.handler(s)
	if err == nil {
		return ch.join(s, nil)
	} else if ok , isOk := err.(*errOkReply); isOk {
		return ch.join(ok.socket, ok.reply)
	} else if err, isErr := err.(*errErrorReply); isErr {
		return err
	} else {
		return fmt.Errorf("join failed: it didn't return either Ok() or Error(), it just crashed with %s", err)
	}
}

func (ch *Channel) join(s *socket, reply interface{}) error {
	s.joined = true

	// TODO subscribe to the broadcasting

	return &errOkReply{
		socket: s,
		reply: reply,
	}
}


// internal info about the handler, mostly used for routing
type joinHandler struct {
	topic string
	handler Handler
}

func (ch *Channel) searchJoinHandler(topic string) *joinHandler {
	// XXX least found is returned, and note that it can still be nil if nothing is found
	// TODO make "*" a special case, so that order doesn't matter when '*' is used after everything else, like: Join("room:private", ...); Join("room:*", ...);
	// NOTE that it would also make the interface more complex to explain and implement
	var found *joinHandler
	for _, h := range ch.joinHandlers {
		if isTopicMatch(h.topic, topic) {
			found = h
		}
	}
	return found
}

// used internally for topic sorting comparison
type topic struct {
	topic string
	subtopic string
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


// returned when a topic is not handled by a channel.
type errNotHandled struct {
	topic string
}

func (err *errNotHandled) Error() string {
	return "topic is not handled by this channel, topic = " + err.topic
}


// returned when a channel.Join method returns c.Error() or c.ErrorReply(reply)
type errErrorReply struct {
	reply interface{}
	socket *socket
}

func (err *errErrorReply) Error() string {
	return fmt.Sprintf("error reply: %s", err.reply)
}


// returned when a channel.Join method returns c.Ok() or c.OkReply(reply)
type errOkReply struct {
	reply interface{}
	socket *socket
}

func (err *errOkReply) Error() string {
	return "ok"
}