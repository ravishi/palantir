package bennu

type broadcaster struct {
	pub       chan *envelope
	subs      map[chan<- *envelope]filterFunc
	subReqs   chan *subRequest
	unsubReqs chan chan<- *envelope
}

type subRequest struct {
	ch     chan<- *envelope
	filter filterFunc
}

type filterFunc func(*envelope) bool

func newBroadcaster() *broadcaster {
	b := &broadcaster{
		pub:       make(chan *envelope),
		subs:      make(map[chan<- *envelope]filterFunc),
		subReqs:   make(chan *subRequest),
		unsubReqs: make(chan chan<- *envelope),
	}
	go b.broadcastForever()
	return b
}

func (b *broadcaster) broadcast(m *envelope) {
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
		case req, ok := <-b.subReqs:
			if ok {
				b.subs[req.ch] = req.filter
			} else {
				return
			}
		case ch := <-b.unsubReqs:
			delete(b.subs, ch)
		case m := <-b.pub:
			b.broadcast(m)
		}
	}
}

func (b *broadcaster) Publish(m *envelope) {
	b.pub <- m
}

func noFilter(*envelope) bool {
	return true
}

func (b *broadcaster) Subscribe(ch chan<- *envelope) {
	b.FilteredSubscribe(ch, noFilter)
}

func (b *broadcaster) FilteredSubscribe(ch chan<- *envelope, filter filterFunc) {
	b.subReqs <- &subRequest{ch: ch, filter: filter}
}

func (b *broadcaster) Unsubscribe(ch chan<- *envelope) {
	b.unsubReqs <- ch
}

func (b *broadcaster) Close() {
	close(b.subReqs)
}