package commands

import (
	"context"
	"encoding/json"
	"fin_notifications/internal/config"
	"fin_notifications/internal/db"
	"fin_notifications/internal/email"
	"fin_notifications/internal/entity"
	"fin_notifications/internal/log"
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"time"
)

func ReadFromQueue(ctx context.Context, cfg *config.Config) error {
	mongoClient := db.GetMongoDbConnection(ctx, cfg)
	defer mongoClient.Disconnect(ctx)
	mongoCollectionReports := mongoClient.Database(cfg.MongoDatabase).Collection(cfg.MongoCollection)

	conn, err := amqp.Dial(cfg.GetRabbitDSN())
	if err != nil {
		log.Error("Failed to connect to RabbitMQ: ", err)
		return err
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Error("Failed to open a channel: ", err)
		return err
	}
	defer ch.Close()

	msgsNotification, err := ch.Consume(
		cfg.RabbitNotificationQueueName,
		"",
		false,
		false,
		false,
		false,
		nil,
	)

	msgsEmailConfirm, err := ch.Consume(
		cfg.RabbitEmailConfirmQueueName,
		"",
		false,
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

			log.Info(fmt.Sprintf("Отчет успешно сохранен в журнал сообщений %s", one.InsertedID))

			// Подтверждение сообщения после успешной обработки
			if err := msg.Ack(false); err != nil {
				log.Error("Failed to ack message: %s", err)
				return err
			}

		case msg := <-msgsEmailConfirm:
			var emailConfirm entity.EmailConfirm
			err := json.Unmarshal(msg.Body, &emailConfirm)
			if err != nil {
				return err
			}

			subject := "Подтверждение email на сайте 'Робот для инвестора'"
			text := fmt.Sprintf("Для подтверждения email перейдите по ссылке: %s", emailConfirm.Url)
			err = email.NotifyByEmail(subject, text, []string{emailConfirm.Email})
			if err != nil {
				return err
			}

			log.Info("Уведомление о подтверждении email успешно отправлено пользователю")

			// Подтверждение сообщения после успешной обработки
			if err := msg.Ack(false); err != nil {
				log.Error("Failed to ack message: %s", err)
				return err
			}

		case <-ctx.Done():
			log.Info("Сервис обработки данных остановлен")
			time.Sleep(5 * time.Second)
			return nil
		}
	}
}
