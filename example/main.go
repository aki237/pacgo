package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"pacgo"
)

func main() {
	// commandline check
	if len(os.Args) != 2 {
		fmt.Println("wrong number of arguments")
		return
	}

	// Read the file
	bs, err := ioutil.ReadFile(os.Args[1])
	if err != nil {
		fmt.Println(err)
		return
	}

	// Parsing...
	p := pacgo.NewParser()
	pb, err := p.Parse(string(bs))
	if err != nil {
		fmt.Println(err)
		return
	}

	// Printing package names.
	for index, name := range pb.Pkgnames {
		fmt.Printf("%d. %s\n", index+1, name)
	}
}
