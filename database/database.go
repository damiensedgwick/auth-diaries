package database

import "database/sql"

var (
	DBConn *sql.DB
)

func NewDBConn() {
	db, err := sql.Open("sqlite3", "auth-diaries.db")
	if err != nil {
		panic(err)
	}

	DBConn = db
}
