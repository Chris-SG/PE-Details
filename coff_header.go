package main

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"reflect"
)

//export GenerateCoffFileHeader
func GenerateCoffFileHeaderJson(b []byte) (j []byte) {
	j, _ = json.Marshal(GenerateResourceDirectoryTree(b))
	return
}

// GenerateResourceDirectoryTree will build a rsrc tree from the provided
// data. Assumes that the data provided contains the entire rsrc data block.
// All offsets will be relative to zero, so no prepended data should be provided.
// Additionally, any use of the returned tree will need calculated offsets from
// the caller to ensure proper data access.
func GenerateCoffFileHeader(b []byte) (cfh CoffFileHeader) {
	br := bytes.NewReader(b)

	cfh = parseAsCoffFileHeader(br)

	return
}

func parseAsCoffFileHeader(br *bytes.Reader) (cfh CoffFileHeader) {
	binary.Read(br, binary.LittleEndian, &cfh)
	return
}

type CoffFileHeader struct {
	Machine uint16
	NumberOfSections uint16
	TimeDateStamp uint32
	PointerToSymbolTable uint32
	NumberOfSymbols uint32
	SizeOfOptionalHeader uint16
	Characteristics uint16
}

func (cfh CoffFileHeader) MarshalJSON() ([]byte, error) {
	b := make([]byte, 0)
	b = append(b, '{')
	v := reflect.ValueOf(cfh)
	for i := 0; i < v.NumField(); i++ {
		b = append(b, fmt.Sprintf("\"%s\":\"%#x\",", v.Type().Field(i).Name, v.Field(i).Interface())...)
	}
	b = append(b[:len(b)-1], '}')
	return b, nil
}