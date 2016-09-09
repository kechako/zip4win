package zip4win

import (
	"bytes"

	"github.com/pkg/errors"
	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"
)

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
