package lfile

import "testing"

func TestReadFile(t *testing.T) {
	fixture := "../../test/fixture.txt"

	file, err := readFile(fixture)
	if err != nil {
		t.Error(err)
	}
	if file == nil {
		t.Error("Failed to open")
	}

	fixture = "../../test/iDONTEXIST.txt"
	file, err = readFile(fixture)
	if err == nil || file != nil {
		t.Error("Error should have failed to open")
	}
}
