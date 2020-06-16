package main

import (
	"flag"
	"fmt"
)

func main() {
	//num := flag.Int("numb", 23, "any int")
	//fmt.Println("numb: ", *num)

	ver := flag.Bool("version", false, "print version")

	numb := flag.Int("numb", -1, "some positive number")

	flag.Parse()

	if *ver {
		fmt.Println("V0.0.1")
	}
	fmt.Println("version:", *ver)

	// testing if a flag has been set ot not is not trivial :/
	if *numb > 0 {
		fmt.Println("numb: ", *numb)
	}

}
