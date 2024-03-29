/*
В языке Go, когда мы говорим о "неявном парсинге", мы обычно имеем в виду
использование глобального набора флагов flag.CommandLine и вызов метода
flag.Parse() без создания отдельного объекта *flag.FlagSet. Это является более
простым и распространенным подходом, особенно для небольших программ.


Теперь, если вы хотите избежать использования flag.CommandLine и вместо этого
явно указать свой набор флагов, вы можете использовать flag.NewFlagSet. Ниже
приведен пример аналогичной программы без использования flag.CommandLine:
*/

//Команда запуска: go run flagNewFlagSet_2.go -name John

package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	// Создаем новый набор флагов
	myFlagSet := flag.NewFlagSet("myFlags", flag.ExitOnError)

	// Определение флага
	var name string
	myFlagSet.StringVar(&name, "name", "Guest", "Specify your name")

	// Парсинг командной строки
	myFlagSet.Parse(os.Args[1:])

	// Вывод результатов
	fmt.Printf("Hello, %s!\n", name)
}

/*
В этом примере мы явно создаем новый набор флагов с помощью
flag.NewFlagSet("myFlags", flag.ExitOnError), добавляем флаг, и парсим аргументы
 командной строки только для этого набора флагов с использованием
 myFlagSet.Parse(os.Args[1:]).

Выбор между явным и неявным парсингом зависит от требований вашей программы.
Для небольших программ часто используется более простой и короткий вариант с
flag.CommandLine и flag.Parse(), но для более крупных проектов или подсистем
может быть полезно использовать явный *flag.FlagSet для лучшей изоляции и
управляемости флагов.
*/
