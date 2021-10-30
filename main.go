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
	command  = flag.String("c", "", "Connect to MySQL Server execute command, and exit")
	help     = flag.Bool("h", false, "Help for mysql")
	db       *sql.DB
)

func printUsage() {
	fmt.Println(`mysql
A mysql client(or terminal) written in Golang

examples:
mysql -u="username" -p="password" -d="awesome_db"
mysql -u="username" -p="password" -d="awesome_db" -c="SELECT NOW() AS t"

options:`)
	flag.PrintDefaults()
}

func reconnect() error {
	var err error
	DSN := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", *user, *passwd, *host, *port, *database)
	db, err = sql.Open("mysql", DSN)
	db.SetMaxIdleConns(2)
	db.SetMaxOpenConns(10)
	db.SetConnMaxIdleTime(time.Minute)
	err = db.Ping()
	return err
}

func main() {
	flag.Parse()
	if *help {
		printUsage()
		return
	}
	if err2 := reconnect(); err2 != nil {
		fmt.Println(err2)
		return
	}
	fmt.Printf("Connected: %s\n", *host)
	fmt.Printf("Server version: %s\n", showServerVersion(db))

	if *command != "" {
		executeCmd(db, *command)
		return
	}

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
