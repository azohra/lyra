package locker

import "testing"

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
