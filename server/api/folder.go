package api

import (
	"net/http"
	"github.com/labstack/echo"
	"github.com/ravishi/palantir/server/api/cx"
)

func CreateFolder(c echo.Context) error {
	cx := cx.From(c)

	s, err := cx.Session()
	if err != nil {
		return err
	}

	path := cx.JsonString("path")

	folder, err := s.CreateFolder(path)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusCreated, c.Echo().URI(ShowFolder, folder.ID))
}

func RemoveFolder(c echo.Context) error {
	id, err := ParseInt64(c.Param("id"))
	if err != nil {
		return err
	}

	err = cx.From(c).PanickingSession().RemoveFolder(id)
	if err != nil {
		return err
	}

	return c.NoContent(http.StatusNoContent)
}

func ShowFolder(c echo.Context) error {
	id, err := ParseInt64(c.Param("id"))
	if err != nil {
		return err
	}

	f, err := cx.From(c).PanickingSession().GetFolder(id)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, f)
}

func ListFolders(c echo.Context) error {
	folders, err := cx.From(c).PanickingSession().FindFolders()
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, folders)
}
