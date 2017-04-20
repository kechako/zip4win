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
