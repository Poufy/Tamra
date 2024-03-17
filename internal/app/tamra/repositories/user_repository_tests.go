package repositories

import (
	"database/sql"
)

var userRepo *UserRepositoryImpl

func TestMain() {
	db, _ := sql.Open("postgres", "user=postgres password=postgres dbname=postgres sslmode=disable")
	defer db.Close()

	// userRepo = NewUserRepository(db)
}
