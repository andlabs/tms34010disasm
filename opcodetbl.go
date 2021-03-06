// 18 march 2013
package main

import (
	"fmt"
)

// TODO type opfunc

type Opcode struct {
	Expected		uint16
	Mask		uint16
	Func			opfunc
}

var opcodes = []Opcode{
	op(o_abs,		0, 0, 0, 0, 0, 0, 1, 1, 1, 0, 0),
	op(o_add,		0, 1, 0, 0, 0, 0, 0),
	op(o_addc,	0, 1, 0, 0, 0, 0, 1),
	op(o_addiw,	0, 0, 0, 0, 1, 0, 1, 1, 0, 0, 0),
	op(o_addil,	0, 0, 0, 0, 1, 0, 1, 1, 0, 0, 1),
	op(o_addk,	0, 0, 0, 1, 0, 0),
	op(o_addxy,	1, 1, 1, 0, 0, 0, 0),
	op(o_and,		0, 1, 0, 1, 0, 0, 0),
	// andi is an alias for andni; TODO provide an option?
	op(o_andn,	0, 1, 0, 1, 0, 0, 1),
	op(o_andni,	0, 0, 0, 0, 1, 0, 1, 1, 1, 0, 0),
	op(o_btst,		0, 0, 0, 1, 1, 1),				// btst immediate,register
	op(o_btst,		0, 1, 0, 0, 1, 0, 1),			// btst register,register
	op(o_call,		0, 0, 0, 0, 1, 0, 0, 1, 0, 0, 1),
	op(o_calla,	0, 0, 0, 0, 1, 1, 0, 1, 0, 1, 0, 1, 1, 1, 1, 1),
	op(o_callr,	0, 0, 0, 0, 1, 1, 0, 1, 0, 0, 1, 1, 1, 1, 1, 1),
	op(o_clr,		0, 1, 0, 1, 0, 1, 1),
	op(o_clrc,		0, 0, 0, 0, 0, 0, 1, 1, 0, 0, 1, 0, 0, 0, 0, 0),
	op(o_cmp,	0, 1, 0, 0, 1, 0, 0),
	op(o_cmpiw,	0, 0, 0, 0, 1, 0, 1, 1, 0, 1, 0),
	op(o_cmpil,	0, 0, 0, 0, 1, 0, 1, 1, 0, 1, 1),
	op(o_cmpxy,	1, 1, 1, 0, 0, 1, 0),
	op(o_cpw,	1, 1, 1, 0, 0, 1, 1),
	op(o_cvxyl,	1, 1, 1, 0, 1, 0, 0),
	op(o_dec,		0, 0, 0, 1, 0, 1, 0, 0, 0, 0, 1),
	op(o_dint,	0, 0, 0, 0, 0, 0, 1, 1, 0, 1, 1, 0, 0, 0, 0, 0),
	op(o_divs,	0, 1, 0, 1, 1, 0, 0),
	op(o_divu,	0, 1, 0, 1, 1, 0, 1),
	op(o_drav,	1, 1, 1, 1, 0, 1, 1),
	op(o_dsj,		0, 0, 0, 0, 1, 1, 0, 1, 1, 0, 0),
	op(o_dsjeq,	0, 0, 0, 0, 1, 1, 0, 1, 1, 0, 1),
	op(o_dsjne,	0, 0, 0, 0, 1, 1, 0, 1, 1, 1, 0),
	op(o_dsjs,	0, 0, 1, 1, 1),
	op(o_eint,		0, 0, 0, 0, 1, 1, 0, 1, 0, 1, 1, 0, 0, 0, 0, 0),
	op(o_emu,	0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0),
	op(o_exgf0,	1, 1, 0, 1, 0, 1, 0, 1, 0, 0, 0),
	op(o_exgf1,	1, 1, 0, 1, 0, 1, 1, 1, 0, 0, 0),
	op(o_exgpc,	0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 1),
	op(o_filll,		0, 0, 0, 0, 1, 1, 1, 1, 1, 1, 0, 0, 0, 0, 0, 0),
	op(o_fillxy,	0, 0, 0, 0, 1, 1, 1, 1, 1, 1, 1, 0, 0, 0, 0, 0),
	op(o_getpc,	0, 0, 0, 0, 0, 0, 0, 1, 0, 1, 0),
	op(o_getst,	0, 0, 0, 0, 0, 0, 0, 1, 1, 0, 0),
	op(o_inc,		0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 1),		// TODO this is just an alias for addk 1,rN... I'll need to sort this table by mask decreasing to have it be handled properly
	ja(o_jauc,		0, 0, 0, 0),
	ja(o_jalo,		1, 0, 0, 0),						// aka jac
	ja(o_jals,		0, 0, 1, 0),
	ja(o_jahi,		0, 0, 1, 1),
	ja(o_jahs,		1, 0, 0, 1),						// aka janc
	ja(o_jaeq,		1, 0, 1, 0),						// aka jaz
	ja(o_jane,		1, 0, 1, 1),						// aka janz
}

func ja(f opfunc, c0, c1, c2, c3 uint16) {
	return op(f,	1, 1, 0, 0, c0, c1, c2, c3, 1, 0, 0, 0, 0, 0, 0, 0)
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

func op(f opfunc, bits ...uint16) (o Opcode) {
	o.Expected = ui16(bits...)
	for i := 0; i < len(bits); i++ {
		bits[i] = 1		// TODO whether or not this should be legal should be in the spec
	}
	o.Mask = ui16(bits...)
	o.Func = f
	return o
}
