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

// Decode decodes a Half-Life image.
func Decode(r io.Reader) (outImage image.Image, err error) {
	var data bytes.Buffer

	if _, err = data.ReadFrom(r); err != nil {
		return
	}

	b := data.Bytes()
	size := len(b)

	if size < 24 {
		err = ErrFormat
		return
	}

	//nameLen := bytes.IndexByte(b[:16], 0)
	//name := string(bytes.ToLower(b[:nameLen]))

	dataOff := int(LittleEndian.Uint32(b[24:]))
	width := int(LittleEndian.Uint32(b[16:]))
	height := int(LittleEndian.Uint32(b[20:]))

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

	palSize := int(LittleEndian.Uint16(b[palOff:]))
	palOff += 2
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

	outImage = &image.Paletted{
		Pix:     b[dataOff : dataOff+width*height],
		Stride:  width,
		Rect:    image.Rect(0, 0, width, height),
		Palette: palette,
	}

	return
}

// DecodeConfig decodes a header of Half-Life image and returns its
// configuration.
func DecodeConfig(r io.Reader) (cfg image.Config, err error) {
	var b [24]byte

	if _, err = r.Read(b[:]); err == nil {
		cfg.Width = int(LittleEndian.Uint32(b[16:]))
		cfg.Height = int(LittleEndian.Uint32(b[20:]))
	}

	return
}

func init() {
	image.RegisterFormat("hltex", "", Decode, DecodeConfig)
}
