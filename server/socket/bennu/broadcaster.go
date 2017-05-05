package bennu

type broadcaster struct {
	pub       chan *message
	subs      map[chan<- *message]filterFunc
	subReqs   chan *subRequest
	unsubReqs chan chan<- *message
}

type subRequest struct {
	ch     chan<- *message
	filter filterFunc
}

type filterFunc func(*message) bool

func createBroadcaster() *broadcaster {
	b := &broadcaster{
		pub:       make(chan *message),
		subs:      make(map[chan<- *message]filterFunc),
		subReqs:   make(chan *subRequest),
		unsubReqs: make(chan chan<- *message),
	}
	go b.broadcastForever()
	return b
}

func (b *broadcaster) broadcast(m *message) {
	for ch := range b.subs {
		if filter := b.subs[ch]; filter == nil || filter(m) {
			ch <- m
		}
	}
}

func (b *broadcaster) broadcastForever() {
	defer close(b.pub)
	defer close(b.unsubReqs)
	for {
		select {
		case m := <-b.pub:
			b.broadcast(m)
		case req, ok := <-b.subReqs:
			if ok {
				b.subs[req.ch] = req.filter
			} else {
				return
			}
		case ch := <-b.unsubReqs:
			delete(b.subs, ch)
		}
	}
}

func (b *broadcaster) Publish(m *message) {
	b.pub <- m
}

func noFilter(*message) bool {
	return true
}

func (b *broadcaster) Subscribe(ch chan<- *message) {
	b.FilteredSubscribe(ch, noFilter)
}

func (b *broadcaster) FilteredSubscribe(ch chan<- *message, filter filterFunc) {
	b.subReqs <- &subRequest{ch: ch, filter: filter}
}

func (b *broadcaster) Unsubscribe(ch chan<- *message) {
	b.unsubReqs <- ch
}

func (b *broadcaster) Close() {
	close(b.subReqs)
}
