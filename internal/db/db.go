package db

import (
	"context"
	"fin_notifications/internal/config"
	"fmt"
	_ "github.com/lib/pq"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log/slog"
)

func GetMongoDbConnection(ctx context.Context, cfg *config.Config) *mongo.Client {
	clientOptions := options.Client().ApplyURI(cfg.GetMongoDSN())

	client, err := mongo.Connect(ctx, clientOptions)

	if err != nil {
		slog.Error(err.Error())
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		slog.Error(err.Error())
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		slog.Error(err.Error())
	}
	fmt.Println("Connected to MongoDB!")

	return client
}
