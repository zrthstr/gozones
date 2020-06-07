package main

import (
	"fmt"
)

func main() {
	rslice := [...]string{"first", "sec", "third", "last"}
	fmt.Println(rslice[0])
	fmt.Println(rslice[len(rslice)-1])
}
