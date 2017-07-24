package server

import (
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/ravishi/palantir/server/socket"
	"fmt"
)

type Server interface {
	Start(address string) error
}

type ServerConfig struct {
	Debug bool
}

func New(config *ServerConfig) Server {
	e := echo.New()

	e.Debug = config.Debug

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// XXX For now...
	e.Use(middleware.CORS())

	s := socket.ChatSocket()


	e.GET("/ws/websocket", func(c echo.Context) error {
		err := s.Handle(c.Response().Writer, c.Request())
		if err != nil {
			fmt.Println("ERROR at e.GET:", err)
		}
		return err
	})

	return e
}
