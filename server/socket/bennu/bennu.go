package bennu

type (
	Beenu struct {
		h *SocketHandler
	}

	SocketHandler struct {
		channels []channelEntry
	}

	channelEntry struct {
		channel *Channel
		topicPattern string
	}
)

func New() *SocketHandler {
	return &SocketHandler{}
}

// TODO? add an options parameter to restrict channel's transport adapters
func (b *Beenu) Channel(topicPattern string, ch *Channel) SocketHandler {
	b.h.channels = append(b.h.channels, &channelEntry{
		channel: ch,
		topicPattern: topicPattern,
	})
	return b
}
