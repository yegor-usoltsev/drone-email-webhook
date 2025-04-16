package main

import (
	"github.com/kelseyhightower/envconfig"
)

const envPrefix = "DRONE"

type Settings struct {
	Secret            string `split_words:"true" required:"true"`
	ServerHost        string `split_words:"true" required:"true" default:"0.0.0.0"`
	ServerPort        uint16 `split_words:"true" required:"true" default:"3000"`
	EmailSMTPHost     string `split_words:"true" required:"true" default:"localhost"`
	EmailSMTPPort     uint16 `split_words:"true" required:"true" default:"25"`
	EmailSMTPUsername string `split_words:"true" required:"false"`
	EmailSMTPPassword string `split_words:"true" required:"false"`
	EmailFrom         string `split_words:"true" required:"true" default:"drone@localhost"`
}

func NewSettingsFromEnv() Settings {
	var settings Settings
	if err := envconfig.Process(envPrefix, &settings); err != nil {
		_ = envconfig.Usage(envPrefix, &settings)
		panic(err)
	}
	return settings
}
