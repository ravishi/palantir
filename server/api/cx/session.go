package cx

import (
	"github.com/ravishi/palantir/server/data"
)

const (
	cxSessionPrefix = "session"
)

func (cx *cxt) Session() (data.Session, error) {
	var err error
	var s data.Session
	val := cx.get(cxSessionPrefix)
	if val != nil {
		var ok bool
		if s, ok = val.(data.Session); !ok {
			s = nil
		}
	}
	if s == nil {
		s, err = data.NewSession("./palantir.db")
		if err == nil {
			cx.set(cxSessionPrefix, s)
		}
	}
	return s, err
}

func (cx *cxt) PanickingSession() data.Session {
	s, err := cx.Session()
	if err != nil {
		panic(err)
	}
	return s
}