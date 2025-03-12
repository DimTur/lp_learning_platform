package config

type Redis struct {
	Host       string `yaml:"host"`
	Port       int    `yaml:"port"`
	AttemptsDB int    `yaml:"attemts_db"`
	Password   string `yaml:"password"`
}
