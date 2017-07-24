package bennu

type (
	Socket interface {
		Topic() string
		Joined() bool
		Payload() interface{}
		Subtopic() string
	}

	JoinSocket interface {
		Socket

		Ok() error
		Nope(reason interface {}) error
		OkReply(reply interface{}) error
	}

	socket struct {
		topic *topic
		joined bool
		channel *Channel
		payload interface{}
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
	return &errOkReply{
		reply: nil,
		socket: s,
	}
}

func (s *socket) OkReply(reply interface{}) error {
	return &errOkReply{
		reply: reply,
		socket: s,
	}
}

func (s *socket) Nope(reply interface{}) error {
	return &errErrorReply{
		reply: reply,
		socket: s,
	}
}
