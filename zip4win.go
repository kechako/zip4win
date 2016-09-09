package zip4win

import (
	"archive/zip"
	"bytes"
	"io"

	"github.com/pkg/errors"
	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"
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
func (w *Writer) Create(name string) (io.Writer, error) {
	sname, err := convertToShiftJIS(name)
	if err != nil {
		return nil, err
	}

	return w.zw.Create(sname)
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
