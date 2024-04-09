package main

import "fmt"

func main() {
	type Circle struct {
		x, y, r float64
	}

	c := Circle{0, 0, 5}

	fmt.Println(c)

}
