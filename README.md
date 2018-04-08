<p align="center">
  <img src="https://cdn.rawgit.com/azohra/lyra/master/docs/assets/logo.svg" height="250" width="250" />
  <p align="center">
    <a href="https://travis-ci.org/azohra/lyra"><img src="https://travis-ci.org/azohra/lyra.svg?branch=master"></a>
    <a href="https://goreportcard.com/report/github.com/azohra/lyra"><img src="https://goreportcard.com/badge/github.com/azohra/lyra"></a>
    <a href="https://github.com/azohra/lyra/blob/master/LICENSE"><img src="https://img.shields.io/badge/License-BSD%203--Clause-blue.svg"></a>
  </p>
</p>

---
Lyra is a lightweight and fast encryption tool that makes protecting your sensitive files easy. 

# Features
* Simple command line encryption and decryption. 
  * Encrypting is as easy as entering: `lyra encrypt file` and to decrypting is just as simple: `lyra decrypt file`
* No need to worry about complex cipher options and configurations.
  * Lyra uses a single strong time tested cipher to protect your data. Your data is encrypted with AES-256-GCM which simultaneously provides data confidentiality, authenticity and integrity (see [authenticated encryption](https://en.wikipedia.org/wiki/Authenticated_encryption)).
  * Values that need to be unique and random (salt and nonce) are generated via a cryptographically secure pseudo random number generator.
* Strong GPU and ASIC resistant KDF being used to protect your passphrase.
  * Lyra uses argon2 to make dictionary attacks and brute force guessing even harder.
* Generate strong memorable passphrases via the diceware method.

# Requirements (if building from source)
* Go 1.9 and above
* [Go Deps](https://golang.github.io/dep/) for dependency management.

# Dependencies
* [gware](https://github.com/brsmsn/gware) for diceware passphrase generation.
* [golang.org/x/crypto](https://github.com/golang/crypto) for argon2 and passphrase terminal.
* [memguard](https://github.com/awnumar/memguard) for handling keys and plaintext secrurely in memory.

# Installation
#### Binaries
* Signed [binaries](#releases) available for windows, linux and macOs available [here](https://github.com/azohra/lyra/releases)
* Simply `mv` the binary to your `$PATH`
#### Installation from source
* If you don't feel comfortable to pipe to `sh` you can alternatevly do:
  * `go get -d github.com/azohra/lyra`
  * `make install`
#### Installation from Brew
* `brew install azohra/tools/lyra`

# Usage
```
Lyra is a lightweight tool used to protect sensitive data

Usage: lyra [Command]

Commands:

	encrypt		Encipher a specified file with inputed passphrase
	decrypt		Decipher a specified file with inputed passphrase
	generate	Generate diceware passphrase(s) via the EFF new worldlist
		
To get more info on commands do: lyra [Command] --help
```

# Releases
Binaries and tags are all signed. The signing key used can be found by searching `Brandon Sam Soon (work key) <brandon.samsoon@loblaw.ca>` with keyid of `5604E4DC6DC74D9B`.
