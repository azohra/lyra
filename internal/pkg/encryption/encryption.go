package encryption

import (
	"github.com/azohra/lyra/pkg/lcrypt"
	"github.com/azohra/lyra/pkg/lfile"
)

// Encrypt encrypts file file and overides the content of file with the ciphertext
// of the specified plaintext file.
func Encrypt(file, saveTo string, passphrase []byte) error {
	ptFile, err := lfile.NewParsedLyraFile(file)
	if err != nil {
		return err
	}

	key, err := lcrypt.NewLKey(passphrase, nil)
	if err != nil {
		return err
	}

	ctFile, err := ptFile.EncipherFile(key)
	if err != nil {
		return err
	}

	switch saveTo {
	case "":
		err = ctFile.Write(file)
	default:
		err = ctFile.Write(saveTo)
	}
	if err != nil {
		return err
	}
	err = key.DestroyKey()
	if err != nil {
		return err
	}

	return nil
}

// Decrypt encrypts file file and overides the content of file with the ciphertext
// of the specified plaintext file.
func Decrypt(file, saveTo string, print bool, passphrase []byte) error {
	ctFile, err := lfile.NewParsedSLFile(file)
	if err != nil {
		return err
	}

	key, err := lcrypt.NewLKey(passphrase, ctFile.RetrieveSalt())
	if err != nil {
		return err
	}

	ptFile, err := ctFile.DecipherFile(key)
	if err != nil {
		return err
	}

	if saveTo == "" && !print {
		err = ptFile.Write(file)
	} else if saveTo != "" && !print {
		err = ptFile.Write(saveTo)
	}

	if err != nil {
		return err
	}

	if print {
		ptFile.PrintLyraFile()
	}

	err = key.DestroyKey()
	if err != nil {
		return err
	}

	return nil
}
