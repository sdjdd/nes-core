package main

import "fmt"

// Mapper of NES
// see http://wiki.nesdev.com/w/index.php/Mapper
type Mapper interface {
	Init(con *Console)
	Read(addr uint16) byte
	Write(addr uint16, val byte)
}

var mappers [768]Mapper

// RegisterMapper - register a mapper by id
func RegisterMapper(id int, mapper Mapper) {
	mappers[id] = mapper
}

// GetMapper - get a mapper by id
func GetMapper(id int) (mapper Mapper, err error) {
	if id < 0 || id >= len(mappers) {
		err = fmt.Errorf("invalid mapper id")
	} else if mapper = mappers[id]; mapper == nil {
		err = fmt.Errorf("mapper %d not implemented", id)
	}
	return
}
