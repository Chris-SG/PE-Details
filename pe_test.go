package main

import (
	"testing"
)

func TestProcess(t *testing.T) {
	f := Process("./test_executables/strings.exe")

	t.Error(f.Print(0))
}