package stats

import (
	"net/http"

	"github.com/labstack/echo"

	log "github.com/Yarikf01/graduatedwork/api/utils"
)

func Assemble(root *echo.Group, m Manager) {
	h := &handler{
		manager: m,
	}

	g := root.Group("/stats")

	g.GET("", h.getStats)
}

// impl

type handler struct {
	manager Manager
}

func (h *handler) getStats(c echo.Context) error {
	ctx := c.Request().Context()
	stats, err := h.manager.GetStats(ctx)
	if err != nil {
		return log.LogicErr(err)
	}

	return c.JSON(http.StatusOK, stats)
}
