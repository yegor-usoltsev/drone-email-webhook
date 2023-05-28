package main

import (
	"fmt"
	"log"
	"time"

	"github.com/kelseyhightower/envconfig"
)

const EnvPrefix = "PLUGIN_"

type Settings struct {
	ServerHost            string        `split_words:"true" default:"0.0.0.0"`
	ServerPort            int           `split_words:"true" default:"8080"`
	ServerMaxHeaderBytes  int           `split_words:"true" default:"4096"`    // 4 * 1024 = 4 KB
	ServerMaxBodyBytes    int           `split_words:"true" default:"1048576"` // 1 * 1024 * 1024 = 1 MB
	ServerReadTimeout     time.Duration `split_words:"true" default:"15s"`
	ServerHandlerTimeout  time.Duration `split_words:"true" default:"10s"`
	ServerWriteTimeout    time.Duration `split_words:"true" default:"15s"`
	ServerIdleTimeout     time.Duration `split_words:"true" default:"120s"`
	ServerShutdownTimeout time.Duration `split_words:"true" default:"15s"`
	EmailSmtpHost         string        `split_words:"true" default:"localhost"`
	EmailSmtpPort         int           `split_words:"true" default:"1025"`
	EmailSmtpUsername     string        `split_words:"true" default:"maildev"`
	EmailSmtpPassword     string        `split_words:"true" default:"maildev"`
	EmailFrom             string        `split_words:"true" default:"Drone <drone@example.com>"`
}

func NewSettingsFromEnv() Settings {
	var settings Settings
	err := envconfig.Process(EnvPrefix, &settings)
	if err != nil {
		_ = envconfig.Usage(EnvPrefix, &settings)
		fmt.Println()
		log.Panicln(err)
	}
	return settings
}
