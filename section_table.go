package main

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"reflect"
)

//export GenerateSectionTable
func GenerateSectionTableJson(b []byte, numberOfSections uint16) (j []byte) {
	j, _ = json.Marshal(GenerateSectionTable(b, numberOfSections))
	return
}

// GenerateResourceDirectoryTree will build a rsrc tree from the provided
// data. Assumes that the data provided contains the entire rsrc data block.
// All offsets will be relative to zero, so no prepended data should be provided.
// Additionally, any use of the returned tree will need calculated offsets from
// the caller to ensure proper data access.
func GenerateSectionTable(b []byte, numberOfSections uint16) (cfh SectionTable) {
	br := bytes.NewReader(b)

	cfh = parseAsSectionTable(br, numberOfSections)

	return
}

func parseAsSectionTable(br *bytes.Reader, numberOfSections uint16) (st SectionTable) {
	for idx := uint16(0); idx < numberOfSections; idx++ {
		sh := SectionHeader{}
		binary.Read(br, binary.LittleEndian, &sh)
		st.SectionHeaders = append(st.SectionHeaders, sh)
	}
	return
}

type SectionTable struct {
	SectionHeaders []SectionHeader
}

type SectionHeader struct {
	Name [8]byte
	VirtualSize uint32
	VirtualAddress uint32
	SizeOfRawData uint32
	PointerToRawData uint32
	PointerToRelocations uint32
	PointerToLineNumbers uint32
	NumberOfRelocations uint16
	NumberOfLineNumbers uint16
	Characteristics uint32
}

func (sh SectionHeader) MarshalJSON() ([]byte, error) {
	b := make([]byte, 0)
	b = append(b, '{')
	v := reflect.ValueOf(sh)
	byteArr := v.Field(0).Interface().([8]byte)
	byteSlice := byteArr[:]
	byteSlice = byteSlice[:bytes.IndexByte(byteSlice, 0)]
	b = append(b, fmt.Sprintf("\"%s\":\"%s\",", v.Type().Field(0).Name, string(byteSlice))...)
	for i := 1; i < v.NumField(); i++ {
		b = append(b, fmt.Sprintf("\"%s\":\"%#x\",", v.Type().Field(i).Name, v.Field(i).Interface())...)
	}
	b = append(b[:len(b)-1], '}')
	return b, nil
}