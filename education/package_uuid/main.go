package main

import (
	"fmt"

	"github.com/google/uuid"
)

func main() {
	// Creating UUID Version 7
	// panic on error
	u1 := uuid.Must(uuid.NewV7())
	fmt.Printf("UUIDv7 через uuid.Must: %s\n", u1)

	// or error handling
	u2, err := uuid.NewV7()
	if err != nil {
		fmt.Printf("Something went wrong: %s", err)
		return
	}
	fmt.Printf("UUIDv7 через uuid.NewV7: %s\n", u2)

	u3 := uuid.New().String()
	fmt.Printf("UUIDv7 через uuid.New().String(): %s\n", u3)

	u4 := uuid.New().String()[:8]
	fmt.Printf("UUIDv7 через uuid.New().String()[:8]: %s\n", u4)

	idString := "f47ac10b-58cc-0372-8567-0e02b2c3d479"

	// Разбор строки в UUID.
	u5, err := uuid.Parse(idString)
	if err != nil {
		fmt.Println("Ошибка разбора UUID:", err)
		return
	}

	// Вывод разобранного UUID.
	fmt.Println("Разобранный UUID:", u5)

}
