package ar

import (
	"io"
	"time"
)

type Writer struct {
	w io.Writer
	globalHeader bool
}

type Header struct {
	FileName string
	ModTime time.Time
	Uid int
	Gid int
	Mode int64
	Size int64
}

func NewWriter(w io.Writer) *Writer { return &Writer{w: w} }

func (aw *Writer) Write(b []byte) (n int, err error) {
	return
}

func (aw *Writer) WriteHeader(hdr *Header) error {
	if !aw.globalHeader {
		_, err := aw.w.Write([]byte("!<arch>\n"))
		if err != nil {
			return err
		}
		aw.globalHeader = true
	}

	return nil
}