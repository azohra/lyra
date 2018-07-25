package main

import (
	"errors"
	"flag"

	"github.com/fvumbaca/lyra/pkg/encryption"
)

const (
	helpstrdec = `
The following examples are all the possible options for the "decrypt" command:

lyra decrypt file

	Decrypts and overides file (user specified file). Users will be asked to provide
	the passphrase via stdin.

lyra decrypt -s file1 file

	Decrypts the contents of file (user specified file) and save the resulting output 
	to file1 (user specified file or specified path). Users will be asked to provide 
	the passphrase via stdin. Original specified file will still remain encrypted
	with the same key.

lyra decrypt -p "mypassphrase" file
	
	Decrypts and overides file (user specified file) with passphrase "mypasshphrase",
	this option will disable stdin interaction.

lyra decrypt --print-only file

	Decrypts content of file and output the resulting plaintext to stdout. The original
	specified file will be left encrypted and passphrase will be asked via stdin.
	
lyra decrypt -p "mypassphrase" -s file1 file

	Decrypts file (user specified file) with passphrase "mypasshphrase" and save the resulting
	output to file1. The original specified file will remain encrypted with the same key. 
	Stdin interaction will also be disabled.

lyra decrypt --print-only -p "mypassphrase" file

	Decrypts content of file with passphrase "mypassphrase" and output the resulting 
	plaintext to stdout. The original specified file will remain encrypted with the 
	same key. 

`

	usagePrint = `Prints the deciphered contents of a specified file to stdout, the original 
file will be unchanged (i.e still encrypted with the same key).
`

	usagePathDec = `Decrypts the contents of file and save the resulting plaintext in a new file. 
The original specified file will be unchanged (i.e still encrypted with the 
same key) if this flag is set.
`
)

type decryptcmd struct {
	path       string
	passphrase string
	printOnly  bool
}

func (cmd *decryptcmd) CName() string {
	return "decrypt"
}

func (cmd *decryptcmd) Help() string {
	return helpstrdec
}

func (cmd *decryptcmd) RegCFlags(fs *flag.FlagSet) {
	fs.StringVar(&cmd.passphrase, "p", "", usagePass)
	fs.StringVar(&cmd.path, "s", "", usagePathDec)
	fs.BoolVar(&cmd.printOnly, "print-only", false, usagePrint)
}

func (cmd *decryptcmd) Run(opt []string) error {
	switch len(opt) {
	case 0:
		return errors.New("You must specify a target file")
	}

	err := cmd.validateInputs()
	if err != nil {
		return err
	}

	if cmd.passphrase == "" {
		cmd.passphrase = string(getPassphrase())
	}

	err = encryption.Decrypt(opt[0], cmd.path, cmd.printOnly, []byte(cmd.passphrase))
	cmd.passphrase = ""

	return err
}

func (cmd *decryptcmd) validateInputs() error {
	if cmd.printOnly && cmd.path != "" {
		return errors.New("Invalid input, -s and --print-only can't be set at the same time")
	}

	return nil
}
