package lfile

import (
	"bytes"
	"io"
	"reflect"
	"testing"

	"github.com/awnumar/memguard"

	"github.com/azohra/lyra/pkg/lcrypt"
)

func TestNewLyraFile(t *testing.T) {
	stc := newLyraFile()
	if stc == nil {
		t.Error("Did not initialize lyraFile")
	}
}

func TestNewParsedLyraFile(t *testing.T) {
	fixture := "../../test/fixture.1.txt"

	stc, err := NewParsedLyraFile(fixture)
	if err != nil {
		t.Error(err)
	}

	ctx := newLyraFile()
	err = ctx.ParseFile(fixture)
	if err != nil {
		t.Error(err)
	}

	if !reflect.DeepEqual(ctx.plaintext.Buffer(), stc.plaintext.Buffer()) {
		t.Error("Plaintext should have been matching")
	}

}

func TestEncipherFile(t *testing.T) {
	file := newLyraFile()
	pt := []byte("Wootwoot")
	var err error
	file.plaintext, err = memguard.NewImmutableFromBytes(pt)
	if err != nil {
		t.Error(err)
	}

	dek := []byte("A password")
	ek, err := lcrypt.NewLKey(dek, nil)
	if err != nil {
		t.Error(err)
	}

	a, err := file.EncipherFile(ek)
	if err != nil {
		t.Error(err)
	}

	if reflect.DeepEqual(a.ciphertext, file.plaintext) {
		t.Error("Encryption did not work should not have been equal")
	}
}

func TestParseFileLyraFile(t *testing.T) {
	fixture := "../../test/fixture.2.txt"
	fixture2 := "../../test/fixture.2.txt"

	ctx := newLyraFile()
	err := ctx.ParseFile(fixture)
	if err != nil {
		t.Error(err)
	}

	getFix := func(file string) []byte {
		hg, _ := readFile(file)
		defer hg.Close()
		buf := bytes.NewBuffer(nil)

		_, _ = io.Copy(buf, hg)

		return buf.Bytes()
	}
	fixText := getFix("../../test/fixture.2.txt")

	if !reflect.DeepEqual(fixText, ctx.plaintext.Buffer()) {
		t.Error("Did not Parse plaintext correctly")
	}

	fixText = getFix("../../test/fixture.1.txt")
	if reflect.DeepEqual(fixText, ctx.plaintext) {
		t.Error("Error should have not parsed correctly")
	}

	err = ctx.ParseFile(fixture2)
	if err != nil {
		t.Error("plaintext should have been destroyed.")
	}
}

func TestWriteLyraFile(t *testing.T) {
	fixture3 := "../../test/fixture.7.go.txt"
	writeTo := "../../test/yoyo.go.txt"

	ctx := newLyraFile()
	err := ctx.ParseFile(fixture3)
	if err != nil {
		t.Error(err)
	}

	getFix := func(file string) []byte {
		hg, _ := readFile(file)
		defer hg.Close()
		buf := bytes.NewBuffer(nil)

		_, _ = io.Copy(buf, hg)

		return buf.Bytes()
	}

	tester := newLyraFile()
	tester.plaintext, err = memguard.Duplicate(ctx.plaintext)
	if err != nil {
		t.Error(err)
	}

	err = ctx.Write(writeTo)
	if err != nil {
		t.Error(err)
	}

	fixText := getFix(writeTo)

	if !reflect.DeepEqual(fixText, tester.plaintext.Buffer()) {
		t.Error("Did not Write plaintxt correctly")
	}
}

func TestDestroyFile(t *testing.T) {
	fixture3 := "../../test/fixture.7.go.txt"

	ctx := newLyraFile()
	err := ctx.ParseFile(fixture3)
	if err != nil {
		t.Error(err)
	}

	err = ctx.DestroyFile()
	if err != nil {
		t.Error(err)
	}
}
