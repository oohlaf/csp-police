// https://github.com/learning-cloud-native-go/myapp/blob/master/config/config.go

package config

import (
	"github.com/joeshaw/envdecode"
	"github.com/rs/zerolog/log"
)

type appConf struct {
	Debug bool `env:"DEBUG,default=false"`
	Web   webConf
	Grpc  grpcConf
	Db	  dbConf
}

type webConf struct {
	Hostname string `env:"WEB_HOSTNAME,default=localhost"`
	Port     int    `env:"WEB_PORT,default=8080"`
}

type grpcConf struct {
	Hostname string `env:"GRPC_HOSTNAME,default=localhost"`
	Port     int    `env:"GRPC_PORT,default=9090"`
}

type dbConf struct {
	Driver string `env:"DB_DRIVER,default=sqlite3"`
	DSN string `env:"DB_DSN,default=file:cspreports.db"`
}

func AppConfig() *appConf {
	var c appConf
	if err := envdecode.StrictDecode(&c); err != nil {
		log.Fatal().Msgf("Failed to decode: %s", err)
	}
	return &c
}
