/*
Package wal provides support for reading Quake2 wal images.
*/
package wal

import (
	"bytes"
	"encoding/binary"
	"errors"
	"image"
	"image/color"
	"io"
)

var Palette color.Palette

var (
	ErrFormat = errors.New("wal: not a valid wal file")
)

// Decode decodes a WAL image.
func Decode(r io.Reader) (outImage image.Image, err error) {
	if w, h, b, loadErr := load(r); loadErr == nil {
		rect := image.Rect(0, 0, w, h)
		outImage = &image.Paletted{
			Pix:     b,
			Stride:  w,
			Rect:    rect,
			Palette: Palette,
		}
	} else {
		err = loadErr
	}

	return
}

// DecodeConfig decodes a header of WAL image and returns its
// configuration.
func DecodeConfig(r io.Reader) (cfg image.Config, err error) {
	if w, h, _, loadErr := load(r); loadErr == nil {
		cfg.Width = w
		cfg.Height = h
	} else {
		err = loadErr
	}

	return
}

func init() {
	Palette = DefaultPalette
	image.RegisterFormat("wal", "", Decode, DecodeConfig)
}

func load(r io.Reader) (w, h int, b []byte, err error) {
	var data bytes.Buffer

	if _, err = data.ReadFrom(r); err != nil {
		return
	}

	b = data.Bytes()

	if len(b) < 100 {
		err = ErrFormat
		return
	}

	w = int(binary.LittleEndian.Uint32(b[32:]))
	h = int(binary.LittleEndian.Uint32(b[36:]))
	b = b[100:]

	if w*h > len(b) {
		err = ErrFormat
	}

	return
}
