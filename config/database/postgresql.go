package database

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
)

// dive deep in here again
func PgOpenConnection() *sql.DB {
	//setup env to psql
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s "+"password=%s dbname=%s sslmode=disable application_name='%s'", os.Getenv("DB_HOST_POSTGRES"), os.Getenv("DB_PORT_POSTGRES"), os.Getenv("DB_USER_POSTGRES"), os.Getenv("DB_PW_POSTGRES"), os.Getenv("DB_NAME_POSTGRES"), os.Getenv("DB_APP_NAME"))

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	return db
}