package main

import (
	"database/sql"
	"fmt"
	"strings"
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
	if cmd == "" {
		return
	}
	rows, err := db.Query(cmd)
	if err != nil {
		fmt.Println(err)
		return
	}

	// cmd `USE $DB`
	_cmd := strings.TrimSuffix(strings.TrimSpace(cmd), ";")
	if strings.Contains(strings.ToUpper(_cmd), "USE ") {
		*database = strings.TrimSpace(_cmd[3:])
		_ = reconnect()
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
		for i := range columnTypes {
			var varcharVal sql.NullString
			dist[i] = &varcharVal
			//switch ct.DatabaseTypeName() {
			//case "CHAR", "VARCHAR", "TEXT", "LONGTEXT", "DATETIME", "DATE", "TIMESTAMP":
			//	var varcharVal sql.NullString
			//	dist[i] = &varcharVal
			//case "INT", "TINYINT", "BIGINT":
			//	var intVal sql.NullInt64
			//	dist[i] = &intVal
			//case "FLOAT", "DECIMAL":
			//	var floatVal sql.NullFloat64
			//	dist[i] = &floatVal
			//}
		}
		err = rows.Scan(dist...)
		if err != nil {
			fmt.Println(err)
			return
		}
		for _, d := range dist {
			switch d.(type) {
			case *sql.NullInt64:
				a := d.(*sql.NullInt64)
				fmt.Printf("%v\t", map[bool]string{
					true:  fmt.Sprintf("%v", a.Int64),
					false: "NULL",
				}[a.Valid])
			case *sql.NullFloat64:
				a := d.(*sql.NullFloat64)
				fmt.Printf("%v\t", map[bool]string{
					true:  fmt.Sprintf("%v", a.Float64),
					false: "NULL",
				}[a.Valid])
			case *sql.NullString:
				a := d.(*sql.NullString)
				fmt.Printf("%v\t", map[bool]string{
					true:  a.String,
					false: "NULL",
				}[a.Valid])
			}
		}
		fmt.Println()
	}
}
