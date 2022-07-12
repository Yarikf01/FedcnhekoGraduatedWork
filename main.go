package main

import (
	"context"
	_ "net/http/pprof"

	"github.com/labstack/echo"
	em "github.com/labstack/echo/middleware"

	"github.com/Yarikf01/bbwy/api"
	"github.com/Yarikf01/bbwy/config"
	"github.com/Yarikf01/bbwy/utils/healthcheck"
	"github.com/Yarikf01/bbwy/utils/log"
)

func main() {
	cfg, err := config.BbwyConfig()
	if err != nil {
		panic(err)
	}

	log.InitLog(cfg.Debug)

	ctx := context.Background()
	logger := log.FromContext(ctx)

	apiManager := api.NewManager(api.Config{})

	e := echo.New()
	e.Use(em.Recover())
	root := e.Group("/api/v1")

	api.Assemble(root, apiManager)
	healthcheck.InternalAssemble(e, logger)

	logger.Fatal(e.Start(":8080"))
}
