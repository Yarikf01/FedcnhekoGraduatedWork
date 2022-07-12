package healthcheck

import (
	"net/http"

	"github.com/labstack/echo"
	"go.uber.org/zap"
)

func InternalAssemble(e *echo.Echo, logger *zap.SugaredLogger) {
	e.GET("/", func(c echo.Context) error {
		logger.Info("BIG BROTHER IS WATCHING YOU (づ｡◕‿‿◕｡)づ) (づ｡◕‿‿◕｡)づ) (づ｡◕‿‿◕｡)づ)")
		return c.JSON(http.StatusOK, map[string]string{"status": "green"})
	})
}
