package main

import "fmt"

func main() {
	x := 5.3
	p := &x
	fmt.Println(p)
	fmt.Println(*p)
}
