package locker

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/fvumbaca/lyra/pkg/lcrypt"
)

const (
	lockedFileExtension = ".locked" // What the extension is for locked files
	commentSequence     = "#"       // What signals a comment in a locker file
)

// Asset represents a Locker file asset.
type Asset struct {
	Filename       string
	LockedFilename string
	IsLocked       bool
	Err            error
}

// ParseLockerFile parses a locker file and returns a list of locker assets
func ParseLockerFile(filename string) ([]Asset, error) {
	jobs := []Asset{}

	lockerFile, err := os.Open(filename)
	defer lockerFile.Close()
	if err != nil {
		return nil, err
	}
	scanner := bufio.NewScanner(lockerFile)
	for scanner.Scan() {
		entry := strings.TrimSpace(scanner.Text())

		// ignore empty lines and comments from locker file
		if entry == "" || strings.HasPrefix(entry, commentSequence) {
			continue
		}

		fileInfo, err := os.Stat(entry)
		// cant be a dir if it is not found but maybe .locked?
		if err == nil && fileInfo.IsDir() {
			fmt.Println("Not handling folders yet")

			filepath.Walk(entry, func(path string, info os.FileInfo, err error) error {
				// newLockerAsset will propagate errors to the top cmd level so we dont need to
				// handle them here
				jobs = append(jobs, newLockerAsset(path))
				return nil
			})

		} else {
			// Creating a locker asset from a missing file is ok. Error will
			// propagate to the cmd level for reporting
			jobs = append(jobs, newLockerAsset(entry))
		}

	}
	return jobs, nil
}

func newLockerAsset(filename string) (l Asset) {
	l.Filename = filename
	l.LockedFilename = filename
	l.IsLocked = true
	l.Err = nil

	if isFilenameLocked(filename) {
		l.Filename = createLockedFilename(filename, false)
	} else {
		l.LockedFilename = createLockedFilename(filename, true)
	}

	_, baseErr := os.Stat(l.Filename)
	_, lockedErr := os.Stat(l.LockedFilename)

	if baseErr == nil && lockedErr != nil && os.IsNotExist(lockedErr) {
		l.IsLocked = false
	} else if baseErr != nil && os.IsNotExist(baseErr) && lockedErr == nil {
		l.IsLocked = true
	} else {
		l.Err = baseErr
	}

	return
}

// Lock encrypts a locker Asset.
func (a Asset) Lock(passphrase []byte) error {
	if !a.IsLocked {
		err := lcrypt.Encrypt(a.Filename, a.LockedFilename, passphrase)
		if err != nil {
			return err
		}
		return os.Remove(a.Filename)
	}
	return nil
}

// Unlock decrypts a locker Asset.
func (a Asset) Unlock(passphrase []byte) error {
	if a.IsLocked {
		err := lcrypt.Decrypt(a.LockedFilename, a.Filename, false, passphrase)
		if err != nil {
			return err
		}
		return os.Remove(a.LockedFilename)
	}
	return nil
}

func createLockedFilename(filename string, lock bool) string {
	if lock && !isFilenameLocked(filename) {
		return filename + lockedFileExtension
	} else if !lock && isFilenameLocked(filename) {
		return strings.TrimSuffix(filename, lockedFileExtension)
	} else {
		return filename
	}
}

func isFilenameLocked(filename string) bool {
	return strings.HasSuffix(filename, lockedFileExtension)
}
