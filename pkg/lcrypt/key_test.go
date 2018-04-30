package lcrypt

import (
	"reflect"
	"testing"
)

func TestNewLKey(t *testing.T) {
	txtk := []byte("A Password")
	salt, err1 := GenSalt()
	if err1 != nil {
		t.Error(err1)
	}

	key, err := NewLKey(txtk, salt)
	if err != nil {
		t.Error(err)
	}

	if key == nil {
		t.Error("Did not init key struct")
	}
}

func TestInitKey(t *testing.T) {
	txtk := []byte("A Password")
	salt, err1 := GenSalt()
	if err1 != nil {
		t.Error(err1)
	}

	key := &LKey{}
	key.initKey(txtk, salt)

	etxtk := []byte("A Password")
	salt2, err1 := GenSalt()
	if err1 != nil {
		t.Error(err1)
	}
	key2 := GenKey(etxtk, salt2)

	if reflect.DeepEqual(key.GetKey(), key2) {
		t.Error("Key init failed, matches an other key that has a unique salt")
	}

	if !reflect.DeepEqual(key.GetKey(), GenKey(etxtk, salt)) {
		t.Error("Key init failed, should match")
	}

	if reflect.DeepEqual(txtk, []byte("A Password")) {
		t.Error("Key wiping failed")
	}

	err := key.key.FillRandomBytes()
	if err == nil {
		t.Error("Write operation should have failed not passed")
	}
}

func TestGetKey(t *testing.T) {
	txtk := []byte("A Password")
	salt, err1 := GenSalt()
	if err1 != nil {
		t.Error(err1)
	}

	key, err := NewLKey(txtk, salt)
	if err != nil {
		t.Error(err)
	}

	txtk1 := []byte("A Password")
	key2 := GenKey(txtk1, salt)
	kk := key.GetKey()

	if !reflect.DeepEqual(kk, key2) {
		t.Error("Failed to fetch key, did not match with generated key")
	}
}

func TestGetSalt(t *testing.T) {
	txtk := []byte("A Password")
	salt, err1 := GenSalt()
	if err1 != nil {
		t.Error(err1)
	}

	key, err := NewLKey(txtk, salt)
	if err != nil {
		t.Error(err)
	}

	txtk1 := []byte("A Password")
	key2 := GenKey(txtk1, salt)
	kk := key.GetKey()

	if !reflect.DeepEqual(kk, key2) {
		t.Error("Failed to fetch key, did not match with generated key")
	}
}
func TestDestroyKey(t *testing.T) {
	txtk := []byte("A Password")
	salt, err1 := GenSalt()
	if err1 != nil {
		t.Error(err1)
	}

	key, err := NewLKey(txtk, salt)
	if err != nil {
		t.Error(err)
	}

	err = key.DestroyKey()
	if err != nil {
		t.Error(err)
	}
}
