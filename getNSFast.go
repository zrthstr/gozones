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

const BUFFERSIZE int = 10000
const CHANCOUNT int = 20

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

	jobs := make(chan Zone, BUFFERSIZE)
	results := make(chan Zone, BUFFERSIZE)

	for c := 0; c < CHANCOUNT; c++ {
		go worker(jobs, results)
	}
	go worker(jobs, results)
	go worker(jobs, results)
	go worker(jobs, results)
	go worker(jobs, results)

	//fmt.Println(zones)
	for k, v := range zones {
		fmt.Println(k, v)
		//getNS(*v)
		jobs <- *v
	}
	for {
		foo := <-results
		fmt.Println(foo)
	}
}

func worker(jobs <-chan Zone, results chan<- Zone) {
	for n := range jobs {
		results <- getNS(n)
	}
}

func getNS(zone Zone) Zone {
	nameserver, err := net.LookupNS(zone.fqdn)
	if err != nil {
		zone.fail = true
	}
	//answer := []string{}
	for _, ns := range nameserver {
		//answer = append(answer, ns.Host)
		zone.ns = append(zone.ns, ns.Host)
	}
	//zone.ns = answer
	return zone
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
