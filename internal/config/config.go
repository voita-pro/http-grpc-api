package config

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

type HTTP struct {
	Host     string        `env:"HOST" envDefault:"0.0.0.0"`
	Port     int           `env:"PORT" envDefault:"8080"`
	Secret   string        `env:"SECRET" envDefault:"019a81c1-26e1-8466-2cca-cfdd550b61e6"`
	ExpToken time.Duration `env:"EXP_TOKEN" envDefault:"1m"`
}

type GRPC struct {
	Host   string `env:"HOST" envDefault:"0.0.0.0"`
	Port   int    `env:"PORT" envDefault:"5001"`
	Secret string `env:"SECRET" envDefault:"2233445511"`
}

type DB struct {
	Host     string `env:"HOST" required:"true"`
	Port     int    `env:"PORT" required:"true"`
	User     string `env:"USER" required:"true"`
	Password string `env:"PASSWORD" required:"true"`
	Name     string `env:"NAME" required:"true"`
}

type Firebase struct {
	ProjectID string `env:"PROJECT_ID"`
	CredsFile string `env:"CREDENTIALS_FILE"`
	WebAPIKey string `env:"WEB_API_KEY"`
}

type Cfg struct {
	Env             string        `env:"ENV" envDefault:".env"`
	Debug           bool          `env:"DEBUG" envDefault:"false"`
	ShutdownTimeout time.Duration `env:"SHUTDOWN_TIMEOUT" envDefault:"5s"`
	Firebase        Firebase      `envPrefix:"FIREBASE_"`
	DB              DB            `envPrefix:"DB_"`
	HTTP            HTTP          `envPrefix:"HTTP_"`
	GRPC            GRPC          `envPrefix:"GRPC_"`
}

var cfg Cfg

func Get() *Cfg {
	return &cfg
}

func Init() (*Cfg, error) {
	if err := load(); err != nil {
		return nil, fmt.Errorf("can't load config: %w", err)
	}

	cfg = Cfg{}
	opts := env.Options{
		Prefix:                "APP_",
		UseFieldNameByDefault: true,
	}

	if err := env.ParseWithOptions(&cfg, opts); err != nil {
		return nil, fmt.Errorf("can't parse config: %w", err)
	}

	return &cfg, nil
}

func load() error {
	cfgEnv := os.Getenv("ENV_FILE")
	if len(cfgEnv) == 0 {
		cfgEnv = ".env"
	}
	err := godotenv.Load(cfgEnv)
	if err != nil {
		log.Panicf("can't parse config: %v", err)
	}
	return nil
}
