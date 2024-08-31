package mongoclient

import (
	"context"
	"fmt"
	"light-backend/config"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func Connect() error {
	uri := fmt.Sprintf("mongodb://%s:%s", config.Config("DB_HOST"), config.Config("DB_PORT"))
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(uri))
	if err != nil {
		return err
	}

	DB = client.Database(config.Config("DB_NAME"))

	if DB == nil {
		return fmt.Errorf("database is not selected")
	}
	return nil
}
