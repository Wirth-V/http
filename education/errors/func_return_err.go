package main

import (
	"errors"
	"fmt"
)

func Vibor(name string) (string, error) {
	if name != "VV" {
		return name, errors.New("Don't correct vibor")
	}
	return name, nil
}

func main() {
	win, err := Vibor("AA")

	if err != nil {
		fmt.Printf("VV\n%s\n", err)
	}

	fmt.Println(win)
}
