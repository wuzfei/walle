package bytes

import (
	"errors"
	"io"
	"sync"
)

var ErrBufferOver = errors.New("write close")

type BufferOver struct {
	mux    sync.RWMutex
	buf    []byte
	isOver bool
}

func NewBufferOver() *BufferOver {
	return &BufferOver{buf: make([]byte, 0, 1024)}
}

func (bo *BufferOver) Write(b []byte) (n int, err error) {
	bo.mux.Lock()
	defer bo.mux.Unlock()
	if bo.isOver {
		return 0, ErrBufferOver
	}
	bo.buf = append(bo.buf, b...)
	return len(b), nil
}

func (bo *BufferOver) ReadAt(b []byte, off int) (n int, err error) {
	bo.mux.RLock()
	defer bo.mux.RUnlock()
	l := len(bo.buf)
	if l > off {
		n = copy(b, bo.buf[off:l])
	}
	if off+n >= l && bo.isOver {
		err = io.EOF
	}
	return
}

func (bo *BufferOver) Len() int {
	bo.mux.RLock()
	defer bo.mux.RUnlock()
	return len(bo.buf)
}

func (bo *BufferOver) WriteOver() {
	bo.mux.Lock()
	defer bo.mux.Unlock()
	bo.isOver = true
}
