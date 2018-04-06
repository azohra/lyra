package main

import "testing"

func TestEncrypt(t *testing.T) {
	file := "../../test/fixture.31.go.txt"
	save := "../../test/fixture.311.go.txt"
	encrypt(file, save, []byte("THIS IS A TEST PASSPHRASE YOU SOULD NOT USE THIS"))

}
