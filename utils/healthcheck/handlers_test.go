package healthcheck_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/appleboy/gofight/v2"
	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"

	"github.com/Yarikf01/bbwy/utils/healthcheck"
	log2 "github.com/Yarikf01/bbwy/utils/log"
)


func TestInternalHealthCheck(t *testing.T) {
	t.Run("happy path for internal", func(t *testing.T) {
		req := gofight.New()
		e := echo.New()
		logger := log2.FromContext(context.TODO())

		healthcheck.InternalAssemble(e, logger)

		req.GET("/").
			SetDebug(true).
			Run(e, func(resp gofight.HTTPResponse, req gofight.HTTPRequest) {
				assert.Equal(t, http.StatusOK, resp.Code)
			})
	})
}
