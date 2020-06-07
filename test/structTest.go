package main

import (
	"fmt"
)

type SomeSt struct {
	count int
	msg   []string
}

func main() {
	randMsg := [...]string{"some", "some", "other"}
	someSt := SomeSt{count: 0}

	for _, m := range randMsg {
		fmt.Println(m)
		someSt.msg = append(someSt.msg, m)
	}

	fmt.Println(someSt)
	fmt.Printf("%+q\n", someSt.msg)
	fmt.Println(len(someSt.msg))
}
