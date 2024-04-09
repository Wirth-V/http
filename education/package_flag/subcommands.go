/*
Подумайте о подкомандах как о основных опциях для команды, которые можно
дополнительно настроить с использованием собственных флагов. Например,
рассмотрим команду go build, где build - это подкоманда. Для настройки
подкоманды используйте метод NewFlagSet. После вызова этого метода вы можете
добавлять флаги к вызову, используя следующий синтаксис:

subCmd := flag.NewFlagSet("sub",flag.ExitOnError)
flagOne := subCmd.Bool("flagKeyword", false, "description")
*/

package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	subOne := flag.NewFlagSet("one", flag.ExitOnError)
	oneCream := subOne.String("cream", "No", "Cream")
	oneSuger := subOne.String("sugar", "No", "Sugar")

	subTwo := flag.NewFlagSet("two", flag.ExitOnError)
	twoWisk := subTwo.String("wiske", "No", "Wiske")

	if len(os.Args) < 2 {
		fmt.Println("expected 'one' or 'two' subcommands")
		os.Exit(1)
	}

	switch os.Args[1] {

	case "one":
		subOne.Parse(os.Args[2:]) // Почему нельзя написать os.Args[2]?
		fmt.Println("Cream:", *oneCream)
		fmt.Println("Sugar:", *oneSuger)
		fmt.Println("Tail:", subOne.Args())

	case "two":
		subTwo.Parse(os.Args[2:])
		fmt.Println("Wiske:", *twoWisk)
		fmt.Println("Tail:", subTwo.Args())

	default:
		fmt.Println("expected 'one' or 'two' subcommands")
		os.Exit(1)
	}

}

/*
Сборка:
go build subcommands.go first.

Запуск:
$ ./subcommands one -cream -sugar=Yes
subcommand 'one'
  Cream: No
  Sugar: Yes
  tail: []

$ ./subcommands two -wisk Yes
subcommand 'two'
  Wisk: Yes
  tail: []
*/
