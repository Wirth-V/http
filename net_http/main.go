package main

import (
	"flag"
	"net_http/modules"
	"os"
)

const (
	One  int = 1
	Two  int = 2
	Size int = 21
)

func main() {

	if len(os.Args) < Two {
		modules.ErrorLog.Println("You comand is not correct: there are not enough commands")
		os.Exit(1)
	}

	req := flag.NewFlagSet(os.Args[One], flag.ExitOnError)
	host := req.String("host", "localhost", "Host")
	port := req.String("port", "8080", "Host")
	username_db := req.String("user_db", "server", "db")
	password_db := req.String("password_db", "198416", "db")
	host_db := req.String("host_db", "localhost", "db")
	port_db := req.String("port_db", "6667", "db")

	db := req.String("db", "server", "db")
	table := req.String("table", "item", "table")

	req.Parse(os.Args[Two:])

	var err error

	switch os.Args[One] {
	case "start":
		err = modules.Server(req, *host, *port, *username_db, *password_db, *host_db, *port_db, *db, *table)
	case "request":
		err = modules.Client(req, *host, *port)
	default:
		modules.ErrorLog.Println("you comand is not correct: non-direct command")
		os.Exit(1)
	}

	if err != nil {
		modules.ErrorLog.Println(err)
	}
}
