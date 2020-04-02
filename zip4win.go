// Package zip4win provides Zip archiver.
package zip4win

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"golang.org/x/text/unicode/norm"
)

const dsStoreName = ".ds_store"

// Writer implements a zip file writer.
type Writer struct {
	zw               *zip.Writer
	Normalizing      bool
	ExcludeDSStore   bool
	ExcludeDotfiles  bool
	UseUTC           bool
	CompressionLevel int

	fwPool sync.Pool
}

// New returns a new Writer wrting a zip file to w with converting file name encoding.
func New(w io.Writer) *Writer {
	writer := &Writer{
		zw:               zip.NewWriter(w),
		Normalizing:      true,
		ExcludeDSStore:   true,
		ExcludeDotfiles:  false,
		UseUTC:           false,
		CompressionLevel: 6,
	}
	writer.init()

	return writer
}

func (w *Writer) init() {
	w.zw.RegisterCompressor(zip.Deflate, zip.Compressor(w.newFlateWriter))
}

// Close finishes writing the zip file by writing the central directory. It does not (and cannot) close the underlying writer.
func (w *Writer) Close() error {
	return w.zw.Close()
}

// create adds a file to zip file using the provided name.
func (w *Writer) create(fi os.FileInfo, name string) (io.Writer, error) {
	h, err := zip.FileInfoHeader(fi)
	if err != nil {
		return nil, err
	}

	h.Method = zip.Deflate

	if filepath.IsAbs(name) {
		// If path is absolute, a entry name is a relative path from root.
		name, err = filepath.Rel(filepath.Clean("/"), name)
		if err != nil {
			return nil, fmt.Errorf("could not get a relative path from root [%s]: %w", name, err)
		}
	}
	name = filepath.ToSlash(filepath.Clean(name))

	if w.Normalizing {
		name = norm.NFC.String(name)
	}

	if fi.IsDir() {
		name = name + "/"
	}

	h.Name = name

	if !w.UseUTC {
		h.Modified = fi.ModTime()
	}

	return w.zw.CreateHeader(h)
}

// WriteEntry add a new entry to zip archive.
func (w *Writer) WriteEntry(path string) error {
	wd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("cound not get the working directory: %w", err)
	}
	fiWd, err := os.Lstat(wd)
	if err != nil {
		return fmt.Errorf("cound not get the working directory: %w", err)
	}

	err = filepath.Walk(path, func(p string, fi os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if os.SameFile(fi, fiWd) {
			return nil
		}
		if w.ExcludeDSStore && strings.ToLower(fi.Name()) == dsStoreName {
			return nil
		}
		if w.ExcludeDotfiles && strings.HasPrefix(fi.Name(), ".") {
			return nil
		}

		return w.writeFile(p, fi)
	})
	if err != nil {
		if pathErr, ok := err.(*os.PathError); ok {
			return fmt.Errorf("no such file or directory [%s]: %w", pathErr.Path, err)
		}

		return err
	}

	return nil
}

// writeFile creates a entry to the zip archive.
func (w *Writer) writeFile(path string, fi os.FileInfo) error {
	fw, err := w.create(fi, path)
	if err != nil {
		return fmt.Errorf("could not create a new file in zip archive: %w", err)
	}

	fmt.Printf("%s\n", path)

	if fi.IsDir() {
		return nil
	}

	fp, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("could not open the file [%s]: %w", path, err)
	}
	defer fp.Close()

	_, err = io.Copy(fw, fp)
	if err != nil {
		return fmt.Errorf("could not write to zip archive: %w", err)
	}

	return nil
}
