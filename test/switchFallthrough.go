package main

import (
	"fmt"
	"strings"
)

func main() {
	str := "foooXXX"
	switch {
	case strings.HasSuffix(str, "X"):
		fmt.Println("ends with X")
		fallthrough
	case strings.HasSuffix(str, "XX"):
		fmt.Println("ends with XX")
		fallthrough
	case strings.HasSuffix(str, "XXX"):
		fmt.Println("ends with XXX")
	}
}
