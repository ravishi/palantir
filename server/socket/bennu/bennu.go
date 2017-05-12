package bennu

type (
	SocketHandler struct {
		channels []*Channel
	}
)

func New() *SocketHandler {
	return &SocketHandler{
		channels: make([]*Channel, 0),
	}
}

func (h *SocketHandler) Channel(topic string) *Channel {
	ch := &Channel{
		topic: topic,
		handler: h,

		joinHandlers: make([]*joinHandler, 0),
	}

	h.channels = append(h.channels, ch)

	return ch
}
