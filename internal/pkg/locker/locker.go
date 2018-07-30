package locker

import (
	"bufio"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/azohra/lyra/internal/pkg/encryption"
)

const (
	lockedFileExtension      = ".locked" // What the extension is for locked files
	commentSequence          = "#"       // What signals a comment in a locker file
	encryptionValidatorRegex = `^\@\![a-f0-9]{32}\@\![a-f0-9]{24}$`
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

			filepath.Walk(entry, func(path string, info os.FileInfo, err error) error {
				if info != nil && !info.IsDir() {
					// newLockerAsset will propagate errors to the top cmd level so we dont need to
					// handle them here
					jobs = append(jobs, NewLockerAsset(path))
				}
				return nil
			})

		} else {
			// Creating a locker asset from a missing file is ok. Error will
			// propagate to the cmd level for reporting
			jobs = append(jobs, NewLockerAsset(entry))
		}

	}
	return jobs, nil
}

// NewLockerAsset creates a new asset record from a filename.
// This includes determining if the file is locked and/or
// even exists
func NewLockerAsset(filename string) (l Asset) {
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
func (a *Asset) Lock(passphrase []byte) error {
	if !a.IsLocked {
		err := encryption.Encrypt(a.Filename, a.LockedFilename, passphrase)
		if err != nil {
			return err
		}
		err = os.Remove(a.Filename)
		if err == nil {
			a.IsLocked = true
		}
		return err
	}
	return nil
}

// Unlock decrypts a locker Asset.
func (a *Asset) Unlock(passphrase []byte) error {
	if a.IsLocked {
		err := encryption.Decrypt(a.LockedFilename, a.Filename, false, passphrase)
		if err != nil {
			return err
		}
		err = os.Remove(a.LockedFilename)
		if err == nil {
			a.IsLocked = false
		}
		return err
	}
	return nil
}

// ValidateLocked validates an asset is locked
// and the file is encrypted - not just
// the extension changed.
func (a *Asset) ValidateLocked() (bool, error) {
	validationRegex := regexp.MustCompile(encryptionValidatorRegex)
	if !a.IsLocked {
		return false, nil
	} else {

		// Make sure the file does not exist in plain text first
		_, err := os.Stat(a.Filename)
		if err == nil || os.IsExist(err) {
			return false, nil
		}

		// Check the headder of the locked file
		file, err := os.Open(a.LockedFilename)
		defer file.Close()
		if err != nil {
			return false, err
		}

		scanner := bufio.NewScanner(file)
		if scanner.Scan() {
			firstLine := scanner.Text()
			return validationRegex.MatchString(firstLine), nil
		}
	}
	return false, nil
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
