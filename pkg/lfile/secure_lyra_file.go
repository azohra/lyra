package lfile

import (
	"bytes"
	"encoding/hex"
	"errors"
	"io"
	"os"
	"strings"

	"github.com/awnumar/memguard"

	"github.com/azohra/lyra/pkg/lcrypt"
)

//SecureLyraFile represents an enciphered file.
type SecureLyraFile struct {
	ciphertext []byte

	//salt used to initialize a key
	salt []byte

	//nonce used for authentication
	nonce []byte
}

//newSecureLyraFile initializes a new SecureLyraFile
func newSecureLyraFile() *SecureLyraFile {
	return &SecureLyraFile{}
}

//NewParsedSLFile returns a newly created securelyrafile that was parsed from file
func NewParsedSLFile(file string) (*SecureLyraFile, error) {
	slf := &SecureLyraFile{}

	err := slf.ParseFile(file)
	if err != nil {
		return nil, err
	}
	return slf, nil
}

//GenerateAuthParams for a new nonce and salt parameters for an encryption of a LyraFile
func (payload *SecureLyraFile) GenerateAuthParams() error {
	payload.nonce = make([]byte, lcrypt.NonceSize)

	_, err := lcrypt.GenNonce(payload.nonce)
	if err != nil {
		return err
	}

	return nil
}

//DecipherFile deciphers a SecureLyraFile to a new LyraFile
func (payload *SecureLyraFile) DecipherFile(key *lcrypt.LKey) (*LyraFile, error) {
	ptxt := newLyraFile()

	plain, err := lcrypt.AesDecrypt(payload.ciphertext, payload.nonce, key.GetKey())
	if err != nil {
		return nil, err
	}

	ptxt.plaintext, err = memguard.NewImmutableFromBytes(plain)
	ptxt.isBeingUsed = true
	return ptxt, err
}

//SeedSalt seeds a salt value to a secureLyraFile
func (payload *SecureLyraFile) SeedSalt(salt []byte) {
	payload.salt = salt
}

//RetrieveSalt returns this securelyrafile's salt value
func (payload *SecureLyraFile) RetrieveSalt() []byte {
	return payload.salt
}

//ParseFile parses encrypted file *os.File to a SecureLyraFile struct
func (payload *SecureLyraFile) ParseFile(file string) error {
	ctxt, err := readFile(file)
	if err != nil {
		return err
	}

	defer ctxt.Close()

	var str []byte
	var authData []byte

	buf := bytes.NewBuffer(nil)
	_, err = io.Copy(buf, ctxt)
	if err != nil {
		return err
	}

	str = buf.Bytes()
	start := 0
	for _, val := range str {
		if val == '\n' {
			break
		}
		authData = append(authData, val)
		start++
	}

	if len(authData) == 0 {
		return errors.New("Could Not Load Authentication Data, the file may be corrupt or has been tempered with")
	}

	payload.ciphertext = str[start+1:]
	payload.salt, payload.nonce, err = parseAuthData(string(authData))
	if err != nil {
		return err
	}

	return nil
}

//Write, writes a SecureLyraFile to a path wd
func (payload *SecureLyraFile) Write(wd string) error {
	authData := Separator + hex.EncodeToString(payload.salt) + Separator + hex.EncodeToString(payload.nonce) + "\n"

	f, err := os.Create(wd)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.Write([]byte(authData))
	if err != nil {
		return err
	}

	_, err = f.Write(payload.ciphertext)
	if err != nil {
		return err
	}

	return nil
}

//parseAuthData parses auth data from data
func parseAuthData(data string) ([]byte, []byte, error) {
	adata := strings.Split(data, Separator)
	if len(adata) != 3 {
		return nil, nil, errors.New("Parsing Failed")
	}

	s, err := hex.DecodeString(adata[1])
	if err != nil {
		return nil, nil, err
	}

	n, err := hex.DecodeString(adata[2])
	if err != nil {
		return nil, nil, err
	}

	return s, n, nil
}
