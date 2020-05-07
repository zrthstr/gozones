package main

import (
	"fmt"
	"io"
	//"io/ioutil"
	"bufio"
	"net/http"
	"os"
	"regexp"
)

// curl $TLDURL ; grep -v '^$' | grep -v '//' | sed 's/.*//' | sort | uniq > out

func main() {
	fileUrl := "https://publicsuffix.org/list/effective_tld_names.dat"
	filePathRaw := "tld_raw.lst"
	filePathClean := "tld_clean.lst"
	if err := GetTLD(fileUrl, filePathRaw); err != nil {
		panic(err)
	}
	if err := CleanTLDFileX(filePathRaw, filePathClean); err != nil {
		panic(err)
	}
}

func GetTLD(url string, filePath string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	out, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	return err
}

func CleanTLDFileX(fileInPath string, fileOutPath string) error {
	// Clean the zone file
	// get rid of trailing wild char '.*'
	// get rid of commented lines i.e. line starting with //
	// get rid of empty lines

	fileIn, err := os.Open(fileInPath)
	if err != nil {
		return err
	}
	defer fileIn.Close()

	fileOut, err := os.Create(fileOutPath)
	if err != nil {
		return err
	}
	defer fileOut.Close()

	scanner := bufio.NewScanner(fileIn)
	wild := regexp.MustCompile(`\*\.`)

	for scanner.Scan() {
		line := scanner.Text()
		line = wild.ReplaceAllString(line, "")
		match, _ := regexp.MatchString("//", line)
		if match {
			continue
		}
		match, _ = regexp.MatchString("^$", line)
		if match {
			continue
		}
		// l, err := f.WriteString("Hello World")
		lineN := line + "\n"
		fileOut.WriteString(lineN)
		//fmt.Println(line)
	}
	if err := scanner.Err(); err != nil {
		return err
	}
	fmt.Println(fileOutPath)
	return nil
}
