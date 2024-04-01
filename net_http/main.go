package main

import (
	"app/moduls"
	"flag"
	"fmt"
	"os"
)

func main() {

	req := flag.NewFlagSet(os.Args[1], flag.ExitOnError)
	host := req.String("host", "localhost", "Host")
	port := req.String("port", "8080", "Host")

	req.Parse(os.Args[2:])

	switch os.Args[1] {
	case "start":
		moduls.Server(req, host, port)
	case "request":
		moduls.Client(req, host, port)
	default:
		fmt.Println("You flag is not correct:")
		os.Exit(1)
	}
}