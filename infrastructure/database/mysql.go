package database

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

// NewMySQLConnection creates a new mysql connection
func NewMySQLConnection() *sql.DB {
	var ds = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s",
		"root",
		os.Getenv("MYSQL_PASSWORD"),
		os.Getenv("MYSQL_HOST"),
		os.Getenv("MYSQL_PORT"),
		os.Getenv("MYSQL_DATABASE"),
	)
	//host := os.Getenv("MYSQL_HOST")
	//database := os.Getenv("MYSQL_DATABASE")
	//username := os.Getenv("MYSQL_USER")
	//password := os.Getenv("MYSQL_PASSWORD")
	//fmt.Sprintf(
	//	"%s:%s@tcp(%s)/%s",
	//	username,
	//	password,
	//	host,
	//	database,
	//)
	db, err := sql.Open("mysql", ds)
	if err != nil {
		panic(err)
	}

	return db
}
