package shutdown

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/labstack/echo"

	"github.com/Yarikf01/graduatedwork/services/utils"
)

func Graceful(ctx context.Context, e *echo.Echo, cleanup func()) {
	// listen on application shutdown signals.
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	log.FromContext(ctx).Infof("service received shutdown signal: %v", <-quit)

	// allow live connections a set period of time to complete.
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		log.WithError(log.FromContext(ctx), err).Error("unable to shutdown the server")
	}
	cleanup()
}
