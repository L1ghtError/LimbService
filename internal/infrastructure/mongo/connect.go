package mongoclient

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/gridfs"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Config struct {
	Host     string
	Port     string
	User     string
	Password string
	Database string
}

type Connection struct {
	Client   *mongo.Client
	Database *mongo.Database
	Bucket   *gridfs.Bucket
}

func NewConnect(ctx context.Context, cfg Config) (*Connection, error) {
	uri := fmt.Sprintf("mongodb://%s:%s", cfg.Host, cfg.Port)
	clientOpts := options.Client().ApplyURI(uri)

	// Add credentials if provided
	if cfg.User != "" && cfg.Password != "" {
		clientOpts.SetAuth(options.Credential{
			Username: cfg.User,
			Password: cfg.Password,
		})
	}

	// Connect
	client, err := mongo.Connect(ctx, clientOpts)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB: %w", err)
	}

	// Ping
	if err := client.Ping(ctx, nil); err != nil {
		return nil, fmt.Errorf("failed to ping MongoDB: %w", err)
	}

	db := client.Database(cfg.Database)
	bucket, err := gridfs.NewBucket(db)
	if err != nil {
		return nil, fmt.Errorf("gridfs init error: %w", err)
	}

	return &Connection{
		Client:   client,
		Database: db,
		Bucket:   bucket,
	}, nil
}
