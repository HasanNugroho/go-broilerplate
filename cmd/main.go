package main

import (
	"net"
	"net/http"

	config "github.com/HasanNugroho/starter-golang/internal/configs"
	"github.com/HasanNugroho/starter-golang/internal/middleware"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

const (
	ProductionEnv = "production"
)

func main() {
	// Initialize configuration
	cfg, err := config.InitConfig()
	if err != nil {
		log.Fatal().Msg("❌ Failed to initialize config: " + err.Error())
	}

	// Set production mode if applicable
	if cfg.AppEnv == ProductionEnv {
		gin.SetMode(gin.ReleaseMode)
	}

	// Initialize Logger
	config.InitLogger(cfg)

	// Initialize Gin router with middlewares
	r := gin.Default()
	r.Use(middleware.SetCORS(cfg), middleware.SecurityMiddleware(cfg))

	// Initialize RDBMS if enabled
	if cfg.Database.Enabled {
		db, err := config.InitDB(&cfg.Database)
		if err != nil {
			config.Logger.Fatal().Msg(err.Error())
			panic(1)
		}
		defer config.ShutdownDB(db) // Ensure database is closed on exit
	}

	// Initialize Redis if enabled
	if cfg.Database.Enabled {
		var err error
		redisClient, err := config.InitRedis(&cfg.Redis)
		if err != nil {
			config.Logger.Fatal().Msg(err.Error())
			panic(1)
		}
		// config.Database.Redis.Client = redisClient
		defer config.ShutdownRedis(redisClient)
	}

	// Initialize Rate Limiter if enabled
	if cfg.Security.RateLimit != "" {
		limiter, err := config.InitRateLimiter(cfg, cfg.Security.RateLimit, cfg.Security.TrustedPlatform)
		if err != nil {
			config.Logger.Fatal().Msg(err.Error())
			panic(1)
		}
		r.Use(middleware.RateLimit(limiter))
	}

	// Define test route
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "pong"})
	})

	// Start server
	serverAddr := net.JoinHostPort(cfg.Server.ServerHost, cfg.Server.ServerPort)
	// config.Logger.Info().Msg("🚀 Server running at " + serverAddr)
	err = r.Run(serverAddr)
	if err != nil {
		log.Fatal().Msg(err.Error())
	}
}
