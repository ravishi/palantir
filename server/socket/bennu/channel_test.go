package bennu

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)


func TestChannelJoin(t *testing.T) {
	plz := require.New(t)

	h := New()

	ch := h.Channel("room:*")

	s := h.newSession(nil, nil)

	e := &envelope{Topic: "room:lobby"}
	r := s.handleJoin(e)
	plz.Nil(r, fmt.Sprintf("ch.handleJoin() should've returned NotHandled(), but instead returned: %s", r))

	ch.Join("room:lobby", func (c JoinSocket) error {
		return c.Reply("hello, there!")
	})

	r = s.handleJoin(e)
	_, isOkReply := r.(*errOkReply)

	plz.Equal(isOkReply, true, fmt.Sprintf("ch.handleJoin() should've returned Ok(), but instead returned: %s", r))

	ch.Join("room:secret", func (c JoinSocket) error {
		if s, ok := c.Payload().(string); ok && s == "pls let me in!?1!!" {
			return c.Reply("ok u in")
		}
		return c.Error("nope")
	})

	e.Topic = "room:secret"
	e.Payload = "pls let me in"
	r = s.handleJoin(e)
	_, isErrorReply := r.(*errErrorReply)

	plz.Equal(isErrorReply, true, fmt.Sprintf("ch.handleJoin() should've returned ErrorReply(), but instead returned: %s", r))

	e.Payload = "pls let me in!?1!!"
	r = s.handleJoin(e)
	_, isOkReply = r.(*errOkReply)

	plz.Equal(isOkReply, true, fmt.Sprintf("ch.handleJoin() should've returned OkReply(), but instead returned: %s", r))
}