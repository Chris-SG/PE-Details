package pe_details

import (
	"encoding/json"
	"io/ioutil"
	"testing"
)

func TestProcess(t *testing.T) {
	b, _ := ioutil.ReadFile("./test_executables/strings.exe")
	f := GeneratePEFileStructure(b)

	t.Error(f.Print(0))
}

func TestRsrc(t *testing.T) {
	b, _ := ioutil.ReadFile("./rsrc.dump")
	f := GenerateResourceDirectoryTree(b)

	j, _ := json.MarshalIndent(f, "", "  ")
	t.Errorf("\n%s\n", string(j))
}