package main

import (
	"bufio"
	"fmt"
	"github.com/miekg/dns"
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
	m := new(dns.Msg)
	for scanner.Scan() {
		line := scanner.Text()
		fmt.Println("testing:", line)
		//m.SetQuestion("miek.nl.", dns.TypeMX)
		resp := m.SetQuestion(dns.Fqdn(line), dns.TypeNS)
		fmt.Println(resp)
	}
	if err := scanner.Err(); err != nil {
		//return err
		fmt.Println(err)
	}
}
