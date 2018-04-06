//Package lcrypt provides crypto operations for lyra
package lcrypt

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
)

//NonceSize defines size of nonce in bytes, NIST specifies for GCM, a nonce size of 12 bytes or 96 bits.
const NonceSize = 12

//AesEncrypt encrypts and authenticates a plaintext via the AES encryption scheme in GCM, all values must
//be in decoded into raw values and can not be in string representation.
func AesEncrypt(plain, nonce, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)

	if err != nil {
		return nil, err
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	ciphertxt := aesgcm.Seal(nil, nonce, plain, nil)

	return ciphertxt, nil
}

//AesDecrypt decrypts and authenticates a ciphertxt, all values must be decoded into raw values
//and can not be in string hex representation.
func AesDecrypt(ciphertxt, nonce, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	plaintext, err := aesgcm.Open(nil, nonce, ciphertxt, nil)
	if err != nil {
		return nil, err
	}

	return plaintext, nil
}

//GenNonce generates a nonce via a cryptographicaly secure number generator to dst or if dst is
//nil returns a nonce of size NonceSize.
func GenNonce(dst []byte) ([]byte, error) {

	if dst != nil {
		_, err := rand.Read(dst)
		if err != nil {
			return nil, err
		}
		return nil, nil
	}

	nonce := make([]byte, NonceSize)
	_, err := rand.Read(nonce)
	if err != nil {
		return nil, err
	}
	return nonce, nil

}
