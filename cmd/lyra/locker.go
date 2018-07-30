package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/azohra/lyra/internal/pkg/locker"
	"github.com/brsmsn/gware/pkg/diceware"
)

const (
	lockerUsage = `
	
Sub commands:

	lock	Lock all assets listed in lyralocker file
	unlock	Unlock all assets listed in lyralocker file
	check	Checks the locker to test for any unencrypted files

`

	// Used to make cool reports
	checkmark = "✓"
	cross     = "✕"

	// LockerConfigFilename defines the filename of the projects lyra locker config file
	LockerConfigFilename = "lyralocker"

	// LockerPassphraseFilename defines the filename of the file that stores the projects passphrase
	LockerPassphraseFilename = ".lockerpass"

	// Might need to add flags to configure these in the future
	lockerPassGenNumWords     = 7
	lockerPassGenNumPhrases   = 1
	lockerPassGenRemoveSpaces = false

	lockerFileTemplate = `# This file is for listing assets to be encrypted

# Uncomment the following line to lock the file
# super-secret-key.txt
`
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

	response := ""
	code := 0

	// Init is a special command that does not error
	// when the lyrafile does not exist
	if opt[0] == "init" {
		response, code = initializeLocker()
	} else {

		// If no passphrase provided to command, pull
		// from project password file
		if cmd.passphrase == "" {
			pass, err := readPassFile()
			if err != nil {
				return err
			}
			cmd.passphrase = pass
		}

		files, err := locker.ParseLockerFile(LockerConfigFilename)
		if err != nil {
			return err
		}

		switch opt[0] {
		case "lock":
			response, code = lock(files, cmd)
			break
		case "unlock":
			response, code = unlock(files, cmd)
			break
		case "check":
			response, code = check(files, cmd)
			break
		default:
			fmt.Println(lockerUsage)
			os.Exit(1)
		}

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

// Initializes the cwd with a password file and a config file
func initializeLocker() (string, int) {
	createdPassFile := true
	passFileExists, err := isFileExists(LockerPassphraseFilename)
	if err != nil {
		return err.Error(), BadExit
	}

	if !passFileExists {

		passFile, err := os.Create(LockerPassphraseFilename)
		defer passFile.Close()
		if err != nil {
			return err.Error(), BadExit
		} else {
			words, err := diceware.GeneratePassphrases(lockerPassGenNumPhrases, lockerPassGenNumWords, diceware.EffWorldList)
			if err != nil {
				return "Could not generate password", BadExit
			}

			passphrase := strings.Join(words, " ")
			if lockerPassGenRemoveSpaces {
				passphrase = strings.Join(words, "")
			}

			// Keep the formatting clean. Will be stripped out when read
			passphrase += "\n"

			_, err = fmt.Fprint(passFile, passphrase)
			if err != nil {
				return "Could not write password to " + LockerPassphraseFilename, BadExit
			}

		}

	}

	configFileCreated := false

	configFileExists, err := isFileExists(LockerConfigFilename)
	if err != nil {
		return err.Error(), BadExit
	}

	if !configFileExists {

		configFile, err := os.Create(LockerConfigFilename)
		defer configFile.Close()
		if err != nil {
			if !os.IsExist(err) {
				configFileCreated = false
			}
		} else {

			_, err := fmt.Fprintln(configFile, lockerFileTemplate)
			if err != nil {
				configFileCreated = true
			}

		}
	}

	reply := ""

	if createdPassFile {
		reply += fmt.Sprintf("Generated passphrase file at %s\n", LockerPassphraseFilename)
	} else {
		reply += fmt.Sprintf("Passphrase file already exists\n")
	}

	if configFileCreated {
		reply += fmt.Sprintf("Generated %s file\n", LockerConfigFilename)
	} else {
		reply += fmt.Sprintf("%s file already exists\n", LockerConfigFilename)
	}

	return reply, GoodExit
}

// Locks all provided assets
func lock(assets []locker.Asset, cmd *lockercmd) (string, int) {
	successCount := 0
	failCount := 0
	for _, asset := range assets {
		if asset.Err != nil {
			if os.IsNotExist(asset.Err) {
				fmt.Fprint(os.Stderr, fmt.Sprintf("File %s does not exist\n", asset.Filename))
			}
			failCount++
		} else {
			err := asset.Lock([]byte(cmd.passphrase))
			if err != nil {
				fmt.Fprint(os.Stderr, err.Error())
				failCount++
			} else {
				if !cmd.quiet {
					fmt.Printf("Locked  %s\n", asset.Filename)
				}
				successCount++
			}
		}
	}

	reply := fmt.Sprintf("%d assets encrypted\n", successCount)
	replyCode := GoodExit

	if failCount > 0 {
		replyCode = BadExit
		reply = fmt.Sprintf("%d could not be encrypted\n", failCount)
	}

	return reply, replyCode
}

// Unlocks locked assets provided
func unlock(assets []locker.Asset, cmd *lockercmd) (string, int) {
	successCount := 0
	failCount := 0
	for _, asset := range assets {
		if asset.Err != nil {
			if os.IsNotExist(asset.Err) {
				fmt.Fprint(os.Stderr, fmt.Sprintf("File %s does not exist\n", asset.Filename))
			}
			failCount++
		} else {
			err := asset.Unlock([]byte(cmd.passphrase))

			if err != nil {
				fmt.Fprint(os.Stderr, err.Error())
				failCount++
			} else {
				if !cmd.quiet {
					fmt.Printf("Unlocked  %s\n", asset.Filename)
				}
				successCount++
			}
		}
	}

	reply := fmt.Sprintf("%d assets decrypted\n", successCount)
	replyCode := GoodExit

	if failCount > 0 {
		replyCode = BadExit
		reply = fmt.Sprintf("%d could not be decrypted\n", failCount)
	}

	return reply, replyCode
}

// Checks the current project for any assets that are in an
// unencrypted state
func check(files []locker.Asset, cmd *lockercmd) (string, int) {
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

	reply := fmt.Sprintf("%d/%d assets secured.\n", successCount, len(files))
	replyCode := GoodExit

	if failCount > 0 {
		replyCode = BadExit
		reply = fmt.Sprintf("%d/%d assets secured. %d files could not be encrypted.\n", successCount, len(files), failCount)
	}

	return reply, replyCode
}

// Reads a password file
func readPassFile() (string, error) {
	contents, err := ioutil.ReadFile(LockerPassphraseFilename)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(contents)), nil
}

// Checks if a file exists and not a directory
func isFileExists(filename string) (bool, error) {
	fileInfo, err := os.Stat(filename)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		if os.IsExist(err) {
			return true, nil
		}
		return false, err
	}
	return !fileInfo.IsDir(), nil
}
