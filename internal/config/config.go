package config

import "time"

type Config struct {
	App      AppConfig
	Database DBConfig
	JWT      JWTConfig
}

type AppConfig struct {
	Port string `env:"APP_PORT,required"`
}

type DBConfig struct {
	URL string `env:"POSTGRES_URL,required"`
}

type JWTConfig struct {
	Secret    string        `env:"JWT_SECRET,required"`
	ExpiresIn time.Duration `env:"JWT_EXPIRES_IN" envDefault:"24h"`
}
