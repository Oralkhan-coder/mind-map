package mongo

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type DbConfig struct {
	Host     string `env:"MONGO_HOST"`
	Port     uint16 `env:"MONGO_PORT"`
	Username string `env:"MONGO_USERNAME"`
	Password string `env:"MONGO_PASSWORD"`
	Database string `env:"MONGO_DATABASE"`
}

type DB struct {
	Client   *mongo.Client
	Database *mongo.Database
}

func NewDB(ctx context.Context, cfg *DbConfig) (*DB, error) {
	uri := fmt.Sprintf(
		"mongodb://%s:%d",
		cfg.Host,
		cfg.Port,
	)

	clientOpts := options.Client().
		ApplyURI(uri).
		SetMaxPoolSize(20).
		SetMinPoolSize(5).
		SetMaxConnIdleTime(5 * time.Minute)

	client, err := mongo.Connect(ctx, clientOpts)
	if err != nil {
		return nil, fmt.Errorf("mongo connect error: %w", err)
	}

	pingCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := client.Ping(pingCtx, nil); err != nil {
		return nil, fmt.Errorf("mongo ping failed: %w", err)
	}

	db := client.Database(cfg.Database)

	return &DB{
		Client:   client,
		Database: db,
	}, nil
}

func (db *DB) Close(ctx context.Context) {
	if err := db.Client.Disconnect(ctx); err != nil {
		log.Fatalf("unable to disconnect from database: %v", err)
	}
}
