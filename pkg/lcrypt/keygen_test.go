package lcrypt

import (
	"encoding/hex"
	"reflect"
	"testing"
)

func TestGenKey(t *testing.T) {

	salt := []byte("randomsalt")
	pass := []byte("some random password that is meant to be somewhat long to stretch")

	//Generated via official argon2 implementation with same params as defined in crypto.go
	testVector := "77a74e9d5b2de72be1b23a40c5e2c246042575ea66983adfd86c28371d9ff6bb"

	testVectKey, err := hex.DecodeString(testVector)

	if err != nil {
		t.Errorf("Encoding err")
	}

	key := GenKey(pass, salt)

	if !reflect.DeepEqual(testVectKey, key) {
		t.Errorf("KeyGen error, key are not the same")
	}
}

func TestGenSalt(t *testing.T) {
	salt := GenSalt()

	if len(salt) != int(SaltSize) || cap(salt) != int(SaltSize) {
		t.Errorf("Failed could not generate random bytes")
		return
	}
}
