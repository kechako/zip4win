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

// writeFile add a new entry to zip archive.
func writeEntry(w *zip4win.Writer, path string) error {
	wd, err := os.Getwd()
	if err != nil {
		return errors.Wrap(err, "Cound not get the working directory.")
	}
	fiWd, err := os.Lstat(wd)
	if err != nil {
		return errors.Wrap(err, "Cound not get the working directory.")
	}

	err = filepath.Walk(path, func(p string, fi os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if os.SameFile(fi, fiWd) {
			return nil
		}

		return writeFile(w, p, fi)
	})
	if err != nil {
		if pathErr, ok := err.(*os.PathError); ok {
			return errors.Wrapf(err, "No such file or directory : %s", pathErr.Path)
		}

		return err
	}

	return nil
}

// writeFile creates a entry to the zip archive.
func writeFile(w *zip4win.Writer, path string, fi os.FileInfo) error {
	fw, err := w.Create(fi, path)
	if err != nil {
		return errors.Wrap(err, "Could not create a new file in zip archive.")
	}

	fmt.Printf("%s\n", path)

	if fi.IsDir() {
		return nil
	}

	fp, err := os.Open(path)
	if err != nil {
		return errors.Wrapf(err, "Could not open the file [%s].", path)
	}
	defer fp.Close()

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
		err = writeEntry(w, path)
		if err != nil {
			printError(err)
		}
	}
}
