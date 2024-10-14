package config

import (
	"log"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

type Config struct {
	Port           int           `env:"PORT" env-default:"8080"`
	Env            string        `env:"ENV" env-default:"local"`
	LogLevel       string        `env:"LOG_LEVEL" env-default:"info"`
	IdleTimeout    time.Duration `env:"IDLE_TIMEOUT" env-default:"15s"`
	RequestTimeout time.Duration `env:"REQ_TIMEOUT" env-default:"5s"`

	AuthServiceEndpoint         string `env:"AUTH_SERVICE_ENDPOINT" env-required:"true"`
	UserServiceEndpoint         string `env:"USER_SERVICE_ENDPOINT" env-required:"true"`
	FeedServiceEndpoint         string `env:"FEED_SERVICE_ENDPOINT" env-required:"true"`
	MessageServiceEndpoint      string `env:"MESSAGE_SERVICE_ENDPOINT" env-required:"true"`
	NotificationServiceEndpoint string `env:"NOTIFICATION_SERVICE_ENDPOINT" env-required:"true"`
	MusicServiceEndpoint        string `env:"MUSIC_SERVICE_ENDPOINT" env-required:"true"`
	VideoServiceEndpoint        string `env:"VIDEO_SERVICE_ENDPOINT" env-required:"true"`
	GroupServiceEndpoint        string `env:"GROUP_SERVICE_ENDPOINT" env-required:"true"`
}

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func Load() Config {
	var config Config

	if err := cleanenv.ReadEnv(&config); err != nil {
		log.Fatal("couldn't bind settings to config")
	}

	return config
}
