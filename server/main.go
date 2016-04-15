package main

import (
	"github.com/labstack/echo/engine/standard"
	"github.com/ravishi/palantir/server/api"
)

func main() {
	e := api.Echo()
	e.Run(standard.New(":5000"))
}
