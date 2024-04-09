// Базовая конструкция для создания фалга
// Но в этом примере мы используем функции flag.StringVar, flag.IntVar
// и flag.BoolVar, чтобы определить флаги, сохраняя значения в переменных
// без использования указателей.

// go run base_flag_2.go -name Vadim -age 22 -admin true
package main

import (
	"flag"
	"fmt"
)

func main() {

	var name string
	var age int
	var isAdmin bool

	flag.StringVar(&name, "name", "Ivan", "Name of user")
	flag.IntVar(&age, "age", 25, "Name of user")
	flag.BoolVar(&isAdmin, "admin", false, "Access rights")

	flag.Parse()

	fmt.Printf("Name: %s\n", name)
	fmt.Printf("Age: %d\n", age)
	fmt.Printf("Admin: %t\n", isAdmin)

}
