package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/azohra/lyra/cmd/lyra/locker"
)

const (
	lockerUsage = `Locker is a tool for project secret management using Lyra's encryption.unlock

	$ lyra locker [Flags] <Command>

Commands:

	lock	Lock all assets listed in lyralocker file
	unlock	Unlock all assets listed in lyralocker file
	shake	Shakes the locker to test for unencrypted files

Flags:

	-q		Quiet mode. Only prints necessary info to the screen
	-p		Enter password
`

	lockerConfig       = "./lyralocker"
	lockerpassFilename = "./.lockerpass"
	checkmark          = "✓"
	cross              = "✕"
)

type lockercmd struct {
	fileRecursionDepth int
	passphrase         string
	quiet              bool
}

func (cmd *lockercmd) CName() string {
	return "locker"
}

func (cmd *lockercmd) Help() string {
	return lockerUsage
}

func (cmd *lockercmd) RegCFlags(fs *flag.FlagSet) {
	fs.BoolVar(&cmd.quiet, "q", false, "Run command in quiet mode")
	fs.StringVar(&cmd.passphrase, "p", "", "Use provided passphrase")
}

func (cmd *lockercmd) Run(opt []string) error {

	if len(opt) < 1 {
		fmt.Println(lockerUsage)
		os.Exit(0)
	}

	// If no passphrase provided to command, pull
	// from project password file
	if cmd.passphrase == "" {
		pass, err := readPassFile()
		if err != nil {
			return err
		}
		cmd.passphrase = pass
	}

	files, err := locker.ParseLockerFile("./lyralocker")
	if err != nil {
		return err
	}

	response := ""
	code := 0

	switch opt[0] {
	case "lock":
		response, code = lock(files, cmd)
		break
	case "unlock":
		response, code = unlock(files, cmd)
		break
	case "shake":
		response, code = shake(files, cmd)
		break
	default:
		fmt.Println(lockerUsage)
		os.Exit(1)
	}

	outputTo := os.Stdout
	if code > 0 {
		outputTo = os.Stderr
	}
	if !cmd.quiet {
		fmt.Fprintf(outputTo, response)
	}
	os.Exit(code)
	return nil // will never be reached since we handle things locally
}

func lock(files []locker.Asset, cmd *lockercmd) (string, int) {
	successCount := 0
	failCount := 0
	for _, f := range files {
		if f.Err != nil {
			if os.IsNotExist(f.Err) {
				fmt.Fprint(os.Stderr, fmt.Sprintf("File %s does not exist\n", f.Filename))
			}
			failCount++
		} else {
			err := f.Lock([]byte(cmd.passphrase))
			if err != nil {
				fmt.Fprint(os.Stderr, err.Error())
				failCount++
			} else {
				if !cmd.quiet {
					fmt.Printf("%s  %s\n", checkmark, f.Filename)
				}
				successCount++
			}
		}
	}

	reply := fmt.Sprintf("%d assets encrypted\n", successCount)
	replyCode := 0

	if failCount > 0 {
		replyCode = 1
		reply = fmt.Sprintf("%d could not be encrypted\n", failCount)
	}

	return reply, replyCode
}

func unlock(files []locker.Asset, cmd *lockercmd) (string, int) {
	successCount := 0
	failCount := 0
	for _, f := range files {
		if f.Err != nil {
			if os.IsNotExist(f.Err) {
				fmt.Fprint(os.Stderr, fmt.Sprintf("File %s does not exist\n", f.Filename))
			}
			failCount++
		} else {
			err := f.Unlock([]byte(cmd.passphrase))

			if err != nil {
				fmt.Fprint(os.Stderr, err.Error())
				failCount++
			} else {
				if !cmd.quiet {
					fmt.Printf("%s  %s\n", checkmark, f.Filename)
				}
				successCount++
			}
		}
	}

	reply := fmt.Sprintf("%d assets decrypted\n", successCount)
	replyCode := 0

	if failCount > 0 {
		replyCode = 1
		reply = fmt.Sprintf("%d could not be decrypted\n", failCount)
	}

	return reply, replyCode
}

func shake(files []locker.Asset, cmd *lockercmd) (string, int) {
	successCount := 0
	failCount := 0
	for _, f := range files {
		if f.Err != nil {
			if os.IsNotExist(f.Err) {
				fmt.Fprint(os.Stderr, fmt.Sprintf("File %s does not exist\n", f.Filename))
			}
			failCount++
		} else {
			validEncryption, err := f.ValidateLocked()

			if err != nil {
				fmt.Fprint(os.Stderr, err.Error())
				failCount++
			} else if validEncryption {
				if !cmd.quiet {
					fmt.Printf("%s  %s\n", checkmark, f.Filename)
				}
				successCount++
			} else {
				fmt.Fprintf(os.Stderr, "%s  %s\n", cross, f.Filename)
				failCount++
			}
		}
	}

	reply := fmt.Sprintf("%d assets secure\n", successCount)
	replyCode := 0

	if failCount > 0 {
		replyCode = 1
		reply = fmt.Sprintf("%d assets are insecure\n", failCount)
	}

	return reply, replyCode
}

func readPassFile() (string, error) {
	contents, err := ioutil.ReadFile(lockerpassFilename)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(contents)), nil
}
