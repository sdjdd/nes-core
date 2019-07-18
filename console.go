package main

import (
	"errors"
)

// Console of NES
type Console struct {
	CPU       *CPU
	Cartridge *Cartridge
	Mapper    Mapper
}

// Connect a device to console
func (con *Console) Connect(device interface{}) error {
	switch device.(type) {
	case *CPU:
		cpu := device.(*CPU)
		cpu.console = con
		cpu.Reset()
	case *Cartridge:
		cart := device.(*Cartridge)
		mapper, err := GetMapper(int(cart.Mapper))
		if err != nil {
			return err
		}
		con.Cartridge = cart
		con.Mapper = mapper
		mapper.Init(con)
	default:
		return errors.New("unknown device")
	}
	return nil
}
