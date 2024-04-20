package main

import (
	"flag"
	"fmt"
	"net_http/modules"
	"os"
)

func main() {

	req := flag.NewFlagSet(os.Args[1], flag.ExitOnError)
	host := req.String("host", "localhost", "Host")
	port := req.String("port", "8080", "Host")
	db := req.String("db", "server", "db")
	table := req.String("table", "item", "table")

	req.Parse(os.Args[2:])

	switch os.Args[1] {
	case "start":
		modules.Server(req, *host, *port, *db, *table) //запустит сервер
	case "request":
		modules.Client(req, *host, *port) //запустит клиент
	default:
		fmt.Println("You flag is not correct:")
		os.Exit(1)
	}
}
