package bytes

import (
	"bytes"
	"io"
	"sync"
)

type Buffer struct {
	mux sync.RWMutex
	buf *bytes.Buffer
	w   io.Writer //复制一份的写入
}

func NewBuffer(w io.Writer) *Buffer {
	return &Buffer{
		mux: sync.RWMutex{},
		buf: &bytes.Buffer{},
		w:   w,
	}
}

func (buf *Buffer) Write(b []byte) (n int, err error) {
	buf.mux.Lock()
	defer buf.mux.Unlock()
	return buf.write(b)
}

func (buf *Buffer) WriteString(s string) (n int, err error) {
	buf.mux.Lock()
	defer buf.mux.Unlock()
	return buf.write([]byte(s))
}

func (buf *Buffer) Bytes() []byte {
	buf.mux.RLock()
	defer buf.mux.RUnlock()
	return buf.buf.Bytes()
}

func (buf *Buffer) String() string {
	buf.mux.RLock()
	defer buf.mux.RUnlock()
	return string(buf.buf.Bytes())
}

func (buf *Buffer) Reset() {
	buf.mux.Lock()
	buf.buf.Reset()
	buf.mux.Unlock()
}

func (buf *Buffer) write(b []byte) (n int, err error) {
	if buf.w != nil {
		n, err = buf.w.Write(b)
		if err != nil {
			return
		}
	}
	return buf.buf.Write(b)
}
