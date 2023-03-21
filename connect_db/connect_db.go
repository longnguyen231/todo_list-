package connect_db

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"log"
)

type DBSql struct {
	db *sql.DB
}

var dbSql *sql.DB

func (d *DBSql) New() {
	dbConnect, err := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/todo_list?parseTime=true")
	if err != nil {

		log.Print(err.Error())
	}
	dbSql = dbConnect
}

func GetDB() *sql.DB {
	return dbSql
}
