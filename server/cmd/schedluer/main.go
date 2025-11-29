package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"schedluer/internal/config"
	"schedluer/internal/container"

	_ "schedluer/docs"
)

// @title           Schedluer API
// @version         1.0
// @description     API для работы с расписанием БГУИР

// @host      localhost:8080
// @BasePath  /api/v1

// @schemes   http https
func main() {
	cfg, err := config.Load()
	if err != nil {
		logrus.Fatalf("Failed to load config: %v", err)
	}

	setupLogger(cfg.Logger.Level)

	ctn, err := container.NewContainer(cfg)
	if err != nil {
		logrus.Fatalf("Failed to create container: %v", err)
	}
	defer func() {
		if err := ctn.Close(); err != nil {
			logrus.Errorf("Failed to close container: %v", err)
		}
	}()

	router := setupRouter(ctn, cfg)

	addr := fmt.Sprintf("%s:%s", cfg.Server.Host, cfg.Server.Port)
	logrus.Infof("Starting server on %s", addr)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := router.Run(addr); err != nil {
			logrus.Fatalf("Failed to start server: %v", err)
		}
	}()

	<-quit
	logrus.Info("Shutting down server...")
}

func setupLogger(level string) {
	logrus.SetFormatter(&logrus.JSONFormatter{})

	switch level {
	case "debug":
		logrus.SetLevel(logrus.DebugLevel)
	case "info":
		logrus.SetLevel(logrus.InfoLevel)
	case "warn":
		logrus.SetLevel(logrus.WarnLevel)
	case "error":
		logrus.SetLevel(logrus.ErrorLevel)
	default:
		logrus.SetLevel(logrus.InfoLevel)
	}
}

func setupRouter(ctn *container.Container, cfg *config.Config) *gin.Engine {
	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowOrigins:     cfg.CORS.AllowedOrigins,
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization"},
		AllowCredentials: true,
		MaxAge:           12 * 60 * 60, // 12 hours
	}))

	router.GET("/health", func(c *gin.Context) {
		if err := ctn.MongoDB.Health(c.Request.Context()); err != nil {
			c.JSON(503, gin.H{
				"status": "unhealthy",
				"error":  err.Error(),
			})
			return
		}

		c.JSON(200, gin.H{
			"status": "ok",
		})
	})

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	ctn.Router.SetupRoutes(router)

	return router
}
