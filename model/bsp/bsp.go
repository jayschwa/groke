package bsp

import (
	"bytes"
	"errors"
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

type Edge [2]Vector3

type Texture struct {
	Name   string
	Width  int
	Height int
	data   []byte
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

type Face struct {
	Edges   []Edge
	Front   bool
	Plane   *Plane
	TexInfo TexInfo
}

type Model struct {
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
}

const (
	PlaneAxialX = iota
	PlaneAxialY
	PlaneAxialZ
	PlaneNonAxialX
	PlaneNonAxialY
	PlaneNonAxialZ
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
	{[]byte{29, 0, 0, 0}, q1BSPRead},
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
		if bytes.Compare(id, bh.id) == 0 {
			var m Model
			err = bh.read(r, flags, &m)
			mp = &m
			break
		}
	}

	return
}
