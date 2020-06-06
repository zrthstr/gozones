package main

import (
	"fmt"
)

type Stat struct {
	count int
	stat  map[string]int
}

func main() {
	//errStat := make(Stat)
	stat := Stat{count: 0, stat: make(map[string]int)}
	errs := [...]string{"aaa", "bbb", "ccc", "ddd", "aaa", "ccc"}
	for _, e := range errs {
		_, exists := stat.stat[e]
		if exists {
			stat.stat[e] += 1
			fmt.Println("DUP")
		} else {
			stat.stat[e] = 1
			fmt.Println("nonDUP")
		}
		//fmt.Println(e)
		//fmt.Println(stat)
	}
}
