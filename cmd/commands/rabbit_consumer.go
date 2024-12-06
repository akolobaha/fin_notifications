package commands

import (
	"context"
	"encoding/json"
	"fin_notifications/internal/config"
	"fin_notifications/internal/db"
	"fin_notifications/internal/email"
	"fin_notifications/internal/entity"
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"log/slog"
)

func ReadFromQueue(ctx context.Context, cfg *config.Config) error {
	mongoClient := db.GetMongoDbConnection(ctx, cfg)
	defer mongoClient.Disconnect(ctx)
	mongoCollectionReports := mongoClient.Database(cfg.MongoDatabase).Collection(cfg.MongoCollection)

	conn, err := amqp.Dial(cfg.GetRabbitDSN())

	if err != nil {
		slog.Error("Failed to connect to RabbitMQ: %s", "error", err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		slog.Error("Failed to open a channel: %s", "error", err)
	}
	defer ch.Close()

	msgsNotification, err := ch.Consume(
		cfg.RabbitNotificationQueueName,
		"",
		true, // auto-ack
		false,
		false,
		false,
		nil,
	)

	msgsEmailConfirm, err := ch.Consume(
		cfg.RabbitEmailConfirmQueueName,
		"",
		true, // auto-ack
		false,
		false,
		false,
		nil,
	)

	for {
		select {
		case msg := <-msgsNotification:
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
			one, err := mongoCollectionReports.InsertOne(ctx, report)
			if err != nil {
				return err
			}

			slog.Info("Отчет успешно сохранен в журнал сообщений", one.InsertedID)

		case msg := <-msgsEmailConfirm:
			var emailConfirm entity.EmailConfirm
			err := json.Unmarshal(msg.Body, &emailConfirm)
			if err != nil {
				return err
			}

			subject := `Подтверждения email на сатйе "Робот для инвестора"`
			text := fmt.Sprintf("Для подтверждения email перейдите по ссылке: %s", emailConfirm.Url)
			err = email.NotifyByEmail(subject, text, []string{emailConfirm.Email})
			if err != nil {
				return err
			}

			slog.Info("Уведомление о подтверждении email успешно отправлено пользователю")

		case <-ctx.Done():
			slog.Info("Сервис обработки данных остановлен")
			return nil
		}
	}
}
