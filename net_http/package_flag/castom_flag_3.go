/*
Если вам нужно определить флаг с пользовательским типом или проверкой, вы можете
реализовать интерфейс flag.Value.  Этот интерфейс имеет два метода:

     String() string — возвращает строковое представление значения флага.
     Set(string) error — устанавливает значение флага из строки.

После реализации этого интерфейса вы можете использовать функцию flag.Var для
создания собственного флага.
*/

//Давайте создадим собственный флаг, который принимает список строк,
// разделенных запятыми:

package main

import (
	"flag"
	"fmt"
	"strings"
)

type StringList []string

func (s *StringList) String() string {
	return strings.Join(*s, ",") // // Объединить элементы с разделителем ", "
}

func (s *StringList) Set(value string) error {
	*s = strings.Split(value, ",") // Bспользуется для разделения строки на
	//подстроки на основе определенного разделителя (',' в данном случии).
	return nil
}

func main() {
	// Определить пользовательский флаг
	var tags StringList
	flag.Var(&tags, "tags", "Comma-separated list of tags")

	// Парсим флаг
	flag.Parse()

	// Используем флаг
	fmt.Printf("Tags: %v\n", tags)
}

//Теперь вы можете использовать флаг `-tags` со списком значений,
//разделенных запятыми:

/*
$ go build main.go
$ ./main -tags=go,programming,tutorial
Tags: [go programming tutorial]
*/
