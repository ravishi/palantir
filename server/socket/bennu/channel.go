package bennu

type (
	Handler func (Socket) error

	JoinHandler func (JoiningSocket) error

	ChannelBuilder interface {
		Join(string, JoinHandler) ChannelBuilder
		HandleIn(string, Handler) ChannelBuilder
		HandleOut(string, Handler) ChannelBuilder

		Build() *Channel
	}

	builder struct {
		ch *Channel
	}

	Channel struct {
		handlers []handler
		joinHandlers []joinHandler

	}

	handler struct {
		event string
		handler Handler
	}

	joinHandler struct {
		topic string
		handler JoinHandler
	}
)

func NewChannelBuilder() ChannelBuilder {
	return &builder{
		ch: &Channel{
			handlers: make([]handler, 0),
			joinHandlers: make([]joinHandler, 0),
		},
	}
}

// Spec:
//   join(topic :: binary, auth_msg :: map, Phoenix.Socket.t) ::
//     {:ok, Phoenix.Socket.t} |
//     {:ok, map, Phoenix.Socket.t} |
//     {:error, map}
//
//
// Special returns:
//
//   JoiningSocket.Ok() -> Like {:ok, Phoenix.Socket.t}
//   JoiningSocket.OkWithReply(interface{}) -> Like {:ok, map, Phoenix.Socket.t}
//   JoiningSocket.Error(interface{}) -> Like {:error, map}
//   error -> same as JoiningSocket.Error(nil)
func (b *builder) Join(topic string, handler JoinHandler) ChannelBuilder {
	b.ch.joinHandlers = append(b.ch.joinHandlers, &joinHandler{
		topic: topic,
		handler: handler,
	})
	return b
}

// Spec:
//
// handle_in(event :: String.t, msg :: map, Phoenix.Socket.t) ::
//   {:noreply, Phoenix.Socket.t} |
//   {:reply, reply, Phoenix.Socket.t} |
//   {:stop, reason :: term, Phoenix.Socket.t} |
//   {:stop, reason :: term, reply, Phoenix.Socket.t}
//
//
// Stopping:
//
//   Stop() ->
//     In this case, the exit won't be logged, there is no restart in transient mode and linked processes won't exit.
//
//   Shutdown(error) ->
//     In this case, the exit won't be logged, there is no restart in transient mode and linked processes exit with the same reason.
//     Error can be nil.
//
//   Shutdown(error) ->
//     In this case, the exit won't be logged, there is no restart in transient mode and linked processes exit with the same reason.
//
//   any other error ->
//     In such cases, the exit will be logged, there are restarts in transient mode, and linked processes exit with the same reason.
func (b *builder) HandleIn(event string, handler Handler) ChannelBuilder {
	b.ch.handlers = append(b.ch.handlers, &handler{
		event: event,
		handler: handler,
	})
	return b
}

func (b *builder) Build() *Channel {
	h := make([]handler, len(b.ch.handlers))
	jh := make([]joinHandler, len(b.ch.joinHandlers))
	copy(h, b.ch.handlers)
	copy(jh, b.ch.joinHandlers)
	return &Channel{
		handlers: h,
		joinHandlers: jh,
	}
}