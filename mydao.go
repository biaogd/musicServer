package main

import (
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

//获得连接对象
func getDB() *sql.DB {
	db, err := sql.Open("mysql", "root:,.Rfb8848/mydb")
	if err != nil {
		log.Panicln(err)
	}
	return db
}
