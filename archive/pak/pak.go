/*
Package pak provides support for reading Quake PAK archives.

Note:

All paths inside a pak archive are converted to lower case and path.Clean is
called on each of them.
*/
package pak

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"
	"io/ioutil"
	"os"
	"path"
)

type File struct {
	Name   string
	Size   uint32
	r      io.ReaderAt
	offset uint32
}

type Reader struct {
	File []*File
}

type ReadCloser struct {
	f *os.File
	Reader
}

const (
	pakEntrySize = 64
)

var (
	ErrFormat = errors.New("pak: not a valid pak file")
)

// NewReader returns a new Reader reading from r, which is assumed to have the
// given size in bytes.
func NewReader(r io.ReaderAt, size int64) (*Reader, error) {
	rc := new(Reader)
	if err := rc.init(r, size); err != nil {
		return nil, err
	}

	return rc, nil
}

// OpenReader will open the pak file specified by name and return a ReadCloser.
func OpenReader(name string) (*ReadCloser, error) {
	f, err := os.Open(name)
	if err != nil {
		return nil, err
	}

	fstat, err := f.Stat()
	if err != nil {
		f.Close()
		return nil, err
	}

	rc := new(ReadCloser)
	if err := rc.init(f, fstat.Size()); err != nil {
		f.Close()
		return nil, err
	}

	rc.f = f

	return rc, nil
}

// Open returns a ReadCloser that provides access to the File's contents.
// Multiple files may be read concurrently.
func (f *File) Open() (io.ReadCloser, error) {
	r := io.NewSectionReader(f.r, int64(f.offset), int64(f.Size))
	return ioutil.NopCloser(r), nil
}

// Close closes the pak file, rendering it unusable for I/O.
func (rc *ReadCloser) Close() error {
	return rc.f.Close()
}

func (rc *Reader) init(r io.ReaderAt, size int64) error {
	rs := io.NewSectionReader(r, 0, size)

	var header struct {
		Id        [4]byte
		DirOffset uint32
		DirSize   uint32
	}

	if err := binary.Read(rs, binary.LittleEndian, &header); err != nil {
		return err
	}

	if header.Id != [4]byte{'P', 'A', 'C', 'K'} {
		return ErrFormat
	}

	if _, err := rs.Seek(int64(header.DirOffset), os.SEEK_SET); err != nil {
		return err
	}

	rc.File = make([]*File, 0, header.DirSize/pakEntrySize)

	for i := 0; i < cap(rc.File); i++ {
		var pakEntry [pakEntrySize]byte

		if _, err := io.ReadFull(rs, pakEntry[:]); err != nil {
			return err
		}

		nameLen := bytes.IndexByte(pakEntry[:], 0)
		if nameLen < 0 || nameLen > 56 {
			nameLen = 56
		}

		name := string(bytes.ToLower(pakEntry[:nameLen]))

		f := &File{
			Name:   path.Clean(name),
			Size:   binary.LittleEndian.Uint32(pakEntry[60:]),
			offset: binary.LittleEndian.Uint32(pakEntry[56:]),
			r:      r,
		}

		rc.File = append(rc.File, f)
	}

	return nil
}
