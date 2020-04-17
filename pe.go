package pe_details

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
)

//export GeneratePEFileStructure
func GeneratePEFileStructureJson(b []byte) (j []byte) {
	j, _ = json.Marshal(GeneratePEFileStructure(b))
	return
}

func GeneratePEFileStructure(data []byte) (peFile PEFile) {
	b := make([]byte, 2)
	br := bytes.NewReader(data)

	br.Read(b)
	if len(b) == 2 && b[0] == 'M' && b[1] == 'Z' {
		b = make([]byte, 4)
		br.Seek(0x3c, io.SeekStart)
		br.Read(b)
		peFile.HasStub = true
		peFile.Offset = PEOffset(binary.LittleEndian.Uint32(b))
		peFile.ExecutableType = Image
	} else if len(b) < 2 {
		peFile.ExecutableType = Unknown
		return
	} else {
		peFile.HasStub = false
		peFile.Offset = PEOffset(0)
		peFile.ExecutableType = Object
	}

	additionalOffset := PEOffset(0)

	if peFile.ExecutableType == Image {
		br.Seek(int64(peFile.Offset), io.SeekStart)
		b = make([]byte, 4)
		br.Read(b)
		exp := []byte{ 'P', 'E', byte(0), byte(0) }
		if bytes.Compare(b, exp) != 0 {
			return
		}
		additionalOffset += 4
	}

	coffHeaderSz := 20
	b = make([]byte, coffHeaderSz)
	br.Seek(int64(peFile.Offset + additionalOffset), io.SeekStart)
	br.Read(b)
	peFile.CoffHeader = createCOFFHeader(b)
	additionalOffset += PEOffset(coffHeaderSz)

	if peFile.ExecutableType == Image {
		b = make([]byte, peFile.CoffHeader.OptionalHeaderSz)
		br.Seek(int64(peFile.Offset + additionalOffset), io.SeekStart)
		br.Read(b)
		peFile.OptionalHeader = createOptionalHeader(b)
		additionalOffset += PEOffset(peFile.CoffHeader.OptionalHeaderSz)
	}

	b = make([]byte, int(peFile.CoffHeader.NumberOfSections) * binary.Size(SectionTable{}))
	br.Seek(int64(peFile.Offset + additionalOffset), io.SeekStart)
	br.Read(b)
	peFile.SectionTables = createSections(b)
	additionalOffset += PEOffset(len(b))

	if peFile.ExecutableType == Object {
		peFile.CoffRelocations = peFile.createCOFFRelocations(br)
	}

	return
}

func createCOFFHeader(b []byte) (h COFFHeader) {
	binary.Read(bytes.NewReader(b), binary.LittleEndian, &h)
	return
}

func createOptionalHeader(b []byte) (h OptionalHeader) {
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

func createSections(b []byte) (sections []SectionTable) {
	br := bytes.NewReader(b)
	for br.Len() >= binary.Size(SectionTable{}) {
		st := SectionTable{}
		binary.Read(br, binary.LittleEndian, &st)
		sections = append(sections, st)
	}
	return
}

func (pf PEFile) createCOFFRelocations(br *bytes.Reader) (coffRelocations []CoffRelocation) {
	for _, section := range pf.SectionTables {
		br.Seek(int64(section.PointerToRawData), io.SeekStart)
		cr := CoffRelocation{}
		binary.Read(br, binary.LittleEndian, &cr)
		coffRelocations = append(coffRelocations, cr)
	}
	return
}

func (pf PEFile) createCOFFSymbolTable(br *bytes.Reader) (coffSymbolTable CoffSymbolTable) {
	if pf.CoffHeader.NumberOfSymbols > 0 {
		br.Seek(int64(pf.CoffHeader.SymbolTablePtr), io.SeekStart)
		for idx := int32(0); idx < pf.CoffHeader.NumberOfSymbols; idx++ {
			coffSymbolTable.Symbols = append(coffSymbolTable.Symbols, createCOFFSymbol(br))
		}
	}
	return
}

func createCOFFSymbol(br *bytes.Reader) (coffSymbol CoffSymbol) {
	coffSymbol.AuxSymbols = make([]AuxSymbol, 0)

	binary.Read(br, binary.LittleEndian, &coffSymbol)
	coffSymbol.AuxSymbols = createAuxSymbols(br, coffSymbol.NumberOfAuxSymbols)

	return
}

func createAuxSymbols(br *bytes.Reader, count byte) (auxSymbols []AuxSymbol) {
	auxSymbols = make([]AuxSymbol, count)
	binary.Read(br, binary.LittleEndian, &auxSymbols)

	return
}