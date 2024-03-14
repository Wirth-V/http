/*

В пакете flag в языке программирования Go существует функция flag.ExitOnError,
которая устанавливает поведение по умолчанию для обработки ошибок в парсинге
флагов. Эта функция настроена таким образом, что при возникновении ошибки в
процессе парсинга флагов она вызывает функцию os.Exit(2) для завершения
программы с кодом возврата 2. Этот код возврата обычно используется для
обозначения ошибок в командной строке.

Вот пример использования flag.ExitOnError вместе с пакетом flag:
*/

package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	// Используем flag.ExitOnError для настройки обработки ошибок
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)

	// Определение флага
	var name string
	flag.StringVar(&name, "name", "Guest", "Specify your name")

	// Парсинг командной строки
	flag.Parse()

	// Вывод результатов
	fmt.Printf("Hello, %s!\n", name)
}

/*

В этом примере flag.CommandLine устанавливается в новый объект flag.FlagSet,
созданный с использованием flag.ExitOnError. Это гарантирует, что любая ошибка
при парсинге флагов вызовет os.Exit(2).

Пример использования программы:
$ go run flagExitOnError.go -name John
Hello, John!

Если бы вы ввели что-то некорректное, например:
$ go run flagExitOnError.go -invalidFlag

Программа завершилась бы с кодом возврата 2 и не вывела бы "Hello, Guest!".

*/
