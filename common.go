/*
Copyright (c) 2013 Blake Smith <blakesmith0@gmail.com>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package ar

import (
	"errors"
	"fmt"
	"os"
	"path"
	"time"
)

const (
	HEADER_BYTE_SIZE = 60
	GLOBAL_HEADER    = "!<arch>\n"
)

type Header struct {
	Name    string
	ModTime time.Time
	Uid     int
	Gid     int
	Mode    int64
	Size    int64
}

type slicer []byte

func (sp *slicer) next(n int) (b []byte) {
	s := *sp
	b, *sp = s[0:n], s[n:]
	return
}

// FileInfo returns an os.FileInfo for the Header.
func (h *Header) FileInfo() os.FileInfo {
	return headerFileInfo{h}
}

// headerFileInfo implements os.FileInfo.
type headerFileInfo struct {
	h *Header
}

func (fi headerFileInfo) Size() int64        { return fi.h.Size }
func (fi headerFileInfo) IsDir() bool        { return false }
func (fi headerFileInfo) ModTime() time.Time { return fi.h.ModTime }
func (fi headerFileInfo) Sys() interface{}   { return fi.h }
func (fi headerFileInfo) Name() string       { return path.Base(fi.h.Name) }

// Mode returns the permission and mode bits for the headerFileInfo.
func (fi headerFileInfo) Mode() (mode os.FileMode) {
	// Set file permission bits.
	mode = os.FileMode(fi.h.Mode).Perm()

	// Set setuid, setgid and sticky bits.
	if fi.h.Mode&c_ISUID != 0 {
		mode |= os.ModeSetuid
	}
	if fi.h.Mode&c_ISGID != 0 {
		mode |= os.ModeSetgid
	}
	if fi.h.Mode&c_ISVTX != 0 {
		mode |= os.ModeSticky
	}

	// Set file mode bits; clear perm, setuid, setgid, and sticky bits.
	switch m := os.FileMode(fi.h.Mode) &^ 07777; m {
	case c_ISDIR:
		mode |= os.ModeDir
	case c_ISFIFO:
		mode |= os.ModeNamedPipe
	case c_ISLNK:
		mode |= os.ModeSymlink
	case c_ISBLK:
		mode |= os.ModeDevice
	case c_ISCHR:
		mode |= os.ModeDevice
		mode |= os.ModeCharDevice
	case c_ISSOCK:
		mode |= os.ModeSocket
	}

	return mode
}

// sysStat, if non-nil, populates h from system-dependent fields of fi.
var sysStat func(fi os.FileInfo, h *Header) error

const (
	// Mode constants from the USTAR spec:
	// See http://pubs.opengroup.org/onlinepubs/9699919799/utilities/pax.html#tag_20_92_13_06
	c_ISUID = 04000 // Set uid
	c_ISGID = 02000 // Set gid
	c_ISVTX = 01000 // Save text (sticky bit)

	// Common Unix mode constants; these are not defined in any common tar standard.
	// Header.FileInfo understands these, but FileInfoHeader will never produce these.
	c_ISDIR  = 040000  // Directory
	c_ISFIFO = 010000  // FIFO
	c_ISREG  = 0100000 // Regular file
	c_ISLNK  = 0120000 // Symbolic link
	c_ISBLK  = 060000  // Block special file
	c_ISCHR  = 020000  // Character special file
	c_ISSOCK = 0140000 // Socket
)

// FileInfoHeader creates a partially-populated Header from fi.
// If fi describes a symlink, FileInfoHeader records link as the link target.
// If fi describes a directory, a slash is appended to the name.
//
// Since os.FileInfo's Name method only returns the base name of
// the file it describes, it may be necessary to modify Header.Name
// to provide the full path name of the file.
func FileInfoHeader(fi os.FileInfo, link string) (*Header, error) {
	if fi == nil {
		return nil, errors.New("archive/tar: FileInfo is nil")
	}
	fm := fi.Mode()
	h := &Header{
		Name:    fi.Name(),
		ModTime: fi.ModTime(),
		Mode:    int64(fm.Perm()), // or'd with c_IS* constants later
	}
	switch {
	case fm.IsRegular():
		h.Size = fi.Size()
	case fi.IsDir():
		return nil, fmt.Errorf("archive/tar: directories not supported")
	case fm&os.ModeSymlink != 0:
		return nil, fmt.Errorf("archive/tar: symlinks not supported")
	case fm&os.ModeDevice != 0:
		return nil, fmt.Errorf("archive/tar: devices not supported")
	case fm&os.ModeNamedPipe != 0:
		return nil, fmt.Errorf("archive/tar: fifo not supported")
	case fm&os.ModeSocket != 0:
		return nil, fmt.Errorf("archive/tar: sockets not supported")
	default:
		return nil, fmt.Errorf("archive/tar: unknown file mode %v", fm)
	}
	if fm&os.ModeSetuid != 0 {
		h.Mode |= c_ISUID
	}
	if fm&os.ModeSetgid != 0 {
		h.Mode |= c_ISGID
	}
	if fm&os.ModeSticky != 0 {
		h.Mode |= c_ISVTX
	}
	// If possible, populate additional fields from OS-specific
	// FileInfo fields.
	if sys, ok := fi.Sys().(*Header); ok {
		// This FileInfo came from a Header (not the OS). Use the
		// original Header to populate all remaining fields.
		h.Uid = sys.Uid
		h.Gid = sys.Gid
		h.ModTime = sys.ModTime
	}
	if sysStat != nil {
		return h, sysStat(fi, h)
	}
	return h, nil
}
