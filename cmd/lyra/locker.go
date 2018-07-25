package main

import (
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/fvumbaca/lyra/cmd/lyra/locker"
)

const lockerConfig = "./lyralocker"
const lockerpassFilename = "./.lockerpass"

type lockercmd struct {
	doEncrypt          bool
	fileRecursionDepth int
	passphrase         string
}

type lockerfile struct {
	filename       string
	lockedFilename string
	isLocked       bool
}

func (cmd *lockercmd) CName() string {
	return "locker"
}

func (cmd *lockercmd) Help() string {
	return "Help TBT"
}

func (cmd *lockercmd) RegCFlags(fs *flag.FlagSet) {
}

func (cmd *lockercmd) Run(opt []string) error {

	if len(opt) < 1 {
		return errors.New("Bad args...")
	}

	switch opt[0] {
	case "lock":
		cmd.doEncrypt = true
		break
	case "unlock":
		cmd.doEncrypt = false
		break
	default:
		return errors.New("Bad args...")
	}

	pass, err := readPassFile()
	if err != nil {
		return err
	}
	cmd.passphrase = pass

	files, err := locker.ParseLockerFile("./lyralocker")
	if err != nil {
		return err
	}

	for _, f := range files {
		fmt.Printf("Parsed file: %+v\n", f)
		if f.Err != nil {
			report(f.Err)
		} else {
			if cmd.doEncrypt {
				f.Lock([]byte(cmd.passphrase))
			} else {
				f.Unlock([]byte(cmd.passphrase))
			}
		}

	}

	return nil
}

func report(err error) {
	fmt.Println("An error occurred: " + err.Error())
}

func readPassFile() (string, error) {
	contents, err := ioutil.ReadFile(lockerpassFilename)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(contents)), nil
}
