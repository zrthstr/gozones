package dns

import (
	"fmt"
	"github.com/miekg/dns"
)

func main(){
	fmt.Println("jooo")
	mx, err := dns.NewRR("miek.nl. 3600 IN MX 10 mx.miek.nl.")
	fmt.Println(mx, err)
}
