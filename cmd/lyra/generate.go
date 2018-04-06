package main

import (
	"errors"
	"flag"
	"fmt"
	"os"

	"github.com/brsmsn/gware/pkg/diceware"
)

const helpstrgen = `
The following exmaples are all the possible options for the "generate" command:

lyra generate --words 7 --phrases 6

	Generate 6 diceware passphrases each containing 7 words.

`

const usageWords = `Specify the number of words that a passphrase will have.
`

const usagePhrases = `Specify the number of passphrases that will be generated.
`

type gencmd struct {
	numWords   int
	numPhrases int
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
