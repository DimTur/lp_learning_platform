package config

import (
	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	GRPCServer GRPCServer `yaml:"grpc_server"`
	Storage    Storage    `yaml:"storage"`
	RabbitMQ   RabbitMQ   `yaml:"rabbit_mq"`
}

type GRPCServer struct {
	Address string `yaml:"address" env-default:":8002"`
}

type Storage struct {
	Type     string `yaml:"type"`
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	DBName   string `yaml:"dbname"`
}

func Parse(s string) (*Config, error) {
	c := &Config{}
	if err := cleanenv.ReadConfig(s, c); err != nil {
		return nil, err
	}

	return c, nil
}
