package main

// NROM mapper 000
// http://wiki.nesdev.com/w/index.php/INES_Mapper_000
type NROM struct {
	console     *Console
	PRGBankSize int
}

func init() {
	RegisterMapper(0, &NROM{})
}

// Init initialize mapper
func (m *NROM) Init(con *Console) {
	m.console = con
	m.PRGBankSize = len(m.console.Cartridge.PRG) / (1024 * 16)
}

func (m *NROM) Read(addr uint16) byte {
	var data byte
	switch {
	case addr < 0x8000:
		data = m.console.Cartridge.SRAM[addr-0x6000]
	case addr < 0xC000:
		data = m.console.Cartridge.PRG[addr-0x8000]
	default:
		addr -= 0x8000
		if m.PRGBankSize > 1 {
			data = m.console.Cartridge.PRG[addr]
		} else {
			data = m.console.Cartridge.PRG[addr-0x4000]
		}
	}
	return data
}

func (m *NROM) Write(addr uint16, val byte) {
	if addr < 0x8000 {
		m.console.Cartridge.SRAM[addr-0x6000] = val
	} else {
		addr -= 0x8000
		if m.PRGBankSize == 1 {
			addr -= 0x4000
		}
		m.console.Cartridge.PRG[addr] = val
	}
}
