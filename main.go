package main

import (
	"flag"

	"github.com/nathanieltruitt/scanner/scanner"
)

// we're defining our flags inside this struct.
type Flags struct {
	Address string
	Start int
	End int
}

// declare flags struct
var flags Flags

func setFlags() {
	flag.StringVar(&flags.Address, "A", "", "specify subnet")
	flag.IntVar(&flags.Start, "s", 1, "specify a starting IP range")
	flag.IntVar(&flags.End, "e", 254, "specify ending IP range")
}


func main() {
	// set our flags
	setFlags()
	// parse our flags
	flag.Parse()

	if flags.Address != "" {
		scanner.PingScan(flags.Address, flags.Start, flags.End)
	}
}