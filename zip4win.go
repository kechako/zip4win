package zip4win

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
	"golang.org/x/text/unicode/norm"
)

const dsStoreName = ".ds_store"

// Writer implements a zip file writer.
type Writer struct {
	zw              *zip.Writer
	Normalizing     bool
	ExcludeDSStore  bool
	ExcludeDotfiles bool
}

// New returns a new Writer wrting a zip file to w with converting file name encoding.
func New(w io.Writer) *Writer {
	return &Writer{
		zw:              zip.NewWriter(w),
		Normalizing:     true,
		ExcludeDSStore:  true,
		ExcludeDotfiles: false,
	}
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

	if filepath.IsAbs(name) {
		// If path is absolute, a entry name is a relative path from root.
		name, err = filepath.Rel(filepath.Clean("/"), name)
		if err != nil {
			return nil, errors.Wrapf(err, "Could not get a relative path from root : %s", name)
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

	// Set UTF-8 Flag
	h.Flags = h.Flags | 0x0800

	return w.zw.CreateHeader(h)
}

// WriteEntry add a new entry to zip archive.
func (w *Writer) WriteEntry(path string) error {
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
			return errors.Wrapf(err, "No such file or directory : %s", pathErr.Path)
		}

		return err
	}

	return nil
}

// writeFile creates a entry to the zip archive.
func (w *Writer) writeFile(path string, fi os.FileInfo) error {
	fw, err := w.create(fi, path)
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
