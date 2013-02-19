package ar

import (
	"io"
	"io/ioutil"
	"os"
	"strconv"
	"time"
)

type Reader struct {
	r io.Reader
	nb int64
}

func NewReader(r io.Reader) *Reader {
	return &Reader{r: r}
}

func (rd *Reader) string(b []byte) string {
	i := len(b)-1
	for i > 0 && b[i] == 32 {
		i--
	}

	return string(b[0:i+1])
}

func (rd *Reader) numeric(b []byte) int64 {
	i := len(b)-1
	for i > 0 && b[i] == 32 {
		i--
	}

	n, _ := strconv.ParseInt(string(b[0:i+1]), 10, 64)

	return n
}

func (rd *Reader) skipUnread() error {
	skip := rd.nb
	rd.nb = 0
	if seeker, ok := rd.r.(io.Seeker); ok {
		_, err := seeker.Seek(skip, os.SEEK_CUR)
		return err
	}

	_, err := io.CopyN(ioutil.Discard, rd.r, skip)
	return err
}

func (rd *Reader) readHeader() (*Header, error) {
	headerBuf := make([]byte, HEADER_BYTE_SIZE)
	if _, err := io.ReadFull(rd.r, headerBuf); err != nil {
		return nil, err
	}

	header := new(Header)
	s := slicer(headerBuf)
	s.next(8) // Skip the global header
	header.Name = rd.string(s.next(16))
	header.ModTime = time.Unix(rd.numeric(s.next(12)), 0)

	return header, nil
}

func (rd *Reader) Next() (*Header, error) {
	err := rd.skipUnread()
	if err != nil {
		return nil, err
	}
	
	return rd.readHeader()
}