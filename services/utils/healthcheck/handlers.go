package healthcheck

import (
	"net/http"

	"github.com/labstack/echo"
)

func ExternalAssemble(e *echo.Echo, version string) {
	e.GET("/healthcheck", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{
			"status":  "green",
			"version": version,
		})
	})
}

func InternalAssemble(e *echo.Echo) {
	e.GET("/", func(c echo.Context) error {
		return c.NoContent(http.StatusOK)
	})
}
