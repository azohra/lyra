package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/brsmsn/gware/pkg/diceware"
)

const (
	helpstrgen = `
The following exmaples are all the possible options for the "generate" command:

lyra generate --words 7 --phrases 6

	Generate 6 diceware passphrases each containing 7 words.

lyra generate --rm-spaces --words 7 --phrases 6

	Generate 6 diceware passphrases each containing 7 words with no spaces

lyra generate --rm-spaces --words 7 

	Generate 1 diceware passphrase containing 7 words with no spaces

lyra generate --rm-spaces --phrases 7 

	Generate 7 diceware passphrase containing 7 words with no spaces

`

	usageWords = `Specify the number of words that a passphrase will have.
`

	usagePhrases = `Specify the number of passphrases that will be generated.
`

	usageRmSpaces = `Specify removal of spaces, this will replace all spaces
with a hyphen as a delimiter.
`
)

type gencmd struct {
	numWords   int
	numPhrases int
	noSpaces   bool
}

func (cmd *gencmd) CName() string {
	return "generate"
}

func (cmd *gencmd) Help() string {
	return helpstrgen
}

func (cmd *gencmd) RegCFlags(fs *flag.FlagSet) {
	fs.IntVar(&cmd.numWords, "words", 7, usageWords)
	fs.IntVar(&cmd.numPhrases, "phrases", 1, usagePhrases)
	fs.BoolVar(&cmd.noSpaces, "rm-spaces", false, usageRmSpaces)
}

func (cmd *gencmd) Run(opt []string) error {
	if len(opt) > 0 {
		return errors.New("Invalid input\n" + helpstrgen)
	}

	err := cmd.validateInputs()
	if err != nil {
		return err
	}

	list, err := diceware.GeneratePassphrases(cmd.numPhrases, cmd.numWords, diceware.EffWorldList)
	if err != nil {
		return err
	}

	for _, val := range list {
		if cmd.noSpaces {
			val = removeSpaces(val)
		}
		fmt.Fprint(os.Stdout, val+"\n")
	}

	return nil
}

func (cmd *gencmd) validateInputs() error {
	if cmd.numWords <= 0 || cmd.numPhrases <= 0 {
		return errors.New("Invalid range number can not be less than 0")
	} else if cmd.numWords < 5 {
		return errors.New("Don't generate any passphrase with less than 5 words, it is insecure")
	}

	return nil
}

//removes spaces from a phrase
func removeSpaces(phrase string) string {
	return strings.Replace(phrase, " ", "-", -1)
}
