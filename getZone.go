package main

import (
	"bufio"
	"fmt"
	"github.com/miekg/dns"
	"log"
	"net"
	"os"
	"path/filepath"
	"strings"
)

type Zone struct {
	fqdn      string
	fail      bool
	ns        []string
	zone      map[string]string
	zoneClean []string
}

type Zones map[string]*Zone

const BUFFERSIZE int = 10000
const WORKERCOUNT int = 200
const DOMAINFILE string = "data/tld_clean.lst"
const OUTDIR string = "data/zones/"
const MAXSORTLEN = 10000

func main() {
	flushOldZones()
	domains := []string{}
	domains, err := fileToList(DOMAINFILE, domains)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	zones := make(Zones)

	for _, domain := range domains {
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
	}
}

func flushOldZones() {
	if _, err := os.Stat(OUTDIR); !os.IsNotExist(err) {
		err := os.RemoveAll(OUTDIR)
		if err != nil {
			log.Println(err)
		}
	}
	err := os.Mkdir(OUTDIR, 0777)
	if err != nil {
		log.Println(err)
	}
}

func writeZone(zone Zone) error {
	if strings.Contains(zone.fqdn, "..") {
		log.Println("Dir traversal detected. skipping:", zone.fqdn)
		return nil
	}
	outArray := []string{""}
	for _, dir := range strings.Split(zone.fqdn, ".") {
		outArray = append([]string{dir}, outArray...)
		//fmt.Println(dir)
	}
	base := strings.Split(OUTDIR, "/")
	outArray = append(base, outArray...)
	fmt.Println("outArray:", outArray)
	//outFile = OUTDIR + strings.TrimSuffix(outFile, "/")
	//outFile := filepath.Split(OUTDIR) + filepath.Join(outArray...)
	outFile := filepath.Join(outArray...)
	outDir := filepath.Dir(outFile)
	fmt.Println("outFile:", outFile)
	fmt.Println("outDir:", outDir)

	err := os.MkdirAll(outDir, 0777)
	if err != nil {
		log.Println(err)
		return err
	}
	f, err := os.Create(outFile)
	if err != nil {
		log.Println(err)
		return err
	}

	//defer f.Close()
	//w := bufio.NewWriter(f)
	//_, err = w.WriteString( v.zoneClean)
	//_, err = w.WriteString(strings.Join(v.zoneClean, "\n"))
	fmt.Println("lenlenlen: ", len(zone.zoneClean))
	n, err := f.WriteString(strings.Join(zone.zoneClean, "\n"))
	if err != nil {
		log.Println(err)
		return err
	}
	log.Println("nnnnnnnn:", n)
	f.Sync()
	//w.Flush()

	//err = w.Err()
	//if err != nil {
	//	log.Println(err)
	//	continue
	//}
	f.Close()
	return err
}

func (zone Zone) String() string {
	out := fmt.Sprintf("[ zone: %s ]\n", zone.fqdn)
	out += fmt.Sprintf("fail.......: %t\n", zone.fail)
	out += fmt.Sprintf("ns.........: %s\n", zone.ns)
	out += fmt.Sprintf("zone.......: %s\n", zone.zone)
	out += fmt.Sprintf("zoneClean..: (%d)\n", len(zone.zoneClean))
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
			log.Println(err)
			continue
		}
		for envelope := range answerChan {
			if envelope.Error != nil {
				log.Println(envelope.Error)
				break
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
		//if len(allLines) > MAXSORTLEN {
		//	log.Println("Long zone detected: ", zone.fqdn, " ", len(allLines))
		//	return []string{"toooo long:", zone.fqdn}
		//}
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

	log.Println("XXXXXXXXXXXXx", len(zone.zoneClean))
	return zone
}

func worker(jobs <-chan Zone, results chan<- Zone) {
	for n := range jobs {
		n = getNS(n)
		if !n.fail {
			n = ZoneTransfer(n)
		}
		_ = writeZone(n)
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
