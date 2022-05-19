package main

import (
	"context"
	"html/template"
	"net/http"
	_ "net/http/pprof"

	"github.com/labstack/echo"
	em "github.com/labstack/echo/middleware"
	"github.com/robfig/cron"

	"github.com/Yarikf01/graduatedwork/api/repo"
	"github.com/Yarikf01/graduatedwork/api/stats"
	log "github.com/Yarikf01/graduatedwork/api/utils"
	"github.com/Yarikf01/graduatedwork/api/utils/healthcheck"
	"github.com/Yarikf01/graduatedwork/job"
	"github.com/Yarikf01/graduatedwork/metric/business"
)

const maxDBConn = 2

func main() {
	cfg, err := AdminConfig()
	if err != nil {
		log.L.Fatalf("failed to read config, %v", err)
	}

	log.InitLog(cfg.Debug)

	ctx := context.Background()
	logger := log.FromContext(ctx)

	logger.With(
		"debug", cfg.Debug,
	).Info("Starting Admin service")

	userInfoTemplate := template.Must(template.ParseGlob("views/*.html")).Lookup("user_info.html")
	if userInfoTemplate == nil {
		logger.Fatal("user_info.html template does not exist")
	}

	db, err := repo.NewDB(ctx, cfg.DBConnString, maxDBConn)
	if err != nil {
		log.WithError(logger, err).Fatal("failed to init database")
	}

	// start web server
	e := echo.New()
	e.Use(em.Recover())

	root := e.Group(stats.Prefix)

	// pprof endpoints
	e.GET("/debug/pprof/*", echo.WrapHandler(http.DefaultServeMux))

	healthcheck.InternalAssemble(e)

	// stats
	statsConfig := stats.Config{StatsDB: db}
	statsManager := stats.NewManager(statsConfig)
	stats.Assemble(root, statsManager)

	if !cfg.Debug {
		c := cron.New()
		metricWriter, metricClose := business.NewMetricWriter(business.MetricConfig{
			ServerURL: cfg.BusinessMetricServer,
			Token:     cfg.BusinessMetricToken,
			Org:       cfg.BusinessMetricOrg,
			Bucket:    cfg.BusinessMetricBucket,
		})
		defer metricClose()
		err = c.AddFunc("@every 30m", func() { job.GetStatsJob(ctx, statsManager, metricWriter) })
		if err != nil {
			log.WithError(logger, err).Fatal("failed to run GetStats cron job")
		}

	}

	// start service
	go func() {
		logger.Info("Admin service started")
		if err := e.Start(cfg.Port); err != nil && err != http.ErrServerClosed {
			log.WithError(logger, err).Fatal("shutting down the server")
		}
	}()

}
