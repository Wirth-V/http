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

	pg := req.String("pg", "postgres://server:198416@localhost:6667/server", "db")
	table := req.String("table", "item", "table")

	req.Parse(os.Args[Two:])

	var err error

	switch os.Args[One] {
	case "start":
		err = modules.Server(req, *host, *port, *pg, *table)
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
