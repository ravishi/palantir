package api

import (
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/ravishi/palantir/server/api/cx"
)

func Echo() *echo.Echo {
	e := echo.New()

	e.Pre(middleware.RemoveTrailingSlash())

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	// This middleware allows us to use Cx(c).Json*(key)
	// to get JSON data from Request's body.
	e.Use(cx.JsonBodyMiddleware)

	library := e.Group("/library")

	folder := library.Group("/folders")
	folder.Get("", ListFolders)
	folder.Get("/:id", ShowFolder)
	folder.Delete("/:id", RemoveFolder)
	folder.Post("/new", CreateFolder)

	return e
}
