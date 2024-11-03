package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/go-sql-driver/mysql"
)

var db *sql.DB

func initDB() error {
	cfg := mysql.Config{
		User:   os.Getenv("DBUSER"),
		Passwd: os.Getenv("DBPASSWD"),
		Net:    "tcp",
		Addr:   "127.0.0.1:3306",
		DBName: "recordings",
	}
	var err error
	db, err = sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		return err
	}
	return db.Ping()
}

func createTables() {
	_, err := db.Exec(`
	CREATE TABLE IF NOT EXISTS album (
		id int NOT NULL AUTO_INCREMENT,
		title varchar(128) COLLATE utf8mb4_general_ci NOT NULL,
		artist varchar(255) COLLATE utf8mb4_general_ci NOT NULL,
		price decimal(5,2) NOT NULL,
		PRIMARY KEY (id)
	) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci
	`)
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	err := initDB()
	if err != nil {
		log.Fatal(err)
	}
	createTables()
	fmt.Printf("complete migration")
}
