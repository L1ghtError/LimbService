package mongoclient

import (
	"context"
	"fmt"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// TODO: Excldude usage of global varible, use functoin instead
// TODO: Let caller pass timeout single varible that represents connection timeout
// TODO Refactor: instead of using Getenv we should rely on config.Get, but for now its low priority
func Connect() error {
	uri := fmt.Sprintf("mongodb://%s:%s", os.Getenv("DB_HOST"), os.Getenv("DB_PORT"))
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(uri))
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err = client.Ping(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to ping MongoDB %w", err)
	}
	DB = client.Database(os.Getenv("DB_NAME"))

	if DB == nil {
		return fmt.Errorf("database is not selected")
	}
	return nil
}
