package db

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

var db, _ = sql.Open("mysql", "root:@/wildcare?parseTime=true")

var Exec = db.Exec
var Query = db.Query
var QueryRow = db.QueryRow
