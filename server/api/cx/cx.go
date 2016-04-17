package cx

import (
	"github.com/labstack/echo"
	"strings"
	"github.com/ravishi/palantir/server/data"
)

const (
	cxPrefix = "cx"
	cxUserDataPrefix = "user.data"
)

type (
	CX interface {
		Get(key string) interface{}
		Set(key string, value interface{})

		JsonBody() map[string]interface{}
		Json(string) interface{}
		JsonString(string) string

		Session() (*data.Session, error)
		PanickingSession() *data.Session
	}

	cxt struct {
		c      echo.Context
		prefix string
	}
)

func From(c echo.Context) CX {
	return newCxt(c, cxPrefix)
}

func (cx *cxt) Get(key string) interface{} {
	return cx.get(cxUserDataPrefix)
}

func (cx *cxt) Set(key string, value interface{}) {
	cx.set(cxUserDataPrefix, value)
}

func newCxt(c echo.Context, namespace string) *cxt {
	return &cxt{c, namespace}
}

func (cx *cxt) get(key string) interface{} {
	return cx.c.Get(cx.path(cx.prefix, key))
}

func (cx *cxt) set(key string, value interface{}) {
	cx.c.Set(cx.path(cx.prefix, key), value)
}

func (cx *cxt) path(s ...string) string {
	return strings.Join(append([]string{cxPrefix}, s...), ".")
}
