/*
В Golang пакет flag предоставляет возможность создания пользовательских типов
флагов с помощью интерфейса flag.Value. Этот интерфейс определен в пакете flag
и включает методы, позволяющие пользовательским типам реализовать собственное поведение при установке и чтении флага.

Интерфейс flag.Value включает следующие методы:

    String() string:
        Этот метод должен возвращать строковое представление значения флага.

    Set(string) error:
        Этот метод вызывается при установке значения флага из строки аргумента
		командной строки. Метод должен проанализировать переданную строку и
		установить соответствующее значение для флага. Если строка не может быть
		корректно разобрана, метод должен вернуть ошибку.

    Get() interface{}:
        Этот метод возвращает текущее значение флага в виде пустого интерфейса.
		Это используется внутренне пакетом flag для доступа к текущему значению флага.

Вот пример создания пользовательского типа флага с использованием интерфейса flag.Value:
*/

// Запуск в консоли: go run castom_flag_1.go -customFlag 5
package main

import (
	"flag"
	"fmt"
	"strconv"
)

// CustomFlagType - пользовательский тип флага
type CustomFlagType int

func (c *CustomFlagType) String() string {
	return strconv.Itoa(int(*c))
}

func (c *CustomFlagType) Set(value string) error {
	// Парсинг строки и установка значения
	val, err := strconv.Atoi(value)
	if err != nil {
		return fmt.Errorf("неверное значение: %s", value)
	}
	*c = CustomFlagType(val)
	return nil
}

func main() {
	var customFlag CustomFlagType

	// Добавление флага с пользовательским типом
	flag.Var(&customFlag, "customFlag", "Пользовательский флаг")

	// Разбор аргументов командной строки
	flag.Parse()

	// Использование значения флага
	fmt.Println("Значение пользовательского флага:", customFlag)
}

/*
В этом примере определен пользовательский тип CustomFlagType, реализующий
интерфейс flag.Value. Затем создается переменная этого типа, и с помощью
flag.Var она связывается с конкретным флагом командной строки. Таким образом,
при использовании программы можно устанавливать значение флага customFlag с
помощью аргумента командной строки.
*/
