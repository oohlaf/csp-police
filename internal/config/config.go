// https://github.com/learning-cloud-native-go/myapp/blob/master/config/config.go

package config

import (
	"github.com/joeshaw/envdecode"
	"github.com/rs/zerolog/log"
)

type appConf struct {
	Debug  bool `env:"DEBUG,default=false"`
	Server serverConf
}

type serverConf struct {
	Hostname string `env:"SERVER_HOSTNAME,default=localhost"`
	Port int `env:"SERVER_PORT,default=9090"`
}

func AppConfig() *appConf {
	var c appConf
	if err := envdecode.StrictDecode(&c); err != nil {
		log.Fatal().Msgf("Failed to decode: %s", err)
	}
	return &c
}
