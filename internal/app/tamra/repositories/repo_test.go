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
	var dbConnectionString string
	// Check if the code is being run locally
	if os.Getenv("CODEBUILD_BUILD_ID") != "" {
		// If it is running inside of codebuild use the db as host name since it's the name of the service in the docker-compose file
		dbConnectionString = "postgresql://postgres:mysecretpassword@db:5432/tamra-postgis-test?sslmode=disable"
	} else {
		dbConnectionString = "postgresql://postgres:mysecretpassword@localhost:5432/tamra-postgis-test?sslmode=disable"
	}

	Db, err = utils.NewDB(dbConnectionString)

	if err != nil {
		panic(err)
	}

	defer Db.Close()

	code := m.Run()

	os.Exit(code)
}
