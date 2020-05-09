package main

import (
	"bufio"
	"fmt"
	//	"github.com/miekg/dns"
	"os"
)

func main() {
	domainFile := "tld_clean.lst"
	domains := []string{}
	domains, err := fileToList(domainFile, domains)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	for i, domain := range domains {
		fmt.Println(i, domain)
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
