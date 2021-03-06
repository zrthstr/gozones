// TODO:
//  make sure when one server gives zone not to ask next for same domain

package main

import (
	"bufio"
	"fmt"
	"github.com/miekg/dns"
	"log"
	"net"
	"os"
	"path/filepath"
	//"reflect"
	"flag"
	"strings"
)

type Zone struct {
	fqdn      string
	nsFail    bool
	axFail    bool
	ns        []string
	zone      map[string]string
	zoneClean []string
	errMsg    []string
	noAx      bool
}

type Zones map[string]*Zone

type ZoneErrors struct {
	total  int
	nsFail int
	axFail int
	count  int
	errMsg map[string]int
}

const BUFFERSIZE int = 10000
const WORKERCOUNT int = 112
const DOMAINFILE string = "data/tld_clean.lst"
const OUTDIR string = "data/zones/"
const MAXSORTLEN = 10000

func main() {
	noAx := flag.Bool("noax", false, "no axfr, just get ns")
	flag.Parse()
	if *noAx {
		fmt.Println("Skipping AXFR", *noAx)
	}

	flushOldZones()
	domains := []string{}
	domains, err := fileToList(DOMAINFILE, domains)
	if err != nil {
		log.Println("0", err)
		os.Exit(1)
	}

	zones := make(Zones)
	zoneErrors := ZoneErrors{total: 0, axFail: 0, nsFail: 0, count: 0, errMsg: make(map[string]int)}

	for _, domain := range domains {
		//zone := &Zone{fqdn: domain, fail: false, errMsg: make([]string)}
		zone := &Zone{fqdn: domain, nsFail: false, axFail: false, noAx: *noAx}
		zones[zone.fqdn] = zone
	}

	jobs := make(chan Zone, BUFFERSIZE)
	results := make(chan Zone, BUFFERSIZE)

	for c := 0; c < WORKERCOUNT; c++ {
		go worker(jobs, results)
	}
	for _, v := range zones {
		jobs <- *v
	}
	for i := 0; i < len(zones); i++ {
		thisZone := <-results
		zones[thisZone.fqdn] = &thisZone
		//log.Println(<-results)
		log.Println(thisZone)
		noteStats(thisZone, &zoneErrors)
	}
	printStats(zoneErrors, *noAx)
}

func noteStats(zone Zone, zoneErrors *ZoneErrors) {
	zoneErrors.total += 1
	if zone.nsFail {
		zoneErrors.nsFail += 1
	} else if zone.axFail {
		zoneErrors.axFail += 1
	}
	for _, e := range zone.errMsg {
		_, exists := zoneErrors.errMsg[e]
		if exists {
			zoneErrors.errMsg[e] += 1
		} else {
			zoneErrors.errMsg[e] = 1
		}
	}
}

func printStats(zoneErrors ZoneErrors, noAx bool) {
	log.Println("\n--------stat--------")
	fmt.Printf("total:  %4d\n", zoneErrors.total)
	fmt.Printf("nsFail: %4d/%d\n", zoneErrors.nsFail, zoneErrors.total)
	if noAx {
		fmt.Println("axFail:  -- / --")
	} else {
		fmt.Printf("axFail: %4d/%d\n", zoneErrors.axFail, zoneErrors.total-zoneErrors.nsFail)
	}
	fmt.Println("--------------------")
	for n, e := range zoneErrors.errMsg {
		fmt.Printf("%5d: %s\n", e, n)
	}
	fmt.Println("--------end.--------")
}

func flushOldZones() {
	if _, err := os.Stat(OUTDIR); !os.IsNotExist(err) {
		err := os.RemoveAll(OUTDIR)
		if err != nil {
			log.Println("1 ...", err)
		}
	}
	err := os.Mkdir(OUTDIR, 0777)
	if err != nil {
		log.Println("2", err)
	}
}

func writeZone(zone Zone) error {
	if len(zone.zoneClean) == 0 {
		log.Println("Zone %s's zoneClean is empty. Not writing to disk", zone.fqdn)
		return nil
	}
	if strings.Contains(zone.fqdn, "..") {
		log.Println("Dir traversal detected. skipping:", zone.fqdn)
		return nil
	}
	outArray := []string{""}
	for _, dir := range strings.Split("zone."+zone.fqdn, ".") {
		outArray = append([]string{dir}, outArray...)
	}
	outArray = append(strings.Split(OUTDIR, "/"), outArray...)
	outFile := filepath.Join(outArray...)
	outDir := filepath.Dir(outFile)

	err := os.MkdirAll(outDir, 0777)
	if err != nil {
		log.Println("3", err)
		return err
	}
	f, err := os.Create(outFile)
	if err != nil {
		log.Println("4", err)
		return err
	}

	fmt.Println("lenlenlen: ", len(zone.zoneClean))
	n, err := f.WriteString(strings.Join(zone.zoneClean, "\n"))
	if err != nil {
		log.Println("5", err)
		return err
	}
	log.Println("Written: ", n, "lines to file: ", outFile)
	f.Sync()

	f.Close()
	return err
}

// TODO:
// use this again
func (zone Zone) String() string {
	out := fmt.Sprintf("[ zone: %s ] \n", zone.fqdn)
	out += fmt.Sprintf("axFail.....: %t\n", zone.axFail)
	out += fmt.Sprintf("nsFail.....: %t\n", zone.nsFail)
	out += fmt.Sprintf("ns.........: %+q\n", zone.ns)
	out += fmt.Sprintf("zone.......: %s\n", zone.zone)
	out += fmt.Sprintf("zoneClean..: (%d)\n", len(zone.zoneClean))
	out += fmt.Sprintf("errMsg.....: %+q\n", zone.errMsg)
	for _, v := range zone.zoneClean {
		out = out + "\t" + v + "\n"
	}
	return out
}

func ZoneTransfer(zone *Zone) {
	zone.zone = make(map[string]string)
	fqdn := dns.Fqdn(zone.fqdn)

	for _, server := range zone.ns {
		msg := new(dns.Msg)
		msg.SetAxfr(fqdn)
		transfer := new(dns.Transfer)
		answerChan, err := transfer.In(msg, net.JoinHostPort(server, "53"))
		if err != nil {
			// dial tcp 168.95.192.10:53: connect: no route to host
			errMsg := err.Error()
			switch {
			case strings.HasSuffix(errMsg, "connect: no route to host"):
				log.Println("6.5", errMsg)
				errMsg = "ns: connect no route to host"
			case strings.HasSuffix(errMsg, "i/o timeout"):
				log.Println("6.6", errMsg)
				errMsg = "ns: dial tcp i/o timeout"
			case strings.HasSuffix(errMsg, ": no such host"):
				log.Println("6.6", errMsg)
				errMsg = "ns: dail tcp lookup no such host"
			default:
				log.Println("6", err)
			}
			zone.errMsg = append(zone.errMsg, errMsg)
			continue
		}
		// https://www.iana.org/assignments/dns-parameters/dns-parameters.xhtml#dns-parameters-6
		for envelope := range answerChan {
			if envelope.Error != nil {
				errMsg := envelope.Error.Error()
				switch {
				case errMsg == "dns: bad xfr rcode: 5":
				case errMsg == "dns: bad xfr rcode: 9":
				case strings.HasSuffix(errMsg, ": i/o timeout"):
					log.Println("7", errMsg)
					errMsg = "axfr: read tcp i/o timeout"
				case strings.HasSuffix(errMsg, "read: connection reset by peer"):
					log.Println("7", errMsg)
					errMsg = "axfr: read tcp connection reset by peer"
				default:
					log.Println("7", errMsg)
				}
				zone.errMsg = append(zone.errMsg, errMsg)
				//break //continue /// why was this break?
			}
			for _, rr := range envelope.RR {
				zone.zone[server] += "\n" + rr.String()
			}
			// break on first zone for domain
			if len(zone.zone[server]) > 1 {
				break
			}
		}
	}
	// for now not needed as we break ob the firt zone for a domain
	//zone.zoneClean = dedupZone(zone)
}

func dedupZone(zone *Zone) []string {
	// deduplicate all answers from different NameServers and store nicely in array
	allLines := ""
	for _, v := range zone.zone {
		allLines += "\n" + v
	}
	if len(allLines) > MAXSORTLEN {
		log.Println("Long zone detected: ", zone.fqdn, " ", len(allLines))
		allLines = allLines[:MAXSORTLEN] ///////////////////////////////////////////////
	}
	var dedupLines []string
	dedupDict := make(map[string]bool)
	for _, line := range strings.Split(allLines, "\n") {
		if line == "" {
			continue
		}
		dedupDict[line] = true
	}
	for k, _ := range dedupDict {
		dedupLines = append(dedupLines, k)
	}
	return dedupLines
}

func worker(jobs <-chan Zone, results chan<- Zone) {
	for n := range jobs {
		getNS(&n)
		if !n.axFail && !n.noAx {
			ZoneTransfer(&n)
		}
		if len(n.zoneClean) > 1 {
			_ = writeZone(n)
		}
		results <- n
	}
}

// TODO
// make sure zone transfer attempt wont happen
func getNS(zone *Zone) {
	nameserver, err := net.LookupNS(zone.fqdn)
	if err != nil {
		zone.axFail = true
		errMsg := err.Error()
		switch {
		case strings.HasSuffix(errMsg, "no such host"):
			fmt.Println(errMsg)
			errMsg = "ns: no such host"
		case strings.HasSuffix(errMsg, "server misbehaving"):
			fmt.Println(errMsg)
			errMsg = "ns: server misbehaving"
		case strings.HasSuffix(errMsg, "i/o timeout"):
			fmt.Println(errMsg)
			errMsg = "ns: i/o timeout"
		}
		zone.errMsg = append(zone.errMsg, errMsg)
		zone.nsFail = true
	} else {
		zone.nsFail = false
		for _, ns := range nameserver {
			zone.ns = append(zone.ns, ns.Host)
		}
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
