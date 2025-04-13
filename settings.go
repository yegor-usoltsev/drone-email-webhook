package main

import (
	"log"

	"github.com/kelseyhightower/envconfig"
)

const envPrefix = "DRONE"

type Settings struct {
	Secret            string `split_words:"true" required:"true"`
	ServerHost        string `split_words:"true" required:"true" default:"0.0.0.0"`
	ServerPort        int    `split_words:"true" required:"true" default:"3000"`
	EmailSMTPHost     string `split_words:"true" required:"true" default:"localhost"`
	EmailSMTPPort     int    `split_words:"true" required:"true" default:"25"`
	EmailSMTPUsername string `split_words:"true"`
	EmailSMTPPassword string `split_words:"true"`
	EmailFrom         string `split_words:"true" required:"true" default:"drone@localhost"`
}

func NewSettingsFromEnv() Settings {
	var settings Settings
	if err := envconfig.Process(envPrefix, &settings); err != nil {
		log.Fatalf("[FATAL] settings: failed to load configuration: %v", err)
	}
	return settings
}
