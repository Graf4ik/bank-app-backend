package server

import (
	"context"
	"go.uber.org/zap"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func GracefulShutdown(ctx context.Context, srv *http.Server, logger *zap.Logger) error {
	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, syscall.SIGINT, syscall.SIGTERM)

	<-sigint

	logger.Info("Shutting down gracefully...")

	ctxTimeout, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctxTimeout); err != nil {
		logger.Error("Server shutdown failed", zap.Error(err))
		return err
	}

	logger.Info("Server stopped gracefully")
	return nil
}
