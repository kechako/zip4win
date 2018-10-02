package zip4win

import (
	"compress/flate"
	"errors"
	"io"
	"sync"
)

func (writer *Writer) newFlateWriter(w io.Writer) (io.WriteCloser, error) {
	fw, ok := writer.fwPool.Get().(*flate.Writer)
	if ok {
		fw.Reset(w)
	} else {
		fw, _ = flate.NewWriter(w, writer.CompressionLevel)
	}
	return &pooledFlateWriter{parent: writer, fw: fw}, nil
}

type pooledFlateWriter struct {
	parent *Writer
	mu     sync.Mutex // guards Close and Write
	fw     *flate.Writer
}

func (w *pooledFlateWriter) Write(p []byte) (n int, err error) {
	w.mu.Lock()
	defer w.mu.Unlock()
	if w.fw == nil {
		return 0, errors.New("Write after Close")
	}
	return w.fw.Write(p)
}

func (w *pooledFlateWriter) Close() error {
	w.mu.Lock()
	defer w.mu.Unlock()
	var err error
	if w.fw != nil {
		err = w.fw.Close()
		w.parent.fwPool.Put(w.fw)
		w.fw = nil
	}
	return err
}
