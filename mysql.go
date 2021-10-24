package main

import (
	"database/sql"
	"fmt"
	"reflect"
)

func showServerVersion(db *sql.DB) string {
	row := db.QueryRow(`select version()`)
	if err := row.Err(); err != nil {
		fmt.Println(err)
		return ""
	}
	ver := ""
	if err := row.Scan(&ver); err != nil {
		fmt.Println(err)
		return ""
	}
	return ver
}

func executeCmd(db *sql.DB, cmd string) {
	rows, err := db.Query(cmd)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer func() { _ = rows.Close() }()
	columns, err := rows.Columns()
	if err != nil {
		fmt.Println(err)
		return
	}
	for _, c := range columns {
		fmt.Printf("%s\t", c)
	}
	fmt.Println()
	columnTypes, err := rows.ColumnTypes()
	if err != nil {
		fmt.Println(err)
		return
	}
	for rows.Next() {
		var dist = make([]interface{}, len(columnTypes))
		for i, ct := range columnTypes {
			switch ct.DatabaseTypeName() {
			case "VARCHAR", "TEXT", "LONGTEXT", "DATETIME", "DATE", "TIMESTAMP":
				var varcharVal string
				dist[i] = &varcharVal
			case "INT", "TINYINT", "BIGINT":
				var intVal int
				dist[i] = &intVal
			}
		}
		err = rows.Scan(dist...)
		if err != nil {
			fmt.Println(err)
			return
		}
		for _, d := range dist {
			fmt.Printf("%v\t", reflect.ValueOf(d).Elem().Interface())
		}
		fmt.Println()
	}
}
