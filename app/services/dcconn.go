package services

import (
	"fmt"
	"io"
	"log"

	"github.com/gopherjs/gopherjs/js"
)

// DCConn ...
type DCConn struct {
	*js.Object
	dst io.ReadCloser
}

// NewDCConn ...
func NewDCConn(dc *js.Object) *DCConn {
	dst, src := io.Pipe()
	c := &DCConn{Object: dc, dst: dst}
	dc.Set("onclose", func() { src.Close() })
	dc.Set("onmessage", func(ev *js.Object) {
		src.Write([]byte(ev.Get("data").String()))
	})
	dc.Set("onerror", func(ev *js.Object) { log.Println(ev) })
	return c
}

func (w *DCConn) Read(b []byte) (int, error) {
	return w.dst.Read(b)
}

func (w *DCConn) Write(b []byte) (int, error) {
	if w.Get("readyState").String() != "open" {
		return 0, fmt.Errorf("does't open")
	}
	w.Call("send", string(b))
	return len(b), nil
}

// Close ...
func (w *DCConn) Close() error {
	if err := w.dst.Close(); err != nil {
		return err
	}
	w.Call("close")
	return nil
}
