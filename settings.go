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
	EmailSmtpPort     int    `split_words:"true" required:"true" default:"1025"`
	EmailSmtpUsername string `split_words:"true" required:"true" default:"maildev"`
	EmailSmtpPassword string `split_words:"true" required:"true" default:"maildev"`
	EmailFrom         string `split_words:"true" required:"true" default:"drone@example.com"`
}

func NewSettingsFromEnv() Settings {
	var settings Settings
	if err := envconfig.Process(envPrefix, &settings); err != nil {
		log.Fatalln(err)
	}
	return settings
}
