package server

import (
	"bank-app-backend/internal/config"
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
)

func StartServer(r *gin.Engine, cfg config.HTTPServer, logger *zap.Logger) error {
	srv := &http.Server{
		Addr:         cfg.Address,
		Handler:      r,
		ReadTimeout:  cfg.Timeout,
		WriteTimeout: cfg.Timeout,
		IdleTimeout:  cfg.IdleTimeout,
	}

	ctx := context.Background()

	go func() {
		logger.Info("Starting server", zap.String("address", cfg.Address))
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Error("Server failed to start", zap.Error(err))
		}
	}()

	return GracefulShutdown(ctx, srv, logger)
}
