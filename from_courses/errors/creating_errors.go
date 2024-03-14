//Встроенные функции для создания ошибок: errors.New и fmt.Errorf

//Обе эти функции позволяют нам указывать настраиваемое сообщение об ошибке,
//которое вы можете отображать вашим пользователям.

package main

import (
	"errors"
	"fmt"
	"strings"
)

var MY_STATIC_ERR = errors.New("Error 1") //Позвалет создавать разные ошибки с одинаковым текстом (так как ссылки)
var MY_STATIC_ERR_2 = errors.New("Error 1")

var MY_DYNAMIC_ERR = fmt.Errorf("Error %d", 2)

// Эти две функции возвращают интерфейс err, что позволяет эти ошибки делить эти ошибки по разным типиам
// Мы можем делать специфичную логику для разных типов ошибок

type MyCustErr struct {
	First  int
	Second int
	Third  int
}

//Прописываем кастомную логику обработки ошибок. Тип (интерфейс) err имеет прописанный метод Error. Вот мы его реализовали.
func (e *MyCustErr) Error() string {
	return fmt.Sprintf("Error %d %d %d", e.First, e.Second, e.Third)
}

func boom() error {
	return MY_STATIC_ERR // Возвращает не строку, а указатель на объект
}

func BOOM() error {
	err := MyCustErr{
		First:  2,
		Second: 3,
		Third:  4,
	} 
	return &err
	//return fmt.Errorf("Error %d %d %d", 2, 3, 4) // Функция fmt.Errorf позволяет динамически создавать сообщение об ошибке.
	//Возвращает новый объект (копию), даже если ты работаешь со ссылками, то будет копия ссылки
}

func main() {
	if err := boom(); err != nil {
		if err == MY_STATIC_ERR {

		} else if strings.Contains(err.Error(), "Error xxx ") {
			fmt.Println(err)
		} 
		
	}

	fmt.Printf("\n")

	if err := BOOM(); err != nil {
		fmt.Println(err)

		switch v := err.(type) { //v - конкретный экземпляр метода, сделанный на основе структуры. 
			case *MyCustErr: 
				v.First = 0
				fmt.Println(err)
			default:
				fmt.Println(err)
			}
			err.(type) == MyCustErr {
			  err.First = 0
			  fmt.Println(err)
			}
	}
}
