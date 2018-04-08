package lfile

import (
	"io"
	"os"

	"github.com/azohra/lyra/pkg/lcrypt"
)

//Separator defines delimitor for Salt and Nonce
var Separator = "@!"

//Encipher defines encrypting operations
type Encipher interface {
	//EncipherFile enciphers a file and returns a pointer to a secureLyraFile.
	EncipherFile(key *lcrypt.LKey) (*SecureLyraFile, error)

	Writer
}

//Authenticater defines Authentication operations
type Authenticater interface {
	//GenerateAuthParams generates a new authentication parameters, this involves
	//generating a new nonce for each enciphering of the plaintext. An
	//error is returned in the event of any issues.
	GenerateAuthParams() error
}

//Sanitizer defines secure destruction of sensitive information
type Sanitizer interface {
	DestroyFile() error
}

//Decipher defines decrypting operations
type Decipher interface {
	//DecipherFile deciphers a file and returns a pointer to a lyraFile.
	DecipherFile(key *lcrypt.LKey) (*LyraFile, error)

	Writer
}

//Parser defines parsing utils
type Parser interface {
	//ParseFile, Parses a file
	ParseFile(file string) error
}

//Writer defines write operations
type Writer interface {
	//Writes a to a specific place
	Write(wd string) error
}

//Printer defines printing operations
type Printer interface {
	//Print str to Reader
	Print(reader io.Reader, str string)
}

//readFile read fileName and return a file pointer.
func readFile(fileName string) (*os.File, error) {
	data, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	return data, nil
}
