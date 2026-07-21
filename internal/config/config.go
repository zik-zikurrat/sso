package config

import (
	"flag"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env         string        `yaml:"env" env-default:"local"`
	StoragePath string        `yaml:"storage_path" env-required:"true"`
	TokenTTL    time.Duration `yaml:"token_ttl" env-required:"true"`
	PG          PGConfig      `yaml:"postgre" env-required:"true"`
	Server      ServerConfig  `yaml:"server" env-required:"true"`
	GRPC        GRPCConfig    `yaml:"grpc"`
	SMTP        SMTPConfig    `yaml:"smtp"`
	Logging     LoggingConfig `yaml:"logging_handler"`
	Cache       CacheConfig   `yaml:"cache"`
}

type LoggingConfig struct {
	Level   string `yaml:"level"`
	Format  string `yaml:"format"`
	Discard bool   `yaml:"discard"`
}
type PGConfig struct {
	PORT     int    `yaml:"port"`
	NAME     string `yaml:"name"`
	USER     string `yaml:"user"`
	SSL      string `yaml:"ssl"`
	HOST     string `yaml:"host"`
	PASSWORD string `yaml:"password"`
}

type ServerConfig struct {
	PORT            int    `yaml:"port"`
	HOST            string `yaml:"host"`
	ReadTimeout     int    `yaml:"read_timeout"`
	WriteTimeout    int    `yaml:"write_timeout"`
	ShutdownTimeout int    `yaml:"shutdown_timeout"`
	Prefork         bool   `yaml:"prefork"`
}

type GRPCConfig struct {
	Port    int           `yaml:"port"`
	Timeout time.Duration `yaml:"timeout"`
}

type SMTPConfig struct {
	Password string `yaml:"smtp_password" env-required:"true"`
	Email    string `yaml:"smtp_email" env-required:"true"`
	Host     string `yaml:"smtp_host" env-default:"smtp.gmail.com"`
	Port     int    `yaml:"smtp_port" env-default:"587"`
}

type CacheConfig struct {
	HOST string `yaml:"host"`
	PORT int    `yaml:"port"`
	DB   int    `yaml:"db"`
}

func MustLoad() *Config {
	path := fetchConfigPath()
	if path == "" {
		panic("config path is empty")
	}
	if _, err := os.Stat(path); os.IsNotExist(err) {
		panic("config file does not exist: " + path)
	}

	var config Config
	if err := cleanenv.ReadConfig(path, &config); err != nil {
		panic("failed to read config: " + err.Error())
	}

	return &config
}

func fetchConfigPath() string {
	var res string

	if !flag.Parsed() {
		flag.StringVar(&res, "config", "", "path to config file")
		flag.Parse()
	}

	if res == "" {
		res = os.Getenv("CONFIG_PATH")
	}
	if res == "" {
		res = "./config/prod.yml"
	}
	return res
}
