//Отрабатываем самую базовую конструкцию побработки ошибок

// Проверка ошибки на пустоту
package main

import "fmt"

func divide(a int, b int) int {
	return a / b
}

func main() {
	var a, b int
	_, err := fmt.Scan(&a)
	if err != nil {
		fmt.Println("Проверьте тип входных параметров")
	} else {
		_, err := fmt.Scan(&b)
		if err != nil || b == 0 {
			fmt.Println("Проверьте тип входных параметров")
		} else {
			fmt.Println(divide(a, b))
		}
	}
}
