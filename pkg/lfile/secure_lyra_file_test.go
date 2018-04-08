package lfile

import (
	"bytes"
	"encoding/hex"
	"io"
	"reflect"
	"testing"

	"github.com/awnumar/memguard"
	"github.com/azohra/lyra/pkg/lcrypt"
)

func TestNewSecureLyraFile(t *testing.T) {
	stc := newSecureLyraFile()
	if stc == nil {
		t.Error("Did not initialize securelyraFile")
	}
}

func TestNewParsedSLFile(t *testing.T) {
	fixture := "../../test/fixture.3.go.txt"

	stc, err := NewParsedSLFile(fixture)
	if err != nil {
		t.Error(err)
	}

	ctx := newSecureLyraFile()
	err = ctx.ParseFile(fixture)
	if err != nil {
		t.Error(err)
	}

	if !reflect.DeepEqual(ctx.ciphertext, stc.ciphertext) {
		t.Error("Plaintext should have been matching")
	}

	if !reflect.DeepEqual(ctx.nonce, stc.nonce) {
		t.Error("Nonce should have been matching")
	}

	if !reflect.DeepEqual(ctx.salt, stc.salt) {
		t.Error("salt should have been matching")
	}

}

func TestGenAuthParameters(t *testing.T) {
	stc := newSecureLyraFile()
	err := stc.GenerateAuthParams()
	if err != nil {
		t.Error(err)
	}

	if len(stc.nonce) != lcrypt.NonceSize || cap(stc.nonce) != lcrypt.NonceSize {
		t.Error("Nonce not initialized")
	}
}

func TestDecipherFile(t *testing.T) {
	var err error
	file := newLyraFile()
	p := []byte("exampleplaintext")
	file.plaintext, err = memguard.NewImmutableFromBytes(p)

	dek := []byte("password")
	ek, err := lcrypt.NewLKey(dek, nil)
	if err != nil {
		t.Error(err)
	}

	a, err := file.EncipherFile(ek)
	if err != nil {
		t.Error(err)
	}

	b, err := a.DecipherFile(ek)
	if err != nil {
		t.Error(err)
	}

	if !reflect.DeepEqual([]byte("exampleplaintext"), b.plaintext.Buffer()) {
		t.Error("Decryption Failed")
	}
}

func TestParseFileSecureLyra(t *testing.T) {
	fixture := "../../test/fixture.1.txt"
	fixture2 := "../../test/fixture.2.txt"
	fixture3 := "../../test/fixture.3.go.txt"

	test := func(fix string, fix2 []byte, num string) {
		ctx := newSecureLyraFile()
		err := ctx.ParseFile(fix)
		if err != nil {
			t.Error(err)
		}

		fixTextProper := fix2
		salt, err := hex.DecodeString("A6C4E5")
		if err != nil {
			t.Error(err)
		}
		nonve, err := hex.DecodeString("B5A6")
		if err != nil {
			t.Error(err)
		}

		if !reflect.DeepEqual(fixTextProper, ctx.ciphertext) {
			t.Error("Did not Parse ciphertxt correctly for " + num)
		}

		if !reflect.DeepEqual(salt, ctx.salt) || !reflect.DeepEqual(nonve, ctx.nonce) {
			t.Error("Auth data did not parse correctly for " + num)
		}
	}

	getFix := func(file string) []byte {
		hg, _ := readFile(file)
		defer hg.Close()
		buf := bytes.NewBuffer(nil)

		_, _ = io.Copy(buf, hg)

		return buf.Bytes()
	}
	aa := getFix("../../test/fixture.21.txt")
	bb := getFix("../../test/fixture.4.go.txt")

	test(fixture, []byte("this is a test\n\n\nOne Line"), "1")

	test(fixture2, aa, "2")
	test(fixture3, bb, "3")
}

func TestParseAuthData(t *testing.T) {
	fixture := "@!A7B3@!C6E5"

	a, b, err := parseAuthData(fixture)
	if err != nil {
		t.Error(err)
	}

	a1, err := hex.DecodeString("A7B3")
	if err != nil {
		t.Error(err)
	}

	b1, err := hex.DecodeString("C6E5")
	if err != nil {
		t.Error(err)
	}

	if !reflect.DeepEqual(a, a1) && reflect.DeepEqual(b, b1) {
		t.Error("Did not parse properly")
	}
}

func TestWriteSecureLyraFile(t *testing.T) {
	fixture3 := "../../test/fixture.3.go.txt"
	writeTo := "../../test/writer.go.txt"

	ctx := newSecureLyraFile()
	err := ctx.ParseFile(fixture3)
	if err != nil {
		t.Error(err)
	}

	err = ctx.Write(writeTo)
	if err != nil {
		t.Error(err)
	}

	ctx2 := newSecureLyraFile()
	err = ctx2.ParseFile(fixture3)
	if err != nil {
		t.Error(err)
	}

	if !reflect.DeepEqual(ctx, ctx2) {
		t.Error("Failed to write")
	}
}
