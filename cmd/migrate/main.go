package main

import (
	"database/sql"
	"os"

	"csp-police/internal/config"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	appConf := config.AppConfig()
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	if appConf.Debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}

	log.Info().Str("driver", appConf.Db.Driver).Str("DSN", appConf.Db.DSN).Msg("Connecting to database")
	db, err := sql.Open(appConf.Db.Driver, appConf.Db.DSN)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to open database connection")
	}
	defer db.Close()

	driver, err := sqlite3.WithInstance(db, &sqlite3.Config{})
	defer driver.Close()

	m, err := migrate.NewWithDatabaseInstance("file://db/migrate-sqlite3/", "sqlite3", driver)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to initialize migrate")
	}
	defer m.Close()
	m.Log = logger{}
	err = m.Up()
	if err != nil {
		if err == migrate.ErrNoChange {
			log.Info().Msg("No change")
		} else {
			log.Fatal().Err(err).Msg("Error during migration")
		}
	}
}

type logger struct{}

func (t logger) Printf(format string, v ...interface{}) {
	log.Info().Msgf(format, v...)
}

func (t logger) Verbose() bool {
	return true
}
