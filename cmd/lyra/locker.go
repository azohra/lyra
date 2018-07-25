package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

const lockerConfig = "./lyralocker"
const lockerpassFilename = "./.lockerpass"

type lockercmd struct {
	doEncrypt          bool
	fileRecursionDepth int
	passphrase         string
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
	fmt.Printf("Read password: %s\n", pass)
	cmd.passphrase = pass

	// TODO support multiple possible names like Docker...
	lockerFilename, err := filepath.Abs(lockerConfig)
	if err != nil {
		return err
	}

	lockerFile, err := os.Open(lockerFilename)
	defer lockerFile.Close()
	if err != nil {
		return err
	}

	lockerScanner := bufio.NewScanner(lockerFile)

	for lockerScanner.Scan() {
		filename, err := filepath.Abs(lockerScanner.Text())
		if err != nil {
			report(err)
			continue
		}
		fmt.Println("Read: " + filename)

		fileInfo, err := os.Stat(filename)
		if err != nil {
			fmt.Printf("Error because exists? %t\n", os.IsExist(err))
			report(err)
			continue
		}

		if fileInfo.IsDir() {
			processLockerFolder(filename, cmd)
		} else {
			err := processLockerFile(filename, cmd)
			if err != nil {
				report(err)
				continue
			}

		}
	}

	return nil
}

func report(err error) {
	// TODO report error
}

func processLockerFolder(foldername string, cmd *lockercmd) {
	filepath.Walk(foldername, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			report(err)
			return err
		}

		if !info.IsDir() {
			return processLockerFile(path, cmd)
		}
		return nil
	})
}

func processLockerFile(filename string, cmd *lockercmd) error {
	fmt.Println("Processing: " + filename)

	if cmd.doEncrypt {
		fmt.Println("Trying to encrypt: " + filename)
		err := encrypt(filename, createLockedFilename(filename), []byte(cmd.passphrase))
		if err != nil {
			return err
		}
		return os.Remove(filename)
	} else {
		fmt.Println("Trying to decrypt: " + filename)
		lockedFilename := createLockedFilename(filename)

		err := decrypt(lockedFilename, filename, false, []byte(cmd.passphrase))
		if err != nil {
			return err
		}
		return os.Remove(lockedFilename)
	}

	// if cmd.doEncrypt != isLocked(filename) {
	// 	if cmd.doEncrypt {
	// 		err := encrypt(filename, createFilenameLock(filename, cmd.doEncrypt), []byte(cmd.passphrase))
	// 		if err != nil {
	// 			return err
	// 		}
	// 		return os.Remove(filename)
	// 	} else {
	// 		err := decrypt(filename, createFilenameLock(filename, cmd.doEncrypt), false, []byte(cmd.passphrase))
	// 		if err != nil {
	// 			return err
	// 		}
	// 		return os.Remove(filename)
	// 	}
	// }

}

func createLockedFilename(filename string) string {
	return filename + ".locked"
}

func readPassFile() (string, error) {
	contents, err := ioutil.ReadFile(lockerpassFilename)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(contents)), nil
}

func isFileLocked(filename string) (bool, error) {
	lockedFilename := createLockedFilename(filename)

	_, infoErr := os.Stat(filename)
	_, lockedInfoErr := os.Stat(lockedFilename)

	if infoErr != nil && lockedInfoErr != nil {
		return false, infoErr
	}

	return os.IsExist(infoErr), nil

}
