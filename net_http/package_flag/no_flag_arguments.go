// Функция flag.Args() в языке программирования Go возвращает аргументы
// командной строки, которые не являются флагами. Это означает, что она
// возвращает все аргументы, которые не были обработаны пакетом flag в рамках
// определенных флагов.

// Синтаксис: func Args() []string
// возвращает аргументы командной строки, которые не были
// использованы для установки значений флагов. Эта функция возвращает срез строк.

// Для запуска из консоли используй команду: go run no_flag_arguments_1.go -name Vadim -age 20 23 23 2345 many
package main

import (
	"flag"
	"fmt"
)

func main() {
	// Определение флагов
	var name string
	var age int

	flag.StringVar(&name, "name", "Guest", "Specify your name")
	flag.IntVar(&age, "age", 25, "Specify your age")

	// Парсинг аргументов командной строки
	flag.Parse()

	// Использование значений флагов
	fmt.Printf("Name: %s\n", name)
	fmt.Printf("Age: %d\n", age)

	fmt.Printf("\n")

	// Аргументы без флагов
	fmt.Println("Non-flag arguments:")
	for _, arg := range flag.Args() {
		fmt.Printf("- %v\n", arg)
	}

	fmt.Printf("\n")

	// Получение аргументов командной строки, не связанных с флагами
	nonFlagArgs := flag.Args()
	fmt.Println("Non-flag arguments:", nonFlagArgs)
}

/*
Обе конструкции выполняют схожую задачу, а именно вывод аргументов командной строки, не связанных с флагами, но есть некоторые различия в их использовании и предназначении.
Конструкция 1:

fmt.Println("Non-flag arguments:")
for _, arg := range flag.Args() {
    fmt.Println(arg)
}
В этой конструкции мы прямо в цикле for перебираем значения, возвращаемые функцией flag.Args(), и выводим каждый аргумент в отдельной строке с использованием fmt.Println. Это явный способ итерации по аргументам и вывода каждого из них.

Конструкция 2:
nonFlagArgs := flag.Args()
fmt.Println("Non-flag arguments:", nonFlagArgs)
В этой конструкции мы сначала сохраняем аргументы, возвращаемые функцией flag.Args(), в переменную nonFlagArgs. Затем мы выводим весь срез аргументов в одной строке с использованием fmt.Println. Этот подход более компактен и выводит все аргументы вместе.

Различия и советы по использованию:
    Подход к выводу:
        Конструкция 1 более явна и подходит, если вы хотите более тонкую настройку вывода каждого аргумента.
        Конструкция 2 более компактна и предпочтительна, если вам необходимо вывести все аргументы вместе.

    Сохранение в переменную:
        Конструкция 2 более удобна, если вы хотите сохранить аргументы в переменной для дальнейшего использования в программе.

    Выбор зависит от задачи:
        Выбор между конструкциями зависит от ваших конкретных требований. Если необходима более тонкая настройка вывода каждого аргумента, используйте конструкцию 1. Если важнее компактность и общий вид всех аргументов, используйте конструкцию 2.

Примеры ниже иллюстрируют различия в выводе:
// Конструкция 1
fmt.Println("Non-flag arguments:")
for _, arg := range flag.Args() {
    fmt.Println("  -", arg)
}

// Конструкция 2
nonFlagArgs := flag.Args()
fmt.Println("Non-flag arguments:", nonFlagArgs)
В результате, если аргументы командной строки - "arg1", "arg2", "arg3", то вывод будет следующим:

markdown

Non-flag arguments:
  - arg1
  - arg2
  - arg3

Non-flag arguments: [arg1 arg2 arg3]
*/
