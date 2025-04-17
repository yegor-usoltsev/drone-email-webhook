package main

import (
	"github.com/kelseyhightower/envconfig"
)

const envPrefix = "DRONE"

type Config struct {
	Secret            string   `split_words:"true" required:"true"`
	ServerHost        string   `split_words:"true" required:"true" default:"0.0.0.0"`
	ServerPort        uint16   `split_words:"true" required:"true" default:"3000"`
	EmailSMTPHost     string   `split_words:"true" required:"true" default:"localhost"`
	EmailSMTPPort     uint16   `split_words:"true" required:"true" default:"25"`
	EmailSMTPUsername string   `split_words:"true" required:"false"`
	EmailSMTPPassword string   `split_words:"true" required:"false"`
	EmailFrom         string   `split_words:"true" required:"true" default:"drone@localhost"`
	EmailCC           []string `split_words:"true" required:"false"`
	EmailBCC          []string `split_words:"true" required:"false"`
}

func NewConfigFromEnv() Config {
	var cfg Config
	if err := envconfig.Process(envPrefix, &cfg); err != nil {
		_ = envconfig.Usage(envPrefix, &cfg)
		panic(err)
	}
	return cfg
}
