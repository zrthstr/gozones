package main

import (
	"bufio"
	"fmt"
	"github.com/miekg/dns"
	"net"
	"os"
)

func main() {
	fmt.Println("jooo")
	mx, err := dns.NewRR("miek.nl. 3600 IN MX 10 mx.miek.nl.")
	fmt.Println(mx, err)

	fileIn, err := os.Open("tld_clean.lst")
	if err != nil {
		//return err
		fmt.Println(err)
	}
	defer fileIn.Close()

	scanner := bufio.NewScanner(fileIn)
	for scanner.Scan() {
		line := scanner.Text()
		fmt.Println("testing:", line)
		nameserver, err := net.LookupNS(line)
		if err != nil {
			//return err
			fmt.Println(err)
			continue
		}
		fmt.Println(nameserver)
	}

	if err := scanner.Err(); err != nil {
		//return err
		fmt.Println(err)
	}

}
