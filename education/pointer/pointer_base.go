package main

import "fmt"

func main() {
	var x int = 4
	var p *int
	p = &x
	fmt.Println(x)
	fmt.Println(p)
	fmt.Println(*p)

}
