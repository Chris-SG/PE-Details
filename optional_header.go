package main

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"reflect"
)

//export GenerateOptionalHeader
func GenerateOptionalHeaderJson(b []byte) (j []byte) {
	j, _ = json.Marshal(GenerateOptionalHeader(b))
	return
}

// GenerateResourceDirectoryTree will build a rsrc tree from the provided
// data. Assumes that the data provided contains the entire rsrc data block.
// All offsets will be relative to zero, so no prepended data should be provided.
// Additionally, any use of the returned tree will need calculated offsets from
// the caller to ensure proper data access.
func GenerateOptionalHeader(b []byte) (oh OptionalHeader) {
	br := bytes.NewReader(b)

	oh = parseAsOptionalHeader(br)

	return
}

func parseAsOptionalHeader(br *bytes.Reader) (oh OptionalHeader) {
	oh.StandardFields = parseOptionalHeaderStandardFields(br)
	oh.WindowsFields = parseOptionalHeaderWindowsFields(br, oh.StandardFields.Magic)
	if oh.StandardFields.Magic == 0x10B {
		oh.DataDirectories = parseOptionalHeaderDataDirectories(br, oh.WindowsFields.(OptionalHeaderWindowsFields32).NumberOfRvaAndSizes)
	} else if oh.StandardFields.Magic == 0x20B {
		oh.DataDirectories = parseOptionalHeaderDataDirectories(br, oh.WindowsFields.(OptionalHeaderWindowsFields64).NumberOfRvaAndSizes)
	}
	return
}

func parseOptionalHeaderStandardFields(br *bytes.Reader) (sf OptionalHeaderStandardFields) {
	binary.Read(br, binary.LittleEndian, &sf)
	if sf.Magic == 0x20B {
		sf.BaseOfData = 0xFFFF
		br.Seek(-4, io.SeekCurrent)
	}
	return
}

func parseOptionalHeaderWindowsFields(br *bytes.Reader, magic uint16) (wf interface{}) {
	if magic == 0x10B {
		tmp := OptionalHeaderWindowsFields32{}
		binary.Read(br, binary.LittleEndian, &tmp)
		wf = tmp
	} else if magic == 0x20B {
		tmp := OptionalHeaderWindowsFields64{}
		binary.Read(br, binary.LittleEndian, &tmp)
		wf = tmp
	} else {
		panic("invalid magic")
	}
	return
}

func parseOptionalHeaderDataDirectories(br *bytes.Reader, numberOfRvaAndSizes uint32) (dd OptionalHeaderDataDirectories) {
	for num := uint32(0); num < numberOfRvaAndSizes; num++ {
		idd := ImageDataDirectory{}
		binary.Read(br, binary.LittleEndian, &idd)
		dd.ImageDataDirectories = append(dd.ImageDataDirectories, idd)
	}
	return
}

type OptionalHeader struct {
	StandardFields OptionalHeaderStandardFields
	WindowsFields interface{}
	DataDirectories OptionalHeaderDataDirectories
}

type OptionalHeaderStandardFields struct {
	Magic uint16
	MajorLinkerVersion uint8
	MinorLinkerVersion uint8
	SizeOfCode uint32
	SizeOfInitializedData uint32
	SizeOfUninitializedData uint32
	AddressOfEntryPoint uint32
	BaseOfCode uint32

	BaseOfData uint32
}

func (sf OptionalHeaderStandardFields) MarshalJSON() ([]byte, error) {
	b := make([]byte, 0)
	b = append(b, '{')
	v := reflect.ValueOf(sf)
	for i := 0; i < v.NumField(); i++ {
		b = append(b, fmt.Sprintf("\"%s\":\"%#x\",", v.Type().Field(i).Name, v.Field(i).Interface())...)
	}
	b = append(b[:len(b)-1], '}')
	return b, nil
}

type OptionalHeaderWindowsFields32 struct {
	ImageBase uint32
	SectionAlignment uint32
	FileAlignment uint32
	MajorOperatingSystemVersion uint16
	MinorOperatingSystemVersion uint16
	MajorImageVersion uint16
	MinorImageVersion uint16
	MajorSubsystemVersion uint16
	MinorSubsystemVersion uint16
	Win32VersionValue uint32
	SizeOfImage uint32
	SizeOfHeaders uint32
	CheckSum uint32
	Subsystem uint16
	DllCharacteristics uint16
	SizeOfStackReserve uint32
	SizeOfStackCommit uint32
	SizeOfHeapReserve uint32
	SizeOfHeapCommit uint32
	LoaderFlags uint32
	NumberOfRvaAndSizes uint32
}

func (wf OptionalHeaderWindowsFields32) MarshalJSON() ([]byte, error) {
	b := make([]byte, 0)
	b = append(b, '{')
	v := reflect.ValueOf(wf)
	for i := 0; i < v.NumField(); i++ {
		b = append(b, fmt.Sprintf("\"%s\":\"%#x\",", v.Type().Field(i).Name, v.Field(i).Interface())...)
	}
	b = append(b[:len(b)-1], '}')
	return b, nil
}

type OptionalHeaderWindowsFields64 struct {
	ImageBase uint64
	SectionAlignment uint32
	FileAlignment uint32
	MajorOperatingSystemVersion uint16
	MinorOperatingSystemVersion uint16
	MajorImageVersion uint16
	MinorImageVersion uint16
	MajorSubsystemVersion uint16
	MinorSubsystemVersion uint16
	Win32VersionValue uint32
	SizeOfImage uint32
	SizeOfHeaders uint32
	CheckSum uint32
	Subsystem uint16
	DllCharacteristics uint16
	SizeOfStackReserve uint64
	SizeOfStackCommit uint64
	SizeOfHeapReserve uint64
	SizeOfHeapCommit uint64
	LoaderFlags uint32
	NumberOfRvaAndSizes uint32
}

func (wf OptionalHeaderWindowsFields64) MarshalJSON() ([]byte, error) {
	b := make([]byte, 0)
	b = append(b, '{')
	v := reflect.ValueOf(wf)
	for i := 0; i < v.NumField(); i++ {
		b = append(b, fmt.Sprintf("\"%s\":\"%#x\",", v.Type().Field(i).Name, v.Field(i).Interface())...)
	}
	b = append(b[:len(b)-1], '}')
	return b, nil
}

type OptionalHeaderDataDirectories struct {
	ImageDataDirectories []ImageDataDirectory
}

func (dd OptionalHeaderDataDirectories) MarshalJSON() ([]byte, error) {
	var names = map[int]string {
		0: "Export Table .edata",
		1: "Import Table .idata",
		2: "Resource Table .rsrc",
		3: "Exception Table .pdata",
		4: "Certificate Table",
		5: "Base Relocation Table .reloc",
		6: "Debug .debug",
		7: "Architecture",
		8: "Global Ptr",
		9: "TLS Table .tls",
		10: "Load Config Table",
		11: "Bound Import",
		12: "IAT",
		13: "DelAy Import Descriptor",
		14: "CLR Runtime",
		15: "Reserved",
	}

	b := make([]byte, 0)
	b = append(b, '[')
	for i, idd := range dd.ImageDataDirectories {
		b = append(b, '{')
		v := reflect.ValueOf(idd)
		b = append(b, fmt.Sprintf("\"%s\":\"%s\",", "Name", names[i])...)
		for j := 0; j < v.NumField(); j++ {
			b = append(b, fmt.Sprintf("\"%s\":\"%#x\",", v.Type().Field(j).Name, v.Field(j).Interface())...)
		}
		b = append(b[:len(b)-1], "},"...)
	}
	b = append(b[:len(b)-1], ']')
	fmt.Println(string(b))
	return b, nil
}

type ImageDataDirectory struct {
	VirtualAddress uint32
	Size uint32
}