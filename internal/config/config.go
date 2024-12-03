package config

import (
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	SmtpHost        string `env:"SMTP_HOST"`
	SmtpPort        int    `env:"SMTP_PORT"`
	SmtpUsername    string `env:"SMTP_USERNAME"`
	SmtpPassword    string `env:"SMTP_PASSWORD"`
	RabbitHost      string `env:"RABBIT_HOST"`
	RabbitPort      int    `env:"RABBIT_PORT"`
	RabbitUsername  string `env:"RABBIT_USERNAME"`
	RabbitPassword  string `env:"RABBIT_PASSWORD"`
	RabbitQueueName string `env:"RABBIT_QUEUE_NAME"`
	MongoUsername   string `env:"MONGO_USERNAME"`
	MongoPassword   string `env:"MONGO_PASSWORD"`
	MongoHost       string `env:"MONGO_HOST"`
	MongoPort       string `env:"MONGO_PORT"`
	MongoDatabase   string `env:"MONGO_DATABASE"`
	MongoCollection string `env:"MONGO_COLLECTION"`
}

func Parse(s string) (*Config, error) {
	c := &Config{}
	if err := cleanenv.ReadConfig(s, c); err != nil {
		return nil, err
	}

	return c, nil
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
