package main

import (
	"fmt"
)

type Foo struct {
	foo  string
	bar  bool
	spam []string
	ham  map[string]string
}

func (f Foo) String() string {
	//return fmt.Println("foo:", f.foo)
	var out string
	out += "==== this is a foo ====\n"
	out += fmt.Sprintf("foo:  %s\n", f.foo)
	out += fmt.Sprintf("bar:  %t\n", f.bar)
	out += fmt.Sprintf("spam: %s\n", f.spam)
	out += fmt.Sprintf("ham:  %s\n", f.ham)
	out += "---- end of this ----"
	return out
}

func main() {
	var foo = Foo{foo: "foofofo",
		bar:  true,
		spam: []string{"a", "b", "c", "d"},
		ham:  map[string]string{"some": "thing", "other": "thing"},
	}
	fmt.Println(foo)
}
