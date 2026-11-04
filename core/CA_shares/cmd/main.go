package main

import (
	"CA_shares/config"
	"CA_shares/core/CA_shares"
	"CA_shares/db"
	"CA_shares/logger"
	"CA_shares/pgService"
	transport "CA_shares/transport/http"
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	cfg, err := config.Load("config/config.yaml")
	if err != nil {
		panic(err)
	}

	logger.Init(logger.Config{
		Level:  cfg.Logger.Level,
		Format: cfg.Logger.Format,
	})
	logger.Info("Config successfully parsed", "dbHost", cfg.DB.Host, "dbName", cfg.DB.DBName, "dbSchema", cfg.DB.Schema)

	dbConn, err := db.Connect(cfg.DB)
	if err != nil {
		logger.Fatal("Fail connect database", "error", err.Error())
	}
	defer func() {
		if err := dbConn.Close(); err != nil {
			logger.Error("Failed to close db", "error", err.Error())
		}
	}()
	logger.Info("Connection successfully established", "host", cfg.DB.Host, "port", cfg.DB.Port)

	pg := pgService.NewPgService(dbConn)
	svc := CA_shares.NewReduceService(pg)
	tp := transport.NewTransport(svc)
	handler := tp.Handler()

	srv := &http.Server{
		Addr:    cfg.Server.Addr,
		Handler: handler,
	}

	go func() {
		logger.Info("HTTP server listening on " + cfg.Server.Addr)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Fatal("Fail to start http server", "error", err.Error())
		}
	}()

	sigCtx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	<-sigCtx.Done()
	logger.Info("Shutdown signal received")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer shutdownCancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		logger.Error("Fail to shutdown http server", "error", err.Error())
	} else {
		logger.Info("Server shutdown successfully")
	}

	logger.Info("Shutdown complete")
}
