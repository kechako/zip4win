package zip4win

import (
	"bytes"
	"testing"
)

func TestNew(t *testing.T) {
	buf := new(bytes.Buffer)
	w := New(buf)
	if w == nil {
		t.Error("New retruned nil.")
	}
	if w.zw == nil {
		t.Error("New returned uninitialized value.")
	}
}

func TestConvertToShiftJIS(t *testing.T) {
	expect := "\x82\xA0\x82\xA2\x82\xA4\x82\xA6\x82\xA8"
	result, err := convertToShiftJIS("あいうえお")
	if err != nil {
		t.Error(err)
	}
	if result != expect {
		t.Errorf("got %s\nwant %s", expect, result)
	}
}
