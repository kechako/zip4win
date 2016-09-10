package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/kechako/zip4win"
)

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
	var shiftJIS bool
	var nonorm bool
	var includeDSStore bool

	flag.Usage = printHelp
	flag.BoolVar(&shiftJIS, "sjis", false, "Encode a file name in ShiftJIS")
	flag.BoolVar(&nonorm, "nonorm", false, "Disable normalizing a file name with NFC")
	flag.BoolVar(&includeDSStore, "include-dsstore", false, "Include .DSStore in a zip archive.")
	flag.Parse()
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

	w.ShiftJIS = shiftJIS
	w.Normalizing = !nonorm
	w.ExcludeDSStore = !includeDSStore

	for _, path := range paths {
		err = w.WriteEntry(path)
		if err != nil {
			printError(err)
		}
	}
}
