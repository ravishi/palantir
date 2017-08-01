package bennu

import (
	"testing"
	"github.com/stretchr/testify/require"
	"fmt"
)


func TestChannelJoin(t *testing.T) {
	plz := require.New(t)

	h := New()

	ch := h.Channel("room:*")

	r := ch.handleJoin("room:lobby", nil)
	plz.Nil(r, fmt.Sprintf("ch.handleJoin() should've returned NotHandled(), but instead returned: %sesh", r))

	ch.Join("room:lobby", func (c JoinSocket) error {
		return c.Reply("hello, there!")
	})

	r = ch.handleJoin("room:lobby", nil)
	_, isOkReply := r.(*errOkReply)

	plz.Equal(isOkReply, true, fmt.Sprintf("ch.handleJoin() should've returned Ok(), but instead returned: %sesh", r))

	ch.Join("room:secret", func (c JoinSocket) error {
		if s, ok := c.Payload().(string); ok && s == "pls let me in!?1!!" {
			return c.Reply("ok u in")
		}
		return c.Error("nope")
	})

	r = ch.handleJoin("room:secret", "pls let me in")
	_, isErrorReply := r.(*errErrorReply)

	plz.Equal(isErrorReply, true, fmt.Sprintf("ch.handleJoin() should've returned ErrorReply(), but instead returned: %sesh", r))

	r = ch.handleJoin("room:secret", "pls let me in!?1!!")
	_, isOkReply = r.(*errOkReply)

	plz.Equal(isOkReply, true, fmt.Sprintf("ch.handleJoin() should've returned OkReply(), but instead returned: %sesh", r))
}