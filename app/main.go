package main

import (
	"fmt"
	"toolkit"
)

func main() {
	var tools toolkit.tools

	s := tools.RandomString(10)
	fmt.Println("Random string:", s)
}
