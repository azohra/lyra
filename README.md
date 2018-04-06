# Lyra [![Build Status](https://travis-ci.org/azohra/lyra.svg?branch=master)](https://travis-ci.org/azohra/lyra)

Lyra is a lightweight encryption tool that makes protecting your sensitive files easy. 

# Features
* Simple command line encryption and decryption. 
  * Encrypting is as easy as entering: `lyra encrypt file` and to decrypting is just as simple: `lyra decrypt file`
* Strong time tested encryption scheme being used.
  * Lyra uses AES-256-GCM to simultaneously provide data confidentiality, authenticity and integrity (see [authenticated encryption](https://en.wikipedia.org/wiki/Authenticated_encryption)).
* Strong GPU resistant KDF being used to protect your passphrase.
  * Lyra uses argon2 to make brute-forcing your passphrase even harder.
* No key files, simply enter a passphrase or get one generated for you.
* Generate strong memorable passphrases via the diceware method.

# Requirements (if building from source)
* Go 1.9.4 and above
* [Go Deps](https://golang.github.io/dep/) for dependency management.

# Dependencies
* [gware](https://github.com/brsmsn/gware) for diceware passphrase generation.
* [golang.org/x/crypto](https://github.com/golang/crypto) for argon2, passphrase terminal.
* [memguard](https://github.com/awnumar/memguard) for handling keys and plaintext secrurely in memory.

# Installation
* Binaries available for windows, linux and macOs available [here]()
* Simply `mv` the binary to your `$PATH`

# Installation from source
* `go get -d `
* `make install` will install lyra to your path.


# Usage
```
Usage: lyra [Command]

Commands:

	encrypt                Encipher a specified file with inputed passphrase
	decrypt                Decipher a specified file with inputed passphrase
	generate               Generate diceware passphrase(s) via the EFF new worldlist
		
To get more info on commands do: lyra [Command] help

Coded with ❤️ by the Azohra team
```
