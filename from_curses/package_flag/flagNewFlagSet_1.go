/*
flag.NewFlagSet в языке программирования Go - это функция, которая создает новый
объект *flag.FlagSet. Тип *flag.FlagSet представляет собой набор флагов,
который можно использовать для определения и обработки собственных флагов
командной строки независимо от глобального flag.CommandLine. Это полезно,
когда вы хотите создать несколько наборов флагов для различных подсистем вашей
программы или когда вам нужна более тонкая настройка обработки ошибок.

Пример использования flag.NewFlagSet:
*/

// Команда запуска: go run flagNewFlagSet_1.go -name John
package main

import (
	"flag"
	"fmt"
)

func main() {
	// Создаем новый объект *flag.FlagSet с именем "myFlags"
	myFlagSet := flag.NewFlagSet("myFlags", flag.ExitOnError)

	// Определение флага в новом наборе флагов
	var name string
	myFlagSet.StringVar(&name, "name", "Guest", "Specify your name")

	// Парсинг командной строки только для нового набора флагов
	myFlagSet.Parse([]string{"-name", "John"}) //если написать myFlagSet.Parse() будет ошибка

	// Вывод результатов
	fmt.Printf("Hello, %s!\n", name)
}

/*
В этом примере мы создаем новый *flag.FlagSet с именем "myFlags" и настраиваем
его для аварийного завершения программы при обнаружении ошибки
(с использованием flag.ExitOnError). Затем мы добавляем строковый флаг в этот
новый набор флагов с помощью метода StringVar. Далее, мы вызываем метод Parse
этого конкретного набора флагов только для переданных аргументов командной
строки.

Использование flag.NewFlagSet позволяет создавать изолированные группы флагов
для конкретных частей программы, обеспечивая четкую структуру и предотвращая
конфликты флагов между различными компонентами программы.
*/

/*!!!!!
В предыдущем коде myFlagSet.Parse([]string{"-name", "John"}) используется для
явного парсинга переданных аргументов командной строки только для созданного
объекта *flag.FlagSet с именем "myFlags". Давайте разберем, что происходит
в этой строке кода:

    myFlagSet: Это созданный ранее объект *flag.FlagSet, представляющий набор
	флагов. Он был настроен с именем "myFlags" и опцией flag.ExitOnError, что
	означает аварийное завершение программы при обнаружении ошибок в процессе
	парсинга флагов.

    .Parse([]string{"-name", "John"}): Этот вызов метода Parse анализирует
	переданный массив строк ([]string{"-name", "John"}) в поиске флагов и их
	значений. В данном случае, он ищет флаг с именем "name" и устанавливает его
	значение в "John". Метод Parse возвращает ошибку, если обнаруживает
	невалидные флаги или аргументы командной строки.

Таким образом, в данном коде мы явно передаем массив строк, представляющих аргументы командной строки, в метод Parse для объекта *flag.FlagSet с именем "myFlags". Это позволяет нам провести парсинг только для этого набора флагов, изолируя его от глобального набора flag.CommandLine
*/