package main

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"reflect"
)

//export GenerateResourceDirectoryTree
func GenerateResourceDirectoryTreeJson(b []byte) (j []byte) {
	j, _ = json.Marshal(GenerateResourceDirectoryTree(b))
	return
}

// GenerateResourceDirectoryTree will build a rsrc tree from the provided
// data. Assumes that the data provided contains the entire rsrc data block.
// All offsets will be relative to zero, so no prepended data should be provided.
// Additionally, any use of the returned tree will need calculated offsets from
// the caller to ensure proper data access.
func GenerateResourceDirectoryTree(b []byte) (rdt ResourceDirectoryTable) {
	br := bytes.NewReader(b)

	rdt = parseAsTable(br)

	return
}

func parseAsTable(br *bytes.Reader) (rdt ResourceDirectoryTable) {
	brPos := bytesRead(br)
	binary.Read(br, binary.LittleEndian, &rdt.Metadata)

	for idx := int16(0); idx < rdt.Metadata.NumberOfNameEntries; idx++ {
		rdt.Entries = append(rdt.Entries, parseAsDirectoryEntry(br))
		br.Seek(8, io.SeekCurrent)
	}

	for idx := int16(0); idx < rdt.Metadata.NumberOfIdEntries; idx++ {
		rdt.Entries = append(rdt.Entries, parseAsDirectoryEntry(br))
		br.Seek(8, io.SeekCurrent)
	}

	br.Seek(brPos, io.SeekStart)
	return
}

func parseAsDirectoryEntry(br *bytes.Reader) (e ResourceDirectoryEntry) {
	brPos := bytesRead(br)
	binary.Read(br, binary.LittleEndian, &e.Metadata)

	if (e.Metadata.Offset >> 31) == 0x1 {
		br.Seek(int64(e.Metadata.Offset&0x7FFFFFFF), io.SeekStart)
		e.Item = parseAsTable(br)
	} else {
		br.Seek(int64(e.Metadata.Offset&0x7FFFFFFF), io.SeekStart)
		e.Item = parseAsDataEntry(br)
	}
	e.Metadata.Offset &= 0x7FFFFFFF

	br.Seek(brPos, io.SeekStart)
	return
}

func parseAsDataEntry(br *bytes.Reader) (e ResourceDataEntry) {
	brPos := bytesRead(br)
	binary.Read(br, binary.LittleEndian, &e)

	br.Seek(brPos, io.SeekStart)
	return
}

func bytesRead(br *bytes.Reader) int64 {
	return br.Size() - int64(br.Len())
}

type ResourceDirectoryTable struct {
	Metadata ResourceDirectoryTableMetadata
	Entries []ResourceDirectoryEntry
}

type ResourceDirectoryTableMetadata struct {
	Characteristics int32
	TimeDate int32
	MajorVersion int16
	MinorVersion int16

	NumberOfNameEntries int16
	NumberOfIdEntries int16
}

type ResourceDirectoryEntry struct {
	Name string
	Metadata ResourceDirectoryEntryMetadata

	Item interface{}
}

type ResourceDirectoryEntryMetadata struct {
	Identifier int32
	Offset uint32
}

type ResourceDataEntry struct {
	DataRva int32
	Size int32
	Codepage int32
	Reserved int32
}

func (m ResourceDirectoryTableMetadata) MarshalJSON() ([]byte, error) {
	b := make([]byte, 0)
	b = append(b, '{')
	v := reflect.ValueOf(m)
	for i := 0; i < v.NumField(); i++ {
		b = append(b, fmt.Sprintf("\"%s\":\"%#x\",", v.Type().Field(i).Name, v.Field(i).Interface())...)
	}
	b = append(b[:len(b)-1], '}')
	return b, nil
}

func (e ResourceDataEntry) MarshalJSON() ([]byte, error) {
	b := make([]byte, 0)
	b = append(b, '{')
	v := reflect.ValueOf(e)
	for i := 0; i < v.NumField(); i++ {
		b = append(b, fmt.Sprintf("\"%s\":\"%#x\",", v.Type().Field(i).Name, v.Field(i).Interface())...)
	}
	b = append(b[:len(b)-1], '}')
	return b, nil
}