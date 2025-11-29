package database

import (
	"context"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type MongoDB struct {
	Client   *mongo.Client
	Database *mongo.Database
}

type Config struct {
	URI      string
	Database string
	Timeout  time.Duration
}

func NewMongoDB(cfg Config) (*MongoDB, error) {
	logrus.Infof("Connecting to MongoDB: %s", maskURI(cfg.URI))

	clientOptions := options.Client().
		ApplyURI(cfg.URI).
		SetMaxPoolSize(100).
		SetMinPoolSize(10).
		SetMaxConnIdleTime(30 * time.Second).
		SetServerSelectionTimeout(30 * time.Second).
		SetConnectTimeout(30 * time.Second)

	client, err := mongo.Connect(clientOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB: %w", err)
	}

	// Используем отдельный контекст с увеличенным таймаутом для Ping
	pingCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	logrus.Info("Pinging MongoDB...")
	if err := client.Ping(pingCtx, nil); err != nil {
		client.Disconnect(context.Background())
		return nil, fmt.Errorf("failed to ping MongoDB: %w", err)
	}

	logrus.Info("Successfully connected to MongoDB")

	db := client.Database(cfg.Database)

	return &MongoDB{
		Client:   client,
		Database: db,
	}, nil
}

func (m *MongoDB) Close(ctx context.Context) error {
	if m.Client != nil {
		if err := m.Client.Disconnect(ctx); err != nil {
			return fmt.Errorf("failed to disconnect from MongoDB: %w", err)
		}
		logrus.Info("MongoDB connection closed")
	}
	return nil
}

func (m *MongoDB) Health(ctx context.Context) error {
	if m.Client == nil {
		return fmt.Errorf("MongoDB client is nil")
	}
	return m.Client.Ping(ctx, nil)
}

// maskURI скрывает пароль в URI для логирования
func maskURI(uri string) string {
	// Простая маскировка пароля в connection string
	// Формат: mongodb+srv://user:password@host
	if idx := len("mongodb+srv://"); len(uri) > idx {
		start := idx
		if atIdx := findAt(uri, start); atIdx > start {
			if colonIdx := findColon(uri, start, atIdx); colonIdx > start {
				return uri[:colonIdx+1] + "***" + uri[atIdx:]
			}
		}
	}
	return uri
}

func findAt(s string, start int) int {
	for i := start; i < len(s); i++ {
		if s[i] == '@' {
			return i
		}
	}
	return -1
}

func findColon(s string, start, end int) int {
	for i := start; i < end && i < len(s); i++ {
		if s[i] == ':' {
			return i
		}
	}
	return -1
}
