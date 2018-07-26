package main

import (
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/fvumbaca/lyra/cmd/lyra/locker"
)

const lockerConfig = "./lyralocker"
const lockerpassFilename = "./.lockerpass"
const checkmark = "✓"
const chironKey = string(0x26B7)

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

	successCount := 0
	failCount := 0

	for _, f := range files {
		if f.Err != nil {
			if os.IsNotExist(f.Err) {
				fmt.Fprint(os.Stderr, fmt.Sprintf("File %s does not exist\n", f.Filename))
			}
			failCount++
		} else {
			if cmd.doEncrypt {
				err := f.Lock([]byte(cmd.passphrase))
				if err != nil {
					fmt.Fprint(os.Stderr, err.Error())
					failCount++
				} else {
					fmt.Printf("%s  %s\n", checkmark, f.Filename)
					successCount++
				}
			} else {
				err := f.Unlock([]byte(cmd.passphrase))
				if err != nil {
					fmt.Fprint(os.Stderr, err.Error())
					failCount++
				} else {
					fmt.Printf("%s  %s\n", checkmark, f.Filename)
					successCount++
				}
			}
		}

	}

	action := "encripted"
	if !cmd.doEncrypt {
		action = "decrypted"
	}

	if failCount > 0 {
		fmt.Fprintf(os.Stderr, "%d assets were unable to be encrypted", failCount)
		os.Exit(1)
	} else {
		fmt.Printf("%d assets %s\n", successCount, action)
		return nil
	}
	return nil
}

func readPassFile() (string, error) {
	contents, err := ioutil.ReadFile(lockerpassFilename)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(contents)), nil
}
