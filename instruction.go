package main

const (
	addrAbsolute  = iota
	addrAbsoluteX // +1 if page crossed
	addrAbsoluteY // +1 if page crossed
	addrAccumulator
	addrImmediate
	addrImplied
	addrIndirect
	addrIndirectX
	addrIndirectY // +1 if page crossed
	addrRelative  // +1 if branch succeeds, +2 if to a new page
	addrZeroPage
	addrZeroPageX
	addrZeroPageY
)

const (
	insADC = iota
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

var (
	instructions     = [256]instruction{}
	instructionSizes = [...]int{3, 3, 3, 1, 2, 1, 3, 2, 2, 2, 2, 2, 2}
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
)

func setIns(ins, opcode, cycles byte, addrMode int) {
	instructions[opcode] = instruction{ins, addrMode, cycles}
}

func init() {
	setIns(insADC, 0x69, 2, addrImmediate)
	setIns(insADC, 0x65, 3, addrZeroPage)
	setIns(insADC, 0x75, 4, addrZeroPageX)
	setIns(insADC, 0x6D, 4, addrAbsolute)
	setIns(insADC, 0x7D, 4, addrAbsoluteX)
	setIns(insADC, 0x79, 4, addrAbsoluteY)
	setIns(insADC, 0x61, 6, addrIndirectX)
	setIns(insADC, 0x71, 5, addrIndirectY)
	setIns(insAND, 0x29, 2, addrImmediate)
	setIns(insAND, 0x25, 3, addrZeroPage)
	setIns(insAND, 0x35, 4, addrZeroPageX)
	setIns(insAND, 0x2D, 4, addrAbsolute)
	setIns(insAND, 0x3D, 4, addrAbsoluteX)
	setIns(insAND, 0x39, 4, addrAbsoluteY)
	setIns(insAND, 0x21, 6, addrIndirectX)
	setIns(insAND, 0x31, 5, addrIndirectY)
	setIns(insASL, 0x0A, 2, addrAccumulator)
	setIns(insASL, 0x06, 5, addrZeroPage)
	setIns(insASL, 0x16, 6, addrZeroPageX)
	setIns(insASL, 0x0E, 6, addrAbsolute)
	setIns(insASL, 0x1E, 7, addrAbsoluteX)
	setIns(insBCC, 0x90, 2, addrRelative)
	setIns(insBCS, 0xB0, 2, addrRelative)
	setIns(insBEQ, 0xF0, 2, addrRelative)
	setIns(insBIT, 0x24, 3, addrZeroPage)
	setIns(insBIT, 0x2C, 4, addrAbsolute)
	setIns(insBMI, 0x30, 2, addrRelative)
	setIns(insBNE, 0xD0, 2, addrRelative)
	setIns(insBPL, 0x10, 2, addrRelative)
	setIns(insBRK, 0x00, 7, addrImplied)
	setIns(insBVC, 0x50, 2, addrRelative)
	setIns(insBVS, 0x70, 2, addrRelative)
	setIns(insCLC, 0x18, 2, addrImplied)
	setIns(insCLD, 0xD8, 2, addrImplied)
	setIns(insCLI, 0x58, 2, addrImplied)
	setIns(insCLV, 0xB8, 2, addrImplied)
	setIns(insCMP, 0xC9, 2, addrImmediate)
	setIns(insCMP, 0xC5, 3, addrZeroPage)
	setIns(insCMP, 0xD5, 4, addrZeroPageX)
	setIns(insCMP, 0xCD, 4, addrAbsolute)
	setIns(insCMP, 0xDD, 4, addrAbsoluteX)
	setIns(insCMP, 0xD9, 4, addrAbsoluteY)
	setIns(insCMP, 0xC1, 6, addrIndirectX)
	setIns(insCMP, 0xD1, 5, addrIndirectY)
	setIns(insCPX, 0xE0, 2, addrImmediate)
	setIns(insCPX, 0xE4, 3, addrZeroPage)
	setIns(insCPX, 0xEC, 4, addrAbsolute)
	setIns(insCPY, 0xC0, 2, addrImmediate)
	setIns(insCPY, 0xC4, 3, addrZeroPage)
	setIns(insCPY, 0xCC, 4, addrAbsolute)
	setIns(insDEC, 0xC6, 5, addrZeroPage)
	setIns(insDEC, 0xD6, 6, addrZeroPageX)
	setIns(insDEC, 0xCE, 6, addrAbsolute)
	setIns(insDEC, 0xDE, 7, addrAbsoluteX)
	setIns(insDEX, 0xCA, 2, addrImplied)
	setIns(insDEY, 0x88, 2, addrImplied)
	setIns(insEOR, 0x49, 2, addrImmediate)
	setIns(insEOR, 0x45, 3, addrZeroPage)
	setIns(insEOR, 0x55, 4, addrZeroPageX)
	setIns(insEOR, 0x4D, 4, addrAbsolute)
	setIns(insEOR, 0x5D, 4, addrAbsoluteX)
	setIns(insEOR, 0x59, 4, addrAbsoluteY)
	setIns(insEOR, 0x41, 6, addrIndirectX)
	setIns(insEOR, 0x51, 5, addrIndirectY)
	setIns(insINC, 0xE6, 5, addrZeroPage)
	setIns(insINC, 0xF6, 6, addrZeroPageX)
	setIns(insINC, 0xEE, 6, addrAbsolute)
	setIns(insINC, 0xFE, 7, addrAbsoluteX)
	setIns(insINX, 0xE8, 2, addrImplied)
	setIns(insINY, 0xC8, 2, addrImplied)
	setIns(insJMP, 0x4C, 3, addrAbsolute)
	setIns(insJMP, 0x6C, 5, addrIndirect)
	setIns(insJSR, 0x20, 6, addrAbsolute)
	setIns(insLDA, 0xA9, 2, addrImmediate)
	setIns(insLDA, 0xA5, 3, addrZeroPage)
	setIns(insLDA, 0xB5, 4, addrZeroPageX)
	setIns(insLDA, 0xAD, 4, addrAbsolute)
	setIns(insLDA, 0xBD, 4, addrAbsoluteX)
	setIns(insLDA, 0xB9, 4, addrAbsoluteY)
	setIns(insLDA, 0xA1, 6, addrIndirectX)
	setIns(insLDA, 0xB1, 5, addrIndirectY)
	setIns(insLDX, 0xA2, 2, addrImmediate)
	setIns(insLDX, 0xA6, 3, addrZeroPage)
	setIns(insLDX, 0xB6, 4, addrZeroPageY)
	setIns(insLDX, 0xAE, 4, addrAbsolute)
	setIns(insLDX, 0xBE, 4, addrAbsoluteY)
	setIns(insLDY, 0xA0, 2, addrImmediate)
	setIns(insLDY, 0xA4, 3, addrZeroPage)
	setIns(insLDY, 0xB4, 4, addrZeroPageX)
	setIns(insLDY, 0xAC, 4, addrAbsolute)
	setIns(insLDY, 0xBC, 4, addrAbsoluteX)
	setIns(insLSR, 0x4A, 2, addrAccumulator)
	setIns(insLSR, 0x46, 5, addrZeroPage)
	setIns(insLSR, 0x56, 6, addrZeroPageX)
	setIns(insLSR, 0x4E, 6, addrAbsolute)
	setIns(insLSR, 0x5E, 7, addrAbsoluteX)
	setIns(insNOP, 0xEA, 2, addrImplied)
	setIns(insORA, 0x09, 2, addrImmediate)
	setIns(insORA, 0x05, 3, addrZeroPage)
	setIns(insORA, 0x15, 4, addrZeroPageX)
	setIns(insORA, 0x0D, 4, addrAbsolute)
	setIns(insORA, 0x1D, 4, addrAbsoluteX)
	setIns(insORA, 0x19, 4, addrAbsoluteY)
	setIns(insORA, 0x01, 6, addrIndirectX)
	setIns(insORA, 0x11, 5, addrIndirectY)
	setIns(insPHA, 0x48, 3, addrImplied)
	setIns(insPHP, 0x08, 3, addrImplied)
	setIns(insPLA, 0x68, 4, addrImplied)
	setIns(insPLP, 0x28, 4, addrImplied)
	setIns(insROL, 0x2A, 2, addrAccumulator)
	setIns(insROL, 0x26, 5, addrZeroPage)
	setIns(insROL, 0x36, 6, addrZeroPageX)
	setIns(insROL, 0x2E, 6, addrAbsolute)
	setIns(insROL, 0x3E, 7, addrAbsoluteX)
	setIns(insROR, 0x6A, 2, addrAccumulator)
	setIns(insROR, 0x66, 5, addrZeroPage)
	setIns(insROR, 0x76, 6, addrZeroPageX)
	setIns(insROR, 0x6E, 6, addrAbsolute)
	setIns(insROR, 0x7E, 7, addrAbsoluteX)
	setIns(insRTI, 0x40, 6, addrImplied)
	setIns(insRTS, 0x60, 6, addrImplied)
	setIns(insSBC, 0xE9, 2, addrImmediate)
	setIns(insSBC, 0xE5, 3, addrZeroPage)
	setIns(insSBC, 0xF5, 4, addrZeroPageX)
	setIns(insSBC, 0xED, 4, addrAbsolute)
	setIns(insSBC, 0xFD, 4, addrAbsoluteX)
	setIns(insSBC, 0xF9, 4, addrAbsoluteY)
	setIns(insSBC, 0xE1, 6, addrIndirectX)
	setIns(insSBC, 0xF1, 5, addrIndirectY)
	setIns(insSEC, 0x38, 2, addrImplied)
	setIns(insSED, 0xF8, 2, addrImplied)
	setIns(insSEI, 0x78, 2, addrImplied)
	setIns(insSTA, 0x85, 3, addrZeroPage)
	setIns(insSTA, 0x95, 4, addrZeroPageX)
	setIns(insSTA, 0x8D, 4, addrAbsolute)
	setIns(insSTA, 0x9D, 5, addrAbsoluteX)
	setIns(insSTA, 0x99, 5, addrAbsoluteY)
	setIns(insSTA, 0x81, 6, addrIndirectX)
	setIns(insSTA, 0x91, 6, addrIndirectY)
	setIns(insSTX, 0x86, 3, addrZeroPage)
	setIns(insSTX, 0x96, 4, addrZeroPageY)
	setIns(insSTX, 0x8E, 4, addrAbsolute)
	setIns(insSTY, 0x84, 3, addrZeroPage)
	setIns(insSTY, 0x94, 4, addrZeroPageX)
	setIns(insSTY, 0x8C, 4, addrAbsolute)
	setIns(insTAX, 0xAA, 2, addrImplied)
	setIns(insTAY, 0xA8, 2, addrImplied)
	setIns(insTSX, 0xBA, 2, addrImplied)
	setIns(insTXA, 0x8A, 2, addrImplied)
	setIns(insTXS, 0x9A, 2, addrImplied)
	setIns(insTYA, 0x98, 2, addrImplied)

	// unofficial instructions
	setIns(insAAC, 0x0B, 2, addrImmediate)
	setIns(insAAC, 0x2B, 2, addrImmediate)
	setIns(insAAX, 0x87, 3, addrZeroPage)
	setIns(insAAX, 0x97, 4, addrZeroPageX)
	setIns(insAAX, 0x83, 6, addrIndirectX)
	setIns(insAAX, 0x8F, 4, addrAbsolute)
	setIns(insARR, 0x6B, 2, addrImmediate)
	setIns(insASR, 0x4B, 2, addrImmediate)
	setIns(insATX, 0xAB, 2, addrImmediate)
	setIns(insAXA, 0x9F, 5, addrAbsoluteY)
	setIns(insAXA, 0x93, 6, addrIndirectY)
	setIns(insAXS, 0xCB, 2, addrImmediate)
	setIns(insDCP, 0xC7, 5, addrZeroPage)
	setIns(insDCP, 0xD7, 6, addrZeroPageX)
	setIns(insDCP, 0xCF, 6, addrAbsolute)
	setIns(insDCP, 0xDF, 7, addrAbsoluteX)
	setIns(insDCP, 0xDB, 7, addrAbsoluteY)
	setIns(insDCP, 0xC3, 8, addrIndirectX)
	setIns(insDCP, 0xD3, 8, addrIndirectY)
	setIns(insDOP, 0x04, 3, addrZeroPage)
	setIns(insDOP, 0x14, 4, addrZeroPageX)
	setIns(insDOP, 0x34, 4, addrZeroPageX)
	setIns(insDOP, 0x44, 3, addrZeroPage)
	setIns(insDOP, 0x54, 4, addrZeroPageX)
	setIns(insDOP, 0x64, 3, addrZeroPage)
	setIns(insDOP, 0x74, 4, addrZeroPageX)
	setIns(insDOP, 0x80, 2, addrImmediate)
	setIns(insDOP, 0x82, 2, addrImmediate)
	setIns(insDOP, 0x89, 2, addrImmediate)
	setIns(insDOP, 0xC2, 2, addrImmediate)
	setIns(insDOP, 0xD4, 4, addrZeroPageX)
	setIns(insDOP, 0xE2, 2, addrImmediate)
	setIns(insDOP, 0xF4, 4, addrZeroPageX)
	setIns(insISC, 0xE7, 5, addrZeroPage)
	setIns(insISC, 0xF7, 6, addrZeroPageX)
	setIns(insISC, 0xEF, 6, addrAbsolute)
	setIns(insISC, 0xFF, 7, addrAbsoluteX)
	setIns(insISC, 0xFB, 7, addrAbsoluteY)
	setIns(insISC, 0xE3, 8, addrIndirectX)
	setIns(insISC, 0xF3, 8, addrIndirectY)
	setIns(insKIL, 0x02, 0, addrImplied)
	setIns(insKIL, 0x12, 0, addrImplied)
	setIns(insKIL, 0x22, 0, addrImplied)
	setIns(insKIL, 0x32, 0, addrImplied)
	setIns(insKIL, 0x42, 0, addrImplied)
	setIns(insKIL, 0x52, 0, addrImplied)
	setIns(insKIL, 0x62, 0, addrImplied)
	setIns(insKIL, 0x72, 0, addrImplied)
	setIns(insKIL, 0x92, 0, addrImplied)
	setIns(insKIL, 0xB2, 0, addrImplied)
	setIns(insKIL, 0xD2, 0, addrImplied)
	setIns(insKIL, 0xF2, 0, addrImplied)
	setIns(insLAR, 0xBB, 4, addrAbsoluteY)
	setIns(insLAX, 0xA7, 3, addrZeroPage)
	setIns(insLAX, 0xB7, 4, addrZeroPageY)
	setIns(insLAX, 0xAF, 4, addrAbsolute)
	setIns(insLAX, 0xBF, 4, addrAbsoluteY)
	setIns(insLAX, 0xA3, 6, addrIndirectX)
	setIns(insLAX, 0xB3, 5, addrIndirectY)
	setIns(insNOP, 0x1A, 2, addrImplied)
	setIns(insNOP, 0x3A, 2, addrImplied)
	setIns(insNOP, 0x5A, 2, addrImplied)
	setIns(insNOP, 0x7A, 2, addrImplied)
	setIns(insNOP, 0xDA, 2, addrImplied)
	setIns(insNOP, 0xFA, 2, addrImplied)
	setIns(insRLA, 0x27, 5, addrZeroPage)
	setIns(insRLA, 0x37, 6, addrZeroPageX)
	setIns(insRLA, 0x2F, 6, addrAbsolute)
	setIns(insRLA, 0x3F, 7, addrAbsoluteX)
	setIns(insRLA, 0x3B, 7, addrAbsoluteY)
	setIns(insRLA, 0x23, 8, addrIndirectX)
	setIns(insRLA, 0x33, 8, addrIndirectY)
	setIns(insRRA, 0x67, 5, addrZeroPage)
	setIns(insRRA, 0x77, 6, addrZeroPageX)
	setIns(insRRA, 0x6F, 6, addrAbsolute)
	setIns(insRRA, 0x7F, 7, addrAbsoluteX)
	setIns(insRRA, 0x7B, 7, addrAbsoluteY)
	setIns(insRRA, 0x63, 8, addrIndirectX)
	setIns(insRRA, 0x73, 8, addrIndirectY)
	setIns(insSBC, 0xEB, 2, addrImmediate)
	setIns(insSLO, 0x07, 5, addrZeroPage)
	setIns(insSLO, 0x17, 6, addrZeroPageX)
	setIns(insSLO, 0x0F, 6, addrAbsolute)
	setIns(insSLO, 0x1F, 7, addrAbsoluteX)
	setIns(insSLO, 0x1B, 7, addrAbsoluteY)
	setIns(insSLO, 0x03, 8, addrIndirectX)
	setIns(insSLO, 0x13, 8, addrIndirectY)
	setIns(insSRE, 0x47, 5, addrZeroPage)
	setIns(insSRE, 0x57, 6, addrZeroPageX)
	setIns(insSRE, 0x4F, 6, addrAbsolute)
	setIns(insSRE, 0x5F, 7, addrAbsoluteX)
	setIns(insSRE, 0x5B, 7, addrAbsoluteY)
	setIns(insSRE, 0x43, 8, addrIndirectX)
	setIns(insSRE, 0x53, 8, addrIndirectY)
	setIns(insSXA, 0x9E, 5, addrAbsoluteY)
	setIns(insSYA, 0x9C, 5, addrAbsoluteX)
	setIns(insTOP, 0x0C, 4, addrAbsolute)
	setIns(insTOP, 0x1C, 4, addrAbsoluteX)
	setIns(insTOP, 0x3C, 4, addrAbsoluteX)
	setIns(insTOP, 0x5C, 4, addrAbsoluteX)
	setIns(insTOP, 0x7C, 4, addrAbsoluteX)
	setIns(insTOP, 0xDC, 4, addrAbsoluteX)
	setIns(insTOP, 0xFC, 4, addrAbsoluteX)
	setIns(insXAA, 0x8B, 2, addrImmediate)
	setIns(insXAS, 0x9B, 5, addrAbsoluteY)
}
