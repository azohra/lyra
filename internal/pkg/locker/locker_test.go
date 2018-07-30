package locker

import (
	"io/ioutil"
	"os"
	"testing"
)

func TestParseLockerFile(t *testing.T) {
	assets, err := ParseLockerFile("../../../lyralocker")
	if err != nil {
		t.Error(err)
		return
	}
	if len(assets) != 1 {
		t.Errorf("Expecting 1 entry in test lyralocker file")
		return
	}
}

func TestValidateLockAndUnlock(t *testing.T) {
	testFile := "../../../test/locker/lockme.txt"
	fileText := "This data is to be locked and unlocked by tests\n"
	passphrase := "forthegoodoftesting"

	// clean anything that may be remaining from previous tests
	os.Remove(testFile)
	os.Remove(testFile + ".locked")

	// overwrite the file if already exists, and schedule it to be cleaned up
	ioutil.WriteFile(testFile, []byte(fileText), 0666)
	defer os.Remove(testFile)
	defer os.Remove(testFile + ".locked") // just incase a premature exit

	asset := NewLockerAsset(testFile)
	asset.Lock([]byte(passphrase))

	isLocked, err := asset.ValidateLocked()
	if err != nil || !isLocked {
		t.Errorf("Failed to lock file. Got: (%v, %v)", isLocked, err)
	}

	contents, err := ioutil.ReadFile(testFile + ".locked")
	if err != nil {
		t.Error(err)
	}
	if string(contents) == fileText {
		t.Errorf("File was not encrypted.")
		return
	}

	asset.Unlock([]byte(passphrase))

	contents, err = ioutil.ReadFile(testFile)
	if err != nil {
		t.Error(err)
	}
	if string(contents) != fileText {
		t.Errorf("File was not decrypted properly.")
	}

}

func TestValidateLocked(t *testing.T) {
	asset1 := NewLockerAsset("../../../test/locker/test-file1.txt")    // just the txt
	asset2 := NewLockerAsset("../../../test/locker/test-file2.txt")    // normal txt and locked files
	asset3 := NewLockerAsset("../../../test/locker/test-file3.txt")    // plain txt marked as locked but not encrypted
	asset4 := NewLockerAsset("../../../test/locker/test-file4.txt")    // properly locked
	asset5 := NewLockerAsset("../../../test/locker/test-file-404.txt") // does not exist

	isLocked, err := asset1.ValidateLocked()
	if err != nil || isLocked {
		t.Errorf("%s is plan text and should not be validated as locked. Got: (%v, %v)", asset1.Filename, isLocked, err)
	}

	isLocked, err = asset2.ValidateLocked()
	if err != nil || isLocked {
		t.Errorf("%s is in both plain text and encrypted. It should not be validated as locked. Got: (%v, %v)", asset2.Filename, isLocked, err)
	}

	isLocked, err = asset3.ValidateLocked()
	if err != nil || isLocked {
		t.Errorf("%s has the extension .locked but is still in plan text. It should not be validated as locked. Got: (%v, %v)", asset3.Filename, isLocked, err)
	}

	isLocked, err = asset4.ValidateLocked()
	if err != nil || !isLocked {
		t.Errorf("%s should be validated as locked. Got: (%v, %v)", asset4.Filename, isLocked, err)
	}

	isLocked, err = asset5.ValidateLocked()
	if err == nil {
		t.Errorf("%s does not exist and should error out. Got: (%v, %v)", asset5.Filename, isLocked, err)
	}
}

func TestIsFilenameLocked(t *testing.T) {
	if isFilenameLocked("somefile.txt") {
		t.Errorf("'somefile.txt' is not a locked name")
	}

	if !isFilenameLocked("someOtherFile.exs.locked") {
		t.Errorf("'someOtherFile.exs.locked' is a locked name")
	}
}

func TestCreateLockedFilename(t *testing.T) {
	test1 := [...]string{"someFile.txt", "someFile.txt.locked"}
	if createLockedFilename(test1[0], true) != test1[1] {
		t.Errorf("'%s' locked filename not correct! Got: %s", test1[0], test1[1])
	}

	test2 := [...]string{"someFile.txt.locked", "someFile.txt"}
	if createLockedFilename(test2[0], false) != test2[1] {
		t.Errorf("'%s' unlocked filename not correct! Got: %s", test2[0], test2[1])
	}

	test3 := [...]string{"someFile.txt", "someFile.txt"}
	if createLockedFilename(test3[0], false) != test3[1] {
		t.Errorf("'%s' unlocked filename not correct! Got: %s", test3[0], test3[1])
	}

	test4 := [...]string{"someFile.txt.locked", "someFile.txt.locked"}
	if createLockedFilename(test4[0], true) != test4[1] {
		t.Errorf("'%s' locked filename not correct! Got: %s", test4[0], test4[1])
	}
}
