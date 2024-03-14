// Любой тег, не определенный в программе, может быть сохранен и доступен с
// gомощью flag.Args() метод. Вы можете указать флаг, передав значение индекса
// i к flag.Args(i).

// Команда для запуска: go build no_flag_arguments_2.go, затем
// ./no_flag_arguments_2 -flavor=chocolate -cream -order=incomplete -flag1 -flag2=true
// или
// ./no_flag_arguments_2 -flavor=chocolate -cream -order=incomplete

/* Ожидаемый вывод:
$ ./coffee -flavor=chocolate -cream -order=incomplete
flavor: chocolate
quantity: 2
cream: true
order: incomplete
tail: []
					и
$ ./coffee -flavor=chocolate -cream -order=incomplete -flag1 -flag2=true
flavor: chocolate
quantity: 2
cream: true
order: incomplete
tail: [flag1 flag2=true]

!!! НО ПРОГРАММА НА ВТОРОМ ВЫВОДЕ ВЫДАЕТ:
./no_flag_arguments_2 -flavor=chocolate -cream -order=incomplete -flag1 -flag2=true
flag provided but not defined: -flag1
Usage of ./no_flag_arguments_2:
  -cream
        decide if you want cream
  -flavor string
        select shot flavor (default "vanilla")
  -order string
        status of order (default "complete")
  -quantity int
        quantity of shots (default 2)

ПОЧЕМУ ОЖИДАЕМЫЙ ВЫВОД ОТЛИЧАЕТСЯ ОТ ФАКТИЧЕСКОГО?
*/

//Ответ:
/*
Попробуй: go run no_flag_arguments_2.go -flavor=chocolate -cream -order=incomplete ser -flag1 -flag2=true
Выведет:
quantity: 2
cream: true
order: incomplete
tail: [ser -flag1 -flag2=true]
Т.к. после аргумента ser. Как только компилятор сталкнулся с лишним аргументом,
все следующие лишние символы он воспронимает как аргументы, а не как флаги
*/

package main

import (
	"flag"
	"fmt"
)

func main() {

	wordPtr := flag.String("flavor", "vanilla", "select shot flavor")
	numbPtr := flag.Int("quantity", 2, "quantity of shots")
	boolPtr := flag.Bool("cream", false, "decide if you want cream")

	var order string
	flag.StringVar(&order, "order", "complete", "status of order")

	flag.Parse()

	fmt.Println("flavor:", *wordPtr)
	fmt.Println("quantity:", *numbPtr)
	fmt.Println("cream:", *boolPtr)
	fmt.Println("order:", order)
	fmt.Println("tail:", flag.Args())

}
