package app

import (
	_ "bank-app-backend/docs"
	"bank-app-backend/internal/config"
	"bank-app-backend/internal/controllers/http"
	"bank-app-backend/internal/controllers/middleware"
	"bank-app-backend/internal/db"
	"bank-app-backend/internal/lib/logger"
	"bank-app-backend/internal/repository"
	"bank-app-backend/internal/server"
	"bank-app-backend/internal/services/authService"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	swaggerFiles "github.com/swaggo/files"
	"github.com/swaggo/gin-swagger"
	"go.uber.org/zap"
)

var (
	requestCount = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "status"},
	)
)

func init() {
	prometheus.MustRegister(requestCount)
}

func Run() {
	cfg := config.MustLoad()

	logger.InitLogger(cfg.Env)
	defer logger.Log.Sync()

	logger.Log.Info("Starting bank-app",
		zap.String("env", cfg.Env),
	)

	logger.Log.Debug("debug messages are enabled")

	database, err := db.InitDB(cfg.Storage)
	if err != nil {
		logger.Log.Error("Failed to connect to database", zap.Error(err))
	}

	r := gin.Default()
	r.Use(middleware.ZapLoggerMiddleware())

	auth := r.Group("/auth")
	auth.Use(middleware.JWTAuthMiddleware([]byte("jwt_access_secret")))

	authRepo := repository.NewAuthRepository(database)
	accountService := authService.NewAuthService(authRepo)
	authHandlers := http.NewAuthHandler(accountService)

	auth.GET("/me", authHandlers.Me)
	r.GET("/metrics", gin.WrapH(promhttp.Handler()))

	r.Use(func(c *gin.Context) {
		c.Next()

		statusCode := c.Writer.Status()
		method := c.Request.Method

		requestCount.WithLabelValues(method, fmt.Sprintf("%d", statusCode)).Inc()
	})

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r.POST("/register", authHandlers.Register)
	r.POST("/login", authHandlers.Login)
	r.POST("/refresh", authHandlers.Refresh)
	r.POST("/logout", authHandlers.Logout)

	if err := server.StartServer(r, cfg.HTTPServer, logger.Log); err != nil {
		logger.Log.Error("Server failed", zap.Error(err))
	}
}
