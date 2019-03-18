package main

const (
	addrAbsolute = uint8(iota)
	addrAbsoluteX
	addrAbsoluteY
	addrAccumulator
	addrImmediate
	addrImplied
	addrIndexedIndirect
	addrIndirect
	addrIndirectIndexed
	addrRelative
	addrZeroPage
	addrZeroPageX
	addrZeroPageY
)
const (
	insADC = uint8(iota)
	insAND
	insASL
	insBCC
	insBCS
	insBEQ
	insBIT
	insBMI
	insBNE
	insBPL
	insBRK
	insBVC
	insBVS
	insCLC
	insCLD
	insCLI
	insCLV
	insCMP
	insCPX
	insCPY
	insDEC
	insDEX
	insDEY
	insEOR
	insINC
	insINX
	insINY
	insJMP
	insJSR
	insLDA
	insLDX
	insLDY
	insLSR
	insNOP
	insORA
	insPHA
	insPHP
	insPLA
	insPLP
	insROL
	insROR
	insRTI
	insRTS
	insSBC
	insSEC
	insSED
	insSEI
	insSTA
	insSTX
	insSTY
	insTAX
	insTAY
	insTSX
	insTXA
	insTXS
	insTYA

	// unofficial instructions
	// see http://nesdev.com/undocumented_opcodes.txt
	insAAC
	insAAX
	insARR
	insASR
	insATX
	insAXA
	insAXS
	insDCP
	insDOP
	insISC
	insKIL
	insLAR
	insLAX
	// insNOP
	insRLA
	insRRA
	//insSBC
	insSLO
	insSRE
	insSXA
	insSYA
	insTOP
	insXAA
	insXAS
)

type instruction struct {
	id       uint8
	cycles   uint8
	exCyc    uint8
	addrMode uint8
}

var (
	instructionSizes = [...]uint8{3, 3, 3, 1, 2, 1, 2, 3, 2, 2, 2, 2, 2}
	instructionNames = [...]string{
		"ADC", "AND", "ASL", "BCC", "BCS", "BEQ", "BIT", "BMI", "BNE", "BPL",
		"BRK", "BVC", "BVS", "CLC", "CLD", "CLI", "CLV", "CMP", "CPX", "CPY",
		"DEC", "DEX", "DEY", "EOR", "INC", "INX", "INY", "JMP", "JSR", "LDA",
		"LDX", "LDY", "LSR", "NOP", "ORA", "PHA", "PHP", "PLA", "PLP", "ROL",
		"ROR", "RTI", "RTS", "SBC", "SEC", "SED", "SEI", "STA", "STX", "STY",
		"TAX", "TAY", "TSX", "TXA", "TXS", "TYA",
		// unofficial
		"AAC", "AAX", "ARR", "ASR", "ATX", "AXA", "AXS", "DCP", "DOP", "ISC",
		"KIL", "LAR", "LAX", "RLA", "RRA", "SLO", "SRE", "SXA", "SYA", "TOP",
		"XAA", "XAS",
	}
	instructions = [256]instruction{
		0x69: instruction{insADC, 2, 0, addrImmediate},
		0x65: instruction{insADC, 3, 0, addrZeroPage},
		0x75: instruction{insADC, 4, 0, addrZeroPageX},
		0x6D: instruction{insADC, 4, 0, addrAbsolute},
		0x7D: instruction{insADC, 4, 1, addrAbsoluteX},
		0x79: instruction{insADC, 4, 1, addrAbsoluteY},
		0x61: instruction{insADC, 6, 0, addrIndexedIndirect},
		0x71: instruction{insADC, 5, 1, addrIndirectIndexed},
		0x29: instruction{insAND, 2, 0, addrImmediate},
		0x25: instruction{insAND, 3, 0, addrZeroPage},
		0x35: instruction{insAND, 4, 0, addrZeroPageX},
		0x2D: instruction{insAND, 4, 0, addrAbsolute},
		0x3D: instruction{insAND, 4, 1, addrAbsoluteX},
		0x39: instruction{insAND, 4, 1, addrAbsoluteY},
		0x21: instruction{insAND, 6, 0, addrIndexedIndirect},
		0x31: instruction{insAND, 5, 1, addrIndirectIndexed},
		0x0A: instruction{insASL, 2, 0, addrAccumulator},
		0x06: instruction{insASL, 5, 0, addrZeroPage},
		0x16: instruction{insASL, 6, 0, addrZeroPageX},
		0x0E: instruction{insASL, 6, 0, addrAbsolute},
		0x1E: instruction{insASL, 7, 0, addrAbsoluteX},
		0x90: instruction{insBCC, 2, 1, addrRelative},
		0xB0: instruction{insBCS, 2, 1, addrRelative},
		0xF0: instruction{insBEQ, 2, 1, addrRelative},
		0x24: instruction{insBIT, 3, 0, addrZeroPage},
		0x2C: instruction{insBIT, 4, 0, addrAbsolute},
		0x30: instruction{insBMI, 2, 1, addrRelative},
		0xD0: instruction{insBNE, 2, 1, addrRelative},
		0x10: instruction{insBPL, 2, 1, addrRelative},
		0x00: instruction{insBRK, 7, 0, addrImplied},
		0x50: instruction{insBVC, 2, 1, addrRelative},
		0x70: instruction{insBVS, 2, 1, addrRelative},
		0x18: instruction{insCLC, 2, 0, addrImplied},
		0xD8: instruction{insCLD, 2, 0, addrImplied},
		0x58: instruction{insCLI, 2, 0, addrImplied},
		0xB8: instruction{insCLV, 2, 0, addrImplied},
		0xC9: instruction{insCMP, 2, 0, addrImmediate},
		0xC5: instruction{insCMP, 3, 0, addrZeroPage},
		0xD5: instruction{insCMP, 4, 0, addrZeroPageX},
		0xCD: instruction{insCMP, 4, 0, addrAbsolute},
		0xDD: instruction{insCMP, 4, 1, addrAbsoluteX},
		0xD9: instruction{insCMP, 4, 1, addrAbsoluteY},
		0xC1: instruction{insCMP, 6, 0, addrIndexedIndirect},
		0xD1: instruction{insCMP, 5, 1, addrIndirectIndexed},
		0xE0: instruction{insCPX, 2, 0, addrImmediate},
		0xE4: instruction{insCPX, 3, 0, addrZeroPage},
		0xEC: instruction{insCPX, 4, 0, addrAbsolute},
		0xC0: instruction{insCPY, 2, 0, addrImmediate},
		0xC4: instruction{insCPY, 3, 0, addrZeroPage},
		0xCC: instruction{insCPY, 4, 0, addrAbsolute},
		0xC6: instruction{insDEC, 5, 0, addrZeroPage},
		0xD6: instruction{insDEC, 6, 0, addrZeroPageX},
		0xCE: instruction{insDEC, 6, 0, addrAbsolute},
		0xDE: instruction{insDEC, 7, 0, addrAbsoluteX},
		0xCA: instruction{insDEX, 2, 0, addrImplied},
		0x88: instruction{insDEY, 2, 0, addrImplied},
		0x49: instruction{insEOR, 2, 0, addrImmediate},
		0x45: instruction{insEOR, 3, 0, addrZeroPage},
		0x55: instruction{insEOR, 4, 0, addrZeroPageX},
		0x4D: instruction{insEOR, 4, 0, addrAbsolute},
		0x5D: instruction{insEOR, 4, 1, addrAbsoluteX},
		0x59: instruction{insEOR, 4, 1, addrAbsoluteY},
		0x41: instruction{insEOR, 6, 0, addrIndexedIndirect},
		0x51: instruction{insEOR, 5, 1, addrIndirectIndexed},
		0xE6: instruction{insINC, 5, 0, addrZeroPage},
		0xF6: instruction{insINC, 6, 0, addrZeroPageX},
		0xEE: instruction{insINC, 6, 0, addrAbsolute},
		0xFE: instruction{insINC, 7, 0, addrAbsoluteX},
		0xE8: instruction{insINX, 2, 0, addrImplied},
		0xC8: instruction{insINY, 2, 0, addrImplied},
		0x4C: instruction{insJMP, 3, 0, addrAbsolute},
		0x6C: instruction{insJMP, 5, 0, addrIndirect},
		0x20: instruction{insJSR, 6, 0, addrAbsolute},
		0xA9: instruction{insLDA, 2, 0, addrImmediate},
		0xA5: instruction{insLDA, 3, 0, addrZeroPage},
		0xB5: instruction{insLDA, 4, 0, addrZeroPageX},
		0xAD: instruction{insLDA, 4, 0, addrAbsolute},
		0xBD: instruction{insLDA, 4, 1, addrAbsoluteX},
		0xB9: instruction{insLDA, 4, 1, addrAbsoluteY},
		0xA1: instruction{insLDA, 6, 0, addrIndexedIndirect},
		0xB1: instruction{insLDA, 5, 1, addrIndirectIndexed},
		0xA2: instruction{insLDX, 2, 0, addrImmediate},
		0xA6: instruction{insLDX, 3, 0, addrZeroPage},
		0xB6: instruction{insLDX, 4, 0, addrZeroPageY},
		0xAE: instruction{insLDX, 4, 0, addrAbsolute},
		0xBE: instruction{insLDX, 4, 1, addrAbsoluteY},
		0xA0: instruction{insLDY, 2, 0, addrImmediate},
		0xA4: instruction{insLDY, 3, 0, addrZeroPage},
		0xB4: instruction{insLDY, 4, 0, addrZeroPageX},
		0xAC: instruction{insLDY, 4, 0, addrAbsolute},
		0xBC: instruction{insLDY, 4, 1, addrAbsoluteX},
		0x4A: instruction{insLSR, 2, 0, addrAccumulator},
		0x46: instruction{insLSR, 5, 0, addrZeroPage},
		0x56: instruction{insLSR, 6, 0, addrZeroPageX},
		0x4E: instruction{insLSR, 6, 0, addrAbsolute},
		0x5E: instruction{insLSR, 7, 0, addrAbsoluteX},
		0xEA: instruction{insNOP, 2, 0, addrImplied},
		0x09: instruction{insORA, 2, 0, addrImmediate},
		0x05: instruction{insORA, 3, 0, addrZeroPage},
		0x15: instruction{insORA, 4, 0, addrZeroPageX},
		0x0D: instruction{insORA, 4, 0, addrAbsolute},
		0x1D: instruction{insORA, 4, 1, addrAbsoluteX},
		0x19: instruction{insORA, 4, 1, addrAbsoluteY},
		0x01: instruction{insORA, 6, 0, addrIndexedIndirect},
		0x11: instruction{insORA, 5, 1, addrIndirectIndexed},
		0x48: instruction{insPHA, 3, 0, addrImplied},
		0x08: instruction{insPHP, 3, 0, addrImplied},
		0x68: instruction{insPLA, 4, 0, addrImplied},
		0x28: instruction{insPLP, 4, 0, addrImplied},
		0x2A: instruction{insROL, 2, 0, addrAccumulator},
		0x26: instruction{insROL, 5, 0, addrZeroPage},
		0x36: instruction{insROL, 6, 0, addrZeroPageX},
		0x2E: instruction{insROL, 6, 0, addrAbsolute},
		0x3E: instruction{insROL, 7, 0, addrAbsoluteX},
		0x6A: instruction{insROR, 2, 0, addrAccumulator},
		0x66: instruction{insROR, 5, 0, addrZeroPage},
		0x76: instruction{insROR, 6, 0, addrZeroPageX},
		0x6E: instruction{insROR, 6, 0, addrAbsolute},
		0x7E: instruction{insROR, 7, 0, addrAbsoluteX},
		0x40: instruction{insRTI, 6, 0, addrImplied},
		0x60: instruction{insRTS, 6, 0, addrImplied},
		0xE9: instruction{insSBC, 2, 0, addrImmediate},
		0xE5: instruction{insSBC, 3, 0, addrZeroPage},
		0xF5: instruction{insSBC, 4, 0, addrZeroPageX},
		0xED: instruction{insSBC, 4, 0, addrAbsolute},
		0xFD: instruction{insSBC, 4, 1, addrAbsoluteX},
		0xF9: instruction{insSBC, 4, 1, addrAbsoluteY},
		0xE1: instruction{insSBC, 6, 0, addrIndexedIndirect},
		0xF1: instruction{insSBC, 5, 1, addrIndirectIndexed},
		0x38: instruction{insSEC, 2, 0, addrImplied},
		0xF8: instruction{insSED, 2, 0, addrImplied},
		0x78: instruction{insSEI, 2, 0, addrImplied},
		0x85: instruction{insSTA, 3, 0, addrZeroPage},
		0x95: instruction{insSTA, 4, 0, addrZeroPageX},
		0x8D: instruction{insSTA, 4, 0, addrAbsolute},
		0x9D: instruction{insSTA, 5, 0, addrAbsoluteX},
		0x99: instruction{insSTA, 5, 0, addrAbsoluteY},
		0x81: instruction{insSTA, 6, 0, addrIndexedIndirect},
		0x91: instruction{insSTA, 6, 0, addrIndirectIndexed},
		0x86: instruction{insSTX, 3, 0, addrZeroPage},
		0x96: instruction{insSTX, 4, 0, addrZeroPageY},
		0x8E: instruction{insSTX, 4, 0, addrAbsolute},
		0x84: instruction{insSTY, 3, 0, addrZeroPage},
		0x94: instruction{insSTY, 4, 0, addrZeroPageX},
		0x8C: instruction{insSTY, 4, 0, addrAbsolute},
		0xAA: instruction{insTAX, 2, 0, addrImplied},
		0xA8: instruction{insTAY, 2, 0, addrImplied},
		0xBA: instruction{insTSX, 2, 0, addrImplied},
		0x8A: instruction{insTXA, 2, 0, addrImplied},
		0x9A: instruction{insTXS, 2, 0, addrImplied},
		0x98: instruction{insTYA, 2, 0, addrImplied},

		// unofficial instructions
		0x0B: instruction{insAAC, 2, 0, addrImmediate},
		0x2B: instruction{insAAC, 2, 0, addrImmediate},
		0x87: instruction{insAAX, 3, 0, addrZeroPage},
		0x97: instruction{insAAX, 4, 0, addrZeroPageY},
		0x83: instruction{insAAX, 6, 0, addrIndexedIndirect},
		0x8F: instruction{insAAX, 4, 0, addrAbsolute},
		0x6B: instruction{insARR, 2, 0, addrImmediate},
		0x4B: instruction{insASR, 2, 0, addrImmediate},
		0xAB: instruction{insATX, 2, 0, addrImmediate},
		0x9F: instruction{insAXA, 5, 0, addrAbsoluteY},
		0x93: instruction{insAXA, 6, 0, addrIndirectIndexed},
		0xCB: instruction{insAXS, 2, 0, addrImmediate},
		0xC7: instruction{insDCP, 5, 0, addrZeroPage},
		0xD7: instruction{insDCP, 6, 0, addrZeroPageX},
		0xCF: instruction{insDCP, 6, 0, addrAbsolute},
		0xDF: instruction{insDCP, 7, 0, addrAbsoluteX},
		0xDB: instruction{insDCP, 7, 0, addrAbsoluteY},
		0xC3: instruction{insDCP, 8, 0, addrIndexedIndirect},
		0xD3: instruction{insDCP, 8, 0, addrIndirectIndexed},
		0x04: instruction{insDOP, 3, 0, addrZeroPage},
		0x14: instruction{insDOP, 4, 0, addrZeroPageX},
		0x34: instruction{insDOP, 4, 0, addrZeroPageX},
		0x44: instruction{insDOP, 3, 0, addrZeroPage},
		0x54: instruction{insDOP, 4, 0, addrZeroPageX},
		0x64: instruction{insDOP, 3, 0, addrZeroPage},
		0x74: instruction{insDOP, 4, 0, addrZeroPageX},
		0x80: instruction{insDOP, 2, 0, addrImmediate},
		0x82: instruction{insDOP, 2, 0, addrImmediate},
		0x89: instruction{insDOP, 2, 0, addrImmediate},
		0xC2: instruction{insDOP, 2, 0, addrImmediate},
		0xD4: instruction{insDOP, 4, 0, addrZeroPageX},
		0xE2: instruction{insDOP, 2, 0, addrImmediate},
		0xF4: instruction{insDOP, 4, 0, addrZeroPageX},
		0xE7: instruction{insISC, 5, 0, addrZeroPage},
		0xF7: instruction{insISC, 6, 0, addrZeroPageX},
		0xEF: instruction{insISC, 6, 0, addrAbsolute},
		0xFF: instruction{insISC, 7, 0, addrAbsoluteX},
		0xFB: instruction{insISC, 7, 0, addrAbsoluteY},
		0xE3: instruction{insISC, 8, 0, addrIndexedIndirect},
		0xF3: instruction{insISC, 8, 0, addrIndirectIndexed},
		0x02: instruction{insKIL, 0, 0, addrImplied},
		0x12: instruction{insKIL, 0, 0, addrImplied},
		0x22: instruction{insKIL, 0, 0, addrImplied},
		0x32: instruction{insKIL, 0, 0, addrImplied},
		0x42: instruction{insKIL, 0, 0, addrImplied},
		0x52: instruction{insKIL, 0, 0, addrImplied},
		0x62: instruction{insKIL, 0, 0, addrImplied},
		0x72: instruction{insKIL, 0, 0, addrImplied},
		0x92: instruction{insKIL, 0, 0, addrImplied},
		0xB2: instruction{insKIL, 0, 0, addrImplied},
		0xD2: instruction{insKIL, 0, 0, addrImplied},
		0xF2: instruction{insKIL, 0, 0, addrImplied},
		0xBB: instruction{insLAR, 4, 1, addrAbsoluteY},
		0xA7: instruction{insLAX, 3, 0, addrZeroPage},
		0xB7: instruction{insLAX, 4, 0, addrZeroPageY},
		0xAF: instruction{insLAX, 4, 0, addrAbsolute},
		0xBF: instruction{insLAX, 4, 1, addrAbsoluteY},
		0xA3: instruction{insLAX, 6, 0, addrIndexedIndirect},
		0xB3: instruction{insLAX, 5, 1, addrIndirectIndexed},
		0x1A: instruction{insNOP, 2, 0, addrImplied},
		0x3A: instruction{insNOP, 2, 0, addrImplied},
		0x5A: instruction{insNOP, 2, 0, addrImplied},
		0x7A: instruction{insNOP, 2, 0, addrImplied},
		0xDA: instruction{insNOP, 2, 0, addrImplied},
		0xFA: instruction{insNOP, 2, 0, addrImplied},
		0x27: instruction{insRLA, 5, 0, addrZeroPage},
		0x37: instruction{insRLA, 6, 0, addrZeroPageX},
		0x2F: instruction{insRLA, 6, 0, addrAbsolute},
		0x3F: instruction{insRLA, 7, 0, addrAbsoluteX},
		0x3B: instruction{insRLA, 7, 0, addrAbsoluteY},
		0x23: instruction{insRLA, 8, 0, addrIndexedIndirect},
		0x33: instruction{insRLA, 8, 0, addrIndirectIndexed},
		0x67: instruction{insRRA, 5, 0, addrZeroPage},
		0x77: instruction{insRRA, 6, 0, addrZeroPageX},
		0x6F: instruction{insRRA, 6, 0, addrAbsolute},
		0x7F: instruction{insRRA, 7, 0, addrAbsoluteX},
		0x7B: instruction{insRRA, 7, 0, addrAbsoluteY},
		0x63: instruction{insRRA, 8, 0, addrIndexedIndirect},
		0x73: instruction{insRRA, 8, 0, addrIndirectIndexed},
		0xEB: instruction{insSBC, 2, 0, addrImmediate},
		0x07: instruction{insSLO, 5, 0, addrZeroPage},
		0x17: instruction{insSLO, 6, 0, addrZeroPageX},
		0x0F: instruction{insSLO, 6, 0, addrAbsolute},
		0x1F: instruction{insSLO, 7, 0, addrAbsoluteX},
		0x1B: instruction{insSLO, 7, 0, addrAbsoluteY},
		0x03: instruction{insSLO, 8, 0, addrIndexedIndirect},
		0x13: instruction{insSLO, 8, 0, addrIndirectIndexed},
		0x47: instruction{insSRE, 5, 0, addrZeroPage},
		0x57: instruction{insSRE, 6, 0, addrZeroPageX},
		0x4F: instruction{insSRE, 6, 0, addrAbsolute},
		0x5F: instruction{insSRE, 7, 0, addrAbsoluteX},
		0x5B: instruction{insSRE, 7, 0, addrAbsoluteY},
		0x43: instruction{insSRE, 8, 0, addrIndexedIndirect},
		0x53: instruction{insSRE, 8, 0, addrIndirectIndexed},
		0x9E: instruction{insSXA, 5, 0, addrAbsoluteY},
		0x9C: instruction{insSYA, 5, 0, addrAbsoluteX},
		0x0C: instruction{insTOP, 4, 0, addrAbsolute},
		0x1C: instruction{insTOP, 4, 1, addrAbsoluteX},
		0x3C: instruction{insTOP, 4, 1, addrAbsoluteX},
		0x5C: instruction{insTOP, 4, 1, addrAbsoluteX},
		0x7C: instruction{insTOP, 4, 1, addrAbsoluteX},
		0xDC: instruction{insTOP, 4, 1, addrAbsoluteX},
		0xFC: instruction{insTOP, 4, 1, addrAbsoluteX},
		0x8B: instruction{insXAA, 2, 0, addrImmediate},
		0x9B: instruction{insXAS, 5, 0, addrAbsoluteY},
	}
)
