package api

import (
	"net/http"

	"github.com/labstack/echo"
)

func Assemble(root *echo.Group, m Manager) {
	h := &handler{
		manager: m,
	}

	g := root.Group("/api")

	g.POST("/upload", h.upload)
	g.GET("/download", h.download)
}

// impl

type handler struct {
	manager Manager
}

func (h *handler) upload(c echo.Context) error {
	ctx := c.Request().Context()

	err := h.manager.Upload(ctx)
	if err != nil {
		return err
	}

	return c.NoContent(http.StatusOK)
}

func (h *handler) download(c echo.Context) error {
	ctx := c.Request().Context()

	err := h.manager.Download(ctx)
	if err != nil {
		return err
	}

	return c.NoContent(http.StatusOK)
}
