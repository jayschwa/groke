package bsp

import (
	"bytes"
	"errors"
	"image"
	"io"
)

type Vector3 [3]float64
type PlaneType uint8
type TexFlags int

type Plane struct {
	N Vector3
	D float64
	T PlaneType
}

type DataSource interface {
	External() bool
}

type dataSourceExternal struct {
}

func (s dataSourceExternal) External() bool {
	return true
}

type dataSourceInternal struct {
	image.Image
}

func (s dataSourceInternal) External() bool {
	return false
}

type Texture struct {
	DataSource
	Name string
	Next *Texture
}

type TexInfo struct {
	S       Vector3
	T       Vector3
	Ds      float64
	Dt      float64
	Texture *Texture
	Flags   TexFlags
}

type Entity map[string]string

type Vert struct {
	Pos Vector3
}

type Face struct {
	Verts   []Vert
	Front   bool
	Plane   *Plane
	TexInfo TexInfo
}

type Model struct {
	Triangle bool // FIXME -- Q3 uses triangles
	Entities []Entity
	Faces    []Face
	TexInfos []TexInfo
	Textures []Texture
}

type bspReader func(io.Reader, int, *Model) error

type bspLoader struct {
	id   []byte
	read bspReader
}

var pt2String = []string{
	"PlaneAxialX",
	"PlaneAxialY",
	"PlaneAxialZ",
	"PlaneNonAxialX",
	"PlaneNonAxialY",
	"PlaneNonAxialZ",
	"PlaneNoType",
}

const (
	PlaneAxialX = iota
	PlaneAxialY
	PlaneAxialZ
	PlaneNonAxialX
	PlaneNonAxialY
	PlaneNonAxialZ
	PlaneNoType
)

const (
	TexAnimated = TexFlags(1 << iota)
)

const (
	EntitiesOnly = 1 << iota
	NoTextures
	NoLightmaps
)

var bspLoaders = []*bspLoader{
	{[]byte{0x1d, 0x00, 0x00, 0x00}, q1BSPRead},
	{[]byte{0x1e, 0x00, 0x00, 0x00}, hlBSPRead},
	{[]byte{'I', 'B', 'S', 'P', 0x26, 0x00, 0x00, 0x00}, q2BSPRead},
	{[]byte{'I', 'B', 'S', 'P', 0x2e, 0x00, 0x00, 0x00}, q3BSPRead},
}

var (
	ErrFormat = errors.New("bsp: invalid bsp format")
)

func (t PlaneType) String() string {
	return pt2String[t]
}

func Read(r io.Reader, flags int) (mp *Model, err error) {
	id := make([]byte, 4)
	if _, err = r.Read(id); err != nil {
		return
	}

	err = ErrFormat

	for _, bh := range bspLoaders {
		if len(bh.id) > len(id) {
			bigger := make([]byte, len(bh.id))
			copy(bigger, id)
			if _, err := r.Read(bigger[len(id):]); err != nil {
				break
			} else {
				id = bigger
			}
		}

		if bytes.Compare(id, bh.id) == 0 {
			var m Model
			err = bh.read(r, flags, &m)
			mp = &m
			break
		}
	}

	return
}
