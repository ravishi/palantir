package cx

import (
	"github.com/ravishi/palantir/server/data"
)

const (
	cxSessionPrefix = "session"
)

func (cx *cxt) Session() (s *data.Session, err error) {
	val := cx.get(cxSessionPrefix)
	if val != nil {
		s = val.(*data.Session)
	}
	if s == nil {
		s, err = data.NewSession("./palantir.db")
		if err == nil {
			cx.set(cxSessionPrefix, s)
		}
	}
	return s, err
}

func (cx *cxt) PanickingSession() *data.Session {
	s, err := cx.Session()
	if err != nil {
		panic(err)
	}
	return s
}