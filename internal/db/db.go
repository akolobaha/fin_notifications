package db

import (
	"context"
	"fin_notifications/internal/config"
	"fin_notifications/internal/log"
	_ "github.com/lib/pq"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func GetMongoDbConnection(ctx context.Context, cfg *config.Config) *mongo.Client {
	clientOptions := options.Client().ApplyURI(cfg.GetMongoDSN())

	client, err := mongo.Connect(ctx, clientOptions)

	if err != nil {
		log.Error("mongo connect: ", err)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Error("mongo connect ping: ", err)
	}

	return client
}
