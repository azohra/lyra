package main

import (
	"errors"
	"flag"
	"fmt"
	"os"

	"github.com/azohra/lyra/internal/pkg/encryption"
	"github.com/brsmsn/gware/pkg/diceware"
)

const (
	helpstrenc = `
The following exmaples are all the possible options for the "encrypt" command:
	
lyra encrypt file
			
	Encrypts and overides file (user specified file). Users will be asked to provide
	the passphrase via stdin.
	
lyra encrypt -s file1 file
	
	Encrypts the contents of file (user specified file) and save the resulting output 
	to file1 (user specified file or specified path). Users will be asked to provide 
	the passphrase via stdin. Original specified file will still remain in plaintext.
	
lyra encrypt --auto-gen file

	Encrypts and overides file (user specified file) with an auto generated passphrase
	and outputs the auto generated passphrase to stdout.
	The auto generated passphrase is a 7 word passphrase generated via the diceware
	method using the EFF new wordlist. It is imperative that the user keep a record
	of the outputted passphrase as there will be no way to decipher the file without it.
	
lyra encrypt --gen-str file
	
	Encrypts and overides file (user specified file) with an auto generated passphrase
	and outputs the auto generated passphrase to stdout.
	Auto generates a 7 word passphrase in kebab case (no spaces).

lyra encrypt -p "mypassphrase" file

	Encrypts and overides file (user specified file) with passphrase "mypassphrase",
	this option will disable stdin interaction. Using this option will also disable
	passphrase checking, therefore it is critical that you do not misspell or forget
	the passphrase. 
	
lyra encrypt --auto-gen -s file1 file
	
	Encrypts file (user specified file) with an auto generated dicewre passphrase and save
	the resulting output to file1 (user specified file). The auto generated passphrase will
	be outputted to stdout and the original specified file will still remain in plaintext
	/decrypted.
	
lyra encrypt --gen-str -s file1 file
	
	Encrypts file (user specified file) with an auto generated a 7 word passphrase without
	spaces and save the resulting output to file1 (user specified file). The auto generated 
	passphrase will be outputted to stdout and the original specified file will still remain 
	in plaintext /decrypted.
	
lyra encrypt -p "mypassphrase" -s file1 file
	
	Encrypts file (user specified file) with passphrase "mypasshphrase" and save the resulting
	output to file1. The original specified file will remain in plaintext/decrypted.
	Stdin interaction will also be disabled.
	
`

	usagePass = `Specify a passphrase used to encrypt/decrypt the specified file, if this flag 
is set, passphrases will not fetched from stdin. 
	
For encryption this flag will disable passphrase verification. Be careful not 
to misspell your passphrase as there will be no way to decrypt your files!
`

	usagePathEnc = `Encrypts the contents of file and save the resulting ciphertext in a new file. 
The original specified file will be unchanged (i.e still decrypted) if this flag is set.
`

	usageGenDice = `Auto generates a single 7 word passphrase that will be used as the key for the
encryption of a specified file. The passphrase is a diceware generated passphrase using
the EFF new wordlist.
`

	usageGenStr = `Auto generate a single 7 word diceware passphrase as a single no-spaced string.
`
)

type encryptcmd struct {
	path        string
	passphrase  string
	autogenDice bool
	autogenStr  bool
}

func (cmd *encryptcmd) CName() string {
	return "encrypt"
}

func (cmd *encryptcmd) Help() string {
	return helpstrenc
}

func (cmd *encryptcmd) RegCFlags(fs *flag.FlagSet) {
	fs.StringVar(&cmd.passphrase, "p", "", usagePass)
	fs.StringVar(&cmd.path, "s", "", usagePathEnc)
	fs.BoolVar(&cmd.autogenDice, "auto-gen", false, usageGenDice)
	fs.BoolVar(&cmd.autogenStr, "gen-str", false, usageGenStr)
}

func (cmd *encryptcmd) Run(opt []string) error {
	switch len(opt) {
	case 0:
		return errors.New("You must specify a target file")
	}

	//check for valid inputs
	err := cmd.validateInputs()
	if err != nil {
		return err
	}

	//if autogenDice was set
	if cmd.autogenDice || cmd.autogenStr {
		cmd.genPass()
		if err != nil {
			return err
		}
		fmt.Fprint(os.Stdout, cmd.passphrase+"\n")
	} else if cmd.passphrase == "" {
		input, err := setPassphrase()
		if err != nil {
			return err
		}
		cmd.passphrase = string(input)
	}

	err = encryption.Encrypt(opt[0], cmd.path, []byte(cmd.passphrase))
	if err != nil {
		return err
	}

	//wiping passphrase before exit
	cmd.passphrase = ""

	return nil
}

func (cmd *encryptcmd) validateInputs() error {
	if (cmd.autogenDice || cmd.autogenStr) && cmd.passphrase != "" {
		return errors.New("Can not specify a passphrase when auto-gen flag has been set")
	} else if cmd.autogenDice && cmd.autogenStr {
		return errors.New("Can not specify --auto-gen and --gen-str at the same time")
	}

	return nil
}

//gen a diceware passphrase using eff long wordlist
func (cmd *encryptcmd) genPass() error {

	phrase, err := diceware.GeneratePassphrases(1, 7, diceware.EffWorldList)
	if err != nil {
		return err
	}

	if cmd.autogenStr {
		cmd.passphrase = removeSpaces(phrase[0])
	} else {
		cmd.passphrase = phrase[0]
	}

	return nil
}
