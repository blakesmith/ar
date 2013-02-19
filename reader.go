package ar

import (
	"io"
	"io/ioutil"
	"os"
)

type Reader struct {
	r io.Reader
	nb int64
}

func NewReader(r io.Reader) *Reader {
	return &Reader{r: r}
}

func (rd *Reader) string(b []byte) string {
	n := len(b)-1
	for n > 0 && b[n] == 32 {
		n--
	}

	return string(b[0:n+1])
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

	return header, nil
}

func (rd *Reader) Next() (*Header, error) {
	err := rd.skipUnread()
	if err != nil {
		return nil, err
	}
	
	return rd.readHeader()
}