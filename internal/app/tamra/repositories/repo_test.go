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
	dbConnectionString := "postgresql://postgres:mysecretpassword@db:5432/tamra-postgis-test?sslmode=disable"

	Db, err = utils.NewDB(dbConnectionString)

	if err != nil {
		panic(err)
	}

	defer Db.Close()

	code := m.Run()

	os.Exit(code)
}
