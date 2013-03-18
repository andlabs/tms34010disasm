// 17 march 2013
package main

import (
	"fmt"
	"os"
	"flag"
	"strconv"
	"io/ioutil"
	"bytes"
	"encoding/binary"
)

var bytes []byte
var words []uint16
var loadaddr uint32
var instructions map[uint32]string
var labels map[uint32]string

var vectorLocs = map[uint32]string{
	0xFFFFFFE0:	"EntryPoint",
	0xFFFFFFC0:	"INT1",
	0xFFFFFFA0:	"INT2",
	0xFFFFFF80:	"Trap3",
	0xFFFFFF60:	"Trap4",
	0xFFFFFF40:	"Trap5",
	0xFFFFFF20:	"Trap6",
	0xFFFFFF00:	"Trap7",
	0xFFFFFEE0:	"NMI",
	0xFFFFFEC0:	"HostInterrupt",
	0xFFFFFEA0:	"DisplayInterrupt",
	0xFFFFFE80:	"WindowViolation",
	0xFFFFFE60:	"Trap12",
	0xFFFFFE40:	"Trap13",
	0xFFFFFE20:	"Trap14",
	0xFFFFFE00:	"Trap15",
	0xFFFFFDE0:	"Trap16",
	0xFFFFFDC0:	"Trap17",
	0xFFFFFDA0:	"Trap18",
	0xFFFFFD80:	"Trap19",
	0xFFFFFD60:	"Trap20",
	0xFFFFFD40:	"Trap21",
	0xFFFFFD20:	"Trap22",
	0xFFFFFD00:	"Trap23",
	0xFFFFFCE0:	"Trap24",
	0xFFFFFCC0:	"Trap25",
	0xFFFFFCA0:	"Trap26",
	0xFFFFFC80:	"Trap27",
	0xFFFFFC60:	"Trap28",
	0xFFFFFC40:	"Trap29",
	0xFFFFFC20:	"IllegalOpcode",
	0xFFFFFC00:	"Trap31",
}

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

	bytes, err = ioutil.ReadFile(flag.Arg(0))
	if err != nil {
		error("error reading input file %s: %v", flag.Arg(0), err)
	}
	if (len(bytes) % 2) != 0 {
		error("given input file %s ends in the middle of a word (this restriction may be removed in the future)", flag.Arg(0))
	}
	if (la + uint64(len(bytes))) < 0x100000000 {
		error("given input file %s does not provide a complete interrupt vector table (this restriction may be removed in the future)", flag.Arg(0))
	}
	words = make([]uint16, len(bytes) / 2)
	// TODO mortal kombat is little endian but are all TMS34010 ROMs? if not I will need to add a command line flag (probably)
	err = binary.Read(bytes.NewReader(b), binary.LittleEndian, &words)
	if err != nil {
		error("error building words array from input byte stream: %v", err)
	}
	// TODO do we just assume the slice was filled properly?

	instructions = map[uint32]string{}
	labels = map[uint32]string{}

	// autoanalyze vectors
	for addr, label := range vectorLocs {
		if labels[addr] != "" {		// if already defined as a different vector, concatenate the labels to make sure everything is represented
			// TODO because this uses a map, it will not be in vector order
			labels[addr] = labels[addr] + "_" + label
		} else {
			labels[addr] = label
		}
		addr /= 2		// words
		pos := (uint32(words[addr + 1]) << 16) | uint32(words[addr])	// little endian longword
		// TODO handle bad addresses here or in disassemble()?
		disassemble(pos)
	}

	// TODO read additional starts from standard input
	// TODO print
}
