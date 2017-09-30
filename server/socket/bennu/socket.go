package bennu

type (
	Socket interface {
		Topic() string
		Joined() bool
		Payload() interface{}
		Subtopic() string
	}

	InSocket interface {
		Socket

		//Stop(reason string, reply interface{}) error
		Ok() error
		Reply(reply interface{}) error
		NoReply() error

		Broadcast(string, interface{})
	}


	JoinSocket interface {
		Socket

		Ok() error
		Reply(reply interface{}) error
		Error(reason interface {}) error
	}

	socket struct {
		event string
		topic *topic
		joined bool
		channel *Channel
		payload interface{}
		session *session
	}
)

func (s *socket) Topic() string {
	return s.topic.topic
}

func (s *socket) Joined() bool {
	return s.joined
}

func (s *socket) Payload() interface{} {
	return s.payload
}

func (s *socket) Subtopic() string {
	return s.topic.subtopic
}

func (s *socket) Ok() error {
	return s.Reply(nil)
}

func (s *socket) Reply(reply interface{}) error {
	return &errOkReply{
		reply: reply,
		socket: s,
	}
}

func (s *socket) NoReply() error {
	return &errNoReply{
		socket: s,
	}
}

func (s *socket) Error(reason interface{}) error {
	return &errErrorReply{
		reason: reason,
		socket: s,
	}
}

func (s *socket) Broadcast(event string, payload interface{}) {
	s.session.handler.bc.Publish(&envelope{
		//JoinRef string
		//Ref string
		Topic: s.Topic() + ":" + s.Subtopic(),
		Event: event,
		Payload: payload,
	})
}
