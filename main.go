package main

import (
	"bufio"
	"database/sql"
	"flag"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"os"
	"time"
)

var (
	host     = flag.String("H", "localhost", "The MySQL Server host name")
	port     = flag.Int("P", 3306, "The MySQL Server port")
	passwd   = flag.String("p", "", "The MySQL Server password")
	user     = flag.String("u", "root", "The MySQL Server username")
	database = flag.String("d", "test", "The MySQL Server database")
	help     = flag.Bool("h", false, "Help for gfm")
)

func printUsage() {
	fmt.Println(`gfm
A mysql client(or terminal) written in Golang

examples:
gfm -u="username" -p="password" -d="awesome_db"

options:`)
	flag.PrintDefaults()
}

func main() {
	flag.Parse()
	if *help {
		printUsage()
		return
	}
	DSN := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", *user, *passwd, *host, *port, *database)
	db, err := sql.Open("mysql", DSN)
	db.SetMaxIdleConns(2)
	db.SetMaxOpenConns(10)
	db.SetConnMaxIdleTime(time.Minute)
	defer func() { _ = db.Close() }()
	if err != nil {
		fmt.Println(err)
		return
	}
	err = db.Ping()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("Connected: %s\n", *host)
	fmt.Printf("Server version: %s\n", showServerVersion(db))
	for {
		fmt.Printf("(%s)>", *database)
		scanner := bufio.NewScanner(os.Stdin)
		_ = scanner.Scan()
		cmd := scanner.Text()
		func() {
			defer func() {
				if re := recover(); re != nil {
					fmt.Println(re)
				}
			}()
			executeCmd(db, cmd)
		}()
	}
}
