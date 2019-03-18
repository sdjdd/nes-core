package main

import (
	"fmt"
	"os"
)

// CPU - MOS6502
// http://wiki.nesdev.com/w/index.php/CPU
type CPU struct {
	cpuRegister
	cpuFlag
	console *Console
	cycles  uint64

	// sample CPU RAM allocation
	// see http://wiki.nesdev.com/w/index.php/Sample_RAM_map
	ram [2048]byte
}

// http://wiki.nesdev.com/w/index.php/CPU_registers
type cpuRegister struct {
	A, X, Y, S byte
	PC         uint16
}

// http://wiki.nesdev.com/w/index.php/Status_flags
type cpuFlag struct {
	C, Z, I, D, N, V byte
}

type stepInfo struct {
	opcode byte
	ins    instruction
	addr   uint16
	PC     uint16
	opnums []byte
}

func (cpu *CPU) step() stepInfo {
	opcode := cpu.read(cpu.PC)
	ins := instructions[opcode]
	size := instructionSizes[ins.addrMode]

	var opnums []byte
	var addr uint16
	if size > 1 {
		opnums = make([]byte, size-1)
		for i := 0; i < int(size-1); i++ {
			opnums[i] = cpu.read(cpu.PC + 1 + uint16(i))
			addr |= uint16(opnums[i]) << uint16(i*8)
		}
	}
	info := stepInfo{opcode, ins, addr, cpu.PC, opnums}

	cpu.PC += uint16(size)
	cpu.cycles += uint64(ins.cycles)

	switch ins.addrMode {
	case addrAbsolute, addrZeroPage, addrImplied, addrAccumulator:
		// do nothing
	case addrAbsoluteX, addrAbsoluteY:
		var offset uint16
		if ins.addrMode == addrAbsoluteX {
			offset = uint16(cpu.X)
		} else {
			offset = uint16(cpu.Y)
		}
		if addr&0xFF00 != (addr+offset)&0xFF00 {
			cpu.cycles += uint64(ins.exCyc)
		}
		addr += offset
	case addrIndexedIndirect:
		addr = cpu.bugRead((addr + uint16(cpu.X)) & 0x00FF)
	case addrIndirect:
		addr = cpu.bugRead(addr)
	case addrIndirectIndexed:
		addr = cpu.bugRead(addr)
		offset := uint16(cpu.Y)
		if addr&0xFF00 != (addr+offset)&0xFF00 {
			cpu.cycles += uint64(ins.exCyc)
		}
		addr += offset
	case addrImmediate:
		addr = cpu.PC - 1
	case addrRelative:
		offset := int8(addr)
		addr = uint16(int32(cpu.PC) + int32(offset))
		pc := cpu.PC
		defer func() {
			if cpu.PC != pc {
				cpu.cycles++ // cycles +1 if branch succeeds
				if pc&0xFF00 != addr&0xFF00 {
					cpu.cycles++ // +2 if to a new page
				}
			}
		}()
	case addrZeroPageX:
		addr = (addr + uint16(cpu.X)) & 0x00FF
	case addrZeroPageY:
		addr = (addr + uint16(cpu.Y)) & 0x00FF
		fmt.Printf("=== %04X ===\n", addr)
	default:
		fmt.Printf("\nunknown address mode: %d, %02X\n", ins.addrMode, opcode)
		os.Exit(0)
	}

	switch ins.id {
	case insNOP:
		// do nothing
	case insADC:
		cpu.adc(addr)
	case insAND:
		cpu.and(addr)
	case insASL:
		cpu.asl(addr, ins.addrMode)
	case insBCC:
		cpu.bcc(addr)
	case insBCS:
		cpu.bcs(addr)
	case insBEQ:
		cpu.beq(addr)
	case insBIT:
		cpu.bit(addr)
	case insBMI:
		cpu.bmi(addr)
	case insBNE:
		cpu.bne(addr)
	case insBPL:
		cpu.bpl(addr)
	case insBVC:
		cpu.bvc(addr)
	case insBVS:
		cpu.bvs(addr)
	case insCLC:
		cpu.C = 0
	case insCLD:
		cpu.D = 0
	case insCLV:
		cpu.V = 0
	case insCMP:
		cpu.cmp(addr)
	case insCPX:
		cpu.cpx(addr)
	case insCPY:
		cpu.cpy(addr)
	case insDEC:
		cpu.dec(addr)
	case insDEX:
		cpu.setValueNZ(&cpu.X, cpu.X-1)
	case insDEY:
		cpu.setValueNZ(&cpu.Y, cpu.Y-1)
	case insEOR:
		cpu.eor(addr)
	case insINC:
		cpu.inc(addr)
	case insINX:
		cpu.setValueNZ(&cpu.X, cpu.X+1)
	case insINY:
		cpu.setValueNZ(&cpu.Y, cpu.Y+1)
	case insJMP:
		cpu.jmp(addr)
	case insJSR:
		cpu.jsr(addr)
	case insLDA:
		cpu.setValueNZ(&cpu.A, cpu.read(addr))
	case insLDX:
		cpu.setValueNZ(&cpu.X, cpu.read(addr))
	case insLDY:
		cpu.setValueNZ(&cpu.Y, cpu.read(addr))
	case insLSR:
		cpu.lsr(addr, ins.addrMode)
	case insORA:
		cpu.ora(addr)
	case insPHA:
		cpu.push(cpu.A)
	case insPHP:
		cpu.push(cpu.flag() | 0x30) // set B flag
	case insPLA:
		cpu.setValueNZ(&cpu.A, cpu.pull())
	case insPLP:
		cpu.setFlags(cpu.pull())
	case insROL:
		cpu.rol(addr, ins.addrMode)
	case insROR:
		cpu.ror(addr, ins.addrMode)
	case insRTI:
		cpu.rti()
	case insRTS:
		cpu.rts()
	case insSBC:
		cpu.sbc(addr)
	case insSEC:
		cpu.C = 1
	case insSED:
		cpu.D = 1 // useless on NES
	case insSEI:
		cpu.I = 1
	case insSTA:
		cpu.write(addr, cpu.A)
	case insSTX:
		cpu.write(addr, cpu.X)
	case insSTY:
		cpu.write(addr, cpu.Y)
	case insTAX:
		cpu.setValueNZ(&cpu.X, cpu.A)
	case insTAY:
		cpu.setValueNZ(&cpu.Y, cpu.A)
	case insTSX:
		cpu.setValueNZ(&cpu.X, cpu.S)
	case insTXA:
		cpu.setValueNZ(&cpu.A, cpu.X)
	case insTXS:
		cpu.S = cpu.X
	case insTYA:
		cpu.setValueNZ(&cpu.A, cpu.Y)

	// unofficial instruction
	case insDOP, insTOP:
		// do nothing
	case insAAX:
		cpu.write(addr, cpu.X&cpu.A)
	case insDCP:
		cpu.dec(addr)
		cpu.cmp(addr)
	case insISC:
		cpu.inc(addr)
		cpu.sbc(addr)
	case insLAX:
		val := cpu.read(addr)
		cpu.setValueNZ(&cpu.A, val)
		cpu.setValueNZ(&cpu.X, val)
	case insRLA:
		cpu.rol(addr, ins.addrMode)
		cpu.and(addr)
	case insRRA:
		cpu.ror(addr, ins.addrMode)
		cpu.adc(addr)
	case insSLO:
		cpu.asl(addr, ins.addrMode)
		cpu.ora(addr)
	case insSRE:
		cpu.lsr(addr, ins.addrMode)
		cpu.eor(addr)
	default:
		panic(fmt.Sprintf("unknown opcode: %02X at %04X", opcode, cpu.PC))
	}

	return info
}

func (f *cpuFlag) setFlags(val byte) {
	f.C = val & 1
	f.Z = val >> 1 & 1
	f.I = val >> 2 & 1
	f.D = val >> 3 & 1
	f.V = val >> 6 & 1
	f.N = val >> 7 & 1
}

func (f *cpuFlag) flag() byte {
	flag := f.N << 7
	flag |= f.V << 6
	flag |= f.D << 3
	flag |= f.I << 2
	flag |= f.Z << 1
	flag |= f.C
	return flag
}

func (f *cpuFlag) setN(test byte) {
	f.N = test >> 7
}

func (f *cpuFlag) setZ(test byte) {
	if test == 0 {
		f.Z = 1
	} else {
		f.Z = 0
	}
}

func (cpu *CPU) read(addr uint16) byte {
	var data byte
	switch {
	case addr < 0x2000:
		data = cpu.ram[addr&0x07FF]
	case addr < 0x4000:
		// NES PPU registers
	case addr < 0x4018:
		// NES APU & I/O registers
	case addr < 0x4020:
		// ignote
	default:
		data = cpu.console.Mapper.Read(addr)
	}
	return data
}

func (cpu *CPU) write(addr uint16, val byte) {
	switch {
	case addr < 0x2000:
		cpu.ram[addr&0x07FF] = val
	case addr < 0x4000:
		// NES PPU registers
	case addr < 0x4018:
		// NES APU & I/O registers
	case addr < 0x4020:
		// ignore
	default:
		cpu.console.Mapper.Write(addr, val)
	}
}

func (cpu *CPU) writeNZ(addr uint16, val byte) {
	cpu.write(addr, val)
	cpu.setN(val)
	cpu.setZ(val)
}

// there is a bug of indirect mode needs to be implemented
// see http://nesdev.com/6502bugs.txt
func (cpu *CPU) bugRead(addr uint16) uint16 {
	lo, hi := uint16(cpu.read(addr)), uint16(0)
	if addr&0x00FF == 0x00FF {
		hi = uint16(cpu.read(addr & 0xFF00))
	} else {
		hi = uint16(cpu.read(addr + 1))
	}
	return hi<<8 | lo
}

func (cpu *CPU) push(val byte) {
	addr := uint16(cpu.S) | 0x0100
	cpu.write(addr, val)
	cpu.S--
}

func (cpu *CPU) pull() byte {
	cpu.S++
	addr := uint16(cpu.S) | 0x0100
	return cpu.read(addr)
}

func (cpu *CPU) setValueNZ(reg *byte, val byte) {
	*reg = val
	cpu.setN(*reg)
	cpu.setZ(*reg)
}

// Reset CPU to initial state
// http://wiki.nesdev.com/w/index.php/CPU_power_up_state
func (cpu *CPU) Reset() {
	cpu.setFlags(0x34)
	cpu.A = 0
	cpu.X = 0
	cpu.Y = 0
	cpu.S = 0xFD
	cpu.write(0x4017, 0)
	cpu.write(0x4015, 0)
	for i := 0x4000; i <= 0x400F; i++ {
		cpu.write(uint16(i), 0)
	}
}

func (cpu *CPU) adc(addr uint16) {
	val := cpu.read(addr)
	t := int16(cpu.A) + int16(val) + int16(cpu.C)
	if t > 0xFF {
		cpu.C = 1
	} else {
		cpu.C = 0
	}
	if (cpu.A^val)&0x80 == 0 && (cpu.A^byte(t))&0x80 != 0 {
		cpu.V = 1
	} else {
		cpu.V = 0
	}
	cpu.setValueNZ(&cpu.A, byte(t))
}

func (cpu *CPU) and(addr uint16) {
	cpu.setValueNZ(&cpu.A, cpu.A&cpu.read(addr))
}

func (cpu *CPU) asl(addr uint16, addrMode uint8) {
	if addrMode == addrAccumulator {
		cpu.C = cpu.A >> 7
		cpu.setValueNZ(&cpu.A, cpu.A<<1)
	} else {
		val := cpu.read(addr)
		cpu.C = val >> 7
		val <<= 1
		cpu.writeNZ(addr, val)
	}
}

func (cpu *CPU) bcc(addr uint16) {
	if cpu.C == 0 {
		cpu.jmp(addr)
	}
}

func (cpu *CPU) bpl(addr uint16) {
	if cpu.N == 0 {
		cpu.jmp(addr)
	}
}

func (cpu *CPU) bvc(addr uint16) {
	if cpu.V == 0 {
		cpu.jmp(addr)
	}
}

func (cpu *CPU) bcs(addr uint16) {
	if cpu.C != 0 {
		cpu.jmp(addr)
	}
}

func (cpu *CPU) beq(addr uint16) {
	if cpu.Z != 0 {
		cpu.jmp(addr)
	}
}

func (cpu *CPU) bit(addr uint16) {
	val := cpu.read(addr)
	cpu.setN(val)
	cpu.setZ(val & cpu.A)
	cpu.V = (val >> 6) & 1
}

func (cpu *CPU) bmi(addr uint16) {
	if cpu.N != 0 {
		cpu.jmp(addr)
	}
}

func (cpu *CPU) bne(addr uint16) {
	if cpu.Z == 0 {
		cpu.jmp(addr)
	}
}

func (cpu *CPU) bvs(addr uint16) {
	if cpu.V != 0 {
		cpu.jmp(addr)
	}
}

func (cpu *CPU) cmp(addr uint16) {
	val := cpu.read(addr)
	cpu.setN(cpu.A - val)
	cpu.setZ(cpu.A - val)
	if cpu.A >= val {
		cpu.C = 1
	} else {
		cpu.C = 0
	}
}

func (cpu *CPU) cpx(addr uint16) {
	val := cpu.read(addr)
	cpu.setN(cpu.X - val)
	cpu.setZ(cpu.X - val)
	if cpu.X >= val {
		cpu.C = 1
	} else {
		cpu.C = 0
	}
}

func (cpu *CPU) cpy(addr uint16) {
	val := cpu.read(addr)
	cpu.setN(cpu.Y - val)
	cpu.setZ(cpu.Y - val)
	if cpu.Y >= val {
		cpu.C = 1
	} else {
		cpu.C = 0
	}
}

func (cpu *CPU) dec(addr uint16) {
	val := cpu.read(addr) - 1
	cpu.writeNZ(addr, val)
}

func (cpu *CPU) eor(addr uint16) {
	cpu.setValueNZ(&cpu.A, cpu.A^cpu.read(addr))
}

func (cpu *CPU) inc(addr uint16) {
	val := cpu.read(addr) + 1
	cpu.writeNZ(addr, val)
}

func (cpu *CPU) jmp(addr uint16) {
	cpu.PC = addr
}

func (cpu *CPU) jsr(addr uint16) {
	cpu.push(byte((cpu.PC - 1) >> 8))
	cpu.push(byte(cpu.PC - 1))
	cpu.jmp(addr)
}

func (cpu *CPU) lsr(addr uint16, addrMode uint8) {
	if addrMode == addrAccumulator {
		cpu.C = cpu.A & 1
		cpu.setValueNZ(&cpu.A, cpu.A>>1)
	} else {
		val := cpu.read(addr)
		cpu.C = val & 1
		val >>= 1
		cpu.writeNZ(addr, val)
	}
}

func (cpu *CPU) ora(addr uint16) {
	cpu.setValueNZ(&cpu.A, cpu.A|cpu.read(addr))
}

func (cpu *CPU) rol(addr uint16, addrMode uint8) {
	c := cpu.C
	if addrMode == addrAccumulator {
		cpu.C = cpu.A >> 7
		cpu.setValueNZ(&cpu.A, cpu.A<<1|c)
	} else {
		val := cpu.read(addr)
		cpu.C = val >> 7
		val = val<<1 | c
		cpu.writeNZ(addr, val)
	}
}

func (cpu *CPU) ror(addr uint16, addrMode uint8) {
	c := cpu.C
	if addrMode == addrAccumulator {
		cpu.C = cpu.A & 1
		cpu.setValueNZ(&cpu.A, cpu.A>>1|c<<7)
	} else {
		val := cpu.read(addr)
		cpu.C = val & 1
		val = val>>1 | c<<7
		cpu.writeNZ(addr, val)
	}
}

func (cpu *CPU) rti() {
	cpu.setFlags(cpu.pull())
	lo := uint16(cpu.pull())
	hi := uint16(cpu.pull())
	cpu.jmp(hi<<8 | lo)
}

func (cpu *CPU) rts() {
	lo := uint16(cpu.pull())
	hi := uint16(cpu.pull())
	cpu.jmp((hi<<8 | lo) + 1)
}

func (cpu *CPU) sbc(addr uint16) {
	val := cpu.read(addr)
	t := int16(cpu.A) - int16(val) - int16(1-cpu.C)
	if t < 0 {
		cpu.C = 0
	} else {
		cpu.C = 1
	}
	if (cpu.A^val)&0x80 != 0 && (cpu.A^byte(t))&0x80 != 0 {
		cpu.V = 1
	} else {
		cpu.V = 0
	}
	cpu.A = byte(t)
	cpu.setN(cpu.A)
	cpu.setZ(cpu.A)
}
