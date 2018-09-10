package main

// create table chat_content(id int primary key auto_increment, uname varchar(10), content varchar(100));

import (
	"fmt"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

func checkErr(err error) {
	if err != nil {
		fmt.Println(err)
	}
}

func main() {
	db_conn, err := sql.Open("mysql", "chat_admin:chat_admin@tcp(192.168.137.128:3306)/chatroom?charset=utf8")
	checkErr(err)
	db_conn.SetMaxOpenConns(2000)
	db_conn.SetMaxIdleConns(1000)
	db_conn.Ping()
	rows, err := db_conn.Query("select * from chat_content")
	checkErr(err)
	columns, _ := rows.Columns()
	fmt.Println(columns)
	scanArgs := make([]interface{}, len(columns))
	values := make([]interface{}, len(columns))
	for i := range scanArgs {
		scanArgs[i] = &values[i]
	}
	// store all rows
	var record []map[string]string
	for rows.Next() {
		rows.Scan(scanArgs...)
		// store one row
		row := make(map[string]string)
		for i, col := range values {
			row[columns[i]] = string(col.([]byte))
		}
		record = append(record, row)
	}
	fmt.Println(record)
}
