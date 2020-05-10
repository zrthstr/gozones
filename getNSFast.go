package main

import (
	"bufio"
	"fmt"
	//	"github.com/miekg/dns"
	"net"
	"os"
)

type Zone struct {
	fqdn string
	fail bool
	ns   []string
	zone []string
}

func main() {
	domainFile := "tld_clean.lst"
	domains := []string{}
	domains, err := fileToList(domainFile, domains)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	zones := make(map[string]*Zone)
	for _, domain := range domains {
		zone := &Zone{fqdn: domain, fail: false}
		zones[zone.fqdn] = zone
	}

	fmt.Println(zones)
	for k, v := range zones {
		fmt.Println(k, v)
	}

	//for i, domain := range domains {
	//	fmt.Println(i, domain)
	//go getNSx(domain)
	//answer, err := getNS(domain)
	//if err != nil {
	//	fmt.Println("err:", err)
	//	continue
	//}
	//fmt.Println(answer, err)
	//}
	//fmt.Scanln()
}

func getNS(domain string) ([]string, error) {
	nameserver, err := net.LookupNS(domain)
	if err != nil {
		return nil, err
	}
	answer := []string{}
	for _, ns := range nameserver {
		//fmt.Println(ns)
		answer = append(answer, ns.Host)
	}
	return answer, nil
}

func getNSx(domain string) {
	nameserver, err := net.LookupNS(domain)
	if err != nil {
		fmt.Println("err:", domain, err)
	}
	answer := []string{}
	for _, ns := range nameserver {
		//fmt.Println(ns)
		answer = append(answer, ns.Host)
	}
}

func fileToList(fileName string, to []string) ([]string, error) {
	fileIn, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer fileIn.Close()

	scanner := bufio.NewScanner(fileIn)
	for scanner.Scan() {
		line := scanner.Text()
		to = append(to, line)
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return to, nil
}
