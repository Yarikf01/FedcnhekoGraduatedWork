package healthcheck_test

import (
	"net/http"
	"testing"

	"github.com/appleboy/gofight/v2"
	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"

	"github.com/Yarikf01/graduatedwork/api/utils/healthcheck"
)

func TestExternalHealthCheck(t *testing.T) {
	t.Run("happy path for external", func(t *testing.T) {
		req := gofight.New()
		e := echo.New()

		healthcheck.ExternalAssemble(e, "v1.2.3")

		req.GET("/healthcheck").
			SetDebug(true).
			Run(e, func(resp gofight.HTTPResponse, req gofight.HTTPRequest) {
				assert.Equal(t, http.StatusOK, resp.Code)
			})
	})
}

func TestInternalHealthCheck(t *testing.T) {
	t.Run("happy path for internal", func(t *testing.T) {
		req := gofight.New()
		e := echo.New()

		healthcheck.InternalAssemble(e)

		req.GET("/").
			SetDebug(true).
			Run(e, func(resp gofight.HTTPResponse, req gofight.HTTPRequest) {
				assert.Equal(t, http.StatusOK, resp.Code)
			})
	})
}
