/*
Package wad provides support for reading Quake/Half-Life WAD archives.

Note:

All paths inside a wad archive are converted to lower case.
*/
package wad

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"
	"io/ioutil"
	"os"
)

type File struct {
	Name   string
	Size   uint32
	Type   byte
	r      io.ReaderAt
	offset uint32
}

type Reader struct {
	File []*File
	Type WadType
}

type ReadCloser struct {
	f *os.File
	Reader
}

type WadType byte

const (
	wadEntrySize = 32
)

const (
	QuakeWad = WadType(iota)
	HalfLifeWad
)

var (
	ErrFormat = errors.New("wad: not a valid wad file")
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

// OpenReader will open the wad file specified by name and return a ReadCloser.
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

// Close closes the wad file, rendering it unusable for I/O.
func (rc *ReadCloser) Close() error {
	return rc.f.Close()
}

func (rc *Reader) init(r io.ReaderAt, size int64) error {
	rs := io.NewSectionReader(r, 0, size)

	var header struct {
		Id        [4]byte
		NumFiles  uint32
		DirOffset uint32
	}

	if err := binary.Read(rs, binary.LittleEndian, &header); err != nil {
		return err
	}

	if header.Id == [4]byte{'W','A','D','2'} {
		rc.Type = QuakeWad
	} else if header.Id == [4]byte{'W','A','D','3'} {
		rc.Type = HalfLifeWad
	} else {
		return ErrFormat
	}

	if _, err := rs.Seek(int64(header.DirOffset), os.SEEK_SET); err != nil {
		return err
	}

	rc.File = make([]*File, 0, header.NumFiles)

	for i := 0; i < cap(rc.File); i++ {
		var wadEntry [wadEntrySize]byte

		if _, err := io.ReadFull(rs, wadEntry[:]); err != nil {
			return err
		}

		nameLen := bytes.IndexByte(wadEntry[16:], 0)
		if nameLen < 0 || nameLen > 16 {
			nameLen = 16
		}

		name := string(bytes.ToLower(wadEntry[16:16+nameLen]))

		f := &File{
			Name:   name,
			Size:   binary.LittleEndian.Uint32(wadEntry[8:]),
			Type:   wadEntry[12],
			offset: binary.LittleEndian.Uint32(wadEntry[:]),
			r:      r,
		}

		rc.File = append(rc.File, f)
	}

	return nil
}
