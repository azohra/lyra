//Package lfile defines a collection of operation and definitions for lyra files
package lfile

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/awnumar/memguard"

	"github.com/azohra/lyra/pkg/lcrypt"
)

//lockedBuffer represents a memguard.LockedBuffer
type lockedBuffer = *memguard.LockedBuffer

//LyraFile represents a deciphered/unencrypted file (also known as a plaintext file).
type LyraFile struct {
	//plaintext
	plaintext lockedBuffer

	//isBeingUsed allows plaintext lockedbuffer to be properly destroyed before
	//allowing any writes to it.
	isBeingUsed bool
}

//NewLyraFile initializes a new Lyrafile
func newLyraFile() *LyraFile {
	return &LyraFile{isBeingUsed: false}
}

//NewParsedLyraFile returns a newly created lyrafile that was parsed from file
func NewParsedLyraFile(file string) (*LyraFile, error) {
	lf := &LyraFile{}

	err := lf.ParseFile(file)
	if err != nil {
		return nil, err
	}

	lf.isBeingUsed = true
	return lf, nil
}

//EncipherFile a LyraFile to a SecureLyraFile and securely destroys the lyrafile upon successful encryption.
func (payload *LyraFile) EncipherFile(key *lcrypt.LKey) (*SecureLyraFile, error) {
	ctxt := newSecureLyraFile()
	err := ctxt.GenerateAuthParams()
	ctxt.SeedSalt(key.GetSalt())

	if err != nil {
		return nil, err
	}

	ctxt.ciphertext, err = lcrypt.AesEncrypt(payload.plaintext.Buffer(), ctxt.nonce, key.GetKey())
	if err != nil {
		return nil, err
	}

	defer payload.DestroyFile()
	return ctxt, nil
}

//ParseFile parses unencrypted file *os.File to a LyraFile struct
func (payload *LyraFile) ParseFile(file string) error {
	if payload.isBeingUsed {
		err := payload.DestroyFile()
		if err != nil {
			return err
		}

		payload.isBeingUsed = false
	}

	ptxt, err := readFile(file)
	defer ptxt.Close()

	buf := bytes.NewBuffer(nil)
	if err != nil {
		return err
	}

	_, err = io.Copy(buf, ptxt)
	if err != nil {
		return err
	}

	if len(buf.Bytes()) == 0 {
		return errors.New("Can not specify an empty file")
	}

	payload.plaintext, err = memguard.NewImmutableFromBytes(buf.Bytes())
	payload.isBeingUsed = true

	return err
}

//Write, writes a LyraFile to a path Wd and destroys the lyrafile upon successful write
func (payload *LyraFile) Write(wd string) error {
	f, err := os.Create(wd)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.Write(payload.plaintext.Buffer())
	if err != nil {
		return err
	}

	defer payload.DestroyFile()
	return nil
}

//DestroyFile safely destroys payload lyrafile, return an error if unable to do so
func (payload *LyraFile) DestroyFile() error {
	payload.plaintext.Destroy()
	if !payload.plaintext.IsDestroyed() {
		return errors.New("Failed to destroy key")
	}
	return nil
}

//PrintLyraFile prints the plaintext of the lyrafile to stdout
func (payload *LyraFile) PrintLyraFile() {
	fmt.Fprintf(os.Stdout, "%s\n", payload.plaintext.Buffer())
}
