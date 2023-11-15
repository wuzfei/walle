package bytes

import (
	"bytes"
	"context"
	"errors"
	"io"
	"testing"
	"time"
)

func TestBufferOver(t *testing.T) {
	tl := NewBufferOver()

	ctx, cancel := context.WithCancel(context.Background())
	wt := time.NewTicker(time.Second / 50)
	defer wt.Stop()
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case <-wt.C:
				_, _ = tl.Write([]byte("123"))
			}
		}
	}()

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case <-wt.C:
				_, _ = tl.Write([]byte("abc"))
			}
		}
	}()

	go func() {
		offset := 0
		buf := bytes.Buffer{}
		for {
			d := make([]byte, 1024)
			n, err := tl.ReadAt(d, offset)
			d = d[:n]
			offset += n
			if n > 0 {
				buf.Write(d)
			}
			if err != nil {
				if !errors.Is(err, io.EOF) {
					t.Fatal("读取失败:", err)
					return
				}
				t.Log("读取结束")
				break
			}
			time.Sleep(time.Second / 15)
		}
		t.Log("获取数据：", buf.String())
		t.Log("获取数据长度：", buf.Len())
	}()

	wOverT := time.NewTimer(time.Second * 1)
	select {
	case <-wOverT.C:
		cancel()
		tl.WriteOver()
	}

	time.Sleep(time.Second * 1)

	t.Log("over")

}
