package main

import (
	"fmt"
)

type Zone struct {
	fqdn string
	ns   []string
	fail bool
	zone string
}

func main() {
	fmt.Println("dd")
	var m map[string]Zone
	m = make(map[string]Zone)

	fmt.Println(m)

	domains := []string{"google.com", "foo.bar", "facebook.com"}
	for c, s := range domains {
		fmt.Println(c, s)
		m[s] = Zone{fqdn: s}
	}

	foo := Zone{fqdn: "google.com."}
	foo.fail = false

	fmt.Println(foo)
	fmt.Println(domains)
	fmt.Println(m)

}
