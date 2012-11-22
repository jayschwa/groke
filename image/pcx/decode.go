/*
Package pcx provides support for reading Quake2 pcx images.
*/
package pcx

import (
	. "encoding/binary"
	"errors"
	"image"
	"image/color"
	"io"
	"io/ioutil"
)

var (
	ErrFormat = errors.New("pcx: unsupported format")
)

func Decode(r io.Reader) (outImage image.Image, err error) {
	b := make([]byte, 128)

	if _, err = io.ReadFull(r, b); err != nil {
		return
	} else if b[2] != 1 || b[3] != 8 {
		err = ErrFormat
		return
	}

	w := int(LittleEndian.Uint16(b[8:])) + 1
	h := int(LittleEndian.Uint16(b[10:])) + 1

	var p []byte
	if p, err = ioutil.ReadAll(r); len(p) < 769 || p[len(p)-769] != 12 {
		err = ErrFormat
		return
	}

	b = p
	p = p[len(p)-768:]

	palette := make(color.Palette, 0)
	for i := 0; i < 256; i++ {
		o := i * 3
		color := color.NRGBA{p[o+0], p[o+1], p[o+2], 0xff}
		if color.R == 0x9f && color.G == 0x5b && color.B == 0x53 {
			color.A = 0
		}
		palette = append(palette, color)
	}

	pix := make([]byte, w*h)
	i := 0
	for y := 0; y < h; y++ {
		for x := 0; x < w; {
			data := b[i]
			i++

			if data&0xc0 == 0xc0 {
				runLen := data & 0x3f
				data = b[i]
				i++
				for ; runLen > 0; runLen-- {
					pix[x+y*w] = data
					x++
				}
			} else {
				pix[x+y*w] = data
				x++
			}
		}
	}

	outImage = &image.Paletted{
		Pix:     pix,
		Stride:  w,
		Rect:    image.Rect(0, 0, w, h),
		Palette: palette,
	}

	return
}

func DecodeConfig(r io.Reader) (cfg image.Config, err error) {
	b := make([]byte, 8)

	if _, err = io.ReadFull(r, b); err != nil {
		return
	}

	w := int(LittleEndian.Uint16(b[8:]) + 1)
	h := int(LittleEndian.Uint16(b[10:]) + 1)

	cfg = image.Config{
		Width:  w,
		Height: h,
	}

	return
}

func init() {
	image.RegisterFormat("\x0a\x05", "", Decode, DecodeConfig)
}
