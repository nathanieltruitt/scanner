package main

import (
	"flag"

	"github.com/nathanieltruitt/scanner/scanner"
)

// we're defining our flags inside this struct.
type Flags struct {
	Address string
}

// declare flags struct
var flags Flags

func setFlags() {
	flag.StringVar(&flags.Address, "A", "", "Attempt to contact an IP")
}


func main() {
	// set our flags
	setFlags()
	// parse our flags
	flag.Parse()

	if flags.Address != "" {
		scanner.IcmpPing(flags.Address)
	}
}