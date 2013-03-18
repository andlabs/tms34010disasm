// 17 march 2013
package main

import (
	"fmt"
	"os"
	"flag"
	"strconv"
	"io/ioutil"
)

var fin []uint16
var loadaddr uint32
var instructions map[uint32]string
var labels map[uint32]string

// command line options
var (
	mirrorAddr = flag.Uint64("mirror", 0xFFFFFFFF, "bit address of ROM mirror (used for games like Mortal Kombat); 0xFFFFFFFF to disable")
)

func error(format string, args ...interface{}) {
	fmg.Fprintf(os.Stderr, format, args...)
	os.Exit(1)
}

func usage() {
	fmt.Fprintf(os.Stderr, "usage: %s [-mirror bitaddress] ROM load-bitaddress"
	flag.PrintDefaults()
	fmt.Fprintf(os.Stderr, "addresses are given in C notation (for instance, 0x12345670, not >12345670)\n")
	os.Exit(1)
}

func main() {
	flag.Usage = usage
	flag.Parse()
	if flag.NArg() != 2 {
		usage()
	}

	if *mirrorAddr != 0xFFFFFFFF {
		if *mirrorAddr > 0xFFFFFFFF {
			error("given mirror address 0x%08X out of range", *mirrorAddr)
		}
		if (*mirrorAddr & 0xF) != 0 {
			error("given mirror address 0x%08X not aligned to word boundary (this restriction may be removed in the future)", *mirrorAddr)
		}
	}

	filename := flag.Arg(0)
	la, err := strconv.ParseUInt(flag.Arg(1), 0, 32)
	if err != nil {
		error("error reading given load address: %v", err)
	}
	loadaddr = uint32(la)
	if (loadaddr & 0xF) != 0 {
		error("given load address 0x%08X not aligned to word boundary (this restriction may be removed in the future)", loadaddr)
	}

	b, err := ioutil.ReadFile(flag.Arg(0))
	if err != nil {
		error("error reading input file %s: %v", flag.Arg(0), err)
	}
	if (len(b) % 2) != 0 {
		error("given input file %s ends in the middle of a word (this restriction may be removed in the future)", flag.Arg(0))
	}
	if (la + uint64(len(b))) < 0x100000000 {
		error("given input file %s does not provide a complete interrupt vector table (this restriction may be removed in the future)", flag.Arg(0))
	}

	instructions = map[uint32]string{}
	labels = map[uint32]string{}

	// TODO handle instruction vectors
	// TODO read additional starts from standard input
	// TODO print
}
