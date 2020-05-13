package main

import (
	"bufio"
	"fmt"
	"github.com/miekg/dns"
	"log"
	"net"
	"os"
	"strings"
)

//func ZoneTransferZ(fqdn string, NSs []string) {
func ZoneTransfer(zone Zone) {
	fqdn := dns.Fqdn(zone.fqdn)

	for _, server := range zone.ns {
		msg := new(dns.Msg)
		msg.SetAxfr(fqdn)

		transfer := new(dns.Transfer)
		answerChan, err := transfer.In(msg, net.JoinHostPort(server, "53"))
		if err != nil {
			log.Println(err)
			continue
		}

		for envelope := range answerChan {
			if envelope.Error != nil {
				log.Println(envelope.Error)
				break
			}

			for _, rr := range envelope.RR {
				switch v := rr.(type) {
				case *dns.A:
					//results.Add(strings.TrimRight(v.Header().Name, "."), v.A.String())
					fmt.Println(strings.TrimRight(v.Header().Name, "."), v.A.String())
				case *dns.AAAA:
					//	results.Add(strings.TrimRight(v.Header().Name, "."), v.AAAA.String())
					fmt.Println(strings.TrimRight(v.Header().Name, "."), v.AAAA.String())
				default:
				}
			}
		}
	}
}

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

	zoneMe := Zone{fqdn: "zonetransfer.me", fail: false}
	zoneMe = getNS(zoneMe)
	fmt.Println(zoneMe)
	fmt.Println(zoneMe.fqdn)
	fmt.Println(zoneMe.ns)

	ZoneTransfer(zoneMe)

	os.Exit(1)

	for c := 0; c < WORKERCOUNT; c++ {
		go worker(jobs, results)
	}

	for k, v := range zones {
		fmt.Println(k, v)
		jobs <- *v
	}
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
