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
	EmailSmtpHost     string `split_words:"true" required:"true" default:"localhost"`
	EmailSmtpPort     int    `split_words:"true" required:"true" default:"25"`
	EmailSmtpUsername string `split_words:"true"`
	EmailSmtpPassword string `split_words:"true"`
	EmailFrom         string `split_words:"true" required:"true" default:"drone@localhost"`
}

func NewSettingsFromEnv() Settings {
	var settings Settings
	if err := envconfig.Process(envPrefix, &settings); err != nil {
		log.Fatalln("settings:", err)
	}
	return settings
}
