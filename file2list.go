package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	fmt.Println("start")
	domains := []string{}

	err, domains := fileToList("tld_clean.lst", domains)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(domains)
}

func fileToList(fileName string, to []string) (error, []string) {
	fileIn, err := os.Open(fileName)
	if err != nil {
		return err, nil
	}
	defer fileIn.Close()

	scanner := bufio.NewScanner(fileIn)
	for scanner.Scan() {
		line := scanner.Text()
		to = append(to, line)
	}

	if err := scanner.Err(); err != nil {
		return err, nil
	}

	return nil, to

}
