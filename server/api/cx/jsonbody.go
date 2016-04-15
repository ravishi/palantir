package cx

import (
	"encoding/json"
	"github.com/labstack/echo"
)

const (
	jsonBodyPrefix = "body.json"
)

/**
	Usage:

		// Considering a POST request with `Content-Type: application/json`
		// and body `{"path": "/tmp", "count": 30}`

		e.Use(cx.JsonBodyMiddleware)

		e.Post("/", func(c echo.Context) error {
			cx := cx.Cx(c)

			p := cx.JsonString("path") // p := "/tmp"
			i := cx.JsonInt("count") // i := 30
		}
 */
func JsonBodyMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		if c.Request().Method() == "POST" {
			if c.Request().Header().Get("Content-Type") == "application/json" {
				var v interface{}
				decoder := json.NewDecoder(c.Request().Body())
				err := decoder.Decode(&v)
				if err == nil {
					newCxt(c, cxPrefix).set(jsonBodyPrefix, v)
				}
				// FIXME silently ignored error should be logged
			}
		}
		return next(c)
	}
}

func (cx *cxt) JsonBody() map[string]interface{} {
	return cx.get(jsonBodyPrefix).(map[string]interface{})
}

func (cx *cxt) Json(key string) interface {} {
	return cx.JsonBody()[key]
}

func (cx *cxt) JsonString(key string) string {
	return cx.Json(key).(string)
}
