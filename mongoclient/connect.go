package mongoclient

import (
	"context"
	"fmt"
	"light-backend/config"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// TODO: Excldude usage of global valibe, use functoin instead
// TODO: Create single varible that represents connection timeout
func Connect() error {
	uri := fmt.Sprintf("mongodb://%s:%s", config.Config("DB_HOST"), config.Config("DB_PORT"))
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
	DB = client.Database(config.Config("DB_NAME"))

	if DB == nil {
		return fmt.Errorf("database is not selected")
	}
	return nil
}
