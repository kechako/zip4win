package zip4win

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

// Writer implements a zip file writer.
type Writer struct {
	zw *zip.Writer
}

// New returns a new Writer wrting a zip file to w with converting file name encoding.
func New(w io.Writer) *Writer {
	return &Writer{
		zw: zip.NewWriter(w),
	}
}

// Close finishes writing the zip file by writing the central directory. It does not (and cannot) close the underlying writer.
func (w *Writer) Close() error {
	return w.zw.Close()
}

// Create adds a file to zip file using the provided name.
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
	name = norm.NFC.String(name)

	if fi.IsDir() {
		name = name + "/"
	}

	h.Name, err = convertToShiftJIS(name)
	if err != nil {
		return nil, err
	}

	return w.zw.CreateHeader(h)
}

// writeFile add a new entry to zip archive.
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

// convertToShiftJIS converts a UTF-8 string to a ShiftJIS string.
func convertToShiftJIS(name string) (string, error) {
	var buf bytes.Buffer
	w := transform.NewWriter(&buf, japanese.ShiftJIS.NewEncoder())
	defer w.Close()

	_, err := w.Write([]byte(name))
	if err != nil {
		return "", errors.Wrap(err, "Could not convert a utf8 string to a sjis string.")
	}

	return buf.String(), nil
}
