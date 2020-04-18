package main

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
)

func main() {
	if len(os.Args) == 2 {
		file, _ := os.Open(os.Args[1])
		defer file.Close()
		f, _ := ioutil.ReadAll(file)
		offset := GetPEOffset(f)
		fmt.Println(offset)
		if offset != 0xFFFFFFFF {
			f = f[offset+4:]
			coff := GenerateCoffFileHeader(f)
			j, _ := json.MarshalIndent(coff, "", "  ")
			fmt.Printf("\n%s\n", string(j))
			f = f[20:]
			opt := GenerateOptionalHeader(f)
			j, _ = json.MarshalIndent(opt, "", "  ")
			fmt.Printf("\n%s\n", string(j))
			optLen := len(opt.DataDirectories.ImageDataDirectories) * 8
			if opt.StandardFields.Magic == 0x10b {
				optLen += 96
			} else {
				optLen += 112
			}
			f = f[optLen:]
			sectTable := GenerateSectionTable(f, coff.NumberOfSections)
			j, _ = json.MarshalIndent(sectTable, "", "  ")
			fmt.Printf("\n%s\n", string(j))
		}
	} else {
		f, _ := ioutil.ReadFile("rsrc.dump")
		fmt.Printf("\n%s\n", string(GenerateResourceDirectoryTreeJson(f)))
	}
}

func GetPEOffset(data []byte) (offset uint32) {
	br := bytes.NewReader(data)
	b := make([]byte, 2)
	binary.Read(br, binary.LittleEndian, &b)
	if bytes.Compare(b, []byte{'M', 'Z'}) == 0 {
		br.Seek(0x3c, io.SeekStart)
		binary.Read(br, binary.LittleEndian, &offset)
		return
	} else if bytes.Compare(b, []byte{'P', 'E'}) == 0 {
		return 0
	}
	return 0xFFFFFFFF
}