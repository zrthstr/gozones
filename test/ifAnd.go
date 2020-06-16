package main

import (
	"fmt"
)

func main() {
	if true && true {
		fmt.Println("0")
	}
	if !false && true {
		fmt.Println("1")
	}
	if !!true && !false {
		fmt.Println("3")
	}
	if true || false {
		fmt.Println("4")
	}
	if true || true {
		fmt.Println("5")
	}
}
