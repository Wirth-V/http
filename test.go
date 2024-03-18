package main

import "fmt"

var globalVariable = 6 // Глобальная переменна

func main() {
	fmt.Println("0:", globalVariable)

	// globalVariable = 10 // Изменяем значение глобальной переменной внутри функции main
	// fmt.Println("10:", globalVariable)

	anotherFunction() // Вызываем другую функцию

	fmt.Println("6:", globalVariable) // Выводим значение глобальной переменной после выполнения функции anotherFunction
}

func anotherFunction() {
	fmt.Println("10", globalVariable)
	globalVariable = 20 // Изменяем значение глобальной переменной внутри другой функции
	fmt.Println("20:", globalVariable)
}
