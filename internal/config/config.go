package config

import (
	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	GRPCServer GRPCServer `yaml:"grpc_server"`
	Storage    Storage    `yaml:"storage"`
}

type GRPCServer struct {
	Address string `yaml:"address" env-default:":8002"`
}

type Storage struct {
	SQLitePath string `yaml:"path" env-default:"db.sql"`
}

func Parse(s string) (*Config, error) {
	c := &Config{}
	if err := cleanenv.ReadConfig(s, c); err != nil {
		return nil, err
	}

	return c, nil
}
