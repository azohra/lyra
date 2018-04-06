//Package lcrypt provides crypto operations for lyra
package lcrypt

import (
	"fmt"
	"os"

	"golang.org/x/crypto/argon2"
)

//Argon2 parameters

//Iterations uint32 number of times argon2 is to be run with args
const Iterations uint32 = 4

//Memory uint32 size of memory cost
const Memory uint32 = 64 * 1024

//Threads uint8 the number of threads to be adjusted
const Threads uint8 = 4

//KeyLen uint32 the size of the key to be generated
const KeyLen uint32 = 32

//SaltSize uint8 is the salt length
const SaltSize uint16 = 16

//GenKey generates a cipher key from pass. Key is generated via argon2id with parameters
//specified via Iterations, Memory, Threads and KeyLen
func GenKey(pass []byte, salt []byte) []byte {
	key := argon2.IDKey(pass, salt, Iterations, Memory, Threads, KeyLen)
	return key
}

//GenSalt generates a random salt of SaltSize
func GenSalt() []byte {
	salt := make([]byte, SaltSize)

	_, err := GenNonce(salt)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to generate Salt")
	}

	return salt
}
