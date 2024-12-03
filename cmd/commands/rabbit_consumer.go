package commands

import (
	"context"
	"encoding/json"
	"fin_notifications/internal/config"
	"fin_notifications/internal/db"
	"fin_notifications/internal/email"
	"fin_notifications/internal/entity"
	amqp "github.com/rabbitmq/amqp091-go"
	"log/slog"
)

func ReadFromQueue(ctx context.Context, cfg *config.Config) error {
	mongoClient := db.GetMongoDbConnection(ctx, cfg)
	defer mongoClient.Disconnect(ctx)
	mongoCollection := mongoClient.Database(cfg.MongoDatabase).Collection(cfg.MongoCollection)

	conn, err := amqp.Dial(cfg.GetRabbitDSN())

	if err != nil {
		slog.Error("Failed to connect to RabbitMQ: %s", err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		slog.Error("Failed to open a channel: %s", err)
	}
	defer ch.Close()

	msgs, err := ch.Consume(
		cfg.RabbitQueueName,
		"",
		true, // auto-ack
		false,
		false,
		false,
		nil,
	)

	one, err := mongoCollection.InsertOne(ctx, entity.TargetUser{})
	if err != nil {
		return err
	}

	slog.Info("сообщение успешно сохранено в журнал сообщений", one.InsertedID)

	for {
		select {
		case msg := <-msgs:
			var targetUser entity.TargetUser
			err := json.Unmarshal(msg.Body, &targetUser)
			if err != nil {
				return err
			}

			subject, text := email.GenerateEmailText(targetUser)
			err = email.NotifyByEmail(subject, text, []string{targetUser.User.Email})
			if err != nil {
				return err
			}

			report := entity.NewReport(targetUser, subject, text)
			one, err := mongoCollection.InsertOne(ctx, report)
			if err != nil {
				return err
			}

			slog.Info("сообщение успешно сохранено в журнал сообщений", one.InsertedID)
		case <-ctx.Done():
			slog.Info("Сервис обработки данных остановлен")
			return nil
		}
	}
}
