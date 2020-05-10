package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	fmt.Println("start")
	domains := []string{}

	domains, err := fileToList("tld_clean.lst", domains)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(domains)
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
