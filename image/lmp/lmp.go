/*
Package lmp provides support for reading Quake lmp images (stored in
WADs).
*/
package lmp

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
	ErrFormat = errors.New("lmp: not a valid lmp file")
)

// Decode decodes a LMP image.
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

// DecodeConfig decodes a header of LMP image and returns its
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
	image.RegisterFormat("lmp", "", Decode, DecodeConfig)
}

func load(r io.Reader) (w, h int, b []byte, err error) {
	var data bytes.Buffer

	if _, err = data.ReadFrom(r); err != nil {
		return
	}

	b = data.Bytes()
	size := len(b)

	if size < 8 {
		err = ErrFormat
		return
	}

	w = int(binary.LittleEndian.Uint32(b))
	h = int(binary.LittleEndian.Uint32(b[4:]))

	if w*h >= size-8 {
		// Quake image with header
		b = b[8 : 8+w*h]
	} else if 128*128 == size {
		// conchars, no header
		w = 128
		h = 128

		// convert all black to transparent
		for i := 0; i < w*h; i++ {
			c := Palette[b[i]].(color.NRGBA)
			if c.R == c.G && c.R == c.B && c.R == 0 {
				b[i] = 255
			}
		}
	} else {
		err = ErrFormat
	}

	return
}
