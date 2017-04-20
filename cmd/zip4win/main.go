package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"

	"github.com/kechako/zip4win"
	"github.com/pkg/errors"
)

var (
	version  = "1.0"
	revision string
)

// printVersion output a version info.
func printVersion() {
	fmt.Printf("%s %s (%s)\n", os.Args[0], version, revision)
}

// printHelp outputs a help message to STDERR.
func printHelp() {
	fmt.Fprintf(os.Stderr, `Usage: %s [options] zipfile file ...

options:
`, os.Args[0])
	flag.PrintDefaults()
	os.Exit(1)
}

// printError outputs a error message to STDERR.
func printError(err error) {
	fmt.Fprintf(os.Stderr, "Error: %s", err)
	os.Exit(2)
}

// entry point
func main() {
	var nonorm bool
	var includeDSStore bool
	var excludeDotfiles bool
	var useUTC bool
	var printVer bool

	flag.Usage = printHelp
	flag.BoolVar(&nonorm, "nonorm", false, "Disable normalizing a file name with NFC.")
	flag.BoolVar(&includeDSStore, "include-dsstore", false, "Include .DSStore in a zip archive.")
	flag.BoolVar(&excludeDotfiles, "exclude-dotfiles", false, "Exclude dotfiles in a zip archive.")
	flag.BoolVar(&useUTC, "utc", false, "Use UTC as mod time. Default use local time.")
	flag.BoolVar(&printVer, "version", false, "Show version info.")

	flag.Parse()

	if printVer {
		printVersion()
		return
	}

	args := flag.Args()
	if len(args) < 2 {
		printHelp()
	}

	zipfile := args[0]
	paths := args[1:]

	fp, err := os.Create(zipfile)
	if err != nil {
		printError(err)
	}
	defer fp.Close()

	w := zip4win.New(fp)
	defer w.Close()

	w.Normalizing = !nonorm
	w.ExcludeDSStore = !includeDSStore
	w.ExcludeDotfiles = excludeDotfiles
	w.UseUTC = useUTC

	if paths[0] == "-" {
		// Input from stdin
		s := bufio.NewScanner(os.Stdin)
		for s.Scan() {
			err = w.WriteEntry(s.Text())
			if err != nil {
				printError(err)
			}
		}
		if err = s.Err(); err != nil {
			printError(errors.Wrap(err, "Could not read from STDIN."))
		}
	} else {
		// Input from parameters
		for _, path := range paths {
			err = w.WriteEntry(path)
			if err != nil {
				printError(err)
			}
		}
	}
}
