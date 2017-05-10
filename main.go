package main

import (
	"fmt"
	"github.com/kardianos/osext"
	"log"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
)

// Makes multiple scripts/executables look like a single command, like git does for extras

// Idea from magic-cli, some code from http://stackoverflow.com/questions/18537257/golang-how-to-get-the-directory-of-the-currently-running-file
// (I wasn't 100% sure how do do that *reliably*)

// Original magic-cli is here: https://github.com/slackhq/magic-cli/blob/master/magic-cli

// Set in main func
var (
	debug       bool
	magicPrefix string
	magicPath   string
)

// The original (Ruby) implementation *starts* by compiling a list of all the possible commands
// This would make some cases slower but maybe some cases faster
// Could be smart and make it compile the list (as a global?) when it first needs it
// (This sounds like a tiny difference to fret over, but we serve this stuff over NFS and Lustre so stat()s are slowish)

// Could do this by doing a Walk and appending each prefix-thing to a thing->full_path string->string map

// *Could* instead of/as well as current same-dir, have a MAGIC_PATH or MAGIC_SEARCH_PATH env var or something similar
// (and *default* to same-dir)

// Could also have a universal prefix that all magic commands treat as theirs
// (Default to "any" but settable by env var?)

func getMagicPath() string {
	path, err := osext.ExecutableFolder()
	if err != nil {
		log.Fatal(err)
	}
	return path
}

func printFilename(path string, info os.FileInfo, err error) error {
	if err != nil {
		log.Print(err)
		return nil
	}
	fmt.Println(info.Name())
	return nil
}

func listCommands() {

	if debug == true {
		fmt.Println("Looking in: ", magicPath)
	}
	commands, err := filepath.Glob(filepath.Join(magicPath, magicPrefix+"-*"))
	if debug == true {
		fmt.Println("Found: ", commands)
	}

	if err != nil {
		panic(err)
	}

	// Might want to do this with Walk instead.
	// (Walk: https://xojoc.pw/justcode/golang-file-tree-traversal.html )
	// This would let you put commands in subdirs as well as same dir as exec
	// (Good for testing, organising)
	for i := 0; i < len(commands); i++ {
		if canUseCommandFile(commands[i]) {
			printCommandInfo(commands[i], magicPrefix)
		}
	}

	// E.g.
	err = filepath.Walk(magicPath, printFilename)
	// Might want to stop at first result, e.g. http://stackoverflow.com/a/36713726/180184
	// You'd want to put the needle you're searching for into WalkFun (printFilename) as a closure
}

func testForCommand(commandName string) bool {
	filename := filepath.Join(magicPath, magicPrefix+"-"+commandName)
	return canUseCommandFile(filename)
}

func canUseCommandFile(filename string) bool {
	info, _ := os.Stat(filename)

	// Can we stat it?
	if info != nil {
		mode := info.Mode()
		// Can we read it?
		canRead := bool((mode & 0444) != 0)

		// Can we execute it?
		canExec := bool((mode & 0111) != 0)

		if canRead && canExec {
			return true
		}
	} else {
		return false
	}
	return false
}

func printCommandInfo(filename string, execBasename string) {
	// Still not 100% sure how to interface this
	// Is it better to read info from files, or run them with an arg?
	// Or something else?

	// Ooooor, test first two bytes for "#!", if present, grab second line as description, otherwise run with -w
	// With magic marker to tell this to run a script with -w anyway

	commandBasename := path.Base(filename)
	commandName := strings.SplitAfterN(commandBasename, execBasename+"-", 2)[1]

	fmt.Printf("  %-10v\t", commandName)

	cmd := exec.Command(filename, "-w")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()

	if err != nil {
		panic(err)
	}
}

// TODO: remove runcommand functionality from printCommandInfo, wrap runCommand
func runCommand(commandName string, args []string) {
	// Note that at this point we should have already tested that this executable exists.
	executable := filepath.Join(magicPath, magicPrefix+"-"+commandName)

	if debug == true {
		log.Println("Trying to run: ", executable, "\n with args: ", strings.Join(args, " "))
	}

	// v- variadic syntax
	cmd := exec.Command(executable, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		panic(err)
	}
}

func usage() {
	listCommands()
}

func main() {
	// v- Package scope
	magicPath = getMagicPath()
	magicPrefix = path.Base(os.Args[0])
	debug = true

	numArgs := len(os.Args)

	// TODO: Messy and manual, needs work.
	if (numArgs < 2) || (numArgs == 2 && ((os.Args[1] == "-h") || (os.Args[1] == "--help"))) {
		usage()
	} else if (numArgs == 2) && ((os.Args[1] == "-l") || (os.Args[1] == "--list")) {
		listCommands()
	} else if testForCommand(os.Args[1]) {
		if debug == true {
			log.Println("Found runnable command: ", magicPath, magicPrefix, os.Args[1])
		}
		runCommand(os.Args[1], os.Args[2:])
	} else {
		panic("Unknown flags or options.")
	}
}
