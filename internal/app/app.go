package app

import (
	_ "bank-app-backend/docs"
	"bank-app-backend/internal/config"
	"bank-app-backend/internal/controllers/http"
	"bank-app-backend/internal/controllers/middleware"
	"bank-app-backend/internal/db"
	lib "bank-app-backend/internal/lib/logger"
	redis "bank-app-backend/internal/lib/redis"
	"bank-app-backend/internal/repository"
	"bank-app-backend/internal/server"
	"bank-app-backend/internal/services"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.uber.org/zap"
	"gorm.io/gorm"
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

	loggerZap := setupLogger(cfg.Env)
	defer loggerZap.Sync()

	database := setupDatabase(cfg.Storage, loggerZap)

	redisClient := redis.NewRedisClient(cfg.Redis.Address, cfg.Redis.Password, cfg.Redis.DB)

	loggerZap.Info("Starting bank-app",
		zap.String("env", cfg.Env),
	)

	loggerZap.Debug("debug messages are enabled")

	r := gin.Default()
	r.Use(middleware.ZapLoggerMiddleware())
	r.Use(middleware.PrometheusMiddleware(requestCount))

	authRepo := repository.NewUsersRepository(database, redisClient)
	authorizationService := services.NewAuthService(authRepo, redisClient)
	authHandlers := http.NewAuthHandler(authorizationService)

	usersRepo := repository.NewUsersRepository(database, redisClient)
	usersService := services.NewUsersService(usersRepo)
	usersHandlers := http.NewUsersHandler(usersService)

	accountsRepo := repository.NewAccountsRepository(database)
	accountsService := services.NewAccountsService(accountsRepo)
	accountsHandlers := http.NewAccountsHandler(accountsService)

	auth := r.Group("/auth")
	auth.Use(middleware.JWTAuthMiddleware([]byte("jwt_access_secret")))

	{
		auth.GET("/me", usersHandlers.Me)
		auth.GET("/accounts", accountsHandlers.GetAllByUser)
		auth.POST("/accounts", accountsHandlers.Create)
		auth.GET("/accounts/:id", accountsHandlers.GetByID)
		auth.PATCH("/accounts/:id", accountsHandlers.CloseAccount)
	}

	users := r.Group("/users")
	{
		users.GET("", usersHandlers.GetAll)
		users.PATCH("/:id", usersHandlers.Update)
	}

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	r.GET("/metrics", gin.WrapH(promhttp.Handler()))

	r.POST("/register", authHandlers.Register)
	r.POST("/login", authHandlers.Login)
	r.POST("/refresh", authHandlers.Refresh)
	r.POST("/logout", authHandlers.Logout)

	if err := server.StartServer(r, cfg.HTTPServer, loggerZap); err != nil {
		loggerZap.Error("Server failed", zap.Error(err))
	}
}

func setupLogger(env string) *zap.Logger {
	lib.InitLogger(env)
	return lib.Log
}

func setupDatabase(path string, logger *zap.Logger) *gorm.DB {
	database, err := db.InitDB(path)
	if err != nil {
		logger.Fatal("Failed to connect to database", zap.Error(err))
	}
	return database
}
