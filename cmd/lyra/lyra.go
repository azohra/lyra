package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"reflect"
	"strings"

	"golang.org/x/crypto/ssh/terminal"
)

const (
	//BadExit represents an exit with an error(s)
	BadExit = 1

	//GoodExit represents an exit with no errors
	GoodExit = 0

	//usage string
	usage = `Lyra is a lightweight tool used to protect sensitive data

Usage: lyra [Command]

Commands:

	encrypt		Encipher a specified file with inputed passphrase
	decrypt		Decipher a specified file with inputed passphrase
	generate	Generate diceware passphrase(s) via the EFF new worldlist
		
To get more info on commands do: lyra [Command] --help
`
	about = `Lyra is a lightweight tool used to protect sensitive data.
	
Coded with ❤️ by the Azohra team and made possible by open source software:

  * gware by brsmsn (BSD-3-clause) https://github.com/brsmsn/gware
    * license available @ https://github.com/brsmsn/gware/blob/master/LICENSE

  * memguard by awnumar (Apache 2.0) https://github.com/awnumar/memguard
    * License available @ https://github.com/awnumar/memguard/blob/master/LICENSE  
		
  * crypto by golang.org (BSD-style) https://github.com/golang/crypto
    * license available @ https://github.com/golang/crypto/blob/master/LICENSE
`
	version = `Version: 1.0.0 (April 2018)
`
)

type command interface {
	CName() string
	Help() string
	RegCFlags(*flag.FlagSet)
	Run([]string) error
}

//app start point
func main() {
	commands := [...]command{
		&encryptcmd{},
		&decryptcmd{},
		&gencmd{},
	}

	versionSet := flag.Bool("version", false, "")
	aboutSet := flag.Bool("about", false, "")

	//for lyra --help
	flag.Usage = func() {
		fmt.Fprint(os.Stderr, usage)
	}

	flag.Parse()
	args := flag.Args()

	if *versionSet {
		fmt.Fprint(os.Stdout, version)
		os.Exit(GoodExit)
	} else if *aboutSet {
		fmt.Fprint(os.Stdout, about)
		os.Exit(GoodExit)
	}

	cmdName, helpc, exit := parseCmd(args)
	if exit {
		fmt.Fprint(os.Stderr, usage)
		os.Exit(BadExit)
	}

	for _, val := range commands {
		if val.CName() == cmdName {
			flags := flag.NewFlagSet(cmdName, flag.ContinueOnError)
			flags.SetOutput(os.Stderr)

			val.RegCFlags(flags)
			if helpc {
				fmt.Fprint(os.Stderr, "Available options/flags - see below to learn how to implement them:\n\n")
				flags.PrintDefaults()
				fmt.Fprint(os.Stderr, val.Help())
				os.Exit(BadExit)
			}

			err := flags.Parse(args[1:])
			if err != nil {
				handleErr(err)
			}

			err = val.Run(flags.Args())
			if err != nil {
				handleErr(err)
			}

			os.Exit(GoodExit)
		}
	}

	fmt.Fprint(os.Stderr, "Command "+"\""+args[0]+"\" not found\n\n"+usage)
	os.Exit(BadExit)
}

func parseCmd(args []string) (name string, helpNd bool, exit bool) {
	helpNeeded := func(opt string) bool {
		str := strings.ToLower(opt)
		return str == "--help" || str == "-h" || str == "-help"
	}

	switch len(args) {
	case 0:
		exit = true
	case 1:
		if helpNeeded(args[0]) {
			exit = true
		} else {
			name = args[0]
		}
	default:
		if helpNeeded(args[1]) {
			name = args[0]
			helpNd = true
		} else {
			name = args[0]
		}
	}

	return name, helpNd, exit
}

func handleErr(err error) {
	if err != nil {
		fmt.Fprint(os.Stderr, "An Error has occurred:\n"+err.Error()+"\n")
		os.Exit(BadExit)
	}
}

func getPassphrase() []byte {
	fmt.Println("Enter passphrase: ")
	input, err := terminal.ReadPassword(0)
	handleErr(err)
	return input
}

func setPassphrase() ([]byte, error) {
	fmt.Println("Enter passphrase: ")
	in1, err := terminal.ReadPassword(0)
	handleErr(err)
	fmt.Println("Enter passphrase again: ")
	in2, err := terminal.ReadPassword(0)
	handleErr(err)

	defer wipe(in2)

	if !reflect.DeepEqual(in1, in2) {
		return nil, errors.New("Inputed passphrases do not match")
	}

	return in1, nil
}

func wipe(ins ...[]byte) {
	for _, val := range ins {
		for k := range val {
			val[k] = 0
		}
	}
}
