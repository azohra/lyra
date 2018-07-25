package lcrypt

import (
	"testing"
)

func TestEncrypt(t *testing.T) {
	file := "../../test/fixture.31.go.txt"
	save := "../../test/fixture.311.go.txt"
	Encrypt(file, save, []byte("THIS IS A TEST PASSPHRASE YOU SHOULD NOT USE THIS"))
}

func TestDecrypt(t *testing.T) {
	file := "../../test/fixture.32.go.txt"
	save := "../../test/fixture.322.go.txt"
	Decrypt(file, save, false, []byte("THIS IS A TEST PASSPHRASE YOU SHOULD NOT USE THIS"))
}
