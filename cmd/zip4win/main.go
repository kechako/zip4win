package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/kechako/zip4win"
	"github.com/pkg/errors"
)

// writeFile creates a entry to the zip archive.
func writeFile(w *zip4win.Writer, path string) error {
	var name string
	if filepath.IsAbs(path) {
		name = filepath.Base(path)
	} else {
		name = filepath.Clean(path)
	}

	fp, err := os.Open(path)
	if err != nil {
		return errors.Wrapf(err, "Could not open the file [%s].", path)
	}
	defer fp.Close()

	fw, err := w.Create(name)
	if err != nil {
		return errors.Wrap(err, "Could not create a new file in zip archive.")
	}

	fmt.Printf("%s => %s\n", path, name)

	_, err = io.Copy(fw, fp)
	if err != nil {
		return errors.Wrap(err, "Could not write to zip archive.")
	}

	return nil
}

// printHelp outputs a help message to STDERR.
func printHelp() {
	fmt.Fprintf(os.Stderr, `Usage: %s zipfile file ...`, os.Args[0])
	os.Exit(1)
}

// printError outputs a error message to STDERR.
func printError(err error) {
	fmt.Fprintf(os.Stderr, "Error: %s", err)
	os.Exit(2)
}

// entry point
func main() {
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

	for _, path := range paths {
		err = writeFile(w, path)
		if err != nil {
			printError(err)
		}
	}
}
