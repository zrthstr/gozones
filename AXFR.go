package main

import (
	"fmt"
	"github.com/miekg/dns"
	"log"
	"net"
	"strings"
	"sync"
)

type Results []Result
type Result struct {
	Domain string   `json:"domain"`
	Addrs  []string `json:"addrs"`
}

func NewResultSet() *ResultSet {
	return &ResultSet{
		results: make(map[string]map[string]struct{}),
	}
}

type ResultSet struct {
	mu      sync.Mutex
	results map[string]map[string]struct{}
}

func main() {
	//Results := ZoneTransfer("heise.de")
	//Results := ZoneTransfer("zonetransfer.me")
	ZoneTransfer("zonetransfer.me")
}

//func ZoneTransfer(domain string) Results {
func ZoneTransfer(domain string) {
	results := NewResultSet()
	fqdn := dns.Fqdn(domain)

	servers, err := net.LookupNS(domain)
	if err != nil {
		log.Fatal(err)
	}

	for _, server := range servers {
		msg := new(dns.Msg)
		msg.SetAxfr(fqdn)

		transfer := new(dns.Transfer)
		answerChan, err := transfer.In(msg, net.JoinHostPort(server.Host, "53"))
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
				//case *dns.AAAA:
				//	results.Add(strings.TrimRight(v.Header().Name, "."), v.AAAA.String())
				default:
				}
			}
		}
	}

	//return results.Results()
	fmt.Println(results)
}
