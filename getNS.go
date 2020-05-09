package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

func main() {
	fileIn, err := os.Open("tld_clean.lst")
	if err != nil {
		fmt.Println(err)
	}
	defer fileIn.Close()

	scanner := bufio.NewScanner(fileIn)
	for scanner.Scan() {
		line := scanner.Text()
		fmt.Println("testing:", line)
		nameserver, err := net.LookupNS(line)
		if err != nil {
			//fmt.Println(err)
			continue
		}
		for _, ns := range nameserver {
			fmt.Println(ns)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Println(err)
	}

}
