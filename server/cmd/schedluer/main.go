package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"schedluer/internal/config"
	"schedluer/internal/container"
)

func main() {
	// Загружаем конфигурацию
	cfg, err := config.Load()
	if err != nil {
		logrus.Fatalf("Failed to load config: %v", err)
	}

	// Настраиваем логгер
	setupLogger(cfg.Logger.Level)

	// Создаем DI-контейнер
	ctn, err := container.NewContainer(cfg)
	if err != nil {
		logrus.Fatalf("Failed to create container: %v", err)
	}
	defer func() {
		if err := ctn.Close(); err != nil {
			logrus.Errorf("Failed to close container: %v", err)
		}
	}()

	// Инициализируем HTTP сервер
	router := setupRouter(ctn)

	// Запускаем сервер
	addr := fmt.Sprintf("%s:%s", cfg.Server.Host, cfg.Server.Port)
	logrus.Infof("Starting server on %s", addr)

	// Graceful shutdown
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

func setupRouter(ctn *container.Container) *gin.Engine {
	router := gin.Default()

	// Health check
	router.GET("/health", func(c *gin.Context) {
		// Проверяем состояние MongoDB
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

	// Настраиваем API роуты
	ctn.Router.SetupRoutes(router)

	return router
}
