// 18 march 2013
package main

import (
	"fmt"
)

type Opcode struct {
	Expected		uint16
	Mask		uint16
	// TODO Func
}

var opcodes = []Opcode{
	// ...
}

func ui16(bits ...uint16) uint16 {
	var x uint16
	var i uint
	
	if len(bits) > 16 {
		panic(fmt.Sprintf("too many bits passed to u16 (given %d)", len(bits)))
	}
	for i = 0; i < uint(len(bits)); i++ {
		x <<= 1
		x |= bits[i]
	}
	x <<= 16 - i
	return x
}

func op(bits ...uint16) (o Opcode) {
	o.Expected = ui16(bits...)
	for i := 0; i < len(bits); i++ {
		bits[i] = 1		// TODO whether or not this should be legal should be in the spec
	}
	o.Mask = ui16(bits...)
	// TODO func
	return o
}
