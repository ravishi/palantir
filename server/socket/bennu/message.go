package bennu

type (
	message struct {
		// The topic or topic:subtopic pair namespace, for example "messages", "messages:123"
		topic string

		// The event name, for example "beenu_join"
		event string

		payload interface{}

		// TODO ref
	}
)

