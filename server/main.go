package main

import (
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/ravishi/palantir/server/socket"
)

func main() {

	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	s := socket.HelloSocket()

	e.GET("/ws", func(c echo.Context) error {
		return s.Handle(c.Response(), c.Request())
	})

	e.Logger.Fatal(e.Start(":8080"))
}
