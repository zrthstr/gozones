package main

import (
	"fmt"
	"strings"
)

func main() {
	foo := "testXXX"
	if strings.HasSuffix(foo, "XX") {
		fmt.Println("true")
	}
}
