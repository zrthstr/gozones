package main

import (
	"fmt"
	"strings"
)

func main() {
	str := "fooooXXX"
	switch {
	case strings.HasSuffix(str, "XXX"):
		fmt.Println("XXX")
	case strings.HasSuffix(str, "YYY"):
		fmt.Println("YYY")
	}
}
