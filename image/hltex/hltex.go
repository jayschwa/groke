package hltex

import (
	"bytes"
	. "encoding/binary"
	"errors"
	"image"
	"image/color"
	"io"
)

var (
	ErrFormat = errors.New("hltex: not a valid Half-Life texture")
)

type HLTex struct {
	image.Image
	Name string
}

// Decode decodes a Half-Life image.
func Decode(r io.Reader) (outImage image.Image, err error) {
	b := make([]byte, 40)

	if _, err = io.ReadFull(r, b); err != nil {
		return
	}

	name := string(bytes.ToLower(b[:bytes.IndexByte(b[:16], 0)]))
	width := int(LittleEndian.Uint32(b[16:]))
	height := int(LittleEndian.Uint32(b[20:]))
	dataOff := int(LittleEndian.Uint32(b[24:]))

	if dataOff == 0 {
		outImage = &HLTex{
			&image.Paletted{
				Pix:     nil,
				Stride:  width,
				Rect:    image.Rect(0, 0, width, height),
				Palette: nil,
			},
			name,
		}
		return
	} else if dataOff < len(b) {
		err = ErrFormat
		return
	}

	var palOff int

	if palOff = int(LittleEndian.Uint32(b[36:])); palOff != 0 {
		palOff += width * height / 64
	} else if palOff = int(LittleEndian.Uint32(b[32:])); palOff != 0 {
		palOff += width * height / 16
	} else if palOff = int(LittleEndian.Uint32(b[28:])); palOff != 0 {
		palOff += width * height / 4
	} else {
		palOff = dataOff + width*height
	}

	dataOff -= len(b)
	palOff -= len(b)

	var size int
	b = make([]byte, palOff+2+256*3)
	if size, err = r.Read(b); err != nil {
		return
	} else if size < palOff+2 {
		err = ErrFormat
		return
	}

	palSize := int(LittleEndian.Uint16(b[palOff:]))
	palOff += 2
	if palSize > 256 || size < palOff+palSize {
		err = ErrFormat
		return
	}

	palette := make(color.Palette, 0)
	for i := 0; i < palSize; i++ {
		o := i * 3
		color := color.NRGBA{b[palOff+o+0], b[palOff+o+1], b[palOff+o+2], 0xff}
		if color.R == color.G && color.G == 0 && color.B == 0xff {
			color.A = 0
		}
		palette = append(palette, color)
	}

	for i := palSize; i < 256; i++ {
		palette = append(palette, color.NRGBA{0, 0, 0, 0})
	}

	outImage = &HLTex{
		&image.Paletted{
			Pix:     b[dataOff : dataOff+width*height],
			Stride:  width,
			Rect:    image.Rect(0, 0, width, height),
			Palette: palette,
		},
		name,
	}

	return
}

// DecodeConfig decodes a header of Half-Life image and returns its
// configuration.
func DecodeConfig(r io.Reader) (cfg image.Config, err error) {
	var b [24]byte

	if _, err = io.ReadFull(r, b[:]); err == nil {
		cfg.Width = int(LittleEndian.Uint32(b[16:]))
		cfg.Height = int(LittleEndian.Uint32(b[20:]))
	}

	return
}

func init() {
	image.RegisterFormat("hltex", "", Decode, DecodeConfig)
}
