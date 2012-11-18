/*
Package wal provides support for reading Quake2 wal images.
*/
package wal

import (
	"bytes"
	. "encoding/binary"
	"errors"
	"image"
	"image/color"
	"io"
)

var Palette color.Palette

var (
	ErrFormat = errors.New("wal: not a valid wal file")
)

// Contents is a type for contents flags.
type Contents uint32

// Visible contents.
const (
	ContentsSolid = Contents(1 << iota)
	ContentsWindow
	ContentsAux
	ContentsLava
	ContentsSlime
	ContentsWater
	ContentsMist
)

// Non-visible contents.
const (
	ContentsAreaportal = Contents(1 << (15 + iota))
	ContentsPlayerClip
	ContentsMonsterClip
	ContentsCurrent0
	ContentsCurrent90
	ContentsCurrent180
	ContentsCurrent270
	ContentsCurrentUp
	ContentsCurrentDown
	ContentsOrigin
	ContentsMonster
	ContentsDeadMonster
	ContentsDetail
	ContentsTranslucent
	ContentsLadder
)

var contentsStr = map[Contents]string{
	ContentsSolid:       "Solid",
	ContentsWindow:      "Window",
	ContentsAux:         "Aux",
	ContentsLava:        "Lava",
	ContentsSlime:       "Slime",
	ContentsWater:       "Water",
	ContentsMist:        "Mist",
	ContentsAreaportal:  "Areaportal",
	ContentsPlayerClip:  "PlayerClip",
	ContentsMonsterClip: "MonsterClip",
	ContentsCurrent0:    "Current0",
	ContentsCurrent90:   "Current90",
	ContentsCurrent180:  "Current180",
	ContentsCurrent270:  "Current270",
	ContentsCurrentUp:   "CurrentUp",
	ContentsCurrentDown: "CurrentDown",
	ContentsOrigin:      "Origin",
	ContentsMonster:     "Monster",
	ContentsDeadMonster: "DeadMonster",
	ContentsDetail:      "Detail",
	ContentsTranslucent: "Translucent",
	ContentsLadder:      "Ladder",
}

// SurFlags is a type for surface flags.
type SurfFlags uint32

const (
	// Value holds light strength.
	SurfLight = SurfFlags(1 << iota)
	// Slick surface affecting game physics.
	SurfSlick
	// Skybox.
	SurfSky
	// Turbulent water warp.
	SurfWarp
	// Translucent 33.3%
	SurfTrans33
	// Translucent 66.6%
	SurfTrans66
	// Scrolls towards angle.
	SurfFlowing
	// Do not draw.
	SurfNoDraw
)

var surfFlagsStr = map[SurfFlags]string{
	SurfLight:   "Light",
	SurfSlick:   "Slick",
	SurfSky:     "Sky",
	SurfWarp:    "Warp",
	SurfTrans33: "Trans33",
	SurfTrans66: "Trans66",
	SurfFlowing: "Flowing",
	SurfNoDraw:  "NoDraw",
}

type WAL struct {
	image.Image
	Name     string
	NextName string
	Flags    SurfFlags
	Contents Contents
	Value    uint32
}

// Decode decodes a WAL image.
func Decode(r io.Reader) (outImage image.Image, err error) {
	var data bytes.Buffer

	if _, err = data.ReadFrom(r); err != nil {
		return
	}

	b := data.Bytes()

	if len(b) < 100 {
		err = ErrFormat
		return
	}

	nameLen := bytes.IndexByte(b, 0)
	if nameLen < 0 || nameLen > 32 {
		nameLen = 32
	}
	name := string(b[:nameLen])

	w := int(LittleEndian.Uint32(b[32:]))
	h := int(LittleEndian.Uint32(b[36:]))
	flags := SurfFlags(LittleEndian.Uint32(b[88:]))
	contents := Contents(LittleEndian.Uint32(b[92:]))
	value := LittleEndian.Uint32(b[96:])

	nextLen := bytes.IndexByte(b[56:], 0)
	if nextLen < 0 || nextLen > 32 {
		nextLen = 32
	}
	nextName := string(b[56 : 56+nextLen])

	b = b[100:]

	if w*h > len(b) {
		err = ErrFormat
	} else {
		rect := image.Rect(0, 0, w, h)
		outImage = &WAL{
			&image.Paletted{
				Pix:     b,
				Stride:  w,
				Rect:    rect,
				Palette: Palette,
			},
			name,
			nextName,
			flags,
			contents,
			value,
		}
	}

	return
}

// DecodeConfig decodes a header of WAL image and returns its
// configuration.
func DecodeConfig(r io.Reader) (cfg image.Config, err error) {
	b := make([]byte, 40)
	var n int

	if n, err = r.Read(b); n < len(b) {
		err = ErrFormat
	} else if err == nil {
		cfg.Width = int(LittleEndian.Uint32(b[32:]))
		cfg.Height = int(LittleEndian.Uint32(b[36:]))
	}

	return
}

func init() {
	Palette = DefaultPalette
	image.RegisterFormat("wal", "", Decode, DecodeConfig)
}

func (c Contents) String() (s string) {
	if c == 0 {
		return "0"
	}

	for i := uint32(0); i < 32 && c != 0; i++ {
		m := Contents(1 << i)
		if c&m != 0 {
			c ^= m
			if mStr, ok := contentsStr[m]; ok {
				if s != "" {
					s += "|"
				}

				s += mStr
			}
		}
	}

	return "(" + s + ")"
}

func (c SurfFlags) String() (s string) {
	if c == 0 {
		return "0"
	}

	for i := uint32(0); i < 32 && c != 0; i++ {
		m := SurfFlags(1 << i)
		if c&m != 0 {
			c ^= m
			if mStr, ok := surfFlagsStr[m]; ok {
				if s != "" {
					s += "|"
				}

				s += mStr
			}
		}
	}

	return "(" + s + ")"
}
