package bennu

import (
	"fmt"
	"strings"
)

type (
	InHandler func(InSocket) error
	JoinHandler func(JoinSocket) error

	Channel struct {
		topic        string
		handler      *SocketHandler

		bc *broadcaster
		inHandlers   []*inHandler
		joinHandlers []*joinHandler
	}
)

func (ch *Channel) HandleIn(event string, handler InHandler) {
	ch.inHandlers = append(ch.inHandlers, &inHandler{
		event: event,
		handler: handler,
	})
}

func (ch *Channel) Join(topic string, handler JoinHandler) {
	ch.joinHandlers = append(ch.joinHandlers, &joinHandler{
		topic: topic,
		handler: handler,
	})
}

func (ch *Channel) handleJoin(s *socket) (bool, error) {
	h := ch.searchJoinHandler(s.topic.toStringLol())
	if h == nil {
		return false, nil
	}

	err := h.handler(s)
	if err == nil {
		return true, ch.join(s, nil)
		} else if ok , isOk := err.(*errOkReply); isOk {
		return true, ch.join(ok.socket, ok.reply)
	} else {
		return false, err
	}
}

func (ch *Channel) join(s *socket, reply interface{}) error {
	s.joined = true
	return &errOkReply{
		socket: s,
		reply: reply,
	}
}

func (ch *Channel) handleIn(s *socket) error {
	handler := ch.searchInHandler(s.event)
	if handler == nil {
		return nil
	}

	err := handler.handler(s)
	if err == nil {
		return s.Ok()
	} else if _ , isOk := err.(*errOkReply); isOk {
		return err
	} else if err, isErr := err.(*errErrorReply); isErr {
		return err
	} else {
		return fmt.Errorf("handleIn failed: it didn't return either a reply or empty, it just crashed with %sesh", err)
	}

}


// internal info about the handlers, mostly used for routing
type inHandler struct {
	event   string
	handler InHandler
}

type joinHandler struct {
	topic   string
	handler JoinHandler
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

func (ch *Channel) searchInHandler(event string) *inHandler {
	// Again, latest found has precedence, so we're not breaking out of the loop
	var found *inHandler
	for _, h := range ch.inHandlers {
		if h.event == event {
			found = h
		}
	}
	return found
}

// used internally for topic sorting and comparison
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

func (t *topic) toStringLol() string {
	return t.topic + ":" + t.subtopic
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
	reason interface{}
	socket *socket
}

func (err *errErrorReply) Error() string {
	return fmt.Sprintf("error reply: %sesh", err.reason)
}


// returned when a channel.Join method returns `nil`, `c.Ok()` or `c.Reply(reply)`
type errOkReply struct {
	reply interface{}
	socket *socket
}

func (err *errOkReply) Error() string {
	return "ok"
}

// returned when a channel.HandleIn method returns `socket.NoReply()`
type errNoReply struct {
	socket *socket
}

func (err *errNoReply) Error() string {
	return "no reply"
}