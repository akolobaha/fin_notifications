package config

import (
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
	"log/slog"
)

type Config struct {
	SmtpHost                    string `env:"SMTP_HOST"`
	SmtpPort                    int    `env:"SMTP_PORT"`
	SmtpUsername                string `env:"SMTP_USERNAME"`
	SmtpPassword                string `env:"SMTP_PASSWORD"`
	RabbitHost                  string `env:"RABBIT_HOST"`
	RabbitPort                  int    `env:"RABBIT_PORT"`
	RabbitUsername              string `env:"RABBIT_USERNAME"`
	RabbitPassword              string `env:"RABBIT_PASSWORD"`
	RabbitNotificationQueueName string `env:"RABBIT_NOTIFICATION_QUEUE_NAME"`
	RabbitEmailConfirmQueueName string `env:"RABBIT_EMAIL_CONFIRM_QUEUE_NAME"`
	MongoUsername               string `env:"MONGO_USERNAME"`
	MongoPassword               string `env:"MONGO_PASSWORD"`
	MongoHost                   string `env:"MONGO_HOST"`
	MongoPort                   string `env:"MONGO_PORT"`
	MongoDatabase               string `env:"MONGO_DATABASE"`
	MongoCollection             string `env:"MONGO_COLLECTION"`
	LogLevel                    string `env:"LOG_LEVEL"`
}

func Parse(s string) (*Config, error) {
	c := &Config{}
	if err := cleanenv.ReadConfig(s, c); err != nil {
		return nil, err
	}
	setLogLevel(c.LogLevel)
	return c, nil
}

func setLogLevel(level string) {
	switch level {
	case "debug":
		slog.SetLogLoggerLevel(-4)
	case "info":
		slog.SetLogLoggerLevel(0)
	case "warn":
		slog.SetLogLoggerLevel(4)
	case "error":
		slog.SetLogLoggerLevel(8)
	default:
		slog.SetLogLoggerLevel(4)
	}
}

func (cfg *Config) GetRabbitDSN() string {
	return fmt.Sprintf(
		"amqp://%s:%s@%s:%d/", cfg.RabbitUsername, cfg.RabbitPassword, cfg.RabbitHost, cfg.RabbitPort,
	)
}

func (cfg *Config) GetMongoDSN() string {
	return fmt.Sprintf(
		"mongodb://%s:%s@%s:%s/admin", cfg.MongoUsername, cfg.MongoPassword, cfg.MongoHost, cfg.MongoPort)
}
