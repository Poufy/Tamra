package repositories

import (
	"Tamra/internal/pkg/utils"
	"database/sql"
	"os"
	"testing"
)

var Db *sql.DB

func TestMain(m *testing.M) {
	var err error
	dbConnectionString := os.Getenv("TEST_DB_CONNECTION_STRING")

	Db, err = utils.NewDB(dbConnectionString)

	if err != nil {
		panic(err)
	}

	defer Db.Close()

	code := m.Run()

	os.Exit(code)
}
