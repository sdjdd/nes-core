package main

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"os"
)

// Cartridge -
type Cartridge struct {
	Trainer   []byte
	PRG       []byte
	Chr       []byte
	SRAM      []byte
	Mapper    byte
	Mirroring byte
	Battery   byte
}

type nesHeader struct {
	Magic   uint32
	PrgSize byte
	ChrSize byte
	Flag6   byte
	Flag7   byte
	Flag8   byte
	Flag9   byte
	Flag10  byte
	_       [5]byte
}

func loadRomFile(path string) (*Cartridge, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	header := new(nesHeader)
	if err = binary.Read(file, binary.LittleEndian, header); err != nil {
		return nil, err
	}
	if header.Magic != 0x1a53454e {
		return nil, errors.New("invalid rom file")
	}

	cart := Cartridge{
		Mapper:    header.Flag6>>4 | header.Flag7&0xf0,
		Mirroring: header.Flag6>>2&0x02 | header.Flag6&0x08,
		Battery:   header.Flag6 & 0x02,
		PRG:       make([]byte, int(header.PrgSize)*1024*16),
		Chr:       make([]byte, int(header.ChrSize)*1024*8),
		SRAM:      make([]byte, 1024*8),
	}

	if header.Flag6&0x04 > 0 {
		if _, err := file.Seek(512, 1); err != nil {
			return nil, err
		}
	}
	if _, err := io.ReadFull(file, cart.PRG); err != nil {
		return nil, err
	}
	if _, err := io.ReadFull(file, cart.Chr); err != nil {
		return nil, err
	}

	return &cart, nil
}

func (c *Cartridge) disassembly() string {
	var (
		buf  bytes.Buffer
		pc   int
		nums [2]byte
	)
	for {
		if pc >= len(c.PRG) {
			break
		}
		var (
			op      = c.PRG[pc]
			ins     = instructions[op]
			insName = instructionNames[ins.id]
			size    = instructionSizes[ins.addrMode]
		)
		for i := 1; i < int(size); i++ {
			if pc+i < len(c.PRG) {
				nums[i-1] = c.PRG[pc+i]
			} else {
				nums[i-1] = 0
			}
		}
		buf.WriteString(insName)
		if size > 1 {
			buf.WriteByte('\t')
		}
		switch ins.addrMode {
		case addrImplied: // do nothing
		case addrAccumulator:
			buf.WriteString("A")
		case addrImmediate:
			buf.WriteString(fmt.Sprintf("#$%02x", nums[0]))
		case addrZeroPage:
			buf.WriteString(fmt.Sprintf("$%02x", nums[0]))
		case addrZeroPageX:
			buf.WriteString(fmt.Sprintf("$%02x,X", nums[0]))
		case addrZeroPageY:
			buf.WriteString(fmt.Sprintf("$%02x,Y", nums[0]))
		case addrRelative:
			buf.WriteString(fmt.Sprintf("*%+d", nums[0]))
		case addrAbsolute:
			buf.WriteString(fmt.Sprintf("$%02x%02x", nums[1], nums[0]))
		case addrAbsoluteX:
			buf.WriteString(fmt.Sprintf("$%02x%02x,X", nums[1], nums[0]))
		case addrAbsoluteY:
			buf.WriteString(fmt.Sprintf("$%02x%02x,Y", nums[1], nums[0]))
		case addrIndirect:
			buf.WriteString(fmt.Sprintf("($%02x%02x)", nums[1], nums[0]))
		case addrIndexedIndirect:
			buf.WriteString(fmt.Sprintf("($%02x,X)", nums[0]))
		case addrIndirectIndexed:
			buf.WriteString(fmt.Sprintf("($%02x),Y", nums[0]))
		}
		buf.WriteByte('\n')
		pc += int(size)
	}
	return buf.String()
}
