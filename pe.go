package main

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"os"
)

func main() {
	peFile := Process(os.Args[1])
	f, _ := json.Marshal(peFile)
	fmt.Printf("%s\n", string(f))
}

func Process(path string) (peFile PEFile) {
	f, err := os.Open(path)
	if err != nil {
		return
	}
	b := make([]byte, 0)
	c := 0

	b = make([]byte, 2)
	if c, err = f.ReadAt(b, 0x0); c == 2 && b[0] == 'M' && b[1] == 'Z' {
		b = make([]byte, 4)
		f.ReadAt(b, 0x3c)
		peFile.HasStub = true
		peFile.Offset = PEOffset(binary.LittleEndian.Uint32(b))
		peFile.ExecutableType = Image
	} else if c < 2 || err != nil {
		peFile.ExecutableType = Unknown
		return
	} else {
		peFile.HasStub = false
		peFile.Offset = PEOffset(0)
		peFile.ExecutableType = Object
	}

	additionalOffset := PEOffset(0)

	if peFile.ExecutableType == Image {
		f.Seek(int64(peFile.Offset), 0)
		b = make([]byte, 4)
		f.Read(b)
		exp := []byte{ 'P', 'E', byte(0), byte(0) }
		if bytes.Compare(b, exp) != 0 {
			return
		}
		additionalOffset += 4
	}

	coffHeaderSz := 20
	b = make([]byte, coffHeaderSz)
	f.Seek(int64(peFile.Offset + additionalOffset), io.SeekStart)
	f.Read(b)
	peFile.CoffHeader = CreateCOFFHeader(b)
	additionalOffset += PEOffset(coffHeaderSz)

	if peFile.ExecutableType == Image {
		b = make([]byte, peFile.CoffHeader.OptionalHeaderSz)
		f.Seek(int64(peFile.Offset + additionalOffset), io.SeekStart)
		f.Read(b)
		peFile.OptionalHeader = CreateOptionalHeader(b)
		additionalOffset += PEOffset(peFile.CoffHeader.OptionalHeaderSz)
	}

	b = make([]byte, int(peFile.CoffHeader.NumberOfSections) * binary.Size(SectionTable{}))
	f.Seek(int64(peFile.Offset + additionalOffset), io.SeekStart)
	f.Read(b)
	peFile.SectionTables = CreateSections(b)
	additionalOffset += PEOffset(len(b))

	return
}

func CreateCOFFHeader(b []byte) (h COFFHeader) {
	binary.Read(bytes.NewReader(b), binary.LittleEndian, &h)
	return
}

func CreateOptionalHeader(b []byte) (h OptionalHeader) {
	br := bytes.NewReader(b)
	magicNumber := int16(0x020b)
	fmt.Println(br.Size() - int64(br.Len()) )
	binary.Read(br, binary.LittleEndian, &h.StandardFields)
	h.Is64Bit = magicNumber == h.StandardFields.Magic
	if h.Is64Bit {
		br.Seek(-4, io.SeekCurrent)
		h.StandardFields.BaseOfData = -1

		windowsFields := OHWindowsFields64{}
		fmt.Println(br.Size() - int64(br.Len()) )
		binary.Read(br, binary.LittleEndian, &windowsFields)
		h.WindowsFields = windowsFields
	} else {
		windowsFields := OHWindowsFields32{}
		fmt.Println(br.Size() - int64(br.Len()) )
		binary.Read(br, binary.LittleEndian, &windowsFields)
		h.WindowsFields = windowsFields
	}

	fmt.Println(br.Size() - int64(br.Len()) )

	for br.Len() >= 8 {
		idd := ImageDataDirectory{}
		binary.Read(br, binary.LittleEndian, &idd)
		h.DataDirectories.ImageDataDirectories = append(h.DataDirectories.ImageDataDirectories, idd)
	}

	return
}

func CreateSections(b []byte) (sections []SectionTable) {
	br := bytes.NewReader(b)
	for br.Len() >= binary.Size(SectionTable{}) {
		st := SectionTable{}
		binary.Read(br, binary.LittleEndian, &st)
		sections = append(sections, st)
	}
	return
}