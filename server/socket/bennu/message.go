package bennu

type (
	message struct {
		topic string
		event string
		payload interface{}
	}
)