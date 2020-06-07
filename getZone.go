package main

import (
	"bufio"
	"fmt"
	"github.com/miekg/dns"
	"log"
	"net"
	"os"
	"path/filepath"
	"reflect"
	"strings"
)

type Zone struct {
	fqdn      string
	fail      bool
	ns        []string
	zone      map[string]string
	zoneClean []string
	errMsg    []string
}

type Zones map[string]*Zone

type ZoneErrors struct {
	okey   int
	count  int
	errMsg map[string]int
}

//func (errMsg string, zoneErrors ZoneErrors) errLog() {
//	zoneErrors.count += 1
//	fmt.Println("errMsg", errMsg)
//	//zoneErrors.errMsg[errMsg] += 1
//}

const BUFFERSIZE int = 10000

//const WORKERCOUNT int = 200
const WORKERCOUNT int = 50
const DOMAINFILE string = "data/tld_clean.lst"
const OUTDIR string = "data/zones/"
const MAXSORTLEN = 10000

func main() {
	println(dns.RcodeNameError)
	flushOldZones()
	domains := []string{}
	domains, err := fileToList(DOMAINFILE, domains)
	if err != nil {
		log.Println("0", err)
		os.Exit(1)
	}

	zones := make(Zones)
	zoneErrors := ZoneErrors{okey: 0, count: 0, errMsg: make(map[string]int)}
	//zoneErrors.errMsg["ddddd"] = 1
	//fmt.Println(zoneErrors)

	for _, domain := range domains {
		//zone := &Zone{fqdn: domain, fail: false, errMsg: make([]string)}
		zone := &Zone{fqdn: domain, fail: false}
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
		//zoneErrors = noteStats(thisZone, &zoneErrors)
		noteStats(thisZone, &zoneErrors)
	}
	printStats(zoneErrors)
}

func noteStats(zone Zone, zoneErrors *ZoneErrors) {
	for _, e := range zone.errMsg {
		fmt.Println("XXXXXXXXXXXXXXXXXXXXXXXXx", e)
		_, exists := zoneErrors.errMsg[e]
		if exists {
			zoneErrors.errMsg[e] += 1
		} else {
			zoneErrors.errMsg[e] = 1
		}
	}
}

func printStats(zoneErrors ZoneErrors) {
	log.Println("--------stat--------")
	for n, e := range zoneErrors.errMsg {
		fmt.Printf("%5d: %s\n", e, n)
	}
	log.Println("--------end.--------")
}

func flushOldZones() {
	if _, err := os.Stat(OUTDIR); !os.IsNotExist(err) {
		err := os.RemoveAll(OUTDIR)
		if err != nil {
			log.Println("1", err)
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
	//base := strings.Split(OUTDIR, "/")
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

func (zone Zone) String() string {
	out := fmt.Sprintf("[ zone: %s ]\n", zone.fqdn)
	out += fmt.Sprintf("fail.......: %t\n", zone.fail)
	out += fmt.Sprintf("ns.........: %+q\n", zone.ns)
	out += fmt.Sprintf("zone.......: %s\n", zone.zone)
	out += fmt.Sprintf("zoneClean..: (%d)\n", len(zone.zoneClean))
	out += fmt.Sprintf("errMsg.....: %+q\n", zone.errMsg)
	for _, v := range zone.zoneClean {
		out = out + "\t" + v + "\n"
	}
	return out
}

//func ZoneTransferZ(fqdn string, NSs []string) {
func ZoneTransfer(zone Zone) Zone {
	zone.zone = make(map[string]string)
	fqdn := dns.Fqdn(zone.fqdn)

	for _, server := range zone.ns {
		msg := new(dns.Msg)
		msg.SetAxfr(fqdn)

		transfer := new(dns.Transfer)
		answerChan, err := transfer.In(msg, net.JoinHostPort(server, "53"))
		if err != nil {
			log.Println("6", err)
			log.Println(reflect.TypeOf(err).String())
			continue
		}
		for envelope := range answerChan {
			if envelope.Error != nil {
				errMsg := envelope.Error.Error()
				//fmt.Println(envelope.Error.Error())
				zone.errMsg = append(zone.errMsg, errMsg)
				switch errMsg {
				case "dns: bad xfr rcode: 5":
					//log.Println("..5")
				case "dns: bad xfr rcode: 9":
					//log.Println("..9")
				default:
					log.Println("7", envelope.Error.Error())
					//log.Println("Other error:", reflect.TypeOf(envelope.Error.Error).String())
					// https://www.iana.org/assignments/dns-parameters/dns-parameters.xhtml#dns-parameters-6
				}
				//break //continue /// why was this break?
			}
			for _, rr := range envelope.RR {
				zone.zone[server] += "\n" + rr.String()
			}
		}
	}
	//return zone
	// deduplicate all answers from different NameServers and store nicely in array
	zone.zoneClean = func(allZones map[string]string) []string {
		allLines := ""
		for _, v := range allZones {
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
		//for _, v := range dedupLines {
		//	fmt.Println(v)
		//}
		//os.Exit(1)
		return dedupLines
	}(zone.zone)
	return zone
}

func worker(jobs <-chan Zone, results chan<- Zone) {
	for n := range jobs {
		n = getNS(n)
		if !n.fail {
			n = ZoneTransfer(n)
		}
		if len(n.zoneClean) > 1 {
			_ = writeZone(n)
		}
		results <- n
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
