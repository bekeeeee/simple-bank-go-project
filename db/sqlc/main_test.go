package db

import (
	"database/sql"
	"github/bekeeeee/simplebank/util"
	"os"
	"testing"

	"github.com/rs/zerolog/log"

	_ "github.com/lib/pq"
)

var testQueries *Queries
var testDB *sql.DB

func TestMain(m *testing.M) {
	config, err := util.LoadConfig("../..")
	if err != nil {
		log.Fatal().Msg("Cannot load config")
	}
	// var err error
	testDB, err = sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal().Msg("cannot connect to db")
	}

	testQueries = New(testDB)

	os.Exit(m.Run())
}
