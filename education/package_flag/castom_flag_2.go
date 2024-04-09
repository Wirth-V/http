// Давайте рассмотрим более практический пример, где пользовательский тип флага
//имеет смысл и выполняет более полезную задачу. Допустим, у нас есть
//пользовательский тип DurationFlag, представляющий временной интервал.
//Мы хотим, чтобы пользователь мог устанавливать этот интервал с помощью
//аргумента командной строки в формате "5s" (5 секунд), "2m" (2 минуты) и так далее.

// Для запуска из консоли: go run castom_flag_2.go -customDuration=5s  или go run castom_flag_2.go -customDuration=5m
package main

import (
	"flag"
	"fmt"
	"time"
)

// DurationFlag - пользовательский тип флага для временных интервалов
type DurationFlag struct {
	duration time.Duration
}

func (d *DurationFlag) String() string {
	return d.duration.String()
}

func (d *DurationFlag) Set(value string) error {
	// Парсинг строки и установка значения
	val, err := time.ParseDuration(value)
	if err != nil {
		return fmt.Errorf("неверное значение: %s", value)
	}
	d.duration = val
	return nil
}

func main() {
	var customDuration DurationFlag

	// Добавление флага с пользовательским типом
	flag.Var(&customDuration, "customDuration", "Пользовательский флаг для временного интервала")

	// Разбор аргументов командной строки
	flag.Parse()

	// Использование значения флага
	fmt.Println("Значение пользовательского флага для временного интервала:", customDuration)
}

/*
В этом примере пользователь может устанавливать значение флага customDuration,
представляющего временной интервал, с использованием аргумента командной строки,
 например, -customDuration=5s для 5 секунд или -customDuration=2m для 2 минут.
 В результате флаг будет корректно разбираться и использоваться в программе,
 предоставляя более высокий уровень абстракции для работы с временем.
*/
