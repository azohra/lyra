package lcrypt

import (
	"errors"

	"github.com/awnumar/memguard"
)

//lockedBuffer is a memguard LockedBuffer
type lockedBuffer = *memguard.LockedBuffer

//LKey represents an immutable encryption key. This key is protected from unintended swap,
//buffer overflows and underflows. However the most important is that it can be securely destroyed.
type LKey struct {
	//key is the encryption key
	key lockedBuffer

	salt []byte
}

//KeyManager implements encryption Key management
type KeyManager interface {
	//initKey initiallizes a key with passphrase and salt, return a bool and error based on success
	initKey(passphrase string, salt []byte) (bool, error)

	//DestroyKey safely destroys a key that has been stored in memory.
	DestroyKey() error

	//GetKey gets a key
	GetKey() []byte
}

//NewLKey returns a new encryption key init to passphrase and salt, if no salt is specified, a new salt will be
//generated for the new key.
func NewLKey(passphrase, salt []byte) (*LKey, error) {
	k := &LKey{}

	if salt == nil {
		k.salt = GenSalt()
	} else {
		k.salt = salt
	}

	err := k.initKey(passphrase, k.salt)
	if err != nil {
		return nil, err
	}
	return k, nil
}

//GetKey returns a reference to the key.
func (k *LKey) GetKey() []byte {
	return k.key.Buffer()
}

//GetSalt returns a reference to the salt of the key.
func (k *LKey) GetSalt() []byte {
	return k.salt
}

//InitKey initializes an encryption key passphrase with salt salt. InitKey initializes a locked
//buffer that will be resistant to unauthorized memory manipulation, additionally it will wipe
//passphrase.
func (k *LKey) initKey(passphrase, salt []byte) error {
	var err error
	k.key, err = memguard.NewImmutableFromBytes(GenKey(passphrase, salt))

	wipeData(passphrase)
	return err
}

//DestroyKey safely destroy the encryption from memory,
func (k *LKey) DestroyKey() error {
	k.key.Destroy()
	if !k.key.IsDestroyed() {
		return errors.New("Failed to destroy key")
	}
	return nil
}

//wipeData will wipe any dst slice with cryptographically secure bytes.
func wipeData(dst []byte) error {
	_, err := GenNonce(dst)
	return err
}
