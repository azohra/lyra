//Package lcrypt provides crypto operations for lyra
package lcrypt

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/hex"
	"reflect"
	"testing"
)

func TestEncipher(t *testing.T) {
	key := "77a74e9d5b2de72be1b23a40c5e2c246042575ea66983adfd86c28371d9ff6bb"
	binkey, _ := hex.DecodeString(key)
	plaintext := []byte("exampleplaintext")
	nonce, _ := hex.DecodeString("64a9433eae7ccceee2fc0eda")

	block, err := aes.NewCipher(binkey)
	if err != nil {
		t.Error(err)
	}
	aesgcm, err := cipher.NewGCM(block)
	ciphertxt := aesgcm.Seal(nil, nonce, plaintext, nil)

	ct, err := AesEncrypt(plaintext, nonce, binkey)
	if err != nil {
		t.Error(err)
	}

	if !reflect.DeepEqual(ciphertxt, ct) {
		t.Error("Encryption test failed")
	}

}

func TestDecipher(t *testing.T) {
	key := "77a74e9d5b2de72be1b23a40c5e2c246042575ea66983adfd86c28371d9ff6bb"
	binkey, _ := hex.DecodeString(key)
	plaintext := []byte("exampleplaintext")
	nonce, _ := hex.DecodeString("64a9433eae7ccceee2fc0eda")

	ct, err := AesEncrypt(plaintext, nonce, binkey)
	if err != nil {
		t.Error(err)
	}

	pt, err := AesDecrypt(ct, nonce, binkey)
	if err != nil {
		t.Error(err)
	}
	if !reflect.DeepEqual(plaintext, pt) {
		t.Error("Decryption test failed")
	}

	// test cases where it is supose to returned failed, where integrity of the file has been tampered with
	pt, err = AesDecrypt(append(ct, byte(0)), nonce, binkey)
	if err == nil {
		t.Error("Decryption test failed should have actually failed, integrity was purposely tampered with")
	}

	nonce, _ = hex.DecodeString("64a9433aae7ccceeb2fc0eda")
	pt, err = AesDecrypt(ct, nonce, binkey)
	if err == nil {
		t.Error("Decryption test failed should have actually failed, nonce was purposely tampered with")
	}

	if reflect.DeepEqual(plaintext, pt) {
		t.Error("Decryption test failed should have actually failed, returned the same plaintext")
	}
}

func TestGenNonce(t *testing.T) {
	dst := make([]byte, NonceSize+12)
	_, err := GenNonce(dst)
	if err != nil {
		t.Error("Failed to set random numbers to dst")
	}

	//For cases when destination is not specified the returning arr wil be len() of NonceSize
	dst2, err := GenNonce(nil)
	if err != nil {
		t.Error("Failed to generate numbers")
	}

	if len(dst) != NonceSize+12 || cap(dst) != NonceSize+12 {
		t.Error("Param slice was not generated properly")
	}

	if len(dst2) != NonceSize || cap(dst2) != NonceSize {
		t.Error("Return slice was not generated properly")
	}
}
