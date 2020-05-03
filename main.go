package main

import (
	"fmt"
	"net/http"
	"os"
	"io"
)

func main(){
	fmt.Println("hello")
	if err := GetTLDS(); err != nil {
		panic(err)
	}
}

func GetTLDS() error {
	fileUrl := "https://publicsuffix.org/list/effective_tld_names.dat"
	filePath := "tld.dat"

	resp, err := http.Get(fileUrl)
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
