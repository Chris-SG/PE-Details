package pe_details

import (
	"fmt"
	"reflect"
)

func (pf PEFile) Print(tabCount int) (ret string) {
	prefix := ""
	for i := 0; i < tabCount; i++ {
		prefix += "\t"
	}
	ret = fmt.Sprintf("%s" +
		"%s[\n" +
		"%s\tExecutable Type: %d\n" +
		"%s\tOffset: %#x\n" +
		"%s\tCOFF Header: %s\n" +
		"%s\tOptional Header: %s\n]" +
		"%s\tSectionTables:\n",
		ret,
		prefix, prefix, pf.ExecutableType, prefix, pf.Offset, prefix, pf.CoffHeader.Print(tabCount+1),
		prefix, pf.OptionalHeader.Print(tabCount+1), prefix)
	for _, s := range pf.SectionTables {
		ret = fmt.Sprintf("%s\n%s\t%s", ret, prefix, s.Print(tabCount+1))
	}
	ret = fmt.Sprintf("%s" +
		"\n%s\tCOFF Relocations:\n",
		ret, prefix)
	for _, s := range pf.CoffRelocations {
		ret = fmt.Sprintf("%s\n%s\t%s", ret, prefix, s.Print(tabCount+1))
	}
	ret = fmt.Sprintf("%s%s]\n", ret, prefix)
	return
}

func (ch COFFHeader) Print(tabCount int) (ret string) {
	prefix := ""
	for i := 0; i < tabCount; i++ {
		prefix += "\t"
	}
	ret = fmt.Sprintf( "%s%s[\n", ret, prefix)

	v := reflect.ValueOf(ch)
	for i := 0; i < v.NumField(); i++ {
		ret = fmt.Sprintf( "%s%s\t%s: %#x\n", ret, prefix, v.Type().Field(i).Name, v.Field(i).Interface())
	}
	ret = fmt.Sprintf( "%s%s]\n", ret, prefix)
	return
}

func (oh OptionalHeader) Print(tabCount int) (ret string) {
	prefix := ""
	for i := 0; i < tabCount; i++ {
		prefix += "\t"
	}
	ret = fmt.Sprintf(
		"%s\n" +
			"%s[\n" +
			"%s\tIs 64bit: %t\n" +
			"%s\tStandard Fields: %s\n" +
			"%s\tWindows Fields: %s\n" +
			"%s\tData Directories: %s\n" +
			"%s]\n",
		ret,
		prefix, prefix, oh.Is64Bit, prefix, oh.StandardFields.Print(tabCount+1),
		prefix, oh.WindowsFields.Print(tabCount+1), prefix, oh.DataDirectories.Print(tabCount+1), prefix)
	return
}

func (sf OHStandardFields) Print(tabCount int) (ret string) {
	prefix := ""
	for i := 0; i < tabCount; i++ {
		prefix += "\t"
	}
	ret = fmt.Sprintf( "%s\n%s[\n", ret, prefix)

	v := reflect.ValueOf(sf)
	for i := 0; i < v.NumField(); i++ {
		ret = fmt.Sprintf( "%s%s\t%s: %#x\n", ret, prefix, v.Type().Field(i).Name, v.Field(i).Interface())
	}
	ret = fmt.Sprintf( "%s%s]\n", ret, prefix)
	return
}

func (wf OHWindowsFields32) Print(tabCount int) (ret string) {
	prefix := ""
	for i := 0; i < tabCount; i++ {
		prefix += "\t"
	}
	ret = fmt.Sprintf( "%s\n%s[\n", ret, prefix)

	v := reflect.ValueOf(wf)
	for i := 0; i < v.NumField(); i++ {
		ret = fmt.Sprintf( "%s%s\t%s: %#x\n", ret, prefix, v.Type().Field(i).Name, v.Field(i).Interface())
	}
	ret = fmt.Sprintf( "%s%s]\n", ret, prefix)
	return
}

func (wf OHWindowsFields64) Print(tabCount int) (ret string) {
	prefix := ""
	for i := 0; i < tabCount; i++ {
		prefix += "\t"
	}
	ret = fmt.Sprintf( "%s\n[\n", prefix)

	v := reflect.ValueOf(wf)
	for i := 0; i < v.NumField(); i++ {
		ret = fmt.Sprintf( "%s%s\t%s: %#x\n", ret, prefix, v.Type().Field(i).Name, v.Field(i).Interface())
	}
	ret = fmt.Sprintf( "%s%s]\n", ret, prefix)
	return
}

func (dd OHDataDirectories) Print(tabCount int) (ret string) {
	prefix := ""
	for i := 0; i < tabCount; i++ {
		prefix += "\t"
	}
	for _, d := range dd.ImageDataDirectories {
		ret = fmt.Sprintf( "%s\n%s[\n" +
			"%s\tVirtual Address: %#x\n" +
			"%s\tSize: %#x\n" +
			"%s]",
			ret, prefix, prefix, d.VirtualAddress, prefix, d.Size, prefix)
	}
	return
}

func (st SectionTable) Print(tabCount int) (ret string) {
	prefix := ""
	for i := 0; i < tabCount; i++ {
		prefix += "\t"
	}
	ret = fmt.Sprintf( "%s%s%s\n" +
		"%s[\n" +
		"%s\tVirtual Size: %#x\n" +
		"%s\tVirtual Address: %#x\n" +
		"%s\tSize of Raw Data: %#x\n" +
		"%s\tPointer to Raw Data: %#x\n" +
		"%s\tPointer to Relocations: %#x\n" +
		"%s\tPointer to Line Numbers%#x\n" +
		"%s\tNumber of Relocations%#x\n" +
		"%s\tNumber of Line Numbers%#x\n" +
		"%s\tCharacteristics: %#x\n" +
		"%s\tCharacteristics String: %s\n" +
		"%s]",
		ret,
		prefix, st.Name, prefix, prefix, st.VirtualSize, prefix, st.VirtualAddress,
		prefix, st.SizeOfRawData, prefix, st.PointerToRawData, prefix, st.PointerToRelocations,
		prefix, st.PointerToLineNumbers, prefix, st.NumberOfRelocations, prefix, st.NumberOfLineNumbers,
		prefix, st.Characteristics, prefix, st.CharacteristicsStringify(), prefix)
	return
}

func (cf CoffRelocation) Print(tabCount int) (ret string) {
	prefix := ""
	for i := 0; i < tabCount; i++ {
		prefix += "\t"
	}
	ret = fmt.Sprintf( "%s\n[\n", prefix)

	v := reflect.ValueOf(cf)
	for i := 0; i < v.NumField(); i++ {
		ret = fmt.Sprintf( "%s%s\t%s: %#x\n", ret, prefix, v.Type().Field(i).Name, v.Field(i).Interface())
	}
	ret = fmt.Sprintf( "%s%s]\n", ret, prefix)
	return
}