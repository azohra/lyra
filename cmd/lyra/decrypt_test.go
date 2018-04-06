package main

import "testing"

func TestDecrypt(t *testing.T) {
	file := "../../test/fixture.32.go.txt"
	save := "../../test/fixture.322.go.txt"
	decrypt(file, save, false, []byte("THIS IS A TEST PASSPHRASE YOU SOULD NOT USE THIS"))
}
