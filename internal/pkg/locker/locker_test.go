package locker

import (
	"testing"
)

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
