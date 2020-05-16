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
const WORKERCOUNT int = 200
const DOMAINFILE string = "tld_clean.lst"

func main() {
	domains := []string{}
	domains, err := fileToList(DOMAINFILE, domains)
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

	for c := 0; c < WORKERCOUNT; c++ {
		go worker(jobs, results)
	}

	//counter := 0
	for k, v := range zones {
		fmt.Println(k, v)
		//if counter > 100 {
		//	break
		//}
		//counter++
		jobs <- *v
	}
	//close(jobs)
	for i := 0; i < len(zones); i++ {
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
	} else {
		zone.fail = false
	}
	for _, ns := range nameserver {
		zone.ns = append(zone.ns, ns.Host)
	}
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
